=======================
logrus-custom-formatter
=======================

Customizable Logrus formatter similar in style to Python's
`logging.Formatter <https://docs.python.org/3.6/library/logging.html#logrecord-attributes>`_.

* Tested with Golang 1.7 on Linux, OS X, and Windows.

📖 Full documentation: https://godoc.org/github.com/Robpol86/logrus-custom-formatter

.. image:: https://img.shields.io/appveyor/ci/Robpol86/logrus-custom-formatter/master.svg?style=flat-square&label=AppVeyor%20CI
    :target: https://ci.appveyor.com/project/Robpol86/logrus-custom-formatter
    :alt: Build Status Windows

.. image:: https://img.shields.io/travis/Robpol86/logrus-custom-formatter/master.svg?style=flat-square&label=Travis%20CI
    :target: https://travis-ci.org/Robpol86/logrus-custom-formatter
    :alt: Build Status

.. image:: https://img.shields.io/codecov/c/github/Robpol86/logrus-custom-formatter/master.svg?style=flat-square&label=Codecov
    :target: https://codecov.io/gh/Robpol86/logrus-custom-formatter
    :alt: Coverage Status

Quickstart
==========

Install:

.. code:: bash

    go get github.com/Robpol86/logrus-custom-formatter

Usage:

.. code:: go

    package main

    import (
        lcf "github.com/Robpol86/logrus-custom-formatter"
        "github.com/Sirupsen/logrus"
    )

    func main() {
        lcf.WindowsEnableNativeANSI(true) // Ignored on non-Windows.
        template := "%[shortLevelName]s[%04[relativeCreated]d] %-45[message]s%[fields]s\n"
        logrus.SetFormatter(lcf.NewFormatter(template, nil))
        logrus.SetLevel(logrus.DebugLevel)

        animal := logrus.Fields{"animal": "walrus", "size": 10}
        logrus.WithFields(animal).Debug("A group of walrus emerges from the ocean")
        logrus.WithFields(animal).Warn("The group's number increased tremendously!")
        number := logrus.Fields{"number": 122, "omg": true}
        logrus.WithFields(number).Info("A giant walrus appears!")
        logrus.Error("Tremendously sized cow enters the ocean.")
        logrus.Fatal("The ice breaks!")
    }

.. changelog-section-start

Changelog
=========

This project adheres to `Semantic Versioning <http://semver.org/>`_.

0.0.1 - 2016-10-27
------------------

* Still in development.

.. changelog-section-end
