//
// aster :: language/vimscript.js
//
//   Copyright (c) 2017-2026 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

'use strict';

var os = require('os');
var path = require('path');
var language = require('language');

var covimerage = exports.covimerage = function() {
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

var primula = exports.primula = function() {
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

primula.annotate = function() {
  return primula.apply(null, ['annotate'].concat(Array.prototype.slice.call(arguments)));
};

primula.combine = function() {
  return primula.apply(null, ['combine'].concat(Array.prototype.slice.call(arguments)));
};

primula.erase = function() {
  return primula.apply(null, ['erase'].concat(Array.prototype.slice.call(arguments)));
};

primula.html = function() {
  return primula.apply(null, ['html'].concat(Array.prototype.slice.call(arguments)));
};

primula.json = function() {
  return primula.apply(null, ['json'].concat(Array.prototype.slice.call(arguments)));
};

primula.lcov = function() {
  return primula.apply(null, ['lcov'].concat(Array.prototype.slice.call(arguments)));
};

primula.report = function() {
  return primula.apply(null, ['report'].concat(Array.prototype.slice.call(arguments)));
};

primula.run = function() {
  return primula.apply(null, ['run'].concat(Array.prototype.slice.call(arguments)));
};

primula.xml = function() {
  return primula.apply(null, ['xml'].concat(Array.prototype.slice.call(arguments)));
};

exports.themis = function() {
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
