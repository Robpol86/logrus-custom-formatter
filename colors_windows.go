package lcf

import (
	"syscall"
)

// EnableVirtualTerminalProcessing indicates the platform supports VT100 color control character sequences.
const EnableVirtualTerminalProcessing = 0x0004

// Does this console window have ENABLE_VIRTUAL_TERMINAL_PROCESSING enabled? Optionally try to enable if not.
func windowsNativeANSI(stderr bool, setMode bool) (enabled bool, err error) {
	// Get the standard device.
	var nStdHandle int
	if stderr {
		nStdHandle = syscall.STD_ERROR_HANDLE
	} else {
		nStdHandle = syscall.STD_OUTPUT_HANDLE
	}

	// Get win32 handle.
	var handle syscall.Handle
	if handle, err = syscall.GetStdHandle(nStdHandle); err != nil {
		return
	}

	// Get console mode.
	var dwMode uint32
	if err = syscall.GetConsoleMode(handle, &dwMode); err != nil {
		return
	}
	enabled = dwMode&EnableVirtualTerminalProcessing != 0
	if enabled || !setMode {
		return
	}

	// Try to enable the feature.
	dwMode |= EnableVirtualTerminalProcessing
	proc := syscall.MustLoadDLL("kernel32").MustFindProc("SetConsoleMode")
	if r1, _, err := proc.Call(uintptr(handle), uintptr(dwMode)); r1 == 0 {
		return false, err
	}
	return true, nil
}
