/*
 * Copyright (c) 2019  InterDigital Communications, Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	sbi "github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-loc-serv/sbi"
	dkm "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr"
	httpLog "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-http-logger"
	clientNotifOMA "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-loc-serv-notification-client"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	redis "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis"

	"github.com/gorilla/mux"
)

const LocServBasePath = "/location/v1/"
const locServKey string = "loc-serv:"
const logModuleLocServ string = "meep-loc-serv"

const typeZone = "zone"
const typeAccessPoint = "accessPoint"
const typeUser = "user"
const typeZonalSubscription = "zonalsubs"
const typeUserSubscription = "usersubs"
const typeZoneStatusSubscription = "zonestatus"

const USER_TRACKING_AND_ZONAL_TRAFFIC = 1
const ZONE_STATUS = 2

type UeUserData struct {
	queryZoneId string
	queryApId   string
	userList    *UserList
}

type ApUserData struct {
	queryInterestRealm string
	apList             *AccessPointList
}

var nextZonalSubscriptionIdAvailable int
var nextUserSubscriptionIdAvailable int
var nextZoneStatusSubscriptionIdAvailable int

var zonalSubscriptionEnteringMap = map[int]string{}
var zonalSubscriptionLeavingMap = map[int]string{}
var zonalSubscriptionTransferringMap = map[int]string{}
var zonalSubscriptionMap = map[int]string{}

var userSubscriptionEnteringMap = map[int]string{}
var userSubscriptionLeavingMap = map[int]string{}
var userSubscriptionTransferringMap = map[int]string{}
var userSubscriptionMap = map[int]string{}

var zoneStatusSubscriptionMap = map[int]*ZoneStatusCheck{}

type ZoneStatusCheck struct {
	ZoneId                 string
	Serviceable            bool
	Unserviceable          bool
	Unknown                bool
	NbUsersInZoneThreshold int
	NbUsersInAPThreshold   int
}

var LOC_SERV_DB = 0
var currentStoreName = ""

var redisAddr string = "meep-redis-master.default.svc.cluster.local:6379"
var influxAddr string = "http://meep-influxdb.default.svc.cluster.local:8086"
var postgisHost string = "meep-postgis.default.svc.cluster.local"
var postgisPort string = "5432"

var rc *redis.Connector
var hostUrl *url.URL
var sandboxName string
var basePath string
var baseKey string

// Init - Location Service initialization
func Init() (err error) {

	sandboxNameEnv := strings.TrimSpace(os.Getenv("MEEP_SANDBOX_NAME"))
	if sandboxNameEnv != "" {
		sandboxName = sandboxNameEnv
	}
	if sandboxName == "" {
		err = errors.New("MEEP_SANDBOX_NAME env variable not set")
		log.Error(err.Error())
		return err
	}
	log.Info("MEEP_SANDBOX_NAME: ", sandboxName)

	// Retrieve Host URL from environment variable
	hostUrl, err = url.Parse(strings.TrimSpace(os.Getenv("MEEP_HOST_URL")))
	if err != nil {
		hostUrl = new(url.URL)
	}
	log.Info("MEEP_HOST_URL: ", hostUrl)

	// Set base path
	basePath = "/" + sandboxName + LocServBasePath

	// Get base storage key
	baseKey = dkm.GetKeyRoot(sandboxName) + locServKey

	// Connect to Redis DB
	rc, err = redis.NewConnector(redisAddr, LOC_SERV_DB)
	if err != nil {
		log.Error("Failed connection to Redis DB. Error: ", err)
		return err
	}
	log.Info("Connected to Redis DB, location service table")

	userTrackingReInit()
	zonalTrafficReInit()
	zoneStatusReInit()

	// Initialize SBI
	sbiCfg := sbi.SbiCfg{
		SandboxName:    sandboxName,
		RedisAddr:      redisAddr,
		PostgisHost:    postgisHost,
		PostgisPort:    postgisPort,
		UserInfoCb:     updateUserInfo,
		ZoneInfoCb:     updateZoneInfo,
		ApInfoCb:       updateAccessPointInfo,
		ScenarioNameCb: updateStoreName,
		CleanUpCb:      cleanUp,
	}
	err = sbi.Init(sbiCfg)
	if err != nil {
		log.Error("Failed initialize SBI. Error: ", err)
		return err
	}
	log.Info("SBI Initialized")

	return nil
}

// Run - Start Location Service
func Run() (err error) {
	return sbi.Run()
}

// Stop - Stop RNIS
func Stop() (err error) {
	return sbi.Stop()
}

func createClient(notifyPath string) (*clientNotifOMA.APIClient, error) {
	// Create & store client for App REST API
	subsAppClientCfg := clientNotifOMA.NewConfiguration()
	subsAppClientCfg.BasePath = notifyPath
	subsAppClient := clientNotifOMA.NewAPIClient(subsAppClientCfg)
	if subsAppClient == nil {
		log.Error("Failed to create Subscription App REST API client: ", subsAppClientCfg.BasePath)
		err := errors.New("Failed to create Subscription App REST API client")
		return nil, err
	}
	return subsAppClient, nil
}

func deregisterZoneStatus(subsIdStr string) {
	subsId, _ := strconv.Atoi(subsIdStr)
	zonalSubscriptionMap[subsId] = ""
}

func registerZoneStatus(zoneId string, nbOfUsersZoneThreshold int32, nbOfUsersAPThreshold int32, opStatus []OperationStatus, subsIdStr string) {

	subsId, _ := strconv.Atoi(subsIdStr)

	var zoneStatus ZoneStatusCheck
	if opStatus != nil {
		for i := 0; i < len(opStatus); i++ {
			switch opStatus[i] {
			case SERVICEABLE:
				zoneStatus.Serviceable = true
			case UNSERVICEABLE:
				zoneStatus.Unserviceable = true
			case OPSTATUS_UNKNOWN:
				zoneStatus.Unknown = true
			default:
			}
		}
	}
	zoneStatus.NbUsersInZoneThreshold = (int)(nbOfUsersZoneThreshold)
	zoneStatus.NbUsersInAPThreshold = (int)(nbOfUsersAPThreshold)
	zoneStatus.ZoneId = zoneId

	zoneStatusSubscriptionMap[subsId] = &zoneStatus
}

func deregisterZonal(subsIdStr string) {
	subsId, _ := strconv.Atoi(subsIdStr)
	zonalSubscriptionMap[subsId] = ""
	zonalSubscriptionEnteringMap[subsId] = ""
	zonalSubscriptionLeavingMap[subsId] = ""
	zonalSubscriptionTransferringMap[subsId] = ""
}

func registerZonal(zoneId string, event []UserEventType, subsIdStr string) {

	subsId, _ := strconv.Atoi(subsIdStr)

	if event != nil {
		for i := 0; i < len(event); i++ {
			switch event[i] {
			case ENTERING:
				zonalSubscriptionEnteringMap[subsId] = zoneId
			case LEAVING:
				zonalSubscriptionLeavingMap[subsId] = zoneId
			case TRANSFERRING:
				zonalSubscriptionTransferringMap[subsId] = zoneId
			default:
			}
		}
	} else {
		zonalSubscriptionEnteringMap[subsId] = zoneId
		zonalSubscriptionLeavingMap[subsId] = zoneId
		zonalSubscriptionTransferringMap[subsId] = zoneId
	}
	zonalSubscriptionMap[subsId] = zoneId
}

func deregisterUser(subsIdStr string) {
	subsId, _ := strconv.Atoi(subsIdStr)
	userSubscriptionMap[subsId] = ""
	userSubscriptionEnteringMap[subsId] = ""
	userSubscriptionLeavingMap[subsId] = ""
	userSubscriptionTransferringMap[subsId] = ""
}

func registerUser(userAddress string, event []UserEventType, subsIdStr string) {

	subsId, _ := strconv.Atoi(subsIdStr)

	if event != nil {
		for i := 0; i < len(event); i++ {
			switch event[i] {
			case ENTERING:
				userSubscriptionEnteringMap[subsId] = userAddress
			case LEAVING:
				userSubscriptionLeavingMap[subsId] = userAddress
			case TRANSFERRING:
				userSubscriptionTransferringMap[subsId] = userAddress
			default:
			}
		}
	} else {
		userSubscriptionEnteringMap[subsId] = userAddress
		userSubscriptionLeavingMap[subsId] = userAddress
		userSubscriptionTransferringMap[subsId] = userAddress
	}
	userSubscriptionMap[subsId] = userAddress
}

func checkNotificationRegistrations(checkType int, param1 string, param2 string, param3 string, param4 string, param5 string) {

	switch checkType {
	case USER_TRACKING_AND_ZONAL_TRAFFIC:
		//params are the following => newZoneId:oldZoneId:newAccessPointId:oldAccessPointId:userAddress
		checkNotificationRegisteredUsers(param1, param2, param3, param4, param5)
		checkNotificationRegisteredZones(param1, param2, param3, param4, param5)
	case ZONE_STATUS:
		//params are the following => zoneId:accessPointId:nbUsersInAP:nbUsersInZone
		checkNotificationRegisteredZoneStatus(param1, param2, param3, param4)
	default:
	}
}

func checkNotificationRegisteredZoneStatus(zoneId string, apId string, nbUsersInAPStr string, nbUsersInZoneStr string) {

	//check all that applies
	for subsId, zoneStatus := range zoneStatusSubscriptionMap {
		if zoneStatus.ZoneId == zoneId {

			nbUsersInZone := 0
			nbUsersInAP := -1
			zoneWarning := false
			apWarning := false
			if nbUsersInZoneStr != "" {
				nbUsersInZone, _ = strconv.Atoi(nbUsersInZoneStr)
				if nbUsersInZone >= zoneStatus.NbUsersInZoneThreshold {
					zoneWarning = true
				}
			}
			if nbUsersInAPStr != "" {
				nbUsersInAP, _ = strconv.Atoi(nbUsersInAPStr)
				if nbUsersInAP >= zoneStatus.NbUsersInAPThreshold {
					apWarning = true
				}
			}

			if zoneWarning || apWarning {
				subsIdStr := strconv.Itoa(subsId)
				jsonInfo, _ := rc.JSONGetEntry(baseKey+typeZoneStatusSubscription+":"+subsIdStr, ".")
				if jsonInfo == "" {
					return
				}

				subscription := convertJsonToZoneStatusSubscription(jsonInfo)

				var zoneStatusNotif clientNotifOMA.ZoneStatusNotification
				zoneStatusNotif.ZoneId = zoneId
				if apWarning {
					zoneStatusNotif.AccessPointId = apId
					zoneStatusNotif.NumberOfUsersInAP = (int32)(nbUsersInAP)
				}
				if zoneWarning {
					zoneStatusNotif.NumberOfUsersInZone = (int32)(nbUsersInZone)
				}
				zoneStatusNotif.Timestamp = time.Now()
				go sendStatusNotification(subscription.CallbackReference.NotifyURL, context.TODO(), subsIdStr, zoneStatusNotif)
				if apWarning {
					log.Info("Zone Status Notification" + "(" + subsIdStr + "): " + "For event in zone " + zoneId + " which has " + nbUsersInAPStr + " users in AP " + apId)
				} else {
					log.Info("Zone Status Notification" + "(" + subsIdStr + "): " + "For event in zone " + zoneId + " which has " + nbUsersInZoneStr + " users in total")
				}
			}

		}
	}
}

func checkNotificationRegisteredUsers(oldZoneId string, newZoneId string, oldApId string, newApId string, userId string) {

	//check all that applies
	for subsId, value := range userSubscriptionMap {
		if value == userId {

			subsIdStr := strconv.Itoa(subsId)
			jsonInfo, _ := rc.JSONGetEntry(baseKey+typeUserSubscription+":"+subsIdStr, ".")
			if jsonInfo == "" {
				return
			}

			subscription := convertJsonToUserSubscription(jsonInfo)

			var zonal clientNotifOMA.TrackingNotification
			zonal.Address = userId
			zonal.Timestamp = time.Now()

			zonal.CallbackData = subscription.ClientCorrelator

			if newZoneId != oldZoneId {
				if userSubscriptionEnteringMap[subsId] != "" && newZoneId != "" {
					zonal.ZoneId = newZoneId
					zonal.CurrentAccessPointId = newApId
					event := new(clientNotifOMA.UserEventType)
					*event = clientNotifOMA.ENTERING_UserEventType
					zonal.UserEventType = event
					go sendNotification(subscription.CallbackReference.NotifyURL, context.TODO(), subsIdStr, zonal)
					log.Info("User Notification" + "(" + subsIdStr + "): " + "Entering event in zone " + newZoneId + " for user " + userId)
				}
				if oldZoneId != "" {
					if userSubscriptionLeavingMap[subsId] != "" {
						zonal.ZoneId = oldZoneId
						zonal.CurrentAccessPointId = oldApId
						event := new(clientNotifOMA.UserEventType)
						*event = clientNotifOMA.LEAVING_UserEventType
						zonal.UserEventType = event
						go sendNotification(subscription.CallbackReference.NotifyURL, context.TODO(), subsIdStr, zonal)
						log.Info("User Notification" + "(" + subsIdStr + "): " + "Leaving event in zone " + oldZoneId + " for user " + userId)
					}
				}
			} else {
				if newApId != oldApId {
					if userSubscriptionTransferringMap[subsId] != "" {
						zonal.ZoneId = newZoneId
						zonal.CurrentAccessPointId = newApId
						zonal.PreviousAccessPointId = oldApId
						event := new(clientNotifOMA.UserEventType)
						*event = clientNotifOMA.TRANSFERRING_UserEventType
						zonal.UserEventType = event
						go sendNotification(subscription.CallbackReference.NotifyURL, context.TODO(), subsIdStr, zonal)
						log.Info("User Notification" + "(" + subsIdStr + "): " + " Transferring event within zone " + newZoneId + " for user " + userId + " from Ap " + oldApId + " to " + newApId)
					}
				}
			}
		}
	}
}

func sendNotification(notifyUrl string, ctx context.Context, subscriptionId string, notification clientNotifOMA.TrackingNotification) {
	startTime := time.Now()

	client, err := createClient(notifyUrl)
	if err != nil {
		log.Error(err)
		return
	}

	jsonNotif, err := json.Marshal(notification)
	if err != nil {
		log.Error(err.Error())
	}

	resp, err := client.NotificationsApi.PostTrackingNotification(ctx, subscriptionId, notification)
	_ = httpLog.LogTx(notifyUrl, "POST", string(jsonNotif), resp, startTime)
	if err != nil {
		log.Error(err)
		return
	}
	defer resp.Body.Close()
}

func sendStatusNotification(notifyUrl string, ctx context.Context, subscriptionId string, notification clientNotifOMA.ZoneStatusNotification) {
	startTime := time.Now()

	client, err := createClient(notifyUrl)
	if err != nil {
		log.Error(err)
		return
	}

	jsonNotif, err := json.Marshal(notification)
	if err != nil {
		log.Error(err.Error())
	}

	resp, err := client.NotificationsApi.PostZoneStatusNotification(ctx, subscriptionId, notification)
	_ = httpLog.LogTx(notifyUrl, "POST", string(jsonNotif), resp, startTime)
	if err != nil {
		log.Error(err)
		return
	}
	defer resp.Body.Close()
}

func checkNotificationRegisteredZones(oldZoneId string, newZoneId string, oldApId string, newApId string, userId string) {

	//check all that applies
	for subsId, value := range zonalSubscriptionMap {

		if value == newZoneId {

			if newZoneId != oldZoneId {

				if zonalSubscriptionEnteringMap[subsId] != "" {
					subsIdStr := strconv.Itoa(subsId)

					jsonInfo, _ := rc.JSONGetEntry(baseKey+typeZonalSubscription+":"+subsIdStr, ".")
					if jsonInfo != "" {
						subscription := convertJsonToZonalSubscription(jsonInfo)

						var zonal clientNotifOMA.TrackingNotification
						zonal.ZoneId = newZoneId
						zonal.CurrentAccessPointId = newApId
						zonal.Address = userId
						event := new(clientNotifOMA.UserEventType)
						*event = clientNotifOMA.ENTERING_UserEventType
						zonal.UserEventType = event
						zonal.Timestamp = time.Now()
						zonal.CallbackData = subscription.ClientCorrelator
						go sendNotification(subscription.CallbackReference.NotifyURL, context.TODO(), subsIdStr, zonal)
						log.Info("Zonal Notify Entering event in zone " + newZoneId + " for user " + userId)
					}
				}
			} else {
				if newApId != oldApId {
					if zonalSubscriptionTransferringMap[subsId] != "" {
						subsIdStr := strconv.Itoa(subsId)

						jsonInfo, _ := rc.JSONGetEntry(baseKey+typeZonalSubscription+":"+subsIdStr, ".")
						if jsonInfo != "" {
							subscription := convertJsonToZonalSubscription(jsonInfo)

							var zonal clientNotifOMA.TrackingNotification
							zonal.ZoneId = newZoneId
							zonal.CurrentAccessPointId = newApId
							zonal.PreviousAccessPointId = oldApId
							zonal.Address = userId
							event := new(clientNotifOMA.UserEventType)
							*event = clientNotifOMA.TRANSFERRING_UserEventType
							zonal.UserEventType = event
							zonal.Timestamp = time.Now()
							zonal.CallbackData = subscription.ClientCorrelator
							go sendNotification(subscription.CallbackReference.NotifyURL, context.TODO(), subsIdStr, zonal)
							log.Info("Zonal Notify Transferring event in zone " + newZoneId + " for user " + userId + " from Ap " + oldApId + " to " + newApId)
						}
					}
				}
			}
		} else {
			if value == oldZoneId {
				if zonalSubscriptionLeavingMap[subsId] != "" {
					subsIdStr := strconv.Itoa(subsId)

					jsonInfo, _ := rc.JSONGetEntry(baseKey+typeZonalSubscription+":"+subsIdStr, ".")
					if jsonInfo != "" {

						subscription := convertJsonToZonalSubscription(jsonInfo)

						var zonal clientNotifOMA.TrackingNotification
						zonal.ZoneId = oldZoneId
						zonal.CurrentAccessPointId = oldApId
						zonal.Address = userId
						event := new(clientNotifOMA.UserEventType)
						*event = clientNotifOMA.LEAVING_UserEventType
						zonal.UserEventType = event
						zonal.Timestamp = time.Now()
						zonal.CallbackData = subscription.ClientCorrelator
						go sendNotification(subscription.CallbackReference.NotifyURL, context.TODO(), subsIdStr, zonal)
						log.Info("Zonal Notify Leaving event in zone " + oldZoneId + " for user " + userId)
					}
				}
			}
		}
	}
}

func usersGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var userData UeUserData

	// Retrieve query parameters
	u, _ := url.Parse(r.URL.String())
	log.Info("url: ", u.RequestURI())
	q := u.Query()
	userData.queryZoneId = q.Get("zoneId")
	userData.queryApId = q.Get("accessPointId")

	// Get user list from DB
	var response ResponseUserList
	var userList UserList
	userList.ResourceURL = hostUrl.String() + basePath + "users"
	response.UserList = &userList
	userData.userList = &userList

	keyName := baseKey + typeUser + ":*"
	err := rc.ForEachJSONEntry(keyName, populateUserList, &userData)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func populateUserList(key string, jsonInfo string, userData interface{}) error {
	// Get query params & userlist from user data
	data := userData.(*UeUserData)
	if data == nil || data.userList == nil {
		return errors.New("userList not found in userData")
	}

	// Retrieve user info from DB
	var userInfo UserInfo
	err := json.Unmarshal([]byte(jsonInfo), &userInfo)
	if err != nil {
		return err
	}

	// Ignore entries with no zoneID or AP ID
	if userInfo.ZoneId == "" || userInfo.AccessPointId == "" {
		return nil
	}

	// Filter using query params
	if data.queryZoneId != "" && userInfo.ZoneId != data.queryZoneId {
		return nil
	}
	if data.queryApId != "" && userInfo.AccessPointId != data.queryApId {
		return nil
	}

	// Add user info to list
	data.userList.User = append(data.userList.User, userInfo)
	return nil
}

func usersGetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	var response ResponseUserInfo
	var userInfo UserInfo
	response.UserInfo = &userInfo

	jsonUserInfo, _ := rc.JSONGetEntry(baseKey+typeUser+":"+vars["userId"], ".")
	if jsonUserInfo == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := json.Unmarshal([]byte(jsonUserInfo), &userInfo)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func zonesByIdGetAps(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var userData ApUserData
	vars := mux.Vars(r)

	// Retrieve query parameters
	u, _ := url.Parse(r.URL.String())
	log.Info("url: ", u.RequestURI())
	q := u.Query()
	userData.queryInterestRealm = q.Get("interestRealm")

	// Get user list from DB
	var response ResponseAccessPointList
	var apList AccessPointList
	apList.ZoneId = vars["zoneId"]
	apList.ResourceURL = hostUrl.String() + basePath + "zones/" + vars["zoneId"] + "/accessPoints"
	response.AccessPointList = &apList
	userData.apList = &apList

	keyName := baseKey + typeZone + ":" + vars["zoneId"] + ":*"
	err := rc.ForEachJSONEntry(keyName, populateApList, &userData)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func zonesByIdGetApsById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	var response ResponseAccessPointInfo
	var apInfo AccessPointInfo
	response.AccessPointInfo = &apInfo

	jsonApInfo, _ := rc.JSONGetEntry(baseKey+typeZone+":"+vars["zoneId"]+":"+typeAccessPoint+":"+vars["accessPointId"], ".")
	if jsonApInfo == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := json.Unmarshal([]byte(jsonApInfo), &apInfo)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func zonesGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response ResponseZoneList
	var zoneList ZoneList
	zoneList.ResourceURL = hostUrl.String() + basePath + "zones"
	response.ZoneList = &zoneList

	keyName := baseKey + typeZone + ":*"
	err := rc.ForEachJSONEntry(keyName, populateZoneList, &zoneList)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func zonesGetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	var response ResponseZoneInfo
	var zoneInfo ZoneInfo
	response.ZoneInfo = &zoneInfo

	jsonZoneInfo, _ := rc.JSONGetEntry(baseKey+typeZone+":"+vars["zoneId"], ".")
	if jsonZoneInfo == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := json.Unmarshal([]byte(jsonZoneInfo), &zoneInfo)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func populateZoneList(key string, jsonInfo string, userData interface{}) error {

	zoneList := userData.(*ZoneList)
	var zoneInfo ZoneInfo

	// Format response
	err := json.Unmarshal([]byte(jsonInfo), &zoneInfo)
	if err != nil {
		return err
	}
	if zoneInfo.ZoneId != "" {
		zoneList.Zone = append(zoneList.Zone, zoneInfo)
	}
	return nil
}

func populateApList(key string, jsonInfo string, userData interface{}) error {
	// Get query params & aplist from user data
	data := userData.(*ApUserData)
	if data == nil || data.apList == nil {
		return errors.New("apList not found in userData")
	}

	// Retrieve AP info from DB
	var apInfo AccessPointInfo
	err := json.Unmarshal([]byte(jsonInfo), &apInfo)
	if err != nil {
		return err
	}

	// Ignore entries with no AP ID
	if apInfo.AccessPointId == "" {
		return nil
	}

	// Filter using query params
	if data.queryInterestRealm != "" && apInfo.InterestRealm != data.queryInterestRealm {
		return nil
	}

	// Add AP info to list
	data.apList.AccessPoint = append(data.apList.AccessPoint, apInfo)
	return nil
}

func userTrackingSubDelById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	err := rc.JSONDelEntry(baseKey+typeUserSubscription+":"+vars["subscriptionId"], ".")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	deregisterUser(vars["subscriptionId"])
	w.WriteHeader(http.StatusNoContent)
}

func userTrackingSubGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response ResponseUserTrackingNotificationSubscriptionList
	var userTrackingSubList UserTrackingNotificationSubscriptionList
	userTrackingSubList.ResourceURL = hostUrl.String() + basePath + "subscriptions/userTracking"
	response.NotificationSubscriptionList = &userTrackingSubList

	keyName := baseKey + typeUserSubscription + "*"
	err := rc.ForEachJSONEntry(keyName, populateUserTrackingList, &userTrackingSubList)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func userTrackingSubGetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	var response ResponseUserTrackingSubscription
	var userTrackingSub UserTrackingSubscription
	response.UserTrackingSubscription = &userTrackingSub

	jsonUserTrackingSub, _ := rc.JSONGetEntry(baseKey+typeUserSubscription+":"+vars["subscriptionId"], ".")
	if jsonUserTrackingSub == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := json.Unmarshal([]byte(jsonUserTrackingSub), &userTrackingSub)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func userTrackingSubPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response ResponseUserTrackingSubscription
	userTrackingSub := new(UserTrackingSubscription)
	response.UserTrackingSubscription = userTrackingSub

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userTrackingSub)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newSubsId := nextUserSubscriptionIdAvailable
	nextUserSubscriptionIdAvailable++
	subsIdStr := strconv.Itoa(newSubsId)

	registerUser(userTrackingSub.Address, userTrackingSub.UserEventCriteria, subsIdStr)
	userTrackingSub.ResourceURL = hostUrl.String() + basePath + "subscriptions/userTracking/" + subsIdStr

	_ = rc.JSONSetEntry(baseKey+typeUserSubscription+":"+subsIdStr, ".", convertUserSubscriptionToJson(userTrackingSub))

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(jsonResponse))
}

func userTrackingSubPutById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	var response ResponseUserTrackingSubscription
	userTrackingSub := new(UserTrackingSubscription)
	response.UserTrackingSubscription = userTrackingSub

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userTrackingSub)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	subsIdStr := vars["subscriptionId"]
	userTrackingSub.ResourceURL = hostUrl.String() + basePath + "subscriptions/userTracking/" + subsIdStr

	_ = rc.JSONSetEntry(baseKey+typeUserSubscription+":"+subsIdStr, ".", convertUserSubscriptionToJson(userTrackingSub))

	deregisterUser(subsIdStr)
	registerUser(userTrackingSub.Address, userTrackingSub.UserEventCriteria, subsIdStr)

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func populateUserTrackingList(key string, jsonInfo string, userData interface{}) error {

	userList := userData.(*UserTrackingNotificationSubscriptionList)
	var userInfo UserTrackingSubscription

	// Format response
	err := json.Unmarshal([]byte(jsonInfo), &userInfo)
	if err != nil {
		return err
	}
	userList.UserTrackingSubscription = append(userList.UserTrackingSubscription, userInfo)
	return nil
}

func zonalTrafficSubDelById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	err := rc.JSONDelEntry(baseKey+typeZonalSubscription+":"+vars["subscriptionId"], ".")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	deregisterZonal(vars["subscriptionId"])
	w.WriteHeader(http.StatusNoContent)
}

func zonalTrafficSubGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response ResponseZonalTrafficNotificationSubscriptionList
	var zonalTrafficSubList ZonalTrafficNotificationSubscriptionList
	zonalTrafficSubList.ResourceURL = hostUrl.String() + basePath + "subscriptions/zonalTraffic"
	response.NotificationSubscriptionList = &zonalTrafficSubList

	keyName := baseKey + typeZonalSubscription + "*"
	err := rc.ForEachJSONEntry(keyName, populateZonalTrafficList, &zonalTrafficSubList)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func zonalTrafficSubGetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	var response ResponseZonalTrafficSubscription
	var zonalTrafficSub ZonalTrafficSubscription
	response.ZonalTrafficSubscription = &zonalTrafficSub

	jsonZonalTrafficSub, _ := rc.JSONGetEntry(baseKey+typeZonalSubscription+":"+vars["subscriptionId"], ".")
	if jsonZonalTrafficSub == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := json.Unmarshal([]byte(jsonZonalTrafficSub), &zonalTrafficSub)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func zonalTrafficSubPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response ResponseZonalTrafficSubscription
	zonalTrafficSub := new(ZonalTrafficSubscription)
	response.ZonalTrafficSubscription = zonalTrafficSub

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&zonalTrafficSub)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newSubsId := nextZonalSubscriptionIdAvailable
	nextZonalSubscriptionIdAvailable++
	subsIdStr := strconv.Itoa(newSubsId)
	/*
		if zonalTrafficSub.Duration > 0 {
			//TODO start a timer mecanism and expire subscription
		}
		//else, lasts forever or until subscription is deleted
	*/
	if zonalTrafficSub.Duration != "" && zonalTrafficSub.Duration != "0" {
		//TODO start a timer mecanism and expire subscription
		log.Info("Non zero duration")
	}
	//else, lasts forever or until subscription is deleted

	zonalTrafficSub.ResourceURL = hostUrl.String() + basePath + "subscriptions/zonalTraffic/" + subsIdStr

	_ = rc.JSONSetEntry(baseKey+typeZonalSubscription+":"+subsIdStr, ".", convertZonalSubscriptionToJson(zonalTrafficSub))

	registerZonal(zonalTrafficSub.ZoneId, zonalTrafficSub.UserEventCriteria, subsIdStr)

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(jsonResponse))
}

