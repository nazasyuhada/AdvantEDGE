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
    // AMD.
    define(['expect.js', '../../src/index'], factory);
  } else if (typeof module === 'object' && module.exports) {
    // CommonJS-like environments that support module.exports, like Node.
    factory(require('expect.js'), require('../../src/index'));
  } else {
    // Browser globals (root is window)
    factory(root.expect, root.MeepControllerRestApi);
  }
}(this, function(expect, MeepControllerRestApi) {
  'use strict';

  var instance;

  beforeEach(function() {
    instance = new MeepControllerRestApi.Deployment();
  });

  var getProperty = function(object, getter, property) {
    // Use getter method if present; otherwise, get the property directly.
    if (typeof object[getter] === 'function')
      return object[getter]();
    else
      return object[property];
  }

  var setProperty = function(object, setter, property, value) {
    // Use setter method if present; otherwise, set the property directly.
    if (typeof object[setter] === 'function')
      object[setter](value);
    else
      object[property] = value;
  }

  describe('Deployment', function() {
    it('should create an instance of Deployment', function() {
      // uncomment below and update the code to test Deployment
      //var instane = new MeepControllerRestApi.Deployment();
      //expect(instance).to.be.a(MeepControllerRestApi.Deployment);
    });

    it('should have the property interDomainLatency (base name: "interDomainLatency")', function() {
      // uncomment below and update the code to test the property interDomainLatency
      //var instane = new MeepControllerRestApi.Deployment();
      //expect(instance).to.be();
    });

    it('should have the property interDomainLatencyVariation (base name: "interDomainLatencyVariation")', function() {
      // uncomment below and update the code to test the property interDomainLatencyVariation
      //var instane = new MeepControllerRestApi.Deployment();
      //expect(instance).to.be();
    });

    it('should have the property interDomainThroughput (base name: "interDomainThroughput")', function() {
      // uncomment below and update the code to test the property interDomainThroughput
      //var instane = new MeepControllerRestApi.Deployment();
      //expect(instance).to.be();
    });

    it('should have the property interDomainPacketLoss (base name: "interDomainPacketLoss")', function() {
      // uncomment below and update the code to test the property interDomainPacketLoss
      //var instane = new MeepControllerRestApi.Deployment();
      //expect(instance).to.be();
    });

    it('should have the property meta (base name: "meta")', function() {
      // uncomment below and update the code to test the property meta
      //var instane = new MeepControllerRestApi.Deployment();
      //expect(instance).to.be();
    });

    it('should have the property userMeta (base name: "userMeta")', function() {
      // uncomment below and update the code to test the property userMeta
      //var instane = new MeepControllerRestApi.Deployment();
      //expect(instance).to.be();
    });

    it('should have the property domains (base name: "domains")', function() {
      // uncomment below and update the code to test the property domains
      //var instane = new MeepControllerRestApi.Deployment();
      //expect(instance).to.be();
    });

  });

}));
