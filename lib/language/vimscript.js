//
// aster :: language/vimscript.js
//
//   Copyright (c) 2017-2020 Akinori Hattori <hattya@gmail.com>
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
