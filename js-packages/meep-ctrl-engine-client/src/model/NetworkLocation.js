/**
 * MEEP Controller REST API
 * Copyright (c) 2019  InterDigital Communications, Inc Licensed under the Apache License, Version 2.0 (the \"License\"); you may not use this file except in compliance with the License. You may obtain a copy of the License at      http://www.apache.org/licenses/LICENSE-2.0  Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an \"AS IS\" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License. 
 *
 * OpenAPI spec version: 1.0.0
 *
 * NOTE: This class is auto generated by the swagger code generator program.
 * https://github.com/swagger-api/swagger-codegen.git
 *
 * Swagger Codegen version: 2.3.1
 *
 * Do not edit the class manually.
 *
 */

(function(root, factory) {
  if (typeof define === 'function' && define.amd) {
    // AMD. Register as an anonymous module.
    define(['ApiClient', 'model/PhysicalLocation'], factory);
  } else if (typeof module === 'object' && module.exports) {
    // CommonJS-like environments that support module.exports, like Node.
    module.exports = factory(require('../ApiClient'), require('./PhysicalLocation'));
  } else {
    // Browser globals (root is window)
    if (!root.MeepControllerRestApi) {
      root.MeepControllerRestApi = {};
    }
    root.MeepControllerRestApi.NetworkLocation = factory(root.MeepControllerRestApi.ApiClient, root.MeepControllerRestApi.PhysicalLocation);
  }
}(this, function(ApiClient, PhysicalLocation) {
  'use strict';




  /**
   * The NetworkLocation model module.
   * @module model/NetworkLocation
   * @version 1.0.0
   */

  /**
   * Constructs a new <code>NetworkLocation</code>.
   * Logical network location object
   * @alias module:model/NetworkLocation
   * @class
   */
  var exports = function() {
    var _this = this;











  };

  /**
   * Constructs a <code>NetworkLocation</code> from a plain JavaScript object, optionally creating a new instance.
   * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
   * @param {Object} data The plain JavaScript object bearing properties of interest.
   * @param {module:model/NetworkLocation} obj Optional instance to populate.
   * @return {module:model/NetworkLocation} The populated <code>NetworkLocation</code> instance.
   */
  exports.constructFromObject = function(data, obj) {
    if (data) {
      obj = obj || new exports();

      if (data.hasOwnProperty('id')) {
        obj['id'] = ApiClient.convertToType(data['id'], 'String');
      }
      if (data.hasOwnProperty('name')) {
        obj['name'] = ApiClient.convertToType(data['name'], 'String');
      }
      if (data.hasOwnProperty('type')) {
        obj['type'] = ApiClient.convertToType(data['type'], 'String');
      }
      if (data.hasOwnProperty('terminalLinkLatency')) {
        obj['terminalLinkLatency'] = ApiClient.convertToType(data['terminalLinkLatency'], 'Number');
      }
      if (data.hasOwnProperty('terminalLinkLatencyVariation')) {
        obj['terminalLinkLatencyVariation'] = ApiClient.convertToType(data['terminalLinkLatencyVariation'], 'Number');
      }
      if (data.hasOwnProperty('terminalLinkThroughput')) {
        obj['terminalLinkThroughput'] = ApiClient.convertToType(data['terminalLinkThroughput'], 'Number');
      }
      if (data.hasOwnProperty('terminalLinkPacketLoss')) {
        obj['terminalLinkPacketLoss'] = ApiClient.convertToType(data['terminalLinkPacketLoss'], 'Number');
      }
      if (data.hasOwnProperty('meta')) {
        obj['meta'] = ApiClient.convertToType(data['meta'], {'String': 'String'});
      }
      if (data.hasOwnProperty('userMeta')) {
        obj['userMeta'] = ApiClient.convertToType(data['userMeta'], {'String': 'String'});
      }
      if (data.hasOwnProperty('physicalLocations')) {
        obj['physicalLocations'] = ApiClient.convertToType(data['physicalLocations'], [PhysicalLocation]);
      }
    }
    return obj;
  }

  /**
   * Unique network location ID
   * @member {String} id
   */
  exports.prototype['id'] = undefined;
  /**
   * Network location name
   * @member {String} name
   */
  exports.prototype['name'] = undefined;
  /**
   * Network location type
   * @member {module:model/NetworkLocation.TypeEnum} type
   */
  exports.prototype['type'] = undefined;
  /**
   * Latency in ms for all terminal links within network location
   * @member {Number} terminalLinkLatency
   */
  exports.prototype['terminalLinkLatency'] = undefined;
  /**
   * Latency variation in ms for all terminal links within network location
   * @member {Number} terminalLinkLatencyVariation
   */
  exports.prototype['terminalLinkLatencyVariation'] = undefined;
  /**
   * The limit of the traffic supported for all terminal links within the network location
   * @member {Number} terminalLinkThroughput
   */
  exports.prototype['terminalLinkThroughput'] = undefined;
  /**
   * Packet lost (in terms of percentage) for all terminal links within the network location
   * @member {Number} terminalLinkPacketLoss
   */
  exports.prototype['terminalLinkPacketLoss'] = undefined;
  /**
   * Key/Value Pair Map (string, string)
   * @member {Object.<String, String>} meta
   */
  exports.prototype['meta'] = undefined;
  /**
   * Key/Value Pair Map (string, string)
   * @member {Object.<String, String>} userMeta
   */
  exports.prototype['userMeta'] = undefined;
  /**
   * @member {Array.<module:model/PhysicalLocation>} physicalLocations
   */
  exports.prototype['physicalLocations'] = undefined;


  /**
   * Allowed values for the <code>type</code> property.
   * @enum {String}
   * @readonly
   */
  exports.TypeEnum = {
    /**
     * value: "POA"
     * @const
     */
    "POA": "POA",
    /**
     * value: "DEFAULT"
     * @const
     */
    "DEFAULT": "DEFAULT"  };


  return exports;
}));


