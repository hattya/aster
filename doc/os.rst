os
==

.. contents::


os.getenv(key)
~~~~~~~~~~~~~~

``os.getenv()`` retrieves the value of the environment variable named by the
``key``. It returns an empty ``String`` if the variable is not present.

key
    ``key`` is a ``String``.


os.getwd()
~~~~~~~~~~

``os.getwd`` returns an absolute path of the current directory. It returns an
empty ``String`` if fails.


os.mkdir(path[, perm=0777])
~~~~~~~~~~~~~~~~~~~~~~~~~~~

``os.mkdir`` creates a directory named ``path``, along with any necessary
parent directories. It returns ``true`` if fails.

path
    ``path`` is a ``String``.

perm
    ``perm`` is a permission bits which are used for all directories that
    ``os.mkdir`` creates.


.. _`os.open`:

os.open(path[, mode='r'])
~~~~~~~~~~~~~~~~~~~~~~~~~

``os.open`` opens a file named ``path`` with ``mode``, and returns an instance
of |os.File|_. It returns ``undefined`` if fails.

path
    ``path`` is a ``String``.

mode
    ``mode`` is a ``String``.

    r
        open for reading (default).

    r+
        open for reading and writing.

    w
        open for writing, truncating file to zero length.

    w+
        open for reading and writing, truncating file to zero length.

    a
        open for writing, appending to the end of file.

    a+
        open for reading and writing, appending to the end of file.


os.remove(path)
~~~~~~~~~~~~~~~

``os.remove`` removes ``path`` and its contents recursively.

path
    ``path`` is a ``String``.


os.rename(src, dst)
~~~~~~~~~~~~~~~~~~~

``os.rename`` renames / moves a file or directory. It returns ``true`` if
fails.

src
    ``src`` is a ``String``.

dst
    ``dst`` is a ``String``.


os.setenv(key, value)
~~~~~~~~~~~~~~~~~~~~~

``os.setenv`` sets the ``value`` of the environment variable named by the ``key``.

key
    ``key`` is a ``String``.

value
    ``value`` is a ``String``.


os.stat(path)
~~~~~~~~~~~~~

``os.stat`` returns an instance of |os.FileInfo|_ which describes the ``path``.
It returns ``undefined`` if fails.

path
    ``path`` is a ``String``.

.. |os.FileInfo| replace:: ``os.FileInfo``
.. _os.FileInfo: `class os.FileInfo`_


os.system(args[, options])
~~~~~~~~~~~~~~~~~~~~~~~~~~~

``os.system`` runs the command specified by ``args``. It returns ``true`` if
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
            The standard output will be split into lines, and added to the
            ``Array``.

    stderr
        ``stderr`` is a ``String``, ``null``, or an ``Array``.

        ``String``
            It is the file name to redirect the standard error. *It will be
            overwritten if exists.*

        ``null``
            The standard error will be discarded.

        ``Array``
            The standard error will be split into lines, and added to the
            ``Array``.


os.whence(name)
~~~~~~~~~~~~~~~

``os.whence`` searches for ``name`` in the directories named by the PATH
environment variable. It returns the path of ``name`` if found, ``undefined``
otherwise.

name
    ``name`` to search.


class os.File
~~~~~~~~~~~~~

File.prototype.close()
""""""""""""""""""""""

``close`` closes the |os.File|_.


File.prototype.name()
"""""""""""""""""""""

``name`` returns the name of the file which specified to |os.open|_.

.. |os.open| replace:: ``os.open``


File.prototype.read(n)
""""""""""""""""""""""

``read`` reads up to ``n`` bytes from the |os.File|_, and returns an
``Object``.

n
    ``n`` is a ``Number``.

Return value
    eof
        It is ``true`` when at the end of the file.

    buffer
        It is a ``String`` which read from the file.


File.prototype.readLine()
"""""""""""""""""""""""""

``readLine`` reads a line from the |os.File|_, and returns an ``Object``.

Return value
    eof
        It is ``true`` when at the end of the file.

    buffer
        It is a ``String`` which read from the file.


File.prototype.write(data)
""""""""""""""""""""""""""

``write`` writes the ``data`` to the |os.File|_.

data
    ``data`` is a ``String``.

.. |os.File| replace:: ``os.File``
.. _os.File: `class os.File`_


class os.FileInfo
~~~~~~~~~~~~~~~~~

FileInfo.name
"""""""""""""

base name of the file.


FileInfo.size
"""""""""""""

file size, in bytes.


FileInfo.mode
"""""""""""""

file mode bits.


FileInfo.mtime
""""""""""""""

time of last modification. It is a ``Date``.


FileInfo.prototype.isDir()
""""""""""""""""""""""""""

``isDir`` reports whether the file is a directory.


FileInfo.prototype.isRegular()
""""""""""""""""""""""""""""""

``isRegular`` reports whether the file is a regular file.


FileInfo.prototype.perm()
"""""""""""""""""""""""""

``perm`` returns the permission bits.
