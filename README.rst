Aster
=====

Aster is a command line tool to handle events on file system modifications. It
is inspired by Guard_.

.. image:: https://semaphoreci.com/api/v1/hattya/aster/branches/master/badge.svg
   :target: https://semaphoreci.com/hattya/aster

.. image:: https://ci.appveyor.com/api/projects/status/qc3luxk7q7jmx2ut/branch/master?svg=true
   :target: https://ci.appveyor.com/project/hattya/aster

.. image:: https://codecov.io/gh/hattya/aster/branch/master/graph/badge.svg
   :target: https://codecov.io/gh/hattya/aster

.. _Guard: http://guardgem.org/


Installation
------------

.. code:: console

   $ go get -u github.com/hattya/aster/cmd/aster


Usage
-----

.. code:: console

   $ aster -g


init
~~~~

.. code:: console

   $ aster init [<template>...]

``aster init`` creates an Asterfile in the current directory if it does not
exist, and add specified template files to it.

Template files are located in:

UNIX
    $XDG_CONFIG_HOME/aster/template/<template>

Windows
    %APPDATA%\\Aster\\template\\<template>


Asterfile
---------

Asterfile is evaluated as JavaScript by otto_.

.. code:: javascript

   var os = require('os');
   var go = require('language/go').go;

   aster.watch(/.+\.go$/, function() {
     // build
     if (go.get('-t', '-v', './...')) {
       return;
     }
     // test
     if (go.test('-v', '-covermode', 'atomic', './...')) {
       return;
     }
   });

.. _otto: https://github.com/robertkrimen/otto


Reference
---------

* `Global Objects <doc/global-objects.rst>`_
* `OS <doc/os.rst>`_
* `Language <doc/language.rst>`_

  * `Go <doc/language.go.rst>`_
  * `JavaScript <doc/language.javascript.rst>`_
  * `reStructuredText <doc/language.restructuredtext.rst>`_


License
-------

Aster is distributed under the terms of the MIT License.
