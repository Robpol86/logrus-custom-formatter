package lcf

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/Sirupsen/logrus"
)

const (
	// Basic formatting just logs the level name, function name, and message.
	Basic = `{{.levelname}}:{{.name}}:{{.message}}\n`

	// Message formatting just logs the message.
	Message = `{{.message}}\n`
)

// TextFormatter is the main formatter for the library.
type TextFormatter struct {
	// Formatting template.
	Template string
}

func get(data logrus.Fields, key string) string {
	if value, ok := data[key]; ok {
		return value.(string)
	}
	return ""
}

// Format is called by logrus and returns the formatted string.
func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var buf bytes.Buffer
	var level string

	// Define level. WARNING needs special attention.
	if entry.Level == logrus.WarnLevel {
		level = "WARN"
	} else {
		level = strings.ToUpper(entry.Level.String())
	}

	// Define values for formatter variables.
	values := map[string]string{
		"levelname": level,
		"message":   entry.Message,
		"name":      get(entry.Data, "lpfFieldName"),
	}

	// Parse entry.
	if t, err := template.New("").Parse(f.Template); err != nil {
		return nil, err
	} else if err := t.Execute(&buf, values); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// NewFormatter creates a new TextFormatter, sets the Format string, and returns its pointer.
// :param format: Log messages will follow this format.
func NewFormatter(format string) *TextFormatter {
	return &TextFormatter{Template: format}
}
