package lcf

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWindowsNativeANSI(t *testing.T) {
	assert := require.New(t)
	assert.False(WindowsNativeANSI())
}

func TestWindowsEnableNativeANSI(t *testing.T) {
	assert := require.New(t)

	if runtime.GOOS == "windows" {
		assert.EqualError(WindowsEnableNativeANSI(true), "The handle is invalid.")
		assert.EqualError(WindowsEnableNativeANSI(false), "The handle is invalid.")
	} else {
		assert.EqualError(WindowsEnableNativeANSI(false), "Not available on this platform.")
	}
}