func zonalTrafficSubPutById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	var response ResponseZonalTrafficSubscription
	zonalTrafficSub := new(ZonalTrafficSubscription)
	response.ZonalTrafficSubscription = zonalTrafficSub

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&zonalTrafficSub)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	subsIdStr := vars["subscriptionId"]
	zonalTrafficSub.ResourceURL = hostUrl.String() + basePath + "subscriptions/zonalTraffic/" + subsIdStr

	_ = rc.JSONSetEntry(baseKey+typeZonalSubscription+":"+subsIdStr, ".", convertZonalSubscriptionToJson(zonalTrafficSub))

	deregisterZonal(subsIdStr)
	registerZonal(zonalTrafficSub.ZoneId, zonalTrafficSub.UserEventCriteria, subsIdStr)

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func populateZonalTrafficList(key string, jsonInfo string, userData interface{}) error {

	zoneList := userData.(*ZonalTrafficNotificationSubscriptionList)
	var zoneInfo ZonalTrafficSubscription

	// Format response
	err := json.Unmarshal([]byte(jsonInfo), &zoneInfo)
	if err != nil {
		return err
	}
	zoneList.ZonalTrafficSubscription = append(zoneList.ZonalTrafficSubscription, zoneInfo)
	return nil
}

func zoneStatusDelById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	err := rc.JSONDelEntry(baseKey+typeZoneStatusSubscription+":"+vars["subscriptionId"], ".")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	deregisterZoneStatus(vars["subscriptionId"])
	w.WriteHeader(http.StatusNoContent)
}

func zoneStatusGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response ResponseZoneStatusNotificationSubscriptionList
	var zoneStatusSubList ZoneStatusNotificationSubscriptionList
	zoneStatusSubList.ResourceURL = hostUrl.String() + basePath + "subscriptions/zoneStatus"
	response.NotificationSubscriptionList = &zoneStatusSubList

	keyName := baseKey + typeZoneStatusSubscription + "*"
	err := rc.ForEachJSONEntry(keyName, populateZoneStatusList, &zoneStatusSubList)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func zoneStatusGetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	var response ResponseZoneStatusSubscription2
	var zoneStatusSub ZoneStatusSubscription
	response.ZoneStatusSubscription = &zoneStatusSub

	jsonZoneStatusSub, _ := rc.JSONGetEntry(baseKey+typeZoneStatusSubscription+":"+vars["subscriptionId"], ".")
	if jsonZoneStatusSub == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := json.Unmarshal([]byte(jsonZoneStatusSub), &zoneStatusSub)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func zoneStatusPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var response ResponseZoneStatusSubscription2
	zoneStatusSub := new(ZoneStatusSubscription)
	response.ZoneStatusSubscription = zoneStatusSub

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&zoneStatusSub)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newSubsId := nextZoneStatusSubscriptionIdAvailable
	nextZoneStatusSubscriptionIdAvailable++
	subsIdStr := strconv.Itoa(newSubsId)

	zoneStatusSub.ResourceURL = hostUrl.String() + basePath + "subscriptions/zoneStatus/" + subsIdStr

	_ = rc.JSONSetEntry(baseKey+typeZoneStatusSubscription+":"+subsIdStr, ".", convertZoneStatusSubscriptionToJson(zoneStatusSub))

	registerZoneStatus(zoneStatusSub.ZoneId, zoneStatusSub.NumberOfUsersZoneThreshold, zoneStatusSub.NumberOfUsersAPThreshold,
		zoneStatusSub.OperationStatus, subsIdStr)

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(jsonResponse))
}

