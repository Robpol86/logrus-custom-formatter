// +build !windows

package lcf

// WindowsNativeANSI will always return false if the compiler builds this file instead of colors_windows.go.
func WindowsNativeANSI() bool {
	return false
}
