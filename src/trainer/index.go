package trainer

import (
	"dk4/winapi"
	"fmt"
	"strings"
)

type Trainer struct {
	process *winapi.Process
	version string
	baseAddr uint64		// 基址地址
}

func (t *Trainer) Init() {
	t.process = nil
	t.baseAddr = 0
	t.version = ""

  processes := winapi.ListProcess()

	for _, p := range processes {
		if strings.ToLower(p.ExecName) == strings.ToLower("DK4HD_sc.exe") {
			t.process = p
			t.version = "sc"		// 简体中文版本
			break
		}
	}
	if t.process == nil {
		return
	}

	p := winapi.GetProcessInfo(t.process.ProcessId)
	if len(p.ModAddrs) == 0 {
		return
	}

	t.baseAddr = uint64(p.ModAddrs[0]) + (uint64(p.ModAddrs[1]) << 32)
	fmt.Printf("进程 %d; 基址： %x\n", t.process.ProcessId, t.baseAddr)
}
