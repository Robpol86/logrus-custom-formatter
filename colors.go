package lcf

import (
	"fmt"

	"github.com/Sirupsen/logrus"
)

// ANSI color codes.
const (
	AnsiReset     = 0
	AnsiRed       = 31
	AnsiHiRed     = 91
	AnsiGreen     = 32
	AnsiHiGreen   = 92
	AnsiYellow    = 33
	AnsiHiYellow  = 93
	AnsiBlue      = 34
	AnsiHiBlue    = 94
	AnsiMagenta   = 35
	AnsiHiMagenta = 95
	AnsiCyan      = 36
	AnsiHiCyan    = 96
	AnsiWhite     = 37
	AnsiHiWhite   = 97
)

// Color colorizes the input string and returns it with ANSI color codes.
func Color(entry *logrus.Entry, formatter *CustomFormatter, s string) string {
	if !formatter.ForceColors && formatter.DisableColors {
		return s
	}

	// Determine color. Default is info.
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel:
		levelColor = formatter.ColorDebug
	case logrus.WarnLevel:
		levelColor = formatter.ColorWarn
	case logrus.ErrorLevel:
		levelColor = formatter.ColorError
	case logrus.PanicLevel:
		levelColor = formatter.ColorPanic
	case logrus.FatalLevel:
		levelColor = formatter.ColorFatal
	default:
		levelColor = formatter.ColorInfo
	}
	if levelColor == AnsiReset {
		return s
	}

	// Colorize.
	return fmt.Sprintf("\033[%dm%s\033[0m", levelColor, s)
}

// WindowsNativeANSI returns true if either the stderr or stdout consoles natively support ANSI color codes. On
// non-Windows platforms this always returns false.
func WindowsNativeANSI() bool {
	enabled, _ := windowsNativeANSI(true, false, nil)
	if enabled {
		return enabled
	}
	enabled, _ = windowsNativeANSI(false, false, nil)
	return enabled
}

// WindowsEnableNativeANSI will attempt to set ENABLE_VIRTUAL_TERMINAL_PROCESSING on a console using SetConsoleMode.
//
// :param stderr: Issue SetConsoleMode win32 API call on stderr instead of stdout handle.
func WindowsEnableNativeANSI(stderr bool) error {
	_, err := windowsNativeANSI(stderr, true, nil)
	return err
}
