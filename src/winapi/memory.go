//go:build windows

package winapi

import (
	"fmt"
	"syscall"
	"unsafe"
)

// https://learn.microsoft.com/en-us/windows/win32/api/memoryapi/nf-memoryapi-readprocessmemory
func readMemory(processHandle uintptr, address uintptr, size uint32) ([]byte) {
	ReadProcessMemory, _ := syscall.GetProcAddress(*kernel32, "ReadProcessMemory")

	buf := []byte{}
	for i := uint32(0); i < size; i++ {
		buf = append(buf, 0)
	}
	var nsize uint32

	ret, _, callErr := syscall.SyscallN(
		uintptr(ReadProcessMemory),
		processHandle,
		address,
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(size),
		uintptr(unsafe.Pointer(&nsize)),
	)

	if ret == 0 {
		panic(fmt.Errorf("ReadProcessMemory failed: %v", callErr))
	}

	return buf[:nsize]
}

// https://learn.microsoft.com/en-us/windows/win32/api/memoryapi/nf-memoryapi-writeprocessmemory
func writeMemory(processHandle uintptr, address uintptr, buf []byte) bool {
	WriteProcessMemory, _ := syscall.GetProcAddress(*kernel32, "WriteProcessMemory")

	var nsize uint32

	ret, _, callErr := syscall.SyscallN(
		uintptr(WriteProcessMemory),
		processHandle,
		address,
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(len(buf)),
		uintptr(unsafe.Pointer(&nsize)),
	)

	if ret == 0 {
		panic(fmt.Errorf("WriteProcessMemory failed: %v", callErr))
	}

	return true
}
