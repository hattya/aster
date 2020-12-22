//
// aster :: language/markdown.js
//
//   Copyright (c) 2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

'use strict';

var os = require('os');
var path = require('path');
var language = require('language');

exports.md2html = function md2html(object) {
  if (!os.whence('md2html')) {
    aster.notify('failure', language.prefix + 'md2html', 'md2html not found!');
    return true;
  }

  var args = ['md2html'];
  if (Array.isArray(object.options)) {
    args = args.concat(object.options);
  }
  args.push(object.src);
  args.push(object.dst || object.src.slice(0, -path.extname(object.src).length) + '.html');
  return language.system({
    args: args,
    title: 'md2html',
    success: object.src,
    failure: object.src + ' failed',
  });
};
