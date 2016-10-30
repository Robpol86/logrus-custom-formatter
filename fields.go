package lcf

import (
	"github.com/Sirupsen/logrus"
)

const fieldPrefix = "lcfField_"

// BuiltInFields collects data from the calling function to expose it to the handlers.
// This function is usually called once per caller function's call regardless of how many log statements are emitted.
// :param formatter: lcf.TextFormatter instance from your logger if different from the one in the logrus StandardLogger.
// :param name: The "logger name". This is the value for the "%[name]s" attribute.
func BuiltInFields(formatter logrus.Formatter, name string) logrus.Fields {
	fields := logrus.Fields{}

	// Handle nil formatter.
	if formatter == nil {
		formatter = logrus.StandardLogger().Formatter
	}

	// Get attributes from formatter.
	var attributes Attributes
	switch formatter.(type) {
	case *TextFormatter:
		attributes = formatter.(*TextFormatter).Attributes
	default:
		return fields
	}

	// Only populate fields that need to be populated.
	if _, ok := attributes["name"]; ok {
		fields[fieldPrefix+"name"] = name
	}

	return fields
}
