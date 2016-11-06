package lcf

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestWindowsNativeANSI(t *testing.T) {
	assert := require.New(t)
	assert.False(WindowsNativeANSI())
}

func TestWindowsEnableNativeANSI(t *testing.T) {
	assert := require.New(t)

	if runtime.GOOS == "windows" {
		assert.EqualError(WindowsEnableNativeANSI(true), "The handle is invalid.")
		assert.EqualError(WindowsEnableNativeANSI(false), "The handle is invalid.")
	} else {
		assert.EqualError(WindowsEnableNativeANSI(false), "Not available on this platform.")
	}
}

func TestColor(t *testing.T) {
	formatter := NewFormatter("", nil)
	formatter.ForceColors = true
	formatter.ColorFatal = AnsiReset // Testing no colors on fatal.
	entry := logrus.NewEntry(logrus.New())

	testCases := map[logrus.Level]string{
		logrus.DebugLevel: "\033[36mTest\033[0m",
		logrus.InfoLevel:  "\033[32mTest\033[0m",
		logrus.WarnLevel:  "\033[33mTest\033[0m",
		logrus.ErrorLevel: "\033[31mTest\033[0m",
		logrus.PanicLevel: "\033[35mTest\033[0m",
		logrus.FatalLevel: "Test",
	}

	for level, expected := range testCases {
		t.Run(fmt.Sprintf("level:%v", level), func(t *testing.T) {
			assert := require.New(t)
			entry.Level = level
			actual := Color(entry, formatter, "Test")
			assert.Equal(expected, actual)
		})
	}
}
