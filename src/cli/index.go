package cli

import (
	"dk4/config"
	"dk4/trainer"
	"dk4/utils"
	"fmt"
)


type ActionFunc func()

var t = &trainer.Trainer{}

var leadSeaman *trainer.Seaman				// 当前玩家角色
var currentOrg *trainer.Organization	// 当前玩家势力
var lockFoodFlag bool = false
var lockFatigueFlag bool = false
var lockShipFlag bool = false


func PrintHelp() {
	fmt.Printf("大航海时代4HD威力加强版修改器(当前版本%s)\n下载最新版本： https://github.com/wusuopu/DK4HDTrainer/releases \n", config.VERSION)
	fmt.Println("输入相应的指令")
	fmt.Println("\t[r] 刷新游戏状态(每次重新开始游戏之后需要刷新一下)")
	fmt.Println("\t[q] 退出")
	fmt.Println("\t[h] 打印此帮助信息")
	fmt.Println("\t[p] 打印游戏信息")
	fmt.Println("\t[0] 金钱加10000")
	fmt.Println("\t[1] 海员增强(属性加满、等级加满)")
	fmt.Println("\t[2] 无限粮食")
	fmt.Println("\t[3] 不会疲劳")
	fmt.Println("\t[4] 锁定战斗船只耐久值和水手数量")
	fmt.Println("\t[5] 武装增强(船只大炮升级到最高)")
	fmt.Println("\t[6] 发现所有补给港")
}
func refresh() {
	t.Init()

	if t.Process != nil {
		printInfo()
	}
}

func printInfo() {
	if t.Process == nil {
		fmt.Println("游戏还未启动")
		return
	}

	fmt.Printf("进程: %d; 版本: %s\n", t.Process.ProcessId, t.Version)

	// 玩家角色
	leadSeaman = nil
	seamans := trainer.ListSeaman(t)
	for _, v := range seamans {
		if v.Id == trainer.CURRENT_LEAD_SEAMAN_NUM {
			leadSeaman = v
			break
		}
	}

	if leadSeaman == nil {
		fmt.Println("游戏还未开始")
		return
	}
	fmt.Printf("当前玩家名: %s, 势力: %s\n", leadSeaman.Name, leadSeaman.OrgName)

	// 玩家势力
	org := &trainer.Organization{}
	org.GetOrganizationById(t, uint64(leadSeaman.OrgId))
	currentOrg = org
	fmt.Println(currentOrg.String())
	fmt.Println("----------")

	// 船只信息
	ships := trainer.ListShip(t)
	for _, v := range ships {
		if v.Valid && v.Id < 100 {
			fmt.Println(v.String())
		}
	}
	fmt.Println("----------")

	// 海员信息
	for _, v := range seamans {
		if uint16(v.OrgId) == currentOrg.Id {
			fmt.Println(v.String())
		}
	}
	fmt.Println("----------")
}

func addMoney() {
	if currentOrg == nil {
		fmt.Println("游戏还未开始")
		return
	}
	currentOrg.GetOrganizationById(t, uint64(currentOrg.Id))
	currentOrg.SetMoney(t, currentOrg.Money + 10000)
	fmt.Println("当前势力资金: ", currentOrg.Money)
}
func enhanceSeaman() {
	if currentOrg == nil {
		fmt.Println("游戏还未开始")
		return
	}

	seamans := trainer.ListSeaman(t)
	for _, v := range seamans {
		if uint16(v.OrgId) == currentOrg.Id {
			v.UpToMaxLevel(t)
		}
	}
	fmt.Println("海员修改完成")
}
func lockFood() {
	lockFoodFlag = !lockFoodFlag
	if lockFoodFlag {
		fmt.Println("食物锁定已开启")
	} else {
		fmt.Println("食物锁定已关闭")
	}
}
func lockFatigue() {
	lockFatigueFlag = !lockFatigueFlag
	if lockFatigueFlag {
		fmt.Println("不会疲劳已开启")
	} else {
		fmt.Println("不会疲劳已关闭")
	}
}
func lockShip() {
	lockShipFlag = !lockShipFlag
	if lockShipFlag {
		fmt.Println("战斗船只锁定已开启")
	} else {
		fmt.Println("战斗船只锁定已关闭")
	}
}

func enhanceShip() {
	if currentOrg == nil {
		fmt.Println("游戏还未开始")
		return
	}
	// 船只信息
	ships := trainer.ListShip(t)
	for _, v := range ships {
		if v.Valid && v.Id < 100 {
			// 船只改为连射炮
			v.SetGun(t, trainer.SHIP_GUN_TYPE_05)
		}
	}
}
func turnOnAllPorts() {
	if currentOrg == nil {
		fmt.Println("游戏还未开始")
		return
	}
	trainer.ToggleOnAllFeedPort(t)
	fmt.Println("修改完成")
}

var Actions = map[string]ActionFunc{
	"h": PrintHelp,
	"r": refresh,
	"p": printInfo,
	"0": addMoney,
	"1": enhanceSeaman,
	"2": lockFood,
	"3": lockFatigue,
	"4": lockShip,
	"5": enhanceShip,
	"6": turnOnAllPorts,
}

// 定时修改游戏数据
func _lockValueTick() {
	if leadSeaman == nil {
		return
	}

	// 船只水粮
	if lockFoodFlag {
		ships := trainer.ListShip(t)
		for _, v := range ships {
			if v.Valid && v.Id < 100 {
				v.LockWaterAndFood(t)
			}
		}
	}
	// 舰队疲劳
	if lockFatigueFlag {
		armadas := trainer.ListArmada(t)
		for _, v := range armadas {
			if uint16(v.LeadSeamanId) == uint16(leadSeaman.LeadId) {
				v.ResetFatigue(t)
			}
		}
	}
	if lockShipFlag {
		trainer.LockAllFight(t)
	}
}
func LockValueTick() {
	err := utils.Try(_lockValueTick)
	if err != nil {
		fmt.Println("游戏已结束")
		if t.Process != nil {
			t.Process.Close()
			t.Process = nil
		}
		leadSeaman = nil
		currentOrg = nil
	}
}