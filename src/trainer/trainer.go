package trainer

import (
	"dk4/winapi"
	"fmt"
	"strings"
)

type Trainer struct {
	Process *winapi.Process
	Version string
	baseAddr uint64		// 基址地址
}

func (t *Trainer) Init() {
	if t.Process != nil && t.Process.Handle != 0 {
		t.Process.Close()
	}

	t.Process = nil
	t.baseAddr = 0
	t.Version = ""

  processes := winapi.ListProcess()

	for _, p := range processes {
		if strings.ToLower(p.ExecName) == strings.ToLower("DK4HD_sc.exe") {
			t.Process = p
			t.Version = "sc"		// 简体中文版本
			break
		}
	}
	if t.Process == nil {
		fmt.Println("游戏还未启动")
		return
	}

	p := winapi.GetProcessInfo(t.Process.ProcessId)
	if len(p.ModAddrs) == 0 {
		return
	}
	t.Process.Inject()

	t.baseAddr = uint64(p.ModAddrs[0]) + (uint64(p.ModAddrs[1]) << 32)
	fmt.Printf("进程 %d; 基址： %x\n", t.Process.ProcessId, t.baseAddr)
}