func zoneStatusPutById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)

	var response ResponseZoneStatusSubscription2
	zoneStatusSub := new(ZoneStatusSubscription)
	response.ZoneStatusSubscription = zoneStatusSub

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&zoneStatusSub)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	subsIdStr := vars["subscriptionId"]
	zoneStatusSub.ResourceURL = hostUrl.String() + basePath + "subscriptions/zoneStatus/" + subsIdStr

	_ = rc.JSONSetEntry(baseKey+typeZoneStatusSubscription+":"+subsIdStr, ".", convertZoneStatusSubscriptionToJson(zoneStatusSub))

	deregisterZoneStatus(subsIdStr)
	registerZoneStatus(zoneStatusSub.ZoneId, zoneStatusSub.NumberOfUsersZoneThreshold, zoneStatusSub.NumberOfUsersAPThreshold,
		zoneStatusSub.OperationStatus, subsIdStr)

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(jsonResponse))
}

func populateZoneStatusList(key string, jsonInfo string, userData interface{}) error {

	zoneList := userData.(*ZoneStatusNotificationSubscriptionList)
	var zoneInfo ZoneStatusSubscription

	// Format response
	err := json.Unmarshal([]byte(jsonInfo), &zoneInfo)
	if err != nil {
		return err
	}
	zoneList.ZoneStatusSubscription = append(zoneList.ZoneStatusSubscription, zoneInfo)
	return nil
}

