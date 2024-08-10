package trainer

import "fmt"

const (
	MAX_ARMADA_COUNT = 69
	ARMADA_SIZE = 0xC8
	ARMADA_OFFSET = 0X42A7B68
)

// 0~62 是预计的，63~68 是自行创建的
var ALL_ARMADA_NAMES [63]string = [63]string{
	"拉斐尔舰队",				// 0
	"赫德拉姆胍队",
	"丽璐舰队",
	"华梅舰队",
	"杏太郎舰队",
	"亚伯拉罕舰队",
	"蒂雅舰队",
	"",							// 7
	"詹姆茲舰队",
	"马丁舰队",
	"迪斯尼舰队",			// 10
	"彼德罗舰队",
	"帕欧罗舰队",
	"赫德姆胍队",
	"巴巴洛沙舰队",
	"菲南舰队",
	"詹洛尼默舰队",
	"海盗格里高利",
	"柴门林舰队",
	"道阿尔泰舰队",
	"安东尼奥舰队",		// 20
	"索静舰队",
	"巴士科斯舰队",
	"迪欧歌舰队",
	"罗伯特舰队",
	"查理舰队",
	"帕罗胍队",
	"舒伯特舰队",
	"西门舰队",
	"阿伦索舰队",
	"恩查舰队",			// 30
	"萨利德舰队",
	"莫沙里舰队",
	"阿纳朱德舰队",
	"伊文舰队",
	"达维易舰队",
	"盖欧古舰队",
	"纳哥兴舰队",
	"赫尔黑舰队",
	"加洛沙舰队",
	"芬舰队",			// 40
	"維路斯舰队",
	"理查舰队",
	"海盗瑜",
	"海盗阿芝莎",
	"海盗格尔哈特",
	"杰克斯舰队",
	"朱里安舰队",
	"海盗尤里斯",
	"海盗柏列特",
	"海盗强",			// 50
	"海盗伽布利厄",
	"海盗雅可布",
	"海盗贝拉肯苏",
	"海盗芝加诺斯",
	"海盗维廉",
	"海盗鄂路南海",
	"盗易斯卡",
	"海盗私设舰队",
	"海盗安吉鲁",
	"海盗斯耐克",	// 60
	"海盗丹第",
	"妖怪",
}

type Armada struct {
  Id uint16
	Name string
	Orgid uint8						`offset:"0x10" size:"1"`			// 所属势力
	OrgName string
	PortId uint8		`offset:"0x0A" size:"1"`		// 根据地港口编号
	PortName string															// 根据地港口名称
	Fatigue uint8				`offset:"0x71" size:"1"`			// 疲劳
	Longitude float32		`offset:"0x68" size:"2"`			// 经度
	Latitude float32		`offset:"0x6C" size:"2"`			// 纬度
}

func ListArmada(t *Trainer) []*Armada {
	if t.Process == nil {
		panic("进程不存在")
	}

	buf := t.Process.ReadMemory(
		uintptr(t.baseAddr + ARMADA_OFFSET),
		ARMADA_SIZE * MAX_ARMADA_COUNT,
	)

	var data []*Armada

	for i := 0; i < MAX_ARMADA_COUNT; i++ {
		o := &Armada{Id: uint16(i), Name: getArmadaName(uint8(i))}
		o.Parse((buf[(i * ARMADA_SIZE):(i * ARMADA_SIZE + ARMADA_SIZE)]))
		data = append(data, o)
	}

	return data
}

func getArmadaName(id uint8) string {
	if int(id) < len(ALL_ARMADA_NAMES) {
		return ALL_ARMADA_NAMES[id]
	}
	return ""
}

func (a *Armada) Parse(buf []byte) {
	if len(buf) != ARMADA_SIZE {
		return
	}
	a.Orgid = buf[0x10]
	a.OrgName = getOrganizationName(a.Orgid)
	a.PortId = buf[0x98]
	a.PortName = getPortCityName(a.PortId)
	a.Fatigue = buf[0x71]
	a.Longitude = parseLongitude(buf[0x69], buf[0x68])
	a.Latitude = parseLatitude(buf[0x6D], buf[0x6C])
}

func (a *Armada) GetArmadaById (t *Trainer, id uint64) *Armada {
	if t.Process == nil {
		panic("进程不存在")
	}

	buf := t.Process.ReadMemory(
		uintptr(t.baseAddr + id * ARMADA_OFFSET),
		ARMADA_SIZE,
	)

	a.Id = uint16(id)
	a.Name = getArmadaName(uint8(id))
	a.Parse((buf))

	return a
}

func (a *Armada) String() string {
	return fmt.Sprintf(
		"舰队:%d %s; 势力:%s; 根据地港口:%s; 疲劳:%d; 坐标：%s %s",
		a.Id, a.Name, a.OrgName, a.PortName, a.Fatigue,
		formatLatitude(a.Latitude), formatLongitude(a.Longitude),
	)
}