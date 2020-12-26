# Aster

Aster is a command line tool to handle events on file system modifications. It
is inspired by [Guard](https://guardgem.org/).

[![pkg.go.dev](https://pkg.go.dev/badge/github.com/hattya/aster.svg)](https://pkg.go.dev/github.com/hattya/aster)
[![GitHub Actions](https://github.com/hattya/aster/workflows/CI/badge.svg)](https://github.com/hattya/aster/actions?query=workflow:CI)
[![Semaphore](https://semaphoreci.com/api/v1/hattya/aster/branches/master/badge.svg)](https://semaphoreci.com/hattya/aster)
[![Appveyor](https://ci.appveyor.com/api/projects/status/qc3luxk7q7jmx2ut/branch/master?svg=true)](https://ci.appveyor.com/project/hattya/aster)
[![Codecov](https://codecov.io/gh/hattya/aster/branch/master/graph/badge.svg)](https://codecov.io/gh/hattya/aster)


## Installation

```console
$ go get -u github.com/hattya/aster/cmd/aster
```


## Usage

```console
$ aster -g
```


### init

```console
$ aster init [<template>...]
```

``aster init`` creates an Asterfile in the current directory if it does not
exist, and add specified template files to it.

Template files are located in:

- UNIX  
  `$XDG_CONFIG_HOME/aster/template/<template>`

- macOS  
  `~/Library/Application Support/Aster/template/<template>`

- Windows  
  `%APPDATA%\Aster\template\<template>`


## Asterfile

Asterfile is evaluated as JavaScript by [otto](https://github.com/robertkrimen/otto).

```javascript
var go = require('language/go').go;

aster.watch(/.+\.go$/, function() {
  // test
  if (go.test('-v', '-covermode', 'atomic', '-coverprofile', 'cover.out', './...')) {
    return;
  }
  // coverage report
  go.tool.cover('-func', 'cover.out');
  go.tool.cover('-html', 'cover.out', '-o', 'coverage.html');
  // vet
  if (go.vet('./...')) {
    return;
  }
});
```


## Reference

- [Global Objects](doc/global-objects.rst)
- [OS](doc/os.rst)
- [Language](doc/language.rst)
  - [Go](doc/language/go.rst)
  - [JavaScript](doc/language/javascript.rst)
  - [Markdown](doc/language/markdown.rst)
  - [Python](doc/language/python.rst)
  - [reStructuredText](doc/language/restructuredtext.rst)
  - [Vim script](doc/language/vimscript.rst)


## License

Aster is distributed under the terms of the MIT License.
