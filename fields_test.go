package lcf

import (
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestBuiltInFieldsManuallyDefined(t *testing.T) {
	assert := require.New(t)
	log := logrus.New()
	log.Formatter = NewFormatter(Basic, nil)
	fields := BuiltInFields(log.Formatter, "testing")
	assert.Equal("testing", fields[fieldPrefix+"name"])
}

func TestBuiltInFieldsAutoValid(t *testing.T) {
	assert := require.New(t)
	defer ResetLogger() // Cleanup after test.
	logrus.SetFormatter(NewFormatter(Basic, nil))
	fields := BuiltInFields(nil, "testing2")
	assert.Equal("testing2", fields[fieldPrefix+"name"])
}

func TestBuiltInFieldsInvalid(t *testing.T) {
	assert := require.New(t)
	defer ResetLogger() // Cleanup after test.
	fields := BuiltInFields(nil, "testing3")
	assert.Nil(fields[fieldPrefix+"name"])
}
