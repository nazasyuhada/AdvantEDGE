/*
 * MEEP Demo App API
 *
 * This is the MEEP Demo App API
 *
 * API version: 0.0.1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package server

// Ue state basic information object
type UeState struct {

	// Duration since the game stated
	Duration int32 `json:"duration,omitempty"`

	// Traffic info for the registered Ue
	TrafficBw int32 `json:"trafficBw,omitempty"`
}