/**
 * Demo iperf transit App API
 * This is the Demo iperf transit App API
 *
 * OpenAPI spec version: 0.0.1
 *
 * NOTE: This class is auto generated by the swagger code generator program.
 * https://github.com/swagger-api/swagger-codegen.git
 *
 * Swagger Codegen version: 2.4.1
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
    factory(root.expect, root.DemoIperfTransitAppApi);
  }
}(this, function(expect, DemoIperfTransitAppApi) {
  'use strict';

  var instance;

  beforeEach(function() {
    instance = new DemoIperfTransitAppApi.IperfInfo();
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

  describe('IperfInfo', function() {
    it('should create an instance of IperfInfo', function() {
      // uncomment below and update the code to test IperfInfo
      //var instance = new DemoIperfTransitAppApi.IperfInfo();
      //expect(instance).to.be.a(DemoIperfTransitAppApi.IperfInfo);
    });

    it('should have the property name (base name: "name")', function() {
      // uncomment below and update the code to test the property name
      //var instance = new DemoIperfTransitAppApi.IperfInfo();
      //expect(instance).to.be();
    });

    it('should have the property app (base name: "app")', function() {
      // uncomment below and update the code to test the property app
      //var instance = new DemoIperfTransitAppApi.IperfInfo();
      //expect(instance).to.be();
    });

    it('should have the property throughput (base name: "throughput")', function() {
      // uncomment below and update the code to test the property throughput
      //var instance = new DemoIperfTransitAppApi.IperfInfo();
      //expect(instance).to.be();
    });

  });

}));