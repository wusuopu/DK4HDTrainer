//go:build windows

package winapi

import (
	"fmt"
	"syscall"
	"unsafe"
)

// https://learn.microsoft.com/en-us/windows/win32/api/memoryapi/nf-memoryapi-readprocessmemory
func readMemory(processHandle uintptr, address uint32, size uint32) ([]byte) {
	ReadProcessMemory, _ := syscall.GetProcAddress(*kernel32, "ReadProcessMemory")

	buf := make([]byte, size)
	var nsize uint32

	ret, _, callErr := syscall.SyscallN(
		uintptr(ReadProcessMemory),
		processHandle,
		uintptr(address),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(len(buf)),
		uintptr(unsafe.Pointer(&nsize)),
	)

	if ret == 0 {
		panic(fmt.Sprintf("ReadProcessMemory failed: %v", callErr))
	}

	return buf[:nsize]
}

// https://learn.microsoft.com/en-us/windows/win32/api/memoryapi/nf-memoryapi-writeprocessmemory
func writeMemory(processHandle uintptr, address uint32, buf []byte) bool {
	WriteProcessMemory, _ := syscall.GetProcAddress(*kernel32, "WriteProcessMemory")

	var nsize uint32

	ret, _, callErr := syscall.SyscallN(
		uintptr(WriteProcessMemory),
		processHandle,
		uintptr(address),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(len(buf)),
		uintptr(unsafe.Pointer(&nsize)),
	)

	if ret == 0 {
		panic(fmt.Sprintf("WriteProcessMemory failed: %v", callErr))
	}

	return true
}
