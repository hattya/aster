//
// aster :: language/python.js
//
//   Copyright (c) 2018-2026 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

'use strict';

var os = require('os');
var language = require('language');

var coverage = exports.coverage = function() {
  if (!os.whence('coverage')) {
    aster.notify('failure', language.prefix + 'coverage', 'coverage not found!');
    return true;
  }

  return language.system({
    args: ['coverage'].concat(Array.prototype.slice.call(arguments)),
    title: 'coverage',
    success: arguments[0] + ' passed',
    failure: arguments[0] + ' failed',
  });
};

coverage.annotate = function() {
  return coverage.apply(null, ['annotate'].concat(Array.prototype.slice.call(arguments)));
};

coverage.combine = function() {
  return coverage.apply(null, ['combine'].concat(Array.prototype.slice.call(arguments)));
};

coverage.erase = function() {
  return coverage.apply(null, ['erase'].concat(Array.prototype.slice.call(arguments)));
};

coverage.html = function() {
  return coverage.apply(null, ['html'].concat(Array.prototype.slice.call(arguments)));
};

coverage.json = function() {
  return coverage.apply(null, ['json'].concat(Array.prototype.slice.call(arguments)));
};

coverage.lcov = function() {
  return coverage.apply(null, ['lcov'].concat(Array.prototype.slice.call(arguments)));
};

coverage.report = function() {
  return coverage.apply(null, ['report'].concat(Array.prototype.slice.call(arguments)));
};

coverage.run = function() {
  return coverage.apply(null, ['run'].concat(Array.prototype.slice.call(arguments)));
};

coverage.xml = function() {
  return coverage.apply(null, ['xml'].concat(Array.prototype.slice.call(arguments)));
};
