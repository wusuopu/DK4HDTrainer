//go:build !windows

package winapi

func Unload() {
	// 空操作
}

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
}
func (p *Process) Close() {
	// 空操作
}
func (p *Process) Inject() {
	// 空操作
}
func (p *Process) ReadMemory (address uintptr, size uint32) {
	// 空操作
}
func (p *Process) WriteMemory (address uintptr, data []byte) {
	// 空操作
}

func GetProcessInfo(pid uint32) *Process {
	// 空操作
	return nil
}
func ListProcess() []*Process {
	// 空操作
	var data []*Process
	return data
}
