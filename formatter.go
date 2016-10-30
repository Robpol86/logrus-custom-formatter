package lcf

import (
	"bytes"
	"fmt"

	"github.com/Sirupsen/logrus"
)

const (
	// Basic formatting just logs the level name, function name, and message.
	Basic = `%[levelName]s:%[name]s:%[message]s\n`

	// Message formatting just logs the message.
	Message = `%[message]s\n`
)

// TextFormatter is the main formatter for the library.
type TextFormatter struct {
	// Post-processed formatting template (e.g. `%[1]s:%[2]s:%[3]s\n`).
	Template string

	// Handler functions whose indexes match up with Template Sprintf explicit argument indexes.
	Handlers []Handler

	// Attribute names (e.g. "levelName") used in pre-processed Template.
	Attributes Attributes

	// Set to true to bypass checking for a TTY before outputting colors.
	ForceColors bool

	// Force disabling colors.
	DisableColors bool
}

// Format is called by logrus and returns the formatted string.
func (f TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Call handlers.
	values := make([]interface{}, len(f.Handlers))
	for i, handler := range f.Handlers {
		value, err := handler(entry, &f)
		if err != nil {
			return nil, err
		}
		values[i] = value
	}

	// Parse template and return.
	parsed := fmt.Sprintf(f.Template, values...)
	return bytes.NewBufferString(parsed).Bytes(), nil
}

// NewFormatter creates a new TextFormatter, sets the Template string, and returns its pointer.
// This function is usually called just once during a running program's lifetime.
// :param template: Pre-processed formatting template (e.g. `%[message]s\n`).
// :param custom: User-defined formatters evaluated before built-in formatters. Keys are attributes to look for in the
// 	formatting string (e.g. `%[myFormatter]s`) and values are formatting functions.
func NewFormatter(template string, custom CustomHandlers) *TextFormatter {
	formatter := TextFormatter{}
	formatter.Template, formatter.Handlers, formatter.Attributes = ParseTemplate(template, custom)
	return &formatter
}
