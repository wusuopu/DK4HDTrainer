package trainer

import (
	"dk4/winapi"
	"fmt"
)

const (
	MAX_SHIP_COUNT = 207
  SHIP_SIZE = 0x80
	SHIP_OFFSET = 0x42B8348
)

type Ship struct {
	Id uint16
	Name string
	Water uint16				`offset:"0x46" size:"2"`
	Food uint16					`offset:"0x48" size:"2"`
	HP uint16						`offset:"0x4A" size:"2"`
	Mariner uint16			`offset:"0x4C" size:"2"`
	Category uint8			`offset:"0x51" size:"1"`			// 船类型, 单独修改无效，值为 0x28 的表示该船只还不存在
	Gun uint8						`offset:"0x52" size:"1"`			// 炮类型, 0~5: 散弹、曲射、加农曲射、加农、重加农、连射
	Valid bool
}

func ListShip (t *Trainer) []*Ship {
	if t.Process == nil {
		panic("进程不存在")
	}

	buf := t.Process.ReadMemory(
		uintptr(t.baseAddr + SHIP_OFFSET),
		SHIP_SIZE * MAX_SHIP_COUNT,
	)

	var data []*Ship

	for i := 0; i < MAX_SHIP_COUNT; i++ {
		o := &Ship{Id: uint16(i)}
		o.Parse((buf[(i * SHIP_SIZE):(i * SHIP_SIZE + SHIP_SIZE)]))
		data = append(data, o)
	}

	return data
}

func (s *Ship) Parse(buf []byte) {
	s.Name, _ = winapi.ByteToString(buf[0x10:0x10+18], "gbk")

	s.Water = winapi.ByteToUInt16(buf[0x46:0x46+2])
	s.Food = winapi.ByteToUInt16(buf[0x48:0x48+2])
	s.HP = winapi.ByteToUInt16(buf[0x4A:0x4A+2])
	s.Mariner = winapi.ByteToUInt16(buf[0x4C:0x4C+2])
	s.Category = buf[0x51]
	s.Gun = buf[0x52]

	s.Valid = s.Category != 0x28
}

func (s *Ship) GetShipById (t *Trainer, id uint64) *Ship {
	if t.Process == nil {
		panic("进程不存在")
	}

	buf := t.Process.ReadMemory(
		uintptr(t.baseAddr + id * SHIP_OFFSET),
		SHIP_SIZE,
	)

	s.Id = uint16(id)
	s.Parse((buf))

	return s
}

func (s *Ship) String() string {
	return fmt.Sprintf(
		"船只:%d %s; 水分:%d 食物:%d 耐久:%d 水手:%d 类别:%d 大炮:%d",
		s.Id, s.Name, s.Water, s.Food, s.HP, s.Mariner, s.Category, s.Gun,
	)
}

// 锁定船只的水和食物
func (s *Ship) LockWaterAndFood(t *Trainer) {
	if !s.Valid{
		return
	}

	s.Water = 200
	s.Food = 200

	addr := t.baseAddr + SHIP_OFFSET + uint64(s.Id) * SHIP_SIZE
	t.Process.WriteByte(uintptr(addr + 0x46), byte(s.Water))
	t.Process.WriteByte(uintptr(addr + 0x48), byte(s.Food))
}