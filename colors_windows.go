package lcf

import (
	"syscall"
)

// EnableVirtualTerminalProcessing indicates the platform supports VT100 color control character sequences.
const EnableVirtualTerminalProcessing = 0x0004

// WindowsNativeANSI returns true if the current Windows console has native support for ANSI color codes.
// Windows 10 Insider since around February 2016 finally introduced support for ANSI colors. Prior versions of Windows
// required issuing win32 API calls to change the next character's foreground and background colors. Windows versions
// that support native ANSI codes have the ENABLE_VIRTUAL_TERMINAL_PROCESSING flag enabled.
func WindowsNativeANSI() bool {
	// Get win32 handle for stdout.
	handle, err := syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
	if err != nil {
		// Try stderr instead.
		if handle, err = syscall.GetStdHandle(syscall.STD_ERROR_HANDLE); err != nil {
			return false
		}
	}

	// Get console mode.
	var dwMode uint32
	if err := syscall.GetConsoleMode(handle, &dwMode); err != nil {
		return false
	}
	return dwMode&EnableVirtualTerminalProcessing != 0
}
