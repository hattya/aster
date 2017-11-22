// Code generated by "modulizer -l std -o std.go lib"; DO NOT EDIT.

package aster

import "github.com/hattya/otto.module"

type stdLoader struct {
}

func (*stdLoader) Load(id string) ([]byte, error) {
	if b, ok := files[id]; ok {
		return b, nil
	}
	return nil, module.ErrModule
}

func (*stdLoader) Resolve(id, _ string) (string, error) {
	for _, ext := range []string{"", ".js", ".json"} {
		k := id + ext
		if _, ok := files[k]; ok {
			return k, nil
		}
	}
	return "", module.ErrModule
}

var files = map[string][]byte{
	"language/go.js": []byte(`//
// aster :: language/go.js
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
var path = require('path');
var language = require('language');

function run(cmd, args, options) {
  // exec
  console.log(language.prompt + args.join(' '));
  var rv = os.system(args, options);
  // notify
  var title = language.prefix + args[0];
  if (!rv) {
    aster.notify('success', title, cmd + ' passed');
  } else {
    aster.notify('failure', title, cmd + ' failed');
  }
  return rv;
}

function parse(a) {
  var i;
  for (i = 0; a[i] === '-'; i++);
  if (0 < i && i < 3) {
    var j = a.indexOf('=', i);
    if (j === -1) {
      return a.slice(i - 1);
    } else if (i < j) {
      return a.slice(i - 1, j + 1);
    }
  }
  return a;
}

var dep = exports.dep = function dep() {
  if (!os.whence('dep')) {
    aster.notify('failure', exports.prefix + 'dep', 'dep command not found!');
    return true;
  }
  return run(arguments[0], ['dep'].concat(Array.prototype.slice.call(arguments)));
};

dep.ensure = function ensure() {
  return dep.apply(null, ['ensure'].concat(Array.prototype.slice.call(arguments)));
};

dep.prune = function prune() {
  return dep.apply(null, ['prune'].concat(Array.prototype.slice.call(arguments)));
};

var go = exports.go = function go() {
  return run(arguments[arguments[0] !== 'tool' ? 0 : 1], ['go'].concat(Array.prototype.slice.call(arguments)));
};

go.generate = function generate() {
  return run('generate', ['go', 'generate'].concat(Array.prototype.slice.call(arguments)));
};

go.get = function get() {
  return run('get', ['go', 'get'].concat(Array.prototype.slice.call(arguments)));
};

go.list = function list() {
  var args = [];
  for (var i = 0; i < arguments.length; i++) {
    var a = arguments[i];
    switch (parse(a)) {
      case '-json':
      case '-json=':
      case '-f=':
        break;
      case '-f':
        i++;
        break;
      default:
        args.push(a);
    }
  }
  var out = [];
  run('list', ['go', 'list', '-f', '{{.Dir}}'].concat(args), { stdout: out });
  return out;
};

go.test = function test() {
  var args = ['go', 'test'];
  for (var i = 0; i < arguments.length; i++) {
    var a = arguments[i];
    switch (parse(a)) {
      case '-race':
      case '-race=':
        if (aster.arch === 'amd64') {
          args.push(a);
        }
        break;
      default:
        args.push(a);
    }
  }
  return run('test', args);
};

go.tool = {
  cover: function cover() {
    var cmd = 'cover';
    var args = ['go', 'tool', 'cover'];
    for (var i = 0; i < arguments.length; i++) {
      var a = arguments[i];
      switch (parse(a)) {
        case '-func':
        case '-func=':
          cmd += ' -func';
          break;
        case '-html':
        case '-html=':
          cmd += ' -html';
          break;
        default:
      }
      args.push(a);
    }
    return run(cmd, args);
  },
};

go.vet = function vet() {
  return run('vet', ['go', 'vet'].concat(Array.prototype.slice.call(arguments)));
};

exports.combine = function combine(args) {
  var out = os.open(args.out, 'w');
  out.write('mode: atomic\n');
  go.list.apply(null, args.packages).forEach(function(p) {
    try {
      var f = os.open(path.join(p, args.profile));
      f.readLine();
      for (;;) {
        var rv = f.readLine();
        if (rv.eof) break;
        out.write(rv.buffer + '\n');
      }
      f.close();
    } catch (ex) {
      // ignore
    }
  });
  out.close();
  return args.out;
};

exports.packagesOf = function packagesOf(files) {
  // changed packages
  var pkgs = files.map(function(f) {
    return ('./' + f.split(/[/\\]+/).slice(0, -1).join('/')).replace(/\/+$/, '');
  }).filter(function(e, i, a) {
    return a.indexOf(e) === i;
  });
  // list packages
  var list = [];
  if (os.system(['go', 'list', '-f', '{{.Dir}}\t{{.ImportPath}}\t{{join .Imports ","}},{{join .TestImports ","}},{{join .XTestImports ","}}', './...'], { stdout: list })) {
    return [];
  }
  var i2p = {};
  var wd = os.getwd();
  list = list.map(function(l) {
    l = l.split(/\t/);
    var p = {
      dir: '.' + l[0].slice(wd.length).replace(/\\/g, '/'),
      importPath: l[1],
      imports: l[2].split(/,/).filter(function(e, i, a) {
        return e && a.indexOf(e) === i;
      }),
    };
    i2p[p.importPath] = p;
    return p;
  });
  // reversed dependencies of packages
  var rdeps = {};
  list.forEach(function(p) {
    p.imports.forEach(function(i) {
      if (i !== p.importPath
          && i in i2p) {
        i = i2p[i];
        if (!(i.dir in rdeps)) {
          rdeps[i.dir] = [];
        }
        rdeps[i.dir].push(p.dir);
      }
    });
  });
  // packages (with dependencies) in order
  var resolve = function(p, seen) {
    var deps = [p];
    if (!seen) {
      seen = {};
    } else if (p in seen) {
      return deps;
    }
    seen[p] = true;
    if (p in rdeps) {
      rdeps[p].forEach(function(d) {
        Array.prototype.push.apply(deps, resolve(d, seen));
      });
    }
    return deps;
  };
  return pkgs.reduce(function(a, b) {
    return list.some(function(p) { return p.dir === b; }) ? a.concat(resolve(b)) : a;
  }, []).filter(function(e, i, a) {
    return a.indexOf(e) === i;
  }).sort(function(a, b) {
    if (a in rdeps
        && rdeps[a].indexOf(b) !== -1) {
      return -1;
    }
    return a === b ? 0 : a < b ? -1 : 1;
  });
};
`),
	"language.js": []byte(`//
// aster :: language.js
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
var path = require('path');

exports.prefix = path.basename(os.getwd()) + ': ';
exports.prompt = '> ';
`),
	"os.js": []byte(`//
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
`),
}
