Aster Changelog
===============

Version 0.5
-----------

* Improve ``language/python`` module.

  * Add ``coverage.json`` function.
  * Add ``coverage.lcov`` function.

* Improve ``language/vimscript`` module.

  * Add ``primula.annotate`` function.
  * Add ``primula.combine`` function.
  * Add ``primula.erase`` function.
  * Add ``primula.html`` function.
  * Add ``primula.json`` function.
  * Add ``primula.lcov`` function.
  * Add ``primula.report`` function.
  * Add ``primula.run`` function.
  * Add ``primula.xml`` function.


Version 0.4
-----------

Release date: 2024-03-13

* Drop Go 1.13 support.
* Drop Go 1.14 support.
* Drop Go 1.15 support.
* Drop Go 1.16 support.
* Drop Go 1.17 support.
* Drop Go 1.18 support.
* Drop Go 1.19 support.


Version 0.3
-----------

Release date: 2020-12-24

* Improve ``language/go`` module.

  * Fix ``dep.prune`` function.
  * Add ``go.build`` function.
  * Add ``go.env`` function.
  * Add ``go.fix`` function.
  * Add ``go.fmt`` function.
  * Add ``go.install`` function.
  * Add ``go.run`` function.
  * Add ``go.mod.download`` function.
  * Add ``go.mod.tidy`` function.
  * Add ``go.mod.vendor`` function.

* Drop Go 1.12 support.
* Add ``language/markdown`` module.


Version 0.2
-----------

Release date: 2018-01-24

* Ignore CVS directory by default.
* Move the aster command to ``github.com/hattya/aster/cmd/aster``.
* Fix deadlock on Windows.
* Introduce the Node.js module loading system.
* Add ``-n`` flag to the aster command.
* ``os`` is now as a module. ``os`` global object is still present for
  backward compatibility.
* Add ``language`` module.
* Add ``language/go`` module.
* Add ``language/javascript`` module.
* Add ``language/python`` module.
* Add ``language/restructuredtext`` module.
* Add ``language/vimscript`` module.


Version 0.1
-----------

Release date: 2017-09-28

* Initial release.
