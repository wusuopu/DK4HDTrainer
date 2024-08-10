//go:build windows

package winapi

import (
	"fmt"
	"runtime"
	"syscall"
	"unsafe"
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


func openProcess(pid uint32, flag uint32) uintptr {
	OpenProcess, _ := syscall.GetProcAddress(*kernel32, "OpenProcess")
	var bInheritHandle = false
	// https://learn.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-openprocess
	hProcess, _, _ := syscall.SyscallN(
		uintptr(OpenProcess),
		uintptr(flag),
		uintptr(unsafe.Pointer(&bInheritHandle)),
		uintptr(pid),
	)

	return hProcess
}
func closeProcess(handle uintptr) uintptr {
	CloseHandle, _ := syscall.GetProcAddress(*kernel32, "CloseHandle")
	ret, _, _ := syscall.SyscallN(
		uintptr(CloseHandle),
		handle,
	)
	return ret
}