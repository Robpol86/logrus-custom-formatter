package lcf

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
)

var _reBracketed = regexp.MustCompile(`%[\d-]*\[(\w+)]\w`)
var _startTime = time.Now()

// Handler is the function signature of formatting attributes such as "levelName" and "message".
type Handler func(*logrus.Entry, *CustomFormatter) (interface{}, error)

// CustomHandlers is a mapping of Handler values to attributes as key names (e.g. "levelName").
type CustomHandlers map[string]Handler

// Attributes is a map used like a "set" to keep track of which formatting attributes are used.
type Attributes map[string]bool

// Contains returns true if attr is present.
func (a Attributes) Contains(attr string) bool {
	_, ok := a[attr]
	return ok
}

// HandlerAscTime returns the formatted timestamp of the entry.
func HandlerAscTime(entry *logrus.Entry, formatter *CustomFormatter) (interface{}, error) {
	return entry.Time.Format(formatter.TimestampFormat), nil
}

// HandlerFields returns the entry's fields (excluding name field if %[name]s is used) colorized according to log level.
// Fields' formatting: key=value key2=value2
func HandlerFields(entry *logrus.Entry, formatter *CustomFormatter) (interface{}, error) {
	var fields string

	// Without sorting no need to get keys from map into a string array.
	if formatter.DisableSorting {
		for key, value := range entry.Data {
			if key == "name" && formatter.Attributes.Contains("name") {
				continue
			}
			fields = fmt.Sprintf("%s %s=%v", fields, Color(entry, formatter, key), value)
		}
		return fields, nil
	}

	// Put keys in a string array and sort it.
	keys := make([]string, len(entry.Data))
	i := 0
	for k := range entry.Data {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	// Do the rest.
	for _, key := range keys {
		if key == "name" && formatter.Attributes.Contains("name") {
			continue
		}
		fields = fmt.Sprintf("%s %s=%v", fields, Color(entry, formatter, key), entry.Data[key])
	}
	return fields, nil
}

// HandlerLevelName returns the entry's long level name (e.g. "WARNING").
func HandlerLevelName(entry *logrus.Entry, formatter *CustomFormatter) (interface{}, error) {
	return Color(entry, formatter, strings.ToUpper(entry.Level.String())), nil
}

// HandlerName returns the name field value set by the user in entry.Data.
func HandlerName(entry *logrus.Entry, _ *CustomFormatter) (interface{}, error) {
	if value, ok := entry.Data["name"]; ok {
		return value.(string), nil
	}
	return "", nil
}

// HandlerMessage returns the unformatted log message in the entry.
func HandlerMessage(entry *logrus.Entry, _ *CustomFormatter) (interface{}, error) {
	return entry.Message, nil
}

// HandlerProcess returns the current process' PID.
func HandlerProcess(_ *logrus.Entry, _ *CustomFormatter) (interface{}, error) {
	return os.Getpid(), nil
}

// HandlerRelativeCreated returns the number of seconds since program start time.
func HandlerRelativeCreated(_ *logrus.Entry, _ *CustomFormatter) (interface{}, error) {
	return int(time.Since(_startTime) / time.Second), nil
}

// HandlerShortLevelName returns the first 4 letters of the entry's level name (e.g. "WARN").
func HandlerShortLevelName(entry *logrus.Entry, formatter *CustomFormatter) (interface{}, error) {
	return Color(entry, formatter, strings.ToUpper(entry.Level.String()[:4])), nil
}

// ParseTemplate parses the template string and prepares it for fmt.Sprintf() and keeps track of which handlers to use.
//
// :param template: Pre-processed formatting template (e.g. `%[message]s\n`).
//
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
			case "ascTime":
				handlers = append(handlers, HandlerAscTime)
			case "fields":
				handlers = append(handlers, HandlerFields)
			case "levelName":
				handlers = append(handlers, HandlerLevelName)
			case "name":
				handlers = append(handlers, HandlerName)
			case "message":
				handlers = append(handlers, HandlerMessage)
			case "process":
				handlers = append(handlers, HandlerProcess)
			case "relativeCreated":
				handlers = append(handlers, HandlerRelativeCreated)
			case "shortLevelName":
				handlers = append(handlers, HandlerShortLevelName)
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