func cleanUp() {
	log.Info("Terminate all")
	rc.DBFlush(baseKey)
	nextZonalSubscriptionIdAvailable = 1
	nextUserSubscriptionIdAvailable = 1
	nextZoneStatusSubscriptionIdAvailable = 1

	zonalSubscriptionEnteringMap = map[int]string{}
	zonalSubscriptionLeavingMap = map[int]string{}
	zonalSubscriptionTransferringMap = map[int]string{}
	zonalSubscriptionMap = map[int]string{}

	userSubscriptionEnteringMap = map[int]string{}
	userSubscriptionLeavingMap = map[int]string{}
	userSubscriptionTransferringMap = map[int]string{}
	userSubscriptionMap = map[int]string{}

	zoneStatusSubscriptionMap = map[int]*ZoneStatusCheck{}

	updateStoreName("")
}

func updateStoreName(storeName string) {
	if currentStoreName != storeName {
		currentStoreName = storeName
		_ = httpLog.ReInit(logModuleLocServ, sandboxName, storeName, redisAddr, influxAddr)
	}
}

func updateUserInfo(address string, zoneId string, accessPointId string, longitude *float32, latitude *float32) {
	var oldZoneId string
	var oldApId string

	// Get User Info from DB
	jsonUserInfo, _ := rc.JSONGetEntry(baseKey+typeUser+":"+address, ".")
	userInfo := convertJsonToUserInfo(jsonUserInfo)

	// Create new user info if necessary
	if userInfo == nil {
		userInfo = new(UserInfo)
		userInfo.Address = address
		userInfo.ResourceURL = hostUrl.String() + basePath + "users/" + address
	} else {
		// Get old zone & AP IDs
		oldZoneId = userInfo.ZoneId
		oldApId = userInfo.AccessPointId
	}
	userInfo.ZoneId = zoneId
	userInfo.AccessPointId = accessPointId

	// Update position
	if longitude == nil || latitude == nil {
		userInfo.LocationInfo = nil
	} else {
		if userInfo.LocationInfo == nil {
			userInfo.LocationInfo = new(LocationInfo)
			userInfo.LocationInfo.Accuracy = 1
		}
		userInfo.LocationInfo.Longitude = *longitude
		userInfo.LocationInfo.Latitude = *latitude
	}

	// Update User info in DB & Send notifications
	_ = rc.JSONSetEntry(baseKey+typeUser+":"+address, ".", convertUserInfoToJson(userInfo))
	checkNotificationRegistrations(USER_TRACKING_AND_ZONAL_TRAFFIC, oldZoneId, zoneId, oldApId, accessPointId, address)
}

