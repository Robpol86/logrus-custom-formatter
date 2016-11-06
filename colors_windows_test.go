package lcf

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type MockSysCall struct {
	mode                uint32
	getStdHandleError   error
	getConsoleModeError error
	setConsoleModeR1    uintptr
	setConsoleModeError error
}

func (s *MockSysCall) getStdHandle() error {
	return s.getStdHandleError
}

func (s *MockSysCall) getConsoleMode(mode *uint32) error {
	*mode = s.mode
	return s.getConsoleModeError
}

func (s *MockSysCall) setConsoleMode(mode uintptr) (uintptr, error) {
	s.mode = uint32(mode)
	return s.setConsoleModeR1, s.setConsoleModeError
}

func Test_windowsNativeANSI_Good(t *testing.T) {
	assert := require.New(t)
	sc := &MockSysCall{3, nil, nil, 1, nil}

	// First read. 3 & 4 == 0 == false.
	enabled, err := windowsNativeANSI(false, false, sc)
	assert.NoError(err)
	assert.False(enabled)

	// Enable feature.
	enabled, err = windowsNativeANSI(false, true, sc)
	assert.NoError(err)
	assert.True(enabled)

	// Is now enabled.
	enabled, err = windowsNativeANSI(false, false, sc)
	assert.NoError(err)
	assert.True(enabled)
}

func Test_windowsNativeANSI_BadSet(t *testing.T) {
	assert := require.New(t)
	sc := &MockSysCall{3, nil, nil, 0, errors.New("The parameter is incorrect.")}

	enabled, err := windowsNativeANSI(false, true, sc)
	assert.EqualError(err, "The parameter is incorrect.")
	assert.False(enabled)
}

func Test_windowsNativeANSI_BadGetConsole(t *testing.T) {
	assert := require.New(t)
	sc := &MockSysCall{3, nil, errors.New("The handle is invalid."), 1, nil}

	enabled, err := windowsNativeANSI(false, false, sc)
	assert.EqualError(err, "The handle is invalid.")
	assert.False(enabled)
}

func Test_windowsNativeANSI_BadGetMode(t *testing.T) {
	assert := require.New(t)
	sc := &MockSysCall{3, errors.New("The handle is invalid."), nil, 1, nil}

	enabled, err := windowsNativeANSI(false, false, sc)
	assert.EqualError(err, "The handle is invalid.")
	assert.False(enabled)
}
