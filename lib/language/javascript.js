//
// aster :: language/javascript.js
//
//   Copyright (c) 2017-2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

'use strict';

var os = require('os');
var language = require('language');

var npm = exports.npm = function() {
  if (!os.whence('npm')) {
    aster.notify('failure', language.prefix + 'npm', 'npm not found!');
    return true;
  }

  var cmd = arguments[0] !== 'run' ? arguments[0] : arguments[1] + ' script';
  return language.system({
    args: ['npm'].concat(Array.prototype.slice.call(arguments)),
    title: 'npm',
    success: cmd + ' passed',
    failure: cmd + ' failed',
  });
};

npm.install = function install() {
  return npm.apply(null, ['install'].concat(Array.prototype.slice.call(arguments)));
};

npm.run = function run() {
  return npm.apply(null, ['run'].concat(Array.prototype.slice.call(arguments)));
};

npm.test = function test() {
  return npm.apply(null, ['test'].concat(Array.prototype.slice.call(arguments)));
};