func updateZoneInfo(zoneId string, nbAccessPoints int, nbUnsrvAccessPoints int, nbUsers int) {
	// Get Zone Info from DB
	jsonZoneInfo, _ := rc.JSONGetEntry(baseKey+typeZone+":"+zoneId, ".")
	zoneInfo := convertJsonToZoneInfo(jsonZoneInfo)

	// Create new zone info if necessary
	if zoneInfo == nil {
		zoneInfo = new(ZoneInfo)
		zoneInfo.ZoneId = zoneId
		zoneInfo.ResourceURL = hostUrl.String() + basePath + "zones/" + zoneId
	}

	// Update info
	if nbAccessPoints != -1 {
		zoneInfo.NumberOfAccessPoints = int32(nbAccessPoints)
	}
	if nbUnsrvAccessPoints != -1 {
		zoneInfo.NumberOfUnservicableAccessPoints = int32(nbUnsrvAccessPoints)
	}
	if nbUsers != -1 {
		zoneInfo.NumberOfUsers = int32(nbUsers)
	}

	// Update Zone info in DB & Send notifications
	_ = rc.JSONSetEntry(baseKey+typeZone+":"+zoneId, ".", convertZoneInfoToJson(zoneInfo))
	checkNotificationRegistrations(ZONE_STATUS, zoneId, "", "", strconv.Itoa(nbUsers), "")
}

