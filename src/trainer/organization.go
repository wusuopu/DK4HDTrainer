package trainer

import (
	"dk4/winapi"
	"fmt"
)

// 最大势力个数
const (
	MAX_ORGANIZATION_COUNT = 25
	ORGANIZATION_SIZE = 0x58
	ORGANIZATION_OFFSET = 0x42B6DB0
)

var ALL_ORGANIZATION_NAMES [25]string = [25]string{
	"卡斯特路商会",		// 0
	"阿歌特商会",
	"柏格斯统商会",
	"李家",
	"克利福德军",
	"舒派亚商会",
	"阿博科魯克军",
	"巴魯迪斯军",
	"陈特利欧商会",
	"巴夏军",
	"海雷丁家",			// 10
	"西魯韦拉商会",
	"埃斯皮诺沙商会",
	"伍丁商会",
	"纳哥普尔商会",
	"普雷依拉商会",
	"库恩商会",
	"来岛家",
	"玛尔德纳德军",
	"埃斯康特军",
	"佐伯家",			// 20
	"昆廷韦拉斯",
	"图鲁维队",
	"李朝水师",
	"私家舰队",
}

// 势力信息 Size: 0x58 Bytes，起始地址: DK4HD_sc.exe+42B6DB0
type Organization struct {
	Id uint16
	Name string
	MasterPortId uint8		`offset:"0x0A" size:"1"`		// 总部港口编号
	MasterPortName string															// 总部港口名称
	Money uint32		`offset:"0x0C" size:"4"`		// 势力资金
	AreaValues [7]uint16	`offset:"0x10" size:"2"`		// 各海域势力值 北海、地中海、非洲、印度洋、东南亚、东亚、新大陆
}

// 获取所有势力信息
func ListOrganization (t *Trainer) []*Organization {
	if t.Process == nil {
		panic("进程不存在")
	}

	buf := t.Process.ReadMemory(
		uintptr(t.baseAddr + ORGANIZATION_OFFSET),
		ORGANIZATION_SIZE * MAX_ORGANIZATION_COUNT,
	)

	var data []*Organization

	for i := 0; i < MAX_ORGANIZATION_COUNT; i++ {
		o := &Organization{Id: uint16(i), Name: ALL_ORGANIZATION_NAMES[i]}
		o.Parse((buf[(i * ORGANIZATION_SIZE):(i * ORGANIZATION_SIZE + ORGANIZATION_SIZE)]))
		data = append(data, o)
	}

	return data
}
func getOrganizationName (id uint8) string {
	if id < MAX_ORGANIZATION_COUNT {
		return ALL_ORGANIZATION_NAMES[id]
	}
	return ""
}

func (o *Organization) Parse(buf []byte) {
	if len(buf) != ORGANIZATION_SIZE {
		return
	}

	o.MasterPortId = buf[0x0A]
	if o.MasterPortId < MAX_ORGANIZATION_COUNT {
		o.MasterPortName = ALL_PORT_CITIES[o.MasterPortId]
	} else {
		o.MasterPortName = ""
	}
	o.Money = winapi.ByteToUInt32(buf[0x0C:0x0C+4])
	for i := 0; i < 7; i++ {
		o.AreaValues[i] = winapi.ByteToUInt16(buf[0x10+i*2:0x10+i*2+2])
	}
}

func (o *Organization) GetOrganizationById (t *Trainer, id uint64) *Organization {
	if t.Process == nil {
		panic("进程不存在")
	}

	buf := t.Process.ReadMemory(
		uintptr(t.baseAddr + id * ORGANIZATION_OFFSET),
		ORGANIZATION_SIZE,
	)

	o.Id = uint16(id)
	o.Name = getOrganizationName(uint8(id))
	o.Parse((buf))

	return o
}

func ( o *Organization) String () string {
  return fmt.Sprintf("势力:%d %s; 总部:%s; 金币:%d %v", o.Id, o.Name, o.MasterPortName, o.Money, o.AreaValues)
}

// 设置金钱
func (o *Organization) SetMoney(t *Trainer, value uint32) {
	o.Money = value

	addr := t.baseAddr + ORGANIZATION_OFFSET + uint64(o.Id) * ORGANIZATION_SIZE
	t.Process.WriteInt32(uintptr(addr + 0x0C), int32(o.Money))
}