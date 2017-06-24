package lcf

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/sirupsen/logrus"
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
		assert.EqualError(WindowsEnableNativeANSI(false), "not available on this platform")
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

func TestCustomFormatter_SprintfNoColors(t *testing.T) {
	assert := require.New(t)

	// len(f.handleColors) == 0
	formatter := NewFormatter("%[shortLevelName]s", nil)
	formatter.ForceColors = true
	assert.Equal("AAA", formatter.Sprintf("AAA"))
	assert.Equal("AA", formatter.Sprintf("AA"))
	assert.Equal("A", formatter.Sprintf("A"))
	assert.Equal("", formatter.Sprintf(""))

	// No \033
	formatter = NewFormatter("%2[shortLevelName]s", nil)
	formatter.ForceColors = true
	assert.Equal("AAA", formatter.Sprintf("AAA"))
	assert.Equal("AA", formatter.Sprintf("AA"))
	assert.Equal(" A", formatter.Sprintf("A"))
	assert.Equal("  ", formatter.Sprintf(""))

	// Multiple no \033
	formatter = NewFormatter("%2[shortLevelName]s %2[shortLevelName]s", nil)
	formatter.ForceColors = true
	assert.Equal("AAA BBB", formatter.Sprintf("AAA", "BBB"))
	assert.Equal("AA BB", formatter.Sprintf("AA", "BB"))
	assert.Equal(" A  B", formatter.Sprintf("A", "B"))
	assert.Equal("     ", formatter.Sprintf("", ""))
}

func TestCustomFormatter_SprintfColors(t *testing.T) {
	assert := require.New(t)

	formatter := NewFormatter("%10[shortLevelName]s", nil)
	formatter.ForceColors = true
	assert.Equal("       \033[31mAAA\033[0m", formatter.Sprintf("\033[31mAAA\033[0m"))
	assert.Equal("        \033[31mAA\033[0m", formatter.Sprintf("\033[31mAA\033[0m"))
	assert.Equal("         \033[31mA\033[0m", formatter.Sprintf("\033[31mA\033[0m"))
	assert.Equal("          \033[31m\033[0m", formatter.Sprintf("\033[31m\033[0m"))
	assert.Equal("          \033[0m", formatter.Sprintf("\033[0m"))

	formatter = NewFormatter("%-10[shortLevelName]s", nil)
	formatter.ForceColors = true
	assert.Equal("\033[31mAAA\033[0m       ", formatter.Sprintf("\033[31mAAA\033[0m"))
	assert.Equal("\033[31mAA\033[0m        ", formatter.Sprintf("\033[31mAA\033[0m"))
	assert.Equal("\033[31mA\033[0m         ", formatter.Sprintf("\033[31mA\033[0m"))
	assert.Equal("\033[31m\033[0m          ", formatter.Sprintf("\033[31m\033[0m"))
	assert.Equal("\033[0m          ", formatter.Sprintf("\033[0m"))

	formatter = NewFormatter("%2[shortLevelName]s", nil)
	formatter.ForceColors = true
	assert.Equal("\033[31mAAA\033[0m", formatter.Sprintf("\033[31mAAA\033[0m"))
	assert.Equal("\033[31mAA\033[0m", formatter.Sprintf("\033[31mAA\033[0m"))
	assert.Equal(" \033[31mA\033[0m", formatter.Sprintf("\033[31mA\033[0m"))
	assert.Equal("  \033[31m\033[0m", formatter.Sprintf("\033[31m\033[0m"))
	assert.Equal("  \033[0m", formatter.Sprintf("\033[0m"))

	formatter = NewFormatter("%-2[shortLevelName]s", nil)
	formatter.ForceColors = true
	assert.Equal("\033[31mAAA\033[0m", formatter.Sprintf("\033[31mAAA\033[0m"))
	assert.Equal("\033[31mAA\033[0m", formatter.Sprintf("\033[31mAA\033[0m"))
	assert.Equal("\033[31mA\033[0m ", formatter.Sprintf("\033[31mA\033[0m"))
	assert.Equal("\033[31m\033[0m  ", formatter.Sprintf("\033[31m\033[0m"))
	assert.Equal("\033[0m  ", formatter.Sprintf("\033[0m"))

	formatter = NewFormatter("%1[shortLevelName]s", nil)
	formatter.ForceColors = true
	assert.Equal("\033[31mAAA\033[0m", formatter.Sprintf("\033[31mAAA\033[0m"))
	assert.Equal("\033[31mAA\033[0m", formatter.Sprintf("\033[31mAA\033[0m"))
	assert.Equal("\033[31mA\033[0m", formatter.Sprintf("\033[31mA\033[0m"))
	assert.Equal(" \033[31m\033[0m", formatter.Sprintf("\033[31m\033[0m"))
	assert.Equal(" \033[0m", formatter.Sprintf("\033[0m"))

	formatter = NewFormatter("%-1[shortLevelName]s", nil)
	formatter.ForceColors = true
	assert.Equal("\033[31mAAA\033[0m", formatter.Sprintf("\033[31mAAA\033[0m"))
	assert.Equal("\033[31mAA\033[0m", formatter.Sprintf("\033[31mAA\033[0m"))
	assert.Equal("\033[31mA\033[0m", formatter.Sprintf("\033[31mA\033[0m"))
	assert.Equal("\033[31m\033[0m ", formatter.Sprintf("\033[31m\033[0m"))
	assert.Equal("\033[0m ", formatter.Sprintf("\033[0m"))

	formatter = NewFormatter("%0[shortLevelName]s", nil)
	formatter.ForceColors = true
	assert.Equal("\033[31mAAA\033[0m", formatter.Sprintf("\033[31mAAA\033[0m"))
	assert.Equal("\033[31mAA\033[0m", formatter.Sprintf("\033[31mAA\033[0m"))
	assert.Equal("\033[31mA\033[0m", formatter.Sprintf("\033[31mA\033[0m"))
	assert.Equal("\033[31m\033[0m", formatter.Sprintf("\033[31m\033[0m"))
	assert.Equal("\033[0m", formatter.Sprintf("\033[0m"))

	formatter = NewFormatter("%-0[shortLevelName]s", nil)
	formatter.ForceColors = true
	assert.Equal("\033[31mAAA\033[0m", formatter.Sprintf("\033[31mAAA\033[0m"))
	assert.Equal("\033[31mAA\033[0m", formatter.Sprintf("\033[31mAA\033[0m"))
	assert.Equal("\033[31mA\033[0m", formatter.Sprintf("\033[31mA\033[0m"))
	assert.Equal("\033[31m\033[0m", formatter.Sprintf("\033[31m\033[0m"))
	assert.Equal("\033[0m", formatter.Sprintf("\033[0m"))

	formatter = NewFormatter("%-[shortLevelName]s", nil)
	formatter.ForceColors = true
	assert.Equal("\033[31mAAA\033[0m", formatter.Sprintf("\033[31mAAA\033[0m"))
	assert.Equal("\033[31mAA\033[0m", formatter.Sprintf("\033[31mAA\033[0m"))
	assert.Equal("\033[31mA\033[0m", formatter.Sprintf("\033[31mA\033[0m"))
	assert.Equal("\033[31m\033[0m", formatter.Sprintf("\033[31m\033[0m"))
	assert.Equal("\033[0m", formatter.Sprintf("\033[0m"))
}
