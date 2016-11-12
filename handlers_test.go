package lcf

import (
	"os"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func handlerOne(_ *logrus.Entry, _ *CustomFormatter) (interface{}, error) {
	return nil, nil
}

func TestParseTemplateCustom(t *testing.T) {
	assert := require.New(t)

	custom := map[string]Handler{"one": handlerOne}
	template, handlers, attributes := ParseTemplate("%[one]d %[two]s", custom)

	assert.Equal("%[1]d %[two]s", template)
	assert.Len(handlers, 1)
	assert.Len(attributes, 1)
	assert.True(attributes["one"])
}

func TestParseTemplateBuiltIn(t *testing.T) {
	assert := require.New(t)

	template, handlers, attributes := ParseTemplate(Basic, nil)

	assert.Equal("%[1]s:%[2]s:%[3]s%[4]s\n", template)
	assert.Len(handlers, 4)
	assert.Len(attributes, 4)
	assert.True(attributes["levelName"])
	assert.True(attributes["name"])
	assert.True(attributes["message"])
	assert.True(attributes["fields"])
}

func TestHandlerAscTime(t *testing.T) {
	assert := require.New(t)

	// Setup.
	formatter := NewFormatter("", nil)
	entry := logrus.NewEntry(logrus.New())
	entry.Level = logrus.ErrorLevel

	// Test.
	fields, err := HandlerAscTime(entry, formatter)
	assert.NoError(err)
	actual := fields.(string)
	assert.Regexp(`^\d{4}-\d\d-\d\d \d\d:\d\d:\d\d\.\d{3}$`, actual)
}

func TestHandlerFields(t *testing.T) {
	assert := require.New(t)

	// Setup.
	formatter := NewFormatter("", nil)
	entry := logrus.NewEntry(logrus.New())
	entry.Level = logrus.ErrorLevel

	// Test with no data.
	fields, err := HandlerFields(entry, formatter)
	assert.NoError(err)
	actual := fields.(string)
	assert.Equal("", actual)

	// Test with sorting and colors.
	formatter.DisableColors = false
	formatter.DisableSorting = false
	formatter.ForceColors = true
	entry.Data["one"] = 1
	entry.Data["two"] = "2"
	entry.Data["3"] = false
	fields, err = HandlerFields(entry, formatter)
	assert.NoError(err)
	actual = fields.(string)
	expected := " \033[31m3\033[0m=false \033[31mone\033[0m=1 \033[31mtwo\033[0m=2"
	assert.Equal(expected, actual)

	// Test with no sorting and no colors.
	formatter.DisableColors = true
	formatter.DisableSorting = true
	formatter.ForceColors = false
	fields, err = HandlerFields(entry, formatter)
	assert.NoError(err)
	actual = fields.(string)
	assert.Contains(actual, "3=false")
	assert.Contains(actual, "one=1")
	assert.Contains(actual, "two=2")
}

func TestHandlerRelativeCreated(t *testing.T) {
	assert := require.New(t)

	// Setup.
	formatter := NewFormatter("", nil)

	// Test.
	var values [2]int
	fields, err := HandlerRelativeCreated(nil, formatter)
	assert.NoError(err)
	values[0] = fields.(int)
	time.Sleep(time.Second * 2)
	fields, err = HandlerRelativeCreated(nil, formatter)
	assert.NoError(err)
	values[1] = fields.(int)
	assert.True(values[0] < values[1])
}

func ExampleCustomHandlers() {
	// Define your own handler for new or to override built-in attributes. Here we'll
	// define LoadAverage() to handle a new %[loadAvg]f attribute.
	LoadAverage := func(e *logrus.Entry, f *CustomFormatter) (interface{}, error) {
		someNumber := 0.3
		return someNumber, nil
	}

	// You can define additional formatting in the template string. Formatting is
	// handled by fmt.Sprintf() after lcf converts keyed indexes to integer indexes.
	template := "[%04[relativeCreated]d] %1.2[loadAvg]f %7[levelName]s %[message]s\n"
	formatter := NewFormatter(template, CustomHandlers{"loadAvg": LoadAverage})

	// Create a new logger or use the standard logger. Here we'll create a new one
	// and configure it.
	log := logrus.New()
	log.Formatter = formatter
	log.Level = logrus.DebugLevel
	log.Out = os.Stdout
	log.Debug("A group of walrus emerges from the ocean")
	log.Warn("The group's number increased tremendously!")
	log.Info("A giant walrus appears!")
	log.Error("Tremendously sized cow enters the ocean.")

	// Output:
	// [0000] 0.30   DEBUG A group of walrus emerges from the ocean
	// [0000] 0.30 WARNING The group's number increased tremendously!
	// [0000] 0.30    INFO A giant walrus appears!
	// [0000] 0.30   ERROR Tremendously sized cow enters the ocean.
}
