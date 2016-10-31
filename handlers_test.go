package lcf

import (
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

	assert.Equal(`%[1]s:%[2]s:%[3]s%[4]s\n`, template)
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
	formatter := CustomFormatter{TimestampFormat: DefaultTimestampFormat}
	entry := logrus.NewEntry(logrus.New())
	entry.Level = logrus.ErrorLevel

	// Test long timestamp.
	fields, err := HandlerAscTime(entry, &formatter)
	assert.NoError(err)
	actual := fields.(string)
	assert.Regexp(`^\d{4}-\d\d-\d\d \d\d:\d\d:\d\d\.\d{3}$`, actual)

	// Test short timestamp.
	formatter.ShortTimestamp = true
	var values [2]int
	fields, err = HandlerAscTime(entry, &formatter)
	assert.NoError(err)
	values[0] = fields.(int)
	time.Sleep(time.Second * 2)
	fields, err = HandlerAscTime(entry, &formatter)
	assert.NoError(err)
	values[1] = fields.(int)
	assert.True(values[0] < values[1])
}

func TestHandlerFields(t *testing.T) {
	assert := require.New(t)

	// Setup.
	formatter := CustomFormatter{}
	entry := logrus.NewEntry(logrus.New())
	entry.Level = logrus.ErrorLevel

	// Test with no data.
	fields, err := HandlerFields(entry, &formatter)
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
	fields, err = HandlerFields(entry, &formatter)
	assert.NoError(err)
	actual = fields.(string)
	expected := " \033[31m3\033[0m=false \033[31mone\033[0m=1 \033[31mtwo\033[0m=2"
	assert.Equal(expected, actual)

	// Test with no sorting and no colors.
	formatter.DisableColors = true
	formatter.DisableSorting = true
	formatter.ForceColors = false
	fields, err = HandlerFields(entry, &formatter)
	assert.NoError(err)
	actual = fields.(string)
	assert.Contains(actual, "3=false")
	assert.Contains(actual, "one=1")
	assert.Contains(actual, "two=2")
}
