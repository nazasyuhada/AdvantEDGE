/*
 * Location API
 *
 * The ETSI MEC ISG MEC012 Location API described using OpenAPI. The API is based on the Open Mobile Alliance's specification RESTful Network API for Zonal Presence
 *
 * API version: 1.1.1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package client

// A type containing zone status subscription.
type ZoneStatusSubscription struct {

	// Uniquely identifies this create subscription request. If there is a communication failure during the request, using the same clientCorrelator when retrying the request allows the operator to avoid creating a duplicate subscription.
	ClientCorrelator string `json:"clientCorrelator,omitempty"`

	// Self referring URL.
	ResourceURL string `json:"resourceURL,omitempty"`

	CallbackReference *UserTrackingSubscriptionCallbackReference `json:"callbackReference"`

	// Identifier of zone
	ZoneId string `json:"zoneId"`

	// Threshold number of users in a zone which if crossed shall cause a notification.
	NumberOfUsersZoneThreshold int32 `json:"numberOfUsersZoneThreshold,omitempty"`

	// Threshold number of users in an access point which if crossed shall cause a notification.
	NumberOfUsersAPThreshold int32 `json:"numberOfUsersAPThreshold,omitempty"`

	// List of operation status values to generate notifications for (these apply to all access points within a zone).
	OperationStatus []OperationStatus `json:"operationStatus,omitempty"`
}