//
// aster :: language/javascript.js
//
//   Copyright (c) 2017 Akinori Hattori <hattya@gmail.com>
//
//   Permission is hereby granted, free of charge, to any person
//   obtaining a copy of this software and associated documentation files
//   (the "Software"), to deal in the Software without restriction,
//   including without limitation the rights to use, copy, modify, merge,
//   publish, distribute, sublicense, and/or sell copies of the Software,
//   and to permit persons to whom the Software is furnished to do so,
//   subject to the following conditions:
//
//   The above copyright notice and this permission notice shall be
//   included in all copies or substantial portions of the Software.
//
//   THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
//   EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
//   MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
//   NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
//   BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
//   ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
//   CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//   SOFTWARE.
//

'use strict';

var os = require('os');
var language = require('language');

var npm = exports.npm = function() {
  if (!os.whence('npm')) {
    aster.notify('failure', exports.prefix + 'npm', 'npm command not found!');
    return true;
  }
  // exec
  var args = ['npm'].concat(Array.prototype.slice.call(arguments));
  console.log(language.prompt + args.join(' '));
  var rv = os.system(args);
  // notify
  var title = language.prefix + args[0];
  var cmd = args[1] !== 'run' ? args[1] : args[2] + ' script';
  if (!rv) {
    aster.notify('success', title, cmd + ' passed');
  } else {
    aster.notify('failure', title, cmd + ' failed');
  }
  return rv;
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