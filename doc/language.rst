language
========

``language`` module provides customizable options for ``language`` submodules.

.. contents::


language.prefix
---------------

``language.prefix`` is a ``String``. It is used as the prefix of notifications
title.

Default value is the name of the current working directory with ``:_``
[#space]_.


language.prompt
---------------

``language.prompt`` is a ``String``. It is used as the prefix of a command.

Default value is ``>_`` [#space]_.


language.system(object)
-----------------------

``language.system`` is a wrapper of ``os.system``.

object
  ``object`` is an ``Object``.

  args
    ``args`` is an ``Array`` of ``String``.

  options
    ``options`` is an ``Object``.

  title
    ``title`` is a ``String``. It is used as a title of notifications.

  success
    ``success`` is a ``String``. It is used as a message of notifications on
    success.

  failure
    ``failure`` is a ``String``. It is used as a message of notifications on
    failure.


.. [#space] An underscore ``_`` represents a space.
