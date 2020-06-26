//
// aster :: language/restructuredtext.js
//
//   Copyright (c) 2017-2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

'use strict';

var os = require('os');
var path = require('path');
var language = require('language');

exports.rst2html = function rst2html(object) {
  var script;
  var ok = ['rst2html5.py', 'rst2html.py'].some(function(s) {
    script = s;
    if (os.whence(script)) {
      return true;
    }
    script = s.slice(0, -3);
    return os.whence(script);
  });
  if (!ok) {
    aster.notify('failure', language.prefix + 'rst2html', 'rst2html not found!');
    return true;
  }

  var args = [script];
  if (Array.isArray(object.options)) {
    args = args.concat(object.options);
  }
  args.push(object.src);
  args.push(object.dst || object.src.slice(0, -path.extname(object.src).length) + '.html');
  return language.system({
    args: args,
    title: 'rst2html',
    success: object.src,
    failure: object.src + ' failed',
  });
};
