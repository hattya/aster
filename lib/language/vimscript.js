//
// aster :: language/vimscript.js
//
//   Copyright (c) 2017-2024 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

'use strict';

var os = require('os');
var path = require('path');
var language = require('language');

var covimerage = exports.covimerage = function covimerage() {
  if (!os.whence('covimerage')) {
    aster.notify('failure', language.prefix + 'covimerage', 'covimerage not found!');
    return true;
  }

  return language.system({
    args: ['covimerage'].concat(Array.prototype.slice.call(arguments)),
    title: 'covimerage',
    success: arguments[0] + ' passed',
    failure: arguments[0] + ' failed',
  });
};

covimerage.report = function() {
  return covimerage.apply(null, ['report'].concat(Array.prototype.slice.call(arguments)));
};

covimerage.run = function() {
  return covimerage.apply(null, ['run'].concat(Array.prototype.slice.call(arguments)));
};

covimerage.write_coverage = function() {
  return covimerage.apply(null, ['write_coverage'].concat(Array.prototype.slice.call(arguments)));
};

covimerage.xml = function() {
  return covimerage.apply(null, ['xml'].concat(Array.prototype.slice.call(arguments)));
};

var primula = exports.primula = function primula() {
  if (!os.whence('primula')) {
    aster.notify('failure', language.prefix + 'primula', 'primula not found!');
    return true;
  }

  return language.system({
    args: ['primula'].concat(Array.prototype.slice.call(arguments)),
    title: 'primula',
    success: arguments[0] + ' passed',
    failure: arguments[0] + ' failed',
  });
};

primula.annotate = function annotate() {
  return primula.apply(null, ['annotate'].concat(Array.prototype.slice.call(arguments)));
};

primula.combine = function combine() {
  return primula.apply(null, ['combine'].concat(Array.prototype.slice.call(arguments)));
};

primula.erase = function erase() {
  return primula.apply(null, ['erase'].concat(Array.prototype.slice.call(arguments)));
};

primula.html = function html() {
  return primula.apply(null, ['html'].concat(Array.prototype.slice.call(arguments)));
};

primula.json = function json() {
  return primula.apply(null, ['json'].concat(Array.prototype.slice.call(arguments)));
};

primula.lcov = function lcov() {
  return primula.apply(null, ['lcov'].concat(Array.prototype.slice.call(arguments)));
};

primula.report = function report() {
  return primula.apply(null, ['report'].concat(Array.prototype.slice.call(arguments)));
};

primula.run = function run() {
  return primula.apply(null, ['run'].concat(Array.prototype.slice.call(arguments)));
};

primula.xml = function xml() {
  return primula.apply(null, ['xml'].concat(Array.prototype.slice.call(arguments)));
};

exports.themis = function themis() {
  var script = 'themis';
  if (!os.whence(script)) {
    var ok = ['.', '..'].some(function(e) {
      script = path.join(e, 'vim-themis', 'bin', 'themis');
      return os.whence(script);
    });
    if (!ok) {
      aster.notify('failure', language.prefix + 'themis', 'themis not found!');
      return true;
    }
  }

  return language.system({
    args: [script].concat(Array.prototype.slice.call(arguments)),
    title: 'themis',
    success: 'passed',
    failure: 'failed',
  });
};
