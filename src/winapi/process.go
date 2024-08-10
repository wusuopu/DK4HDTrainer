//go:build windows

package winapi

import (
	"dk4/config"
	"fmt"
	"path"
	"strings"
	"syscall"
	"unsafe"
)

// 进程权限相关：https://docs.microsoft.com/en-us/windows/desktop/ProcThread/process-security-and-access-rights
func GetProcessInfo (pid uint32) *Process {
	// https://docs.microsoft.com/en-us/windows/desktop/psapi/enumerating-all-modules-for-a-process
	hProcess := openProcess(pid, PROCESS_QUERY_INFORMATION | PROCESS_VM_READ)
	if hProcess == 0 {
		fmt.Printf("OpenProcess %d error\n", pid)
		return nil
	}

	p := Process{
		ProcessId: pid,
	}
	p.Handle = hProcess

	defer p.Close()

	EnumProcessModules, _ := syscall.GetProcAddress(*psapi, "EnumProcessModules")
	count := 100
	var hMods [100]uint32
	var cbNeeded uint32
	ret, _, callErr := syscall.SyscallN(
		uintptr(EnumProcessModules),
		hProcess,
		uintptr(unsafe.Pointer(&hMods)),
		uintptr(DWORD_SIZE * count),
		uintptr(unsafe.Pointer(&cbNeeded)),
	)
	if ret == 0 {
		if config.DEBUG {
			fmt.Printf("EnumProcessModules %d error %v\n", pid, callErr)
		}
		return nil
	}

	GetModuleBaseNameA, _ := syscall.GetProcAddress(*psapi, "GetModuleBaseNameA")

	for i := 0; i < count; i++ {
		p.ModAddrs = append(p.ModAddrs, hMods[i])

		var lpBaseName [100]byte
		ret, _, callErr = syscall.SyscallN(
			uintptr(GetModuleBaseNameA),
			hProcess,
			uintptr(hMods[i]),
			uintptr(unsafe.Pointer(&lpBaseName)),
			uintptr(1 * count),
		)

		if ret == 0 {
			if config.DEBUG {
				fmt.Printf("GetModuleBaseNameA %d %d error %v\n", pid, i, callErr)
			}
			p.ModNames = append(p.ModNames, "")
			continue
		}

		name := string(lpBaseName[:ret])
		if config.DEBUG {
			fmt.Printf("pid %d mod %d %x name %s\n", pid, i, hMods[i], name)
		}
		p.ModNames = append(p.ModNames, name)
	}

	return &p
}
func getProcessName (pid uint32) *Process {
	hProcess := openProcess(pid, PROCESS_QUERY_INFORMATION | PROCESS_VM_READ)
	if hProcess == 0 {
		// fmt.Printf("OpenProcess %d error\n", pid)
		return nil
	}

	p := Process{
		ProcessId: pid,
	}
	p.Handle = hProcess

	defer p.Close()

	GetProcessImageFileNameA, _ := syscall.GetProcAddress(*psapi, "GetProcessImageFileNameA")
	var lpImageFileName [2024]byte
	ret, _, callErr := syscall.SyscallN(
		uintptr(GetProcessImageFileNameA),
		hProcess,
		uintptr(unsafe.Pointer(&lpImageFileName)),
		uintptr(2024),
	)

	if ret == 0 {
		fmt.Printf("GetProcessImageFileNameA %d error: %v\n", pid, callErr)
	} else {
		_, execName := path.Split(strings.TrimSpace(string(lpImageFileName[:ret])))
		names := strings.Split(execName, "\\")
		p.ExecName = names[len(names)-1]
	}

	return &p
}


func ListProcess() []*Process {
	EnumProcesses, _ := syscall.GetProcAddress(*psapi, "EnumProcesses")

	count := 1024
	// lpidProcess := make([]uint32, count)
	var lpidProcess [1024]uint32
	var lpcbNeeded uint32

	// https://docs.microsoft.com/en-us/windows/desktop/api/psapi/nf-psapi-enumprocesses
	ret, _, callErr := syscall.SyscallN(
		uintptr(EnumProcesses),
		uintptr(unsafe.Pointer(&lpidProcess)),
		uintptr(DWORD_SIZE * count),
		uintptr(unsafe.Pointer(&lpcbNeeded)),
	)

	if ret != 1 {
		panic(fmt.Sprintf("ListProcess error: %v", callErr))
	}

	num := int(lpcbNeeded / uint32(DWORD_SIZE))
	var data []*Process

	for i := 0; i < num; i++ {
		p := getProcessName(lpidProcess[i])
		// p := GetProcessInfo(lpidProcess[i])
		if p == nil {
			continue
		}
		data = append(data, p)
	}

	return data
}
