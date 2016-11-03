package lcf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCallerName(t *testing.T) {
	assert := require.New(t)

	assert.Equal("TestCallerName", CallerName(1))
	assert.Equal("TestCallerName", func() string { return CallerName(2) }())
}
