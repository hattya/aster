//
// aster :: os.js
//
//   Copyright (c) 2014-2017 Akinori Hattori <hattya@gmail.com>
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
  getwd: os.getwd,
  mkdir: os.mkdir,
  open: open,
  remove: os.remove,
  rename: os.rename,
  stat: stat,
  system: os.system,
  whence: os.whence,
};
