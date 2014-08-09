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


License
-------

aster is distributed under the terms of the MIT License.
