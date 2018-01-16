//
// aster :: language/restructuredtext.js
//
//   Copyright (c) 2017-2018 Akinori Hattori <hattya@gmail.com>
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
  // exec
  var args = [script].concat(object.options);
  args.push(object.src);
  args.push(object.dst || object.src.slice(0, -path.extname(object.src).length) + '.html');
  console.log(language.prompt + args.join(' '));
  var rv = os.system(args);
  // notify
  var title = language.prefix + 'rst2html';
  if (!rv) {
    aster.notify('success', title, object.src);
  } else {
    aster.notify('failure', title, object.src + ' failed');
  }
  return rv;
};
