language/go
===========

.. contents::


go.dep(...args)
---------------

``go.dep`` represents ``dep``, and its commands are defined as properties of
this function.

Requirements
  - `dep <https://github.com/golang/dep>`_


dep.ensure(...args)
~~~~~~~~~~~~~~~~~~~

``dep.ensure`` represents ``dep ensure``.


dep.prune(...args)
~~~~~~~~~~~~~~~~~~

``dep.prune`` represents ``dep prune``.


go.go(...args)
--------------

``go.go`` represents ``go``, and its commands are defined as properties of this
function.


go.build(...args)
~~~~~~~~~~~~~~~~~

``go.build`` represents ``go build``.


go.env(...args)
~~~~~~~~~~~~~~~

``go.env`` represents ``go env``, and returns the output.


go.fix(...args)
~~~~~~~~~~~~~~~

``go.fix`` represents ``go fix``.


go.fmt(...args)
~~~~~~~~~~~~~~~

``go.fmt`` represents ``go fmt``.


go.generate(...args)
~~~~~~~~~~~~~~~~~~~~

``go.generate`` represents ``go generate``.


go.get(...args)
~~~~~~~~~~~~~~~

``go.get`` represents ``go get``.


go.install(...args)
~~~~~~~~~~~~~~~~~~~

``go.install`` represents ``go install``.


go.list(...args)
~~~~~~~~~~~~~~~~

``go.list`` represents ``go list``, and returns ``-f {{.Dir}}`` of the
specified packages.

It eliminates ``-f`` flag and ``-json`` flag.


go.mod.download(...args)
~~~~~~~~~~~~~~~~~~~~~~~~

``go.mod download`` represents ``go mod download``.


go.mod.tidy(...args)
~~~~~~~~~~~~~~~~~~~~

``go.mod tidy`` represents ``go mod tidy``.


go.mod.vendor(...args)
~~~~~~~~~~~~~~~~~~~~~~

``go.mod vendor`` represents ``go mod vendor``.


go.run(...args)
~~~~~~~~~~~~~~~

``go.run`` represents ``go run``.


go.test(...args)
~~~~~~~~~~~~~~~~

``go.test`` represents ``go test``.

It eliminates ``-race`` flag when the runtime architecture is not x64.


go.tool.cover(...args)
~~~~~~~~~~~~~~~~~~~~~~

``go.tool.cover`` represents ``go tool cover``.


go.vet(...args)
~~~~~~~~~~~~~~~

``go.vet`` represents ``go vet``.


go.combine(object)
------------------

``go.combine`` combines the specified coverage profiles, and returns the name
of the combined coverage profile.

object
  ``object`` is an ``Object``.

  packages
    ``packages`` is an ``Array`` of package names.

  profile
    ``profile`` is a name of the coverage profile to be search.

  out
    ``out`` is a name of the combined coverage profile. It will be overwritten
    if exists.


go.packagesOf(...files)
-----------------------

``go.packagesOf`` returns an ``Array`` of package names and its dependencies
in order of dependency.

files
  ``files`` is an ``Array`` of ``Stirng``.
