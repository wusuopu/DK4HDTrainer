//go:build windows

package winapi

import (
	"fmt"
	"runtime"
	"syscall"
)

var psapi *syscall.Handle
var kernel32 *syscall.Handle

func init() {
	fmt.Println("init win-api")
	if runtime.GOOS != "windows" {
		fmt.Println("current is not windows")
		return
	}
	_psapi, _ := syscall.LoadLibrary("psapi.dll")
	_kernel32, _ := syscall.LoadLibrary("kernel32.dll")

	psapi = &_psapi
	kernel32 = &_kernel32
}

func Unload() {
	if psapi != nil {
		syscall.FreeLibrary(*psapi)
	}
	if kernel32 != nil {
		syscall.FreeLibrary(*kernel32)
	}
}