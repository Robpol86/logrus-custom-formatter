package lcf

import (
	"fmt"

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
	} else if logrus.IsTerminal() {
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
