/*
 * Metrics Engine Service API
 *
 * This is Metrics Engine Services API
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package server

type DataResponseData struct {

	// Number of packets received since last log event
	Rx int32 `json:"rx,omitempty"`

	// Number of bytes received since last log event
	RxBytes int32 `json:"rxBytes,omitempty"`

	// Throughput measured between 2 pods in Mbits/seconds
	Throughput float32 `json:"throughput,omitempty"`

	// Number of packets loss between2 pods as a percentage
	PacketLoss string `json:"packet-loss,omitempty"`

	// Latency measured betwen 2 pods in ms
	Latency int32 `json:"latency,omitempty"`
}
