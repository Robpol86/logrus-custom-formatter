package formatter

import (
	"github.com/Sirupsen/logrus"
)

const (
	// Basic formatter just logs the level name, function name, and message.
	Basic = "%(levelname)s:%(name)s:%(message)s"
)

// TextFormatter is the main formatter for the library.
type TextFormatter struct {
	Formatting string
}

// Format is called by logrus and returns the formatted string.
func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return nil, nil
}

// NewFormatter creates a new TextFormatter, sets the Format string, and returns its pointer.
// :param format: Log messages will follow this format.
func NewFormatter(format string) *TextFormatter {
	return &TextFormatter{Formatting: format}
}
