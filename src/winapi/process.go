//go:build windows

package winapi

import (
	"fmt"
	"path"
	"strings"
	"syscall"
	"unsafe"
)

// 进程权限相关：https://docs.microsoft.com/en-us/windows/desktop/ProcThread/process-security-and-access-rights
const (
	DWORD_SIZE = 4				// 双字 4 字节
	PROCESS_QUERY_INFORMATION = 0x0400
	PROCESS_VM_OPERATION = 0x0008
	PROCESS_VM_READ = 0x0010
	PROCESS_VM_WRITE = 0x0020
)

type Process struct {
	ProcessId uint32
	ExecName string
	Handle uintptr
	ModAddrs []uint32
	ModNames []string
}
func (p *Process) Close() {
	if p.Handle == 0 {
		return
	}
	CloseHandle, _ := syscall.GetProcAddress(*kernel32, "CloseHandle")
	ret, _, _ := syscall.SyscallN(
		uintptr(CloseHandle),
		p.Handle,
	)
	if ret > 0 {
		p.Handle = 0
	}
}
func (p *Process) Inject() {
	if p.ProcessId == 0 {
		fmt.Printf("PID is zero\n")
		return
	}
	if p.Handle != 0 {
		fmt.Printf("Process has opened\n")
		return
	}

	// 注入某个进程
	OpenProcess, _ := syscall.GetProcAddress(*kernel32, "OpenProcess")

	var bInheritHandle = false
	hProcess, _, callErr := syscall.SyscallN(
		uintptr(OpenProcess),
		uintptr(PROCESS_QUERY_INFORMATION|PROCESS_VM_READ|PROCESS_VM_OPERATION|PROCESS_VM_WRITE),
		uintptr(unsafe.Pointer(&bInheritHandle)),
		uintptr(p.ProcessId),
	)

	if hProcess == 0 {
		fmt.Printf("OpenProcess %d error %v\n", p.ProcessId, callErr)
		return
	}
	p.Handle = hProcess
}
func (p *Process) ReadMemory (address uintptr, size uint32) {
  // https://learn.microsoft.com/en-us/windows/win32/api/memoryapi/nf-memoryapi-readprocessmemory
	if p.Handle == 0 {
		panic("Process has not opened")
	}
}
func (p *Process) WriteMemory (address uintptr, data []byte) {
	if p.Handle == 0 {
		panic("Process has not opened")
	}
}

// ================================================================================
func GetProcessInfo (pid uint32) *Process {
	// https://docs.microsoft.com/en-us/windows/desktop/psapi/enumerating-all-modules-for-a-process
	OpenProcess, _ := syscall.GetProcAddress(*kernel32, "OpenProcess")
	
	var bInheritHandle = false
	// https://learn.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-openprocess
	hProcess, _, callErr := syscall.SyscallN(
		uintptr(OpenProcess),
		uintptr(PROCESS_QUERY_INFORMATION | PROCESS_VM_READ),
		uintptr(unsafe.Pointer(&bInheritHandle)),
		uintptr(pid),
	)

	if hProcess == 0 {
		fmt.Printf("OpenProcess %d error %v\n", pid, callErr)
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
		fmt.Printf("EnumProcessModules %d error %v\n", pid, callErr)
		return nil
	}

	GetModuleBaseNameA, _ := syscall.GetProcAddress(*psapi, "GetModuleBaseNameA")

	for i := 0; i < count; i++ {
		p.ModAddrs = append(p.ModAddrs, hMods[i])

		var lpBaseName [100]byte
		ret, _, _ = syscall.SyscallN(
			uintptr(GetModuleBaseNameA),
			hProcess,
			uintptr(hMods[i]),
			uintptr(unsafe.Pointer(&lpBaseName)),
			uintptr(1 * count),
		)

		if ret == 0 {
			fmt.Printf("GetModuleBaseNameA %d %d error %v\n", pid, i, callErr)
			p.ModNames = append(p.ModNames, "")
			continue
		}

		name := string(lpBaseName[:ret])
		fmt.Printf("pid %d mod %d %x name %s\n", pid, i, hMods[i], name)
		p.ModNames = append(p.ModNames, name)
	}

	return &p
}
func getProcessName (pid uint32) *Process {
	OpenProcess, _ := syscall.GetProcAddress(*kernel32, "OpenProcess")
	var bInheritHandle = false
	// https://learn.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-openprocess
	hProcess, _, _ := syscall.SyscallN(
		uintptr(OpenProcess),
		uintptr(PROCESS_QUERY_INFORMATION | PROCESS_VM_READ),
		uintptr(unsafe.Pointer(&bInheritHandle)),
		uintptr(pid),
	)

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
