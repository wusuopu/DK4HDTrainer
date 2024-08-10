package trainer

import (
	"dk4/winapi"
	"fmt"
)

const (
	MAX_FIGHT_COUNT = 5
  FIGHT_SIZE = 0x30
	FIGHT_OFFSET = 0x42AB330
)

// 战斗船只
type Fight struct {
  Id uint16
	HP uint16					`offset:"0x26" size:"2"`			// 耐久
	Mariner uint16		`offset:"0x28" size:"2"`			// 水手
	Valid bool
}


func ListFight(t *Trainer) []*Fight {
	if t.Process == nil {
		panic(fmt.Errorf("进程不存在"))
	}

	buf := t.Process.ReadMemory(
		uintptr(t.baseAddr + FIGHT_OFFSET),
		FIGHT_SIZE * MAX_FIGHT_COUNT,
	)

	var data []*Fight

	for i := 0; i < MAX_FIGHT_COUNT; i++ {
		o := &Fight{Id: uint16(i)}
		o.Parse((buf[(i * FIGHT_SIZE):(i * FIGHT_SIZE + FIGHT_SIZE)]))
		data = append(data, o)
	}

	return data
}

// 锁定所有战斗船只耐久值和水手数量
func LockAllFight(t *Trainer) {
  data := ListFight(t)
	for i := 0; i < len(data); i++ {
		f := data[i]
		if !f.Valid {
			continue
		}
		f.HP = 900
		f.Mariner = 200
		addr := t.baseAddr + FIGHT_OFFSET + uint64(i) * FIGHT_SIZE

		t.Process.WriteInt16(uintptr(addr + 0x26), int16(f.HP))
		t.Process.WriteInt16(uintptr(addr + 0x28), int16(f.Mariner))
	}
}

func (f *Fight) Parse (buf []byte) {
	f.HP = winapi.ByteToUInt16(buf[0x26:0x26+2])
	f.Mariner = winapi.ByteToUInt16(buf[0x28:0x28+2])
	f.Valid = buf[0x21] != 0			// 船只有效
}
func (f *Fight) String() string {
	return fmt.Sprintf("Id: %d, 耐久: %d, 水手: %d, 有效: %t\n", f.Id, f.HP, f.Mariner, f.Valid)
}
