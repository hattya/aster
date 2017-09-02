Aster
=====

Aster is a command line tool to handle events on file system modifications. It
is inspired by Guard_.

.. image:: https://semaphoreci.com/api/v1/hattya/aster/branches/master/badge.svg
   :target: https://semaphoreci.com/hattya/aster

.. image:: https://ci.appveyor.com/api/projects/status/qc3luxk7q7jmx2ut/branch/master?svg=true
   :target: https://ci.appveyor.com/project/hattya/aster

.. _Guard: http://guardgem.org/


Installation
------------

.. code:: console

   $ go get -u github.com/hattya/aster


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

   aster.watch(/.+\.go$/, function(files) {
     // build
     if (os.system('go get -t -v ./...'.split(/\s+/))) {
       aster.notify('failure', 'build', 'failure');
       return;
     }
     aster.notify('success', 'build', 'success');

     // test
     if (os.system('go test -v -cover -coverprofile cover.out ./...'.split(/\s+/))) {
       aster.notify('failure', 'test', 'failure');
       return;
     }
     aster.notify('success', 'test', 'success');

     // coverage
     os.system('go tool cover -func cover.out'.split(/\s+/));
     os.system('go tool cover -html cover.out -o coverage.html'.split(/\s+/));
   });

.. _otto: https://github.com/robertkrimen/otto


Reference
---------

aster.watch(pattern, callback)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

``aster.watch`` defines which files should be watched by Aster.

pattern
    ``pattern`` is a ``RegExp``.

callback
    ``callback`` is a ``Function``. It is invoked on each file system
    modifications when ``pattern`` is matched.

    ``callback`` is invoked with one argument:

    * ``Array`` of paths


aster.notify(name, title, text)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

``aster.notify`` sends a GNTP_ notification. It does nothing when ``-g`` flag
is not specified to Aster.

name
    ``name`` is the name (type) of a GNTP_ notification and which is either
    ``success`` or ``failure``.

title
    ``title`` is the title of a GNTP_ notification.

text
    ``text`` is the text of a GNTP_ notification.


aster.title(title)
~~~~~~~~~~~~~~~~~~

``aster.title`` sets the title of an XTerm, or the title of the console window
on Windows.


aster.arch
~~~~~~~~~~

``aster.arch`` is a ``String``. It is a synonym of |runtime.GOARCH|_.

.. |runtime.GOARCH| replace:: ``runtime.GOARCH``
.. _runtime.GOARCH: runtime_


aster.os
~~~~~~~~

``aster.os`` is a ``String``. It is a synonym of |runtime.GOOS|_.

.. |runtime.GOOS| replace:: ``runtime.GOOS``
.. _runtime.GOOS: runtime_


aster.ignore
~~~~~~~~~~~~

``aster.ignore`` is an ``Array`` of ``RegExp``. It will be ignored recursively
by Aster when a directory is matched to any of ``aster.ignore``.

A path to be matched is a relative path from where the Asterfile exists.


os.getwd()
~~~~~~~~~~

``os.getwd`` returns an absolute path of the current directory, or an empty
``String`` if fails.


os.mkdir(path[, perm=0777])
~~~~~~~~~~~~~~~~~~~~~~~~~~~

``os.mkdir`` creates a directory named ``path``, along with any necessary
parent directories, and returns ``true`` if fails.

path
    ``path`` is a ``String``.

perm
    ``perm`` is a permission bits which are used for all directories that
    ``os.mkdir`` creates.


os.remove(path)
~~~~~~~~~~~~~~~

``os.remove`` removes ``path`` and its contents recursively.

path
    ``path`` is a ``String``.


os.rename(src, dst)
~~~~~~~~~~~~~~~~~~~

``os.rename`` renames / moves a file or directory.

src
    ``src`` is a ``String``.

dst
    ``dst`` is a ``String``.


os.stat(path)
~~~~~~~~~~~~~

``os.stat`` returns a ``os.FileInfo`` which describes the path.


os.system(args[, options])
~~~~~~~~~~~~~~~~~~~~~~~~~~~

``os.system`` runs the command specified by ``args``, and returns ``true`` if
fails.

args
    ``args`` is an ``Array`` of ``String``.

options
    ``options`` is an ``Object``.

    dir
        ``dir`` is the working directory of the command.

    stdout
        ``stdout`` is a ``String``, ``null``, or an ``Array``.

        ``String``
            It is the file name to redirect the standard output. *It will be
            overwritten if exists.*

        ``null``
            The standard output will be discarded.

        ``Array``
            The standard output will be splitted into lines, and added to the
            ``Array``.

    stderr
        ``stderr`` is a ``String``, ``null``, or an ``Array``.

        ``String``
            It is the file name to redirect the standard error. *It will be
            overwritten if exists.*

        ``null``
            The standard error will be discarded.

        ``Array``
            The standard error will be splitted into lines, and added to the
            ``Array``.


os.whence(name)
~~~~~~~~~~~~~~~

``os.whence`` searches for ``name`` in the directories named by the PATH
environment variable. It returns the path of ``name`` if found, ``undefined``
otherwise.

name
    ``name`` to search.


class os.FileInfo
~~~~~~~~~~~~~~~~~

name
    base name of the file.

size
    file size, in bytes.

mode
    file mode bits.

mtime
    time of last miodification. It is a ``Date``.

isDir()
    ``FileInfo.isDir`` reports whether the file is a directory.

isRegular()
    ``FileInfo.isRegular`` reports whether the file is a regular file.

perm()
    ``FileInfo.perm`` returns the permission bits.

.. _GNTP: http://growl.info/documentation/developer/gntp.php
.. _runtime: https://golang.org/pkg/runtime/#pkg-constants


License
-------

Aster is distributed under the terms of the MIT License.
