package lcf

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Sirupsen/logrus"
)

var _reBracketed = regexp.MustCompile(`%\[(\w+)][\w\d-]`)

// Handler is the function signature of formatting attributes such as "levelName" and "message".
type Handler func(*logrus.Entry, *CustomFormatter) (interface{}, error)

// CustomHandlers is a mapping of Handler values to attributes as key names (e.g. "levelName").
type CustomHandlers map[string]Handler

// Attributes is a map used like a "set" to keep track of which formatting attributes are used.
type Attributes map[string]bool

// HandlerLevelName returns the entry's long level name (e.g. "WARNING").
func HandlerLevelName(entry *logrus.Entry, formatter *CustomFormatter) (interface{}, error) {
	return Color(entry, formatter, strings.ToUpper(entry.Level.String())), nil
}

// HandlerName returns the "logger name" set by the user at the beginning of their function's call.
func HandlerName(entry *logrus.Entry, _ *CustomFormatter) (interface{}, error) {
	if value, ok := entry.Data[fieldPrefix+"name"]; ok {
		return value.(string), nil
	}
	return "", nil
}

// HandlerMessage returns the unformatted log message in the entry.
func HandlerMessage(entry *logrus.Entry, _ *CustomFormatter) (interface{}, error) {
	return entry.Message, nil
}

// ParseTemplate parses the template string and prepares it for fmt.Sprintf() and keeps track of which handlers to use.
// :param template: Pre-processed formatting template (e.g. `%[message]s\n`).
// :param custom: User-defined formatters evaluated before built-in formatters. Keys are attributes to look for in the
func ParseTemplate(template string, custom CustomHandlers) (string, []Handler, Attributes) {
	attributes := make(Attributes)
	var handlers []Handler
	var positions [][2]int

	// Find attribute names to replace and with what handler function to map them to.
	for _, idxs := range _reBracketed.FindAllStringSubmatchIndex(template, -1) {
		attribute := template[idxs[2]:idxs[3]]
		if f, ok := custom[attribute]; ok {
			handlers = append(handlers, f)
		} else {
			switch attribute {
			case "levelName":
				handlers = append(handlers, HandlerLevelName)
			case "name":
				handlers = append(handlers, HandlerName)
			case "message":
				handlers = append(handlers, HandlerMessage)
			default:
				continue
			}
		}
		attributes[attribute] = true
		positions = append(positions, [...]int{idxs[2], idxs[3]})
	}

	// Substitute attribute names with Handler indexes in reverse.
	for i := len(positions) - 1; i >= 0; i-- {
		pos := positions[i]
		template = fmt.Sprintf("%s%d%s", template[:pos[0]], i+1, template[pos[1]:])
	}

	return template, handlers, attributes
}
