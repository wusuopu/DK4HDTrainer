package trainer

import (
	"dk4/winapi"
	"fmt"
)


const (
	MAX_PORT_CITY_COUNT = 97
	PORT_CITY_SIZE = 0x78
	PORT_CITY_OFFSET = 0x42B3E88
)
var ALL_PORT_CITIES [97]string = [97]string{
	"伦敦",					// 0
	"布利斯特",
	"阿姆斯特丹",
	"布鲁日",
	"南特",
	"汉堡",
	"卢比克",
	"斯德哥尔摩",
	"奥斯陆",
	"哥本哈根",
	"里加",				// 10
	"里斯本",
	"休达",
	"塞尔维亚",
	"瓦伦西亚",
	"热那亚",
	"马赛",
	"西拉库萨",
	"威尼斯",
	"雅典",
	"克里特",		// 20
	"塞普勒斯",
	"伊斯坦堡",
	"拉古扎",
	"贝鲁特",
	"亚才山卓",
	"的黎波里",
	"阿尔及尔",
	"突尼斯",
	"圣乔治",
	"马德拉",		// 30
	"拉斯帕马斯",
	"绿角",
	"卢安达",
	"索法拉",
	"开普敦",
	"莫三比克",
	"摩加迪休",
	"巴斯拉",
	"亚丁",
	"马斯喀特",	// 40
	"荷姆兹",
	"卡利亥特",
	"哥亚",
	"锡兰",
	"加尔各达",
	"安曼",
	"麻六甲",
	"汶莱",
	"马尼拉",
	"雅加达",		// 50
	"巴邻旁",
	"德尔纳特",
	"安勃那",
	"杭州",
	"泉州",
	"澳门",
	"京城",
	"长崎",
	"大阪",
	"那霸",		// 60
	"哈瓦那",
	"圣多明尼加",
	"圣约翰",
	"牙买加",
	"委拉克路斯",
	"美利达",
	"波多韦罗",
	"马拉开波",
	"伯南布哥",
	"特卢希忧",	// 70
	"卡恩内",
	"希拉雷奥湟",
	"圣多美",
	"马达加斯加",
	"蒙巴萨",
	"索哥德拉",
	"喀拉蚩",
	"马德拉斯",
	"马斯利巴丹",
	"阿镇",		// 80
	"吉阿丁",
	"汉杰鲁马辰",
	"马加撒",
	"泗水",
	"沂州",
	"朋沙科拉",
	"阿布哈兹",
	"雷利史塔特",
	"淡水",
	"圣马罗",	// 90
	"拿坡里",
	"班加西",
	"马纳多",
	"釜山",
	"江户",
	"卡拉卡斯",
}


type PortCity struct {
	Id uint16
	Name string
	DevelopmentValue uint16		`offset:"0x1E" size:"2"`
	ArmamentValue uint16			`offset:"0x20" size:"2"`
	Org1 uint8						`offset:"0x18" size:"1"`
	OrgName1 string
	Possession1 uint8			`offset:"0x19" size:"1"`
	Org2 uint8						`offset:"0x1A" size:"1"`
	OrgName2 string
	Possession2 uint8			`offset:"0x1B" size:"1"`
	Org3 uint8						`offset:"0x1C" size:"1"`
	OrgName3 string
	Possession3 uint8			`offset:"0x1D" size:"1"`
}

func ListPortCity (t *Trainer) []*PortCity {
	if t.Process == nil {
		panic("进程不存在")
	}

	buf := t.Process.ReadMemory(
		uintptr(t.baseAddr + PORT_CITY_OFFSET),
		PORT_CITY_SIZE * MAX_PORT_CITY_COUNT,
	)

	var data []*PortCity

	for i := 0; i < MAX_PORT_CITY_COUNT; i++ {
		o := &PortCity{Id: uint16(i), Name: ALL_PORT_CITIES[i]}
		o.Parse((buf[(i * PORT_CITY_SIZE):(i * PORT_CITY_SIZE + PORT_CITY_SIZE)]))
		data = append(data, o)
	}

	return data
}
func getPortCityName (id uint8) string {
	if id < MAX_ORGANIZATION_COUNT {
		return ALL_PORT_CITIES[id]
	}
	return ""
}

func (p *PortCity) Parse (buf []byte) {
	if len(buf) != PORT_CITY_SIZE {
		return
	}
	p.DevelopmentValue = winapi.ByteToUInt16(buf[0x1E:0x1E+2])
	p.ArmamentValue = winapi.ByteToUInt16(buf[0x20:0x20+2])

	p.Org1 = buf[0x18]
	p.Possession1 = buf[0x19]
	p.OrgName1 = getOrganizationName(p.Org1)

	p.Org2 = buf[0x1A]
	p.Possession2 = buf[0x1B]
	p.OrgName2 = getOrganizationName(p.Org2)

	p.Org3 = buf[0x1C]
	p.Possession3 = buf[0x1D]
	p.OrgName3 = getOrganizationName(p.Org3)
}

func (p *PortCity) GetPortCityById (t *Trainer, id uint64) *PortCity {
	if t.Process == nil {
		panic("进程不存在")
	}

	buf := t.Process.ReadMemory(
		uintptr(t.baseAddr + id * PORT_CITY_OFFSET),
		PORT_CITY_SIZE,
	)

	p.Id = uint16(id)
	p.Name = getPortCityName(uint8(id))
	p.Parse((buf))

	return p
}

func (p *PortCity) String() string {
		return fmt.Sprintf(
			"港口：%d %s 发展：%d 武装：%d； 占有率：%s %d, %s %d, %s %d",
			p.Id, p.Name, p.DevelopmentValue, p.ArmamentValue,
			p.OrgName1, p.Possession1,
			p.OrgName2, p.Possession2,
			p.OrgName3, p.Possession3,
		)
}