func updateAccessPointInfo(zoneId string, apId string, conTypeStr string, opStatusStr string, nbUsers int, longitude *float32, latitude *float32) {
	// Get AP Info from DB
	jsonApInfo, _ := rc.JSONGetEntry(baseKey+typeZone+":"+zoneId+":"+typeAccessPoint+":"+apId, ".")
	apInfo := convertJsonToAccessPointInfo(jsonApInfo)

	// Create new AP info if necessary
	if apInfo == nil {
		apInfo = new(AccessPointInfo)
		apInfo.AccessPointId = apId
		apInfo.ResourceURL = hostUrl.String() + basePath + "zones/" + zoneId + "/accessPoints/" + apId
		conType := convertStringToConnectionType(conTypeStr)
		apInfo.ConnectionType = &conType
	}

	// Update info
	if opStatusStr != "" {
		opStatus := convertStringToOperationStatus(opStatusStr)
		apInfo.OperationStatus = &opStatus
	}
	if nbUsers != -1 {
		apInfo.NumberOfUsers = int32(nbUsers)
	}

	// Update position
	if longitude == nil || latitude == nil {
		apInfo.LocationInfo = nil
	} else {
		if apInfo.LocationInfo == nil {
			apInfo.LocationInfo = new(LocationInfo)
			apInfo.LocationInfo.Accuracy = 1
		}
		apInfo.LocationInfo.Longitude = *longitude
		apInfo.LocationInfo.Latitude = *latitude
	}

	// Update AP info in DB & Send notifications
	_ = rc.JSONSetEntry(baseKey+typeZone+":"+zoneId+":"+typeAccessPoint+":"+apId, ".", convertAccessPointInfoToJson(apInfo))
	checkNotificationRegistrations(ZONE_STATUS, zoneId, apId, strconv.Itoa(nbUsers), "", "")
}

