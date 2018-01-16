Global Objects
==============

.. contents::


aster
-----

aster.arch
~~~~~~~~~~

``aster.arch`` is a ``String``. It is a synonym of |runtime.GOARCH|_.

.. |runtime.GOARCH| replace:: ``runtime.GOARCH``
.. _runtime.GOARCH: runtime_


aster.ignore
~~~~~~~~~~~~

``aster.ignore`` is an ``Array`` of ``RegExp``. It will be ignored recursively
by Aster when a directory is matched to any of ``aster.ignore``.

A path to be matched is a relative path from where the Asterfile exists.


aster.os
~~~~~~~~

``aster.os`` is a ``String``. It is a synonym of |runtime.GOOS|_.

.. |runtime.GOOS| replace:: ``runtime.GOOS``
.. _runtime.GOOS: runtime_


aster.notify(event, title, body)
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

``aster.notify`` sends a notification. It does nothing when ``-g`` flag and/or
``-n`` flag are not specified to Aster.

event
  ``event`` is the event name of a notification, and which is either
  ``success`` or ``failure``.

title
  ``title`` is the title of a notification.

body
  ``body`` is the body text of a notification.


aster.title(title)
~~~~~~~~~~~~~~~~~~

``aster.title`` sets the title of an XTerm, or the title of the console window
on Windows.


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


.. _runtime: https://golang.org/pkg/runtime/#pkg-constants
