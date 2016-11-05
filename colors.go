package lcf

import (
	"fmt"
	"runtime"

	"github.com/Sirupsen/logrus"
)

// ANSI color codes.
const (
	Red    = 31
	Yellow = 33
	Blue   = 34
	Gray   = 37
)

// Color colorizes the input string and returns it with ANSI color codes.
func Color(entry *logrus.Entry, formatter *CustomFormatter, s string) string {
	// Determine if colors should be shown. Large if statement block for easier human readability.
	isColored := false
	if formatter.ForceColors {
		isColored = true // ForceColors takes precedent.
	} else if formatter.DisableColors {
		// false.
	} else if formatter.isTerminal && (runtime.GOOS != "windows" || formatter.isWindowsNativeAnsi) {
		isColored = true
	}

	// Bail if colors are disabled.
	if !isColored {
		return s
	}

	// Determine color.
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel:
		levelColor = Gray
	case logrus.WarnLevel:
		levelColor = Yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = Red
	default:
		levelColor = Blue
	}

	// Colorize.
	return fmt.Sprintf("\033[%dm%s\033[0m", levelColor, s)
}

// WindowsNativeANSI returns true if either the stderr or stdout consoles natively support ANSI color codes. On
// non-Windows platforms this always returns false.
func WindowsNativeANSI() bool {
	enabled, _ := windowsNativeANSI(true, false)
	if enabled {
		return enabled
	}
	enabled, _ = windowsNativeANSI(false, false)
	return enabled
}

// WindowsEnableNativeANSI will attempt to set ENABLE_VIRTUAL_TERMINAL_PROCESSING on a console using SetConsoleMode.
//
// :param stderr: Issue SetConsoleMode win32 API call on stderr instead of stdout handle.
func WindowsEnableNativeANSI(stderr bool) error {
	_, err := windowsNativeANSI(stderr, true)
	return err
}
