package trainer

import (
	"dk4/winapi"
	"fmt"
)


const (
	MAX_SEAMAN_COUNT = 202
	SEAMAN_SIZE = 0x30
	SEAMAN_OFFSET = 0x42b11f0
	CURRENT_LEAD_SEAMAN_NUM = 235			// 当前玩家海员编号
)

var ALL_SEAMAN_NAMES [202]string = [202]string{
	"拉斐尔•卡斯特路",									// 0
	"赫德拉姆•柏格斯统",
	"丽璐•阿歌特",
	"华梅•李",
	"杏太郎•佐伯",
	"亚伯拉罕•伍丁",
	"蒂雅•恰斯卡",
	"杰拿斯•帕沙",
	"库拉乌迪•马奈乌斯",
	"弗利奥•埃涅科",
	"克莉丝汀娜•埃涅科",								// 10
	"阿尔加迪斯•歐多西斯",
	"卡米尔•奥芬埃西",
	"希恩.杨",
	"詹姆•魯德韦",
	"行久•白木",
	"铁礼列•腾尼",
	"埃米利奥•菲陛",
	"安杰洛•普契尼",
	"格尔哈特 阿迪肯",
	"阿尔•西恩",										// 20
	"查理•洛雪弗",
	"科鲁罗•西奈特",
	"费南德•迪阿斯",
	"易安•杜可夫",
	"塞维•汉",
	"曼努埃尔•阿尔米达",
	"塞拉•夏尔巴拉茲",
	"乙凤•宋",
	"尤里安•罗佩斯",
	"阿芝莎•努连纳哈尔",						// 30
	"林•森",
	"阿米娜•安奈富",
	"德尼雅•伊蒂哈德",
	"法娣玛•哈涅",
	"哈希姆•阿尔奈迪尔",
	"柳科•西萨",
	"谢尔•内迪姆",
	"塞西莉雅•梅卡德",
	"詹姆茲•克利福德",
	"马丁•舒派亚",								// 40
	"迪斯尼•阿博科鲁克",
	"彼德罗•巴鲁迪斯",
	"帕欧罗•陈特利欧",
	"赫德姆•巴夏",
	"巴巴洛沙•海雷丁",
	"菲南•西魯韦拉",
	"詹洛尼默•埃斯皮诺沙",
	"格里高利•图鲁维",
	"柴门林•纳哥普尔",
	"道阿尔泰•普雷依拉",					// 50
	"安东尼奥•库恩",
	"索静•来岛",
	"巴士科斯•玛尔德纳德",
	"迪欧歌•埃斯康特",
	"瑜•文",
	"罗伯特•史特科",
	"查理•吉尔邦",
	"帕罗.莫拉依斯",
	"舒伯特•格拉斯",
	"西门•李纳列斯",						// 60
	"阿伦索•罗莎",
	"恩查•列吉欧尼",
	"萨利德•阿加儒",
	"莫沙里•希得拉",
	"阿纳朱德•贾乌米",
	"伊文•尼耶迪",
	"达维易•弗拉芬",
	"盖欧古•札图吉塔",
	"纳哥兴•卡莫路",
	"赫尔黑•阿本思",						// 70
	"加洛沙•奥毕厄徳",
	"芬•布兰科",
	"维路斯•格拉纳特",
	"杰克斯•勃鲁姆",
	"朱里安•菲米尔",
	"尤里斯•休根",
	"柏列特•佩罗",
	"强•拉姆吉欧",
	"伽布利厄•卡尔杜齐",
	"雅可布•勃特朗",					// 80
	"贝拉肯苏•阿基列",
	"芝加诺斯•贝",
	"維廉•克莱布",
	"鄂路南•贝利欧",
	"易斯卡•雅沙尔",
	"理查•回森",
	"汉斯•勒茨",
	"米哈易尔•勒茨",
	"阿科布•达迈",
	"谢乌德•埃米",					// 90
	"米瓦尔•根茨",
	"夏洛特•米列",
	"桑亨•韩",
	"克丽雅•波福特",
	"爱斯法妮雅•梅卡德",
	"特欧贝尔德•梅卡德",
	"宗九郎•泷田",					// 97
	"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",			// 101 个空占位符
	"安吉鲁",							// 199
	"斯耐克",
	"丹第",								// 201
}

