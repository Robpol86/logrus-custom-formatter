package lcf

import (
	"github.com/Sirupsen/logrus"
)

// FieldPrefix is the string prefix in field keys in entry.Data used internally by this library. It allows us to
// differentiate between user-defined fields and lcf-defined fields since all fields are in the same namespace.
const FieldPrefix = "lcfField_"

// BuiltInFields collects data from the calling function to expose it to the handlers.
// This function is usually called once per caller function's call regardless of how many log statements are emitted.
// :param formatter: lcf.CustomFormatter instance from your logger if different from the one in logrus StandardLogger.
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
	case *CustomFormatter:
		attributes = formatter.(*CustomFormatter).Attributes
	default:
		return fields
	}

	// Only populate fields that need to be populated.
	if attributes.Contains("name") {
		fields[FieldPrefix+"name"] = name
	}

	return fields
}
