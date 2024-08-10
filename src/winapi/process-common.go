package winapi

import (
	"encoding/binary"
	"fmt"
)

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
	ret := closeProcess(p.Handle)
	if ret > 0 {
		p.Handle = 0
	}
}
func (p *Process) Inject() error {
	if p.ProcessId == 0 {
		err := fmt.Errorf("PID is zero\n")
		return err
	}
	if p.Handle != 0 {
		err := fmt.Errorf("Process has opened\n")
		return err
	}

	// 注入某个进程
	hProcess := openProcess(p.ProcessId, PROCESS_QUERY_INFORMATION|PROCESS_VM_READ|PROCESS_VM_OPERATION|PROCESS_VM_WRITE)
	if hProcess == 0 {
		err := fmt.Errorf("OpenProcess %d error\n", p.ProcessId)
		return err
	}
	p.Handle = hProcess

	return nil
}

// 读取数据
func (p *Process) ReadMemory (address uintptr, size uint32) []byte {
  // https://learn.microsoft.com/en-us/windows/win32/api/memoryapi/nf-memoryapi-readprocessmemory
	if p.Handle == 0 {
		panic(fmt.Errorf("Process has not opened"))
	}

	data := readMemory(p.Handle, address, size)
	return data
}
func (p *Process) ReadByte (address uintptr) byte {
	return p.ReadMemory(address, 1)[0]
}
func (p *Process) ReadInt16 (address uintptr) int16 {
	return int16(binary.LittleEndian.Uint16(p.ReadMemory(address, 2)))
}
func (p *Process) ReadInt32 (address uintptr) int32 {
	return int32(binary.LittleEndian.Uint32(p.ReadMemory(address, 4)))
}
func (p *Process) ReadInt64 (address uintptr) int64 {
	return int64(binary.LittleEndian.Uint64(p.ReadMemory(address, 8)))
}
func (p *Process) ReadString (address uintptr, size uint32, coding string) string {
	buf := p.ReadMemory(address, size)
	data, err := ByteToString(buf, coding)
	if err != nil {
		panic(err)
	}

	return string(data)
}

func (p *Process) WriteMemory (address uintptr, data []byte) bool {
	if p.Handle == 0 {
		panic(fmt.Errorf("Process has not opened"))
	}
	return writeMemory(p.Handle, address, data)
}
func (p *Process) WriteByte (address uintptr, data byte) bool {
	buf := []byte{data}
	return p.WriteMemory(address, buf)
}
func (p *Process) WriteInt16 (address uintptr, data int16) bool {
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, uint16(data))
	return p.WriteMemory(address, buf)
}
func (p *Process) WriteInt32 (address uintptr, data int32) bool {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(data))
	return p.WriteMemory(address, buf)
}
func (p *Process) WriteInt64 (address uintptr, data int64) bool {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(data))
	return p.WriteMemory(address, buf)
}
func (p *Process) WriteString (address uintptr, data string) bool {
	buf := []byte(data)
	return p.WriteMemory(address, buf)
}
