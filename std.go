// Code generated by "modulizer -l std -o std.go lib"; DO NOT EDIT.

package aster

import "github.com/hattya/otto.module"

type stdLoader struct {
}

func (l *stdLoader) Load(id string) ([]byte, error) {
	for _, ext := range []string{"", ".js", ".json"} {
		if b, ok := files[id+ext]; ok {
			return b, nil
		}
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
//   Copyright (c) 2017-2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

'use strict';

var os = require('os');
var path = require('path');
var language = require('language');

function system(cmd, args, options) {
  return language.system({
    args: args,
    options: options,
    title: args[0],
    success: cmd + ' passed',
    failure: cmd + ' failed',
  });
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
    aster.notify('failure', language.prefix + 'dep', 'dep not found!');
    return true;
  }
  return system(arguments[0], ['dep'].concat(Array.prototype.slice.call(arguments)));
};

dep.ensure = function ensure() {
  return dep.apply(null, ['ensure'].concat(Array.prototype.slice.call(arguments)));
};

dep.prune = function prune() {
  return dep.apply(null, ['prune'].concat(Array.prototype.slice.call(arguments)));
};

var go = exports.go = function go() {
  var cmd;
  switch (arguments[0]) {
    case 'mod':
      cmd = arguments[0] + ' ' + arguments[1];
      break;
    case 'tool':
      cmd = arguments[1];
      break;
    default:
      cmd = arguments[0];
  }
  return system(cmd, ['go'].concat(Array.prototype.slice.call(arguments)));
};

go.build = function build() {
  return go.apply(null, ['build'].concat(Array.prototype.slice.call(arguments)));
};

go.env = function env() {
  var out = [];
  system('env', ['go', 'env'].concat(Array.prototype.slice.call(arguments)), { stdout: out });
  return out;
};

go.fix = function fix() {
  return go.apply(null, ['fix'].concat(Array.prototype.slice.call(arguments)));
};

go.fmt = function fmt() {
  return go.apply(null, ['fmt'].concat(Array.prototype.slice.call(arguments)));
};

go.generate = function generate() {
  return go.apply(null, ['generate'].concat(Array.prototype.slice.call(arguments)));
};

go.get = function get() {
  return go.apply(null, ['get'].concat(Array.prototype.slice.call(arguments)));
};

go.install = function install() {
  return go.apply(null, ['install'].concat(Array.prototype.slice.call(arguments)));
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
  system('list', ['go', 'list', '-f', '{{.Dir}}'].concat(args), { stdout: out });
  return out;
};

go.mod = {
  download: function download() {
    return go.apply(null, ['mod', 'download'].concat(Array.prototype.slice.call(arguments)));
  },

  tidy: function tidy() {
    return go.apply(null, ['mod', 'tidy'].concat(Array.prototype.slice.call(arguments)));
  },

  vendor: function vendor() {
    return go.apply(null, ['mod', 'vendor'].concat(Array.prototype.slice.call(arguments)));
  },
};

go.run = function run() {
  return go.apply(null, ['run'].concat(Array.prototype.slice.call(arguments)));
};

go.test = function test() {
  var args = ['test'];
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
  return go.apply(null, args);
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
    return system(cmd, args);
  },
};

go.vet = function vet() {
  return go.apply(null, ['vet'].concat(Array.prototype.slice.call(arguments)));
};

exports.combine = function combine(object) {
  var out = os.open(object.out, 'w');
  out.write('mode: atomic\n');
  go.list.apply(null, object.packages).forEach(function(p) {
    try {
      var f = os.open(path.join(p, object.profile));
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
  return object.out;
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
	"language/javascript.js": []byte(`//
// aster :: language/javascript.js
//
//   Copyright (c) 2017-2020 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

'use strict';

var os = require('os');
var language = require('language');

var npm = exports.npm = function() {
  if (!os.whence('npm')) {
    aster.notify('failure', language.prefix + 'npm', 'npm not found!');
    return true;
  }

  var cmd = arguments[0] !== 'run' ? arguments[0] : arguments[1] + ' script';
  return language.system({
    args: ['npm'].concat(Array.prototype.slice.call(arguments)),
    title: 'npm',
    success: cmd + ' passed',
    failure: cmd + ' failed',
  });
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
`),
	"language/markdown.js": []byte(`//
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
`),
	"language/python.js": []byte(`//
// aster :: language/python.js
//
//   Copyright (c) 2018-2024 Akinori Hattori <hattya@gmail.com>
//
//   SPDX-License-Identifier: MIT
//

'use strict';

var os = require('os');
var language = require('language');

var coverage = exports.coverage = function coverage() {
  if (!os.whence('coverage')) {
    aster.notify('failure', language.prefix + 'coverage', 'coverage not found!');
    return true;
  }

  return language.system({
    args: ['coverage'].concat(Array.prototype.slice.call(arguments)),
    title: 'coverage',
    success: arguments[0] + ' passed',
    failure: arguments[0] + ' failed',
  });
};

coverage.annotate = function annotate() {
  return coverage.apply(null, ['annotate'].concat(Array.prototype.slice.call(arguments)));
};

coverage.combine = function combine() {
  return coverage.apply(null, ['combine'].concat(Array.prototype.slice.call(arguments)));
};

coverage.erase = function erase() {
  return coverage.apply(null, ['erase'].concat(Array.prototype.slice.call(arguments)));
};

coverage.html = function html() {
  return coverage.apply(null, ['html'].concat(Array.prototype.slice.call(arguments)));
};

coverage.json = function json() {
  return coverage.apply(null, ['json'].concat(Array.prototype.slice.call(arguments)));
};

coverage.lcov = function lcov() {
  return coverage.apply(null, ['lcov'].concat(Array.prototype.slice.call(arguments)));
};

coverage.report = function report() {
  return coverage.apply(null, ['report'].concat(Array.prototype.slice.call(arguments)));
};

coverage.run = function run() {
  return coverage.apply(null, ['run'].concat(Array.prototype.slice.call(arguments)));
};

coverage.xml = function xml() {
  return coverage.apply(null, ['xml'].concat(Array.prototype.slice.call(arguments)));
};
`),
	"language/restructuredtext.js": []byte(`//
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
`),
	"language/vimscript.js": []byte(`//
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
`),
	"language.js": []byte(`//
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
`),
	"os.js": []byte(`//
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
`),
}
