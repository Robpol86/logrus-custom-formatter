package formatter

import (
	"strings"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestNewFormatterBasic(t *testing.T) {
	assert := require.New(t)
	defer ResetLogger() // Cleanup after test.

	_, stderr, err := WithCapSys(func() {
		ResetLogger()
		logrus.SetFormatter(NewFormatter(Basic))
		logrus.SetLevel(logrus.DebugLevel)
		LogMsgs(true)
	})
	assert.NoError(err)

	actual := strings.Split(stderr, `\n`)
	expected := []string{
		"DEBUG:LogMsgs:Sample debug 1.",
		"INFO:LogMsgs:Sample info 1.",
		"WARN:LogMsgs:Sample warn 1.",
		"ERROR:LogMsgs:Sample error 1.",
		"",
	}
	assert.Equal(expected, actual)
}
