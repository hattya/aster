//
// aster :: os.js
//
//   Copyright (c) 2014-2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

'use strict';

var os = process.binding('os');

function File(impl) {
  this._impl = impl;
}

File.prototype.close = function close() {
  return this._impl.Close.apply(this, arguments);
};

File.prototype.name = function name() {
  return this._impl.Name.apply(this, arguments);
};

File.prototype.read = function read() {
  return this._impl.Read.apply(this, arguments);
};

File.prototype.readLine = function readLine() {
  return this._impl.ReadLine.apply(this, arguments);
};

File.prototype.write = function write() {
  return this._impl.Write.apply(this, arguments);
};

function open() {
  return os.open.apply(new File(), arguments);
}

function FileInfo(name, size, mode, mtime) {
  this.name = name;
  this.size = size;
  this.mode = mode;
  this.mtime = mtime;
}

FileInfo.prototype.isDir = function isDir() {
  return (this.mode & os.MODE_DIR) !== 0;
};

FileInfo.prototype.isRegular = function isRegular() {
  return (this.mode & os.MODE_TYPE) === 0;
};

FileInfo.prototype.perm = function perm() {
  return this.mode & os.MODE_PERM;
};

function stat() {
  return os.stat.apply(new FileInfo(), arguments);
}

module.exports = {
  getenv: process.env.__get__,
  getwd: os.getwd,
  mkdir: os.mkdir,
  open: open,
  remove: os.remove,
  rename: os.rename,
  setenv: process.env.__set__,
  stat: stat,
  system: os.system,
  whence: os.whence,
};
