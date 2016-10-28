package formatter

import (
	"bytes"
	"io"
	"os"

	"github.com/Sirupsen/logrus"
)

// WithCapSys temporarily redirects stdout/stderr pipes to capture the output while the function runs. Returns them as
// strings.
func WithCapSys(function func()) (string, string, error) {
	var writeStdout *os.File
	var writeStderr *os.File
	chanStdout := make(chan string)
	chanStderr := make(chan string)

	// Prepare new streams.
	if read, write, err := os.Pipe(); err == nil {
		writeStdout = write
		go func() { var buf bytes.Buffer; io.Copy(&buf, read); chanStdout <- buf.String() }()
		if read, write, err := os.Pipe(); err == nil {
			writeStderr = write
			go func() { var buf bytes.Buffer; io.Copy(&buf, read); chanStderr <- buf.String() }()
		} else {
			return "", "", err
		}
	} else {
		return "", "", err
	}

	// Patch streams.
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	defer func() { os.Stdout = oldStdout; os.Stderr = oldStderr }()
	os.Stdout = writeStdout
	os.Stderr = writeStderr

	// Run.
	function()

	// Collect and return.
	writeStdout.Close()
	writeStderr.Close()
	stdout := <-chanStdout
	stderr := <-chanStderr
	return stdout, stderr, nil
}

// ResetLogger re-initializes the global logrus logger so stdout/stderr changes are applied to it.
// Otherwise after patching the streams logrus still points to the original file descriptor.
func ResetLogger() {
	*logrus.StandardLogger() = *logrus.New()
}

// Log sample messages to logrus.
// :param withLibFields: Creates logger with fields specific to this library.
func LogMsgs(withLibFields bool) {
	var log *logrus.Entry
	if withLibFields {
		log = logrus.WithField("loggerName", "LogMsgs") // TODO reflection method to automate name.
	} else {
		log = logrus.NewEntry(logrus.StandardLogger())
	}

	log.Debug("Sample debug 1.")
	// log.WithFields(logrus.Fields{"a": "b", "c": 10}).Debug("Sample debug 2.")
	log.Info("Sample info 1.")
	// log.WithFields(logrus.Fields{"a": "b", "c": 10}).Info("Sample info 2.")
	log.Warn("Sample warn 1.")
	// log.WithFields(logrus.Fields{"a": "b", "c": 10}).Warn("Sample warn 2.")
	log.Error("Sample error 1.")
	// log.WithFields(logrus.Fields{"a": "b", "c": 10}).Error("Sample error 2.")
}
