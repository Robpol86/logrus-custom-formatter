package lcf

import (
	"testing"

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

	assert.Equal(`%[1]s:%[2]s:%[3]s\n`, template)
	assert.Len(handlers, 3)
	assert.Len(attributes, 3)
	assert.True(attributes["levelName"])
	assert.True(attributes["name"])
	assert.True(attributes["message"])
}
