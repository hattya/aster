//
// aster :: language/python.js
//
//   Copyright (c) 2018-2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

'use strict';

var os = require('os');
var language = require('language');

var coverage = exports.coverage = function coverage() {
  if (!os.whence('coverage')) {
    aster.notify('failure', exports.prefix + 'coverage', 'coverage not found!');
    return true;
  }

  return language.system({
    args: ['coverage'].concat(Array.prototype.slice.call(arguments)),
    title: 'coverage',
    success: arguments[0] + ' passed',
    failure: arguments[0] + ' failed',
  });
};

coverage.annotate = function annotate() {
  return coverage.apply(null, ['annotate'].concat(Array.prototype.slice.call(arguments)));
};

coverage.combine = function combine() {
  return coverage.apply(null, ['combine'].concat(Array.prototype.slice.call(arguments)));
};

coverage.erase = function erase() {
  return coverage.apply(null, ['erase'].concat(Array.prototype.slice.call(arguments)));
};

coverage.html = function html() {
  return coverage.apply(null, ['html'].concat(Array.prototype.slice.call(arguments)));
};

coverage.report = function report() {
  return coverage.apply(null, ['report'].concat(Array.prototype.slice.call(arguments)));
};

coverage.run = function run() {
  return coverage.apply(null, ['run'].concat(Array.prototype.slice.call(arguments)));
};

coverage.xml = function xml() {
  return coverage.apply(null, ['xml'].concat(Array.prototype.slice.call(arguments)));
};
