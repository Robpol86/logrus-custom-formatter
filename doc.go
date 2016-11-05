/*
Package lcf (logrus-custom-formatter) is a customizable formatter for https://github.com/Sirupsen/logrus that lets you
choose which columns to include in your log outputs.

Windows support tested on Windows 10 after May 2016 with native ANSI color support. Previous versions of Windows won't
display actual colors unless os.Stdout/err is intercepted and win32 API calls are made by another library. More info:
https://github.com/Robpol86/colorclass/blob/c7ed6d/colorclass/windows.py#L113

Example Usage

Below is a simple example program that uses lcf with logrus:

	package main

	import (
		lcf "github.com/Robpol86/logrus-custom-formatter"
		"github.com/Sirupsen/logrus"
	)

	func main() {
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

And its output is:

	DEBU[0000] A group of walrus emerges from the ocean      animal=walrus size=10
	WARN[0000] The group's number increased tremendously!    animal=walrus size=10
	INFO[0000] A giant walrus appears!                       number=122 omg=true
	ERRO[0000] Tremendously sized cow enters the ocean.
	FATA[0000] The ice breaks!
	exit status 1

Built-In Attributes

These attributes are provided by lcf and can be specified in your template string:

	%[ascTime]s		Timestamp with formatting defined in CustomFormatter.TimestampFormat.
	%[fields]s		Logrus fields formatted as "key1=value key2=value". Keys sorted unless
				CustomFormatter.DisableSorting is true.
	%[levelName]s		The capitalized log level name (e.g. INFO, WARNING, ERROR).
	%[message]s		The log message.
	%[name]s		The value of the "name" field. If this attribute it used "name" will be
				omitted in %[fields]s.
	%[process]d		The current PID of the running process emitting log statements.
	%[relativeCreated]d	Number of seconds since the program has started.
	%[shortLevelName]s	Like %[levelName]s except WARNING is shown as "WARN".
*/
package lcf
