/*
 * Copyright (c) 2020  InterDigital Communications, Inc
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

package sbi

import (
	dataModel "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	mod "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model"
	mq "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq"
)

const moduleName string = "meep-rnis-sbi"
const redisAddr string = "meep-redis-master.default.svc.cluster.local:6379"

type RnisSbi struct {
	sandboxName          string
	mqLocal              *mq.MsgQueue
	handlerId            int
	activeModel          *mod.Model
	updateUeEcgiInfoCB   func(string, string, string, string)
	updateAppEcgiInfoCB  func(string, string, string, string)
	updateScenarioNameCB func(string)
	cleanUpCB            func()
}

var sbi *RnisSbi

// Init - RNI Service SBI initialization
func Init(sandboxName string,
	updateUeEcgiInfo func(string, string, string, string),
	updateAppEcgiInfo func(string, string, string, string),
	updateScenarioName func(string),
	cleanUp func()) (err error) {

	// Create new SBI instance
	sbi = new(RnisSbi)
	sbi.sandboxName = sandboxName

	// Create message queue
	sbi.mqLocal, err = mq.NewMsgQueue(mq.GetLocalName(sandboxName), moduleName, sandboxName, redisAddr)
	if err != nil {
		log.Error("Failed to create Message Queue with error: ", err)
		return err
	}
	log.Info("Message Queue created")

	// Create new active scenario model
	modelCfg := mod.ModelCfg{
		Name:      "activeScenario",
		Namespace: sandboxName,
		Module:    moduleName,
		UpdateCb:  nil,
		DbAddr:    redisAddr,
	}
	sbi.activeModel, err = mod.NewModel(modelCfg)
	if err != nil {
		log.Error("Failed to create model: ", err.Error())
		return err
	}

	sbi.updateUeEcgiInfoCB = updateUeEcgiInfo
	sbi.updateAppEcgiInfoCB = updateAppEcgiInfo
	sbi.updateScenarioNameCB = updateScenarioName
	sbi.cleanUpCB = cleanUp

	// Initialize service
	processActiveScenarioUpdate()

	return nil
}

// Run - MEEP RNIS execution
func Run() (err error) {

	// Register Message Queue handler
	handler := mq.MsgHandler{Handler: msgHandler, UserData: nil}
	sbi.handlerId, err = sbi.mqLocal.RegisterHandler(handler)
	if err != nil {
		log.Error("Failed to register message queue handler: ", err.Error())
		return err
	}

	return nil
}

// Message Queue handler
func msgHandler(msg *mq.Msg, userData interface{}) {
	switch msg.Message {
	case mq.MsgScenarioActivate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		processActiveScenarioUpdate()
	case mq.MsgScenarioUpdate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		processActiveScenarioUpdate()
	case mq.MsgScenarioTerminate:
		log.Debug("RX MSG: ", mq.PrintMsg(msg))
		processActiveScenarioTerminate()
	default:
		log.Trace("Ignoring unsupported message: ", mq.PrintMsg(msg))
	}
}

func processActiveScenarioTerminate() {
	log.Debug("processActiveScenarioTerminate")

	// Sync with active scenario store
	sbi.activeModel.UpdateScenario()

	sbi.cleanUpCB()
}

func processActiveScenarioUpdate() {
	log.Debug("processActiveScenarioUpdate")

	// Update scenario Name that needs to be accessed by the NBI
	scenarioName := sbi.activeModel.GetScenarioName()
	sbi.updateScenarioNameCB(scenarioName)

	// Update UE info
	ueNameList := sbi.activeModel.GetNodeNames("UE")
	for _, name := range ueNameList {

		ueParent := sbi.activeModel.GetNodeParent(name)
		if poa, ok := ueParent.(*dataModel.NetworkLocation); ok {
			poaParent := sbi.activeModel.GetNodeParent(poa.Name)
			if zone, ok := poaParent.(*dataModel.Zone); ok {
				zoneParent := sbi.activeModel.GetNodeParent(zone.Name)
				if domain, ok := zoneParent.(*dataModel.Domain); ok {
					mnc := ""
					mcc := ""
					cellId := ""
					if domain.CellularDomainConfig != nil {
						mnc = domain.CellularDomainConfig.Mnc
						mcc = domain.CellularDomainConfig.Mcc
					}
					if poa.CellularPoaConfig != nil {
						if poa.CellularPoaConfig.CellId != "" {
							cellId = poa.CellularPoaConfig.CellId
						} else {
							cellId = domain.CellularDomainConfig.DefaultCellId
						}
					} else {
						if domain.CellularDomainConfig != nil {
							cellId = domain.CellularDomainConfig.DefaultCellId
						}
					}

					sbi.updateUeEcgiInfoCB(name, mnc, mcc, cellId)
				}
			}
		}
	}

	// Update Edge App info
	meAppNameList := sbi.activeModel.GetNodeNames("EDGE-APP")
	ueAppNameList := sbi.activeModel.GetNodeNames("UE-APP")
	var appNameList []string
	appNameList = append(appNameList, meAppNameList...)
	appNameList = append(appNameList, ueAppNameList...)
	for _, meAppName := range appNameList {
		meAppParent := sbi.activeModel.GetNodeParent(meAppName)
		if pl, ok := meAppParent.(*dataModel.PhysicalLocation); ok {
			plParent := sbi.activeModel.GetNodeParent(pl.Name)
			if nl, ok := plParent.(*dataModel.NetworkLocation); ok {
				//nl can be either POA for {FOG or UE} or Zone Default for {Edge
				nlParent := sbi.activeModel.GetNodeParent(nl.Name)
				if zone, ok := nlParent.(*dataModel.Zone); ok {
					zoneParent := sbi.activeModel.GetNodeParent(zone.Name)
					if domain, ok := zoneParent.(*dataModel.Domain); ok {
						mnc := ""
						mcc := ""
						cellId := ""
						if domain.CellularDomainConfig != nil {
							mnc = domain.CellularDomainConfig.Mnc
							mcc = domain.CellularDomainConfig.Mcc
						}
						if nl.CellularPoaConfig != nil {
							if nl.CellularPoaConfig.CellId != "" {
								cellId = nl.CellularPoaConfig.CellId
							} else {
								cellId = domain.CellularDomainConfig.DefaultCellId
							}
						} else {
							if domain.CellularDomainConfig != nil {
								cellId = domain.CellularDomainConfig.DefaultCellId
							}
						}

						sbi.updateAppEcgiInfoCB(meAppName, mnc, mcc, cellId)
					}
				}
			}
		}
	}

}