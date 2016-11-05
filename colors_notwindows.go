// +build !windows

package lcf

func windowsNativeANSI(_ bool, _ bool) (bool, error) {
	return false, nil
}
