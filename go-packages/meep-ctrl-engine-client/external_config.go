/*
 * MEEP Controller REST API
 *
 * Copyright (c) 2019  InterDigital Communications, Inc Licensed under the Apache License, Version 2.0 (the \"License\"); you may not use this file except in compliance with the License. You may obtain a copy of the License at      http://www.apache.org/licenses/LICENSE-2.0  Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an \"AS IS\" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package client

// External Process configuration. NOTE: Only valid if 'isExternal' is set.
type ExternalConfig struct {
	IngressServiceMap []IngressService `json:"ingressServiceMap,omitempty"`

	EgressServiceMap []EgressService `json:"egressServiceMap,omitempty"`
}
