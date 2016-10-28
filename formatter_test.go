package formatter

import (
	"bytes"
	"io"
	"os"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/require"
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

// resetLogger re-initializes the global logrus logger so stdout/stderr changes are applied to it.
// Otherwise after patching the streams logrus still points to the original file descriptor.
func resetLogger() {
	*log.StandardLogger() = *log.New()
}

func logMsgs() {
	log.Debug("Sample debug 1.")
	log.WithFields(log.Fields{"a": "b", "c": 10}).Debug("Sample debug 2.")
	log.Info("Sample info 1.")
	log.WithFields(log.Fields{"a": "b", "c": 10}).Info("Sample info 2.")
	log.Warn("Sample warn 1.")
	log.WithFields(log.Fields{"a": "b", "c": 10}).Warn("Sample warn 2.")
	log.Error("Sample error 1.")
	log.WithFields(log.Fields{"a": "b", "c": 10}).Error("Sample error 2.")
}

func TestNewFormatterBasic(t *testing.T) {
	assert := require.New(t)
	defer resetLogger() // Cleanup after test.

	_, stderr, err := WithCapSys(func() {
		resetLogger()
		log.SetFormatter(NewFormatter(Basic))
		log.SetLevel(log.DebugLevel)
		logMsgs()
	})
	assert.NoError(err)

	assert.Equal("TODO", stderr)
}