type Seaman struct {
	Id uint16
	Name string
	OrgId uint8					`offset:"0x0C" size:"1"`			// 所属势力
	OrgName string
	Metrics [6]uint8		`offset:"0x0F" size:"1"`			// 6个属性
	Exps [2]uint32			`offset:"0x18" size:"4"`			// 2个经验值 最大值 30W(0x0493e0)
	LeadId uint8				`offset:"0x28" size:"1"`			// 当前提督海员编号
}

func ListSeaman (t *Trainer) []*Seaman {
	if t.Process == nil {
		panic("进程不存在")
	}

	buf := t.Process.ReadMemory(
		uintptr(t.baseAddr + SEAMAN_OFFSET),
		SEAMAN_SIZE * MAX_SEAMAN_COUNT,
	)

	var data []*Seaman

	for i := 0; i < MAX_SEAMAN_COUNT; i++ {
		o := &Seaman{Id: uint16(i), Name: getSeamanName(uint8(i))}
		o.Parse((buf[(i * SEAMAN_SIZE):(i * SEAMAN_SIZE + SEAMAN_SIZE)]))
		data = append(data, o)
	}

	o := GetCurrentLeadSeaman(t)
	data[o.LeadId] = o
	return data
}

func getSeamanName (id uint8) string {
	if id < MAX_SEAMAN_COUNT {
		return ALL_SEAMAN_NAMES[id]
	}
	return ""
}

func GetCurrentLeadSeaman (t *Trainer) *Seaman {
	s := &Seaman{Id: CURRENT_LEAD_SEAMAN_NUM}
	s.GetSeamanById(t, CURRENT_LEAD_SEAMAN_NUM)
	name1 := t.Process.ReadString(uintptr(t.baseAddr + SEAMAN_OFFSET + CURRENT_LEAD_SEAMAN_NUM * SEAMAN_SIZE + 0x30), 17, "gbk")
	name2 := t.Process.ReadString(uintptr(t.baseAddr + SEAMAN_OFFSET + CURRENT_LEAD_SEAMAN_NUM * SEAMAN_SIZE + 0x41), 17, "gbk")
	name3 := t.Process.ReadString(uintptr(t.baseAddr + SEAMAN_OFFSET + CURRENT_LEAD_SEAMAN_NUM * SEAMAN_SIZE + 0x52), 17, "gbk")
	name := fmt.Sprintf("%s•%s•%s", name1, name2, name3)
	s.Name = name
	return s
}

func (s *Seaman) Parse(buf []byte) {
	if len(buf) != SEAMAN_SIZE {
		return
	}
	s.OrgId = buf[0x0C]
	s.OrgName = getOrganizationName(s.OrgId)
	for i := 0; i < 6; i++ {
		s.Metrics[i] = buf[0x0F+i]
	}
	for i := 0; i < 2; i++ {
		s.Exps[i] = winapi.ByteToUInt32(buf[0x18+i*4:0x18+i*4+4])
	}

	s.LeadId = buf[0x28]
}

func (s *Seaman) GetSeamanById (t *Trainer, id uint64) *Seaman {
	if t.Process == nil {
		panic("进程不存在")
	}

	buf := t.Process.ReadMemory(
		uintptr(t.baseAddr + SEAMAN_OFFSET + id * SEAMAN_SIZE),
		SEAMAN_SIZE,
	)

	s.Id = uint16(id)
	s.Name = getSeamanName(uint8(id))
	s.Parse((buf))

	return s
}

func (s *Seaman) String () string{
	return fmt.Sprintf(
		"海员:%d %s; 势力:%s; 六维:%v; 经验：%v",
		s.Id, s.Name, s.OrgName, s.Metrics, s.Exps,
	)
}