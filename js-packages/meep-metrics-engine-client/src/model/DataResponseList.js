/**
 * Metrics Engine Service API
 * This is Metrics Engine Services API
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
    define(['ApiClient', 'model/DataResponse'], factory);
  } else if (typeof module === 'object' && module.exports) {
    // CommonJS-like environments that support module.exports, like Node.
    module.exports = factory(require('../ApiClient'), require('./DataResponse'));
  } else {
    // Browser globals (root is window)
    if (!root.MetricsEngineServiceApi) {
      root.MetricsEngineServiceApi = {};
    }
    root.MetricsEngineServiceApi.DataResponseList = factory(root.MetricsEngineServiceApi.ApiClient, root.MetricsEngineServiceApi.DataResponse);
  }
}(this, function(ApiClient, DataResponse) {
  'use strict';




  /**
   * The DataResponseList model module.
   * @module model/DataResponseList
   * @version 1.0.0
   */

  /**
   * Constructs a new <code>DataResponseList</code>.
   * @alias module:model/DataResponseList
   * @class
   */
  var exports = function() {
    var _this = this;


  };

  /**
   * Constructs a <code>DataResponseList</code> from a plain JavaScript object, optionally creating a new instance.
   * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
   * @param {Object} data The plain JavaScript object bearing properties of interest.
   * @param {module:model/DataResponseList} obj Optional instance to populate.
   * @return {module:model/DataResponseList} The populated <code>DataResponseList</code> instance.
   */
  exports.constructFromObject = function(data, obj) {
    if (data) {
      obj = obj || new exports();

      if (data.hasOwnProperty('dataResponse')) {
        obj['dataResponse'] = ApiClient.convertToType(data['dataResponse'], [DataResponse]);
      }
    }
    return obj;
  }

  /**
   * @member {Array.<module:model/DataResponse>} dataResponse
   */
  exports.prototype['dataResponse'] = undefined;



  return exports;
}));


