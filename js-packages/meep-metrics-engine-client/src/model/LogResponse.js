/**
 * MEEP Metrics Engine Service REST API
 * Copyright (c) 2019 InterDigital Communications, Inc. All rights reserved. The information provided herein is the proprietary and confidential information of InterDigital Communications, Inc. 
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
    define(['ApiClient', 'model/LogResponseData'], factory);
  } else if (typeof module === 'object' && module.exports) {
    // CommonJS-like environments that support module.exports, like Node.
    module.exports = factory(require('../ApiClient'), require('./LogResponseData'));
  } else {
    // Browser globals (root is window)
    if (!root.MeepMetricsEngineServiceRestApi) {
      root.MeepMetricsEngineServiceRestApi = {};
    }
    root.MeepMetricsEngineServiceRestApi.LogResponse = factory(root.MeepMetricsEngineServiceRestApi.ApiClient, root.MeepMetricsEngineServiceRestApi.LogResponseData);
  }
}(this, function(ApiClient, LogResponseData) {
  'use strict';




  /**
   * The LogResponse model module.
   * @module model/LogResponse
   * @version 1.0.0
   */

  /**
   * Constructs a new <code>LogResponse</code>.
   * @alias module:model/LogResponse
   * @class
   * @param dest {String} Pod where the log message is taken from
   * @param dataType {String} Pod where the log message is taken from
   * @param src {String} Pod that originated the metrics logged in the message
   * @param timestamp {String} System time at which the metric was logged
   */
  var exports = function(dest, dataType, src, timestamp) {
    var _this = this;

    _this['dest'] = dest;
    _this['dataType'] = dataType;
    _this['src'] = src;
    _this['timestamp'] = timestamp;

  };

  /**
   * Constructs a <code>LogResponse</code> from a plain JavaScript object, optionally creating a new instance.
   * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
   * @param {Object} data The plain JavaScript object bearing properties of interest.
   * @param {module:model/LogResponse} obj Optional instance to populate.
   * @return {module:model/LogResponse} The populated <code>LogResponse</code> instance.
   */
  exports.constructFromObject = function(data, obj) {
    if (data) {
      obj = obj || new exports();

      if (data.hasOwnProperty('dest')) {
        obj['dest'] = ApiClient.convertToType(data['dest'], 'String');
      }
      if (data.hasOwnProperty('dataType')) {
        obj['dataType'] = ApiClient.convertToType(data['dataType'], 'String');
      }
      if (data.hasOwnProperty('src')) {
        obj['src'] = ApiClient.convertToType(data['src'], 'String');
      }
      if (data.hasOwnProperty('timestamp')) {
        obj['timestamp'] = ApiClient.convertToType(data['timestamp'], 'String');
      }
      if (data.hasOwnProperty('data')) {
        obj['data'] = LogResponseData.constructFromObject(data['data']);
      }
    }
    return obj;
  }

  /**
   * Pod where the log message is taken from
   * @member {String} dest
   */
  exports.prototype['dest'] = undefined;
  /**
   * Pod where the log message is taken from
   * @member {String} dataType
   */
  exports.prototype['dataType'] = undefined;
  /**
   * Pod that originated the metrics logged in the message
   * @member {String} src
   */
  exports.prototype['src'] = undefined;
  /**
   * System time at which the metric was logged
   * @member {String} timestamp
   */
  exports.prototype['timestamp'] = undefined;
  /**
   * @member {module:model/LogResponseData} data
   */
  exports.prototype['data'] = undefined;



  return exports;
}));

