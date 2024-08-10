//go:build !windows

package winapi

func Unload() {
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

// ==================================== nop process ====================================
func openProcess(pid uint32, flag uint32) uintptr {
	return 0
}
func closeProcess(handle uintptr) uintptr {
	return 0
}

// ==================================== nop memory ====================================
func readMemory(processHandle uintptr, address uintptr, size uint32) ([]byte) {
}
func writeMemory(processHandle uintptr, address uintptr, buf []byte) bool {
}