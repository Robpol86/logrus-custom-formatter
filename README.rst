=======================
logrus-custom-formatter
=======================

Customizable Logrus formatter similar in style to Python's
`logging.Formatter <https://docs.python.org/3.6/library/logging.html#logrecord-attributes>`_.

Windows support tested on Windows 10 after May 2016 with native ANSI color support. Previous versions of Windows won't
display actual colors unless os.Stdout/err is intercepted and win32 API calls are made by another library. More info:
https://github.com/Robpol86/colorclass/blob/c7ed6d/colorclass/windows.py#L113

* Tested with Golang 1.7 on Linux, OS X, and Windows.

.. image:: https://img.shields.io/appveyor/ci/Robpol86/logrus-python-formatter/master.svg?style=flat-square&label=AppVeyor%20CI
    :target: https://ci.appveyor.com/project/Robpol86/logrus-python-formatter
    :alt: Build Status Windows

.. image:: https://img.shields.io/travis/Robpol86/logrus-python-formatter/master.svg?style=flat-square&label=Travis%20CI
    :target: https://travis-ci.org/Robpol86/logrus-python-formatter
    :alt: Build Status

.. image:: https://img.shields.io/codecov/c/github/Robpol86/logrus-python-formatter/master.svg?style=flat-square&label=Codecov
    :target: https://codecov.io/gh/Robpol86/logrus-python-formatter
    :alt: Coverage Status

Quickstart
==========

Install:

.. code:: bash

    go get https://github.com/Robpol86/logrus-python-formatter

Usage:

.. code:: go

    // TODO

.. changelog-section-start

Changelog
=========

This project adheres to `Semantic Versioning <http://semver.org/>`_.

0.0.1 - 2016-10-27
------------------

* Still in development.

.. changelog-section-end