func zoneStatusReInit() {
	//reusing the object response for the get multiple zoneStatusSubscription
	var zoneList ZoneStatusNotificationSubscriptionList

	keyName := baseKey + typeZoneStatusSubscription + "*"
	_ = rc.ForEachJSONEntry(keyName, populateZoneStatusList, &zoneList)

	maxZoneStatusSubscriptionId := 0
	for _, zone := range zoneList.ZoneStatusSubscription {
		resourceUrl := strings.Split(zone.ResourceURL, "/")
		subscriptionId, err := strconv.Atoi(resourceUrl[len(resourceUrl)-1])
		if err != nil {
			log.Error(err)
		} else {
			if subscriptionId > maxZoneStatusSubscriptionId {
				maxZoneStatusSubscriptionId = subscriptionId
			}

			var zoneStatus ZoneStatusCheck
			opStatus := zone.OperationStatus
			if opStatus != nil {
				for i := 0; i < len(opStatus); i++ {
					switch opStatus[i] {
					case SERVICEABLE:
						zoneStatus.Serviceable = true
					case UNSERVICEABLE:
						zoneStatus.Unserviceable = true
					case OPSTATUS_UNKNOWN:
						zoneStatus.Unknown = true
					default:
					}
				}
			}
			zoneStatus.NbUsersInZoneThreshold = (int)(zone.NumberOfUsersZoneThreshold)
			zoneStatus.NbUsersInAPThreshold = (int)(zone.NumberOfUsersAPThreshold)
			zoneStatus.ZoneId = zone.ZoneId
			zoneStatusSubscriptionMap[subscriptionId] = &zoneStatus
		}
	}
	nextZoneStatusSubscriptionIdAvailable = maxZoneStatusSubscriptionId + 1
}

func zonalTrafficReInit() {
	//reusing the object response for the get multiple zonalSubscription
	var zoneList ZonalTrafficNotificationSubscriptionList

	keyName := baseKey + typeZonalSubscription + "*"
	_ = rc.ForEachJSONEntry(keyName, populateZonalTrafficList, &zoneList)

	maxZonalSubscriptionId := 0
	for _, zone := range zoneList.ZonalTrafficSubscription {
		resourceUrl := strings.Split(zone.ResourceURL, "/")
		subscriptionId, err := strconv.Atoi(resourceUrl[len(resourceUrl)-1])
		if err != nil {
			log.Error(err)
		} else {
			if subscriptionId > maxZonalSubscriptionId {
				maxZonalSubscriptionId = subscriptionId
			}

			for i := 0; i < len(zone.UserEventCriteria); i++ {
				switch zone.UserEventCriteria[i] {
				case ENTERING:
					zonalSubscriptionEnteringMap[subscriptionId] = zone.ZoneId
				case LEAVING:
					zonalSubscriptionLeavingMap[subscriptionId] = zone.ZoneId
				case TRANSFERRING:
					zonalSubscriptionTransferringMap[subscriptionId] = zone.ZoneId
				default:
				}
			}
			zonalSubscriptionMap[subscriptionId] = zone.ZoneId
		}
	}
	nextZonalSubscriptionIdAvailable = maxZonalSubscriptionId + 1
}

func userTrackingReInit() {
	//reusing the object response for the get multiple zonalSubscription
	var userList UserTrackingNotificationSubscriptionList

	keyName := baseKey + typeUserSubscription + "*"
	_ = rc.ForEachJSONEntry(keyName, populateUserTrackingList, &userList)

	maxUserSubscriptionId := 0
	for _, user := range userList.UserTrackingSubscription {
		resourceUrl := strings.Split(user.ResourceURL, "/")
		subscriptionId, err := strconv.Atoi(resourceUrl[len(resourceUrl)-1])
		if err != nil {
			log.Error(err)
		} else {
			if subscriptionId > maxUserSubscriptionId {
				maxUserSubscriptionId = subscriptionId
			}

			for i := 0; i < len(user.UserEventCriteria); i++ {
				switch user.UserEventCriteria[i] {
				case ENTERING:
					userSubscriptionEnteringMap[subscriptionId] = user.Address
				case LEAVING:
					userSubscriptionLeavingMap[subscriptionId] = user.Address
				case TRANSFERRING:
					userSubscriptionTransferringMap[subscriptionId] = user.Address
				default:
				}
			}
			userSubscriptionMap[subscriptionId] = user.Address
		}
	}
	nextUserSubscriptionIdAvailable = maxUserSubscriptionId + 1
}
