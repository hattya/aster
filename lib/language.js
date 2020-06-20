//
// aster :: language.js
//
//   Copyright (c) 2017-2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

'use strict';

var os = require('os');
var path = require('path');

exports.prefix = path.basename(os.getwd()) + ': ';
exports.prompt = '> ';

exports.system = function system(object) {
  // exec
  console.log(exports.prompt + object.args.join(' '));
  var rv = os.system(object.args, object.options);
  // notify
  var title = exports.prefix + object.title;
  if (!rv) {
    aster.notify('success', title, object.success);
  } else {
    aster.notify('failure', title, object.failure);
  }
  return rv;
};
