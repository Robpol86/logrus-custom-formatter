package lcf

import (
	"bytes"
	"io"
	"os"

	"github.com/sirupsen/logrus"
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
func LogMsgs() {
	logrus.Debug("Sample debug 1.")
	logrus.WithFields(logrus.Fields{"name": CallerName(1), "a": "b", "c": 10}).Debug("Sample debug 2.")
	logrus.Info("Sample info 1.")
	logrus.WithFields(logrus.Fields{"name": CallerName(1), "a": "b", "c": 10}).Info("Sample info 2.")
	logrus.Warn("Sample warn 1.")
	logrus.WithFields(logrus.Fields{"name": CallerName(1), "a": "b", "c": 10}).Warn("Sample warn 2.")
	logrus.Error("Sample error 1.")
	logrus.WithFields(logrus.Fields{"name": CallerName(1), "a": "b", "c": 10}).Error("Sample error 2.")
}
