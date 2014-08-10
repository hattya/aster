aster
=====

aster is a command line tool to handle events on file system modifications. It
is inspired by Guard_.

.. _Guard: http://guardgem.org/


Install
-------

.. code:: console

   $ go get -u github.com/hattya/aster


Usage
-----

.. code:: console

   $ aster -g


Asterfile
---------

Asterfile is evaluated as JavaScript by otto_.

.. code:: javascript

   aster.watch(/.+\.go$/, function(files) {
     // build
     if (os.system("go", "get", "-t", "-v", "./...")) {
       aster.notify("failure", "build", "failure");
       return;
     }
     aster.notify("success", "build", "success");

     // test
     if (os.system("go", "test", "-v", "-cover", "-coverprofile cover.out", "./...")) {
       aster.notify("failure", "test", "failure");
       return;
     }
     aster.notify("success", "test", "success");

     // coverage
     os.system("go", "tool", "cover", "-func cover.out");
     os.system("go", "tool", "cover", "-html cover.out", "-o coverage.html");
   });

.. _otto: https://github.com/robertkrimen/otto


Reference
---------

aster.watch(pattern, callback)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

``aster.watch`` defines which files should be watched by aster.

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
is not specified to aster.

name
    ``name`` is a name (type) of a GNTP_ notification and which is either
    ``success`` or ``failure``.

title
    ``title`` is a title of a GNTP_ notification.

text
    ``text`` is a text of a GNTP_ notification.


os.system(name[, ... argN])
~~~~~~~~~~~~~~~~~~~~~~~~~~~

``os.system`` spawns a new process and returns ``true`` when it is failed.


.. _GNTP: http://growl.info/documentation/developer/gntp.php


License
-------

aster is distributed under the terms of the MIT License.
