package ui

import (
	"dk4/config"
	"dk4/trainer"
	"dk4/utils"
	"fmt"
	"os"
	"path"

	"github.com/valyala/fastjson"
)

type jsFunc func(string) string

var jsFuncMaps = map[string]jsFunc{
	"readTemplateFile": readTemplateFile,
	"refreshStatus":    refreshStatus,
	"getStatus":    		getStatus,
	"getPlayerInfo":    getPlayerInfo,
	"listOrganization": listOrganization,
	"listArmada": 			listArmada,
	"listSeaman": 			listSeaman,
	"listShip": 				listShip,
	"listPort": 				listPort,
	"addMoney":					addMoney,
	"minusOrgMoney":		minusOrgMoney,
	"enhanceSeaman":		enhanceSeaman,
	"enhanceShip":			enhanceShip,
	"turnOnAllPorts":		turnOnAllPorts,
	"toggleLockFlag":		toggleLockFlag,
}

type GameInfo struct {
	ProcessId int			`json:"processId"`
	Lang string				`json:"lang"`
	Version string		`json:"version"`
}

// =============================================
// 读取静态的模板文件
func readTemplateFile(payload string) string {
	if viewDir == "" {
		return makeErrorResponse("viewDir is not set", 500)
	}

	data, _ := fastjson.Parse(payload)
	file, _ := data.StringBytes()

	name := path.Join(viewDir, string(file))
	fmt.Println("read template file:", name)
	content, err := os.ReadFile(name)
	if err != nil {
		return makeErrorResponse("cannot open file", 400)
	}

	resp := makeResponse(200)
	resp.Data = string(content)
	return formatResponse(&resp)
}


func _getPlayerInfo() bool {
	if t.Process == nil {
		fmt.Println("游戏还未启动")
		return false
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
		return false
	}
	fmt.Printf("当前玩家名: %s, 势力: %s\n", leadSeaman.Name, leadSeaman.OrgName)

	// 玩家势力
	org := &trainer.Organization{}
	org.GetOrganizationById(t, uint64(leadSeaman.OrgId))
	currentOrg = org
	fmt.Println(currentOrg.String())
	fmt.Println("----------")
	return true
}
func getPlayerInfo(payload string) string {
	_getPlayerInfo()

	resp := makeResponse(200)
	resp.Data = map[string]interface{}{
		"leadSeaman": leadSeaman,
		"currentOrg": currentOrg,
	}

	return formatResponse(&resp)
}

func refreshStatus(payload string) string {
	t.Init()

	if t.Process != nil {
		_getPlayerInfo()
	}

	resp := makeResponse(200)
	game := GameInfo{ Lang: t.Version, Version: config.VERSION, }

	if t.Process != nil {
		game.ProcessId = int(t.Process.ProcessId)
	}
	resp.Data = game
	return formatResponse(&resp)
}
func getStatus(payload string) string {
	resp := makeResponse(200)
	game := GameInfo{ Lang: t.Version, Version: config.VERSION, }

	if t.Process != nil {
		game.ProcessId = int(t.Process.ProcessId)
	}
	resp.Data = game
	return formatResponse(&resp)
}

func listOrganization(payload string) string {
	resp := makeResponse(200)
	data := trainer.ListOrganization(t)

	resp.Data = data
	return formatResponse(&resp)
}

func listArmada(payload string) string {
	resp := makeResponse(200)
	data := trainer.ListArmada(t)

	resp.Data = data
	return formatResponse(&resp)
}

func listSeaman(payload string) string {
	resp := makeResponse(200)
	data := trainer.ListSeaman(t)

	resp.Data = data
	return formatResponse(&resp)
}

func listShip(payload string) string {
	resp := makeResponse(200)
	data := trainer.ListShip(t)

	resp.Data = data
	return formatResponse(&resp)
}

func listPort(payload string) string {
	resp := makeResponse(200)
	data := trainer.ListPortCity(t)

	resp.Data = data
	return formatResponse(&resp)
}


func addMoney(payload string) string {
	if currentOrg == nil {
		return makeErrorResponse("", 400)
	}

	currentOrg.GetOrganizationById(t, uint64(currentOrg.Id))
	currentOrg.SetMoney(t, currentOrg.Money + 10000)

	resp := makeResponse(200)
	return formatResponse(&resp)
}
func minusOrgMoney(payload string) string {
	if currentOrg == nil {
		return makeErrorResponse("", 400)
	}

	data, _ := fastjson.Parse(payload)
	id, err := utils.GetJSONInt64(data, "id")
	if err != nil {
		return makeErrorResponse("", 400)
	}

	org := &trainer.Organization{}
	org.GetOrganizationById(t, uint64(id))
	org.SetMoney(t, 5000)

	resp := makeResponse(200)
	return formatResponse(&resp)
}

func enhanceSeaman(payload string) string {
	if currentOrg == nil {
		return makeErrorResponse("", 400)
	}

	seamans := trainer.ListSeaman(t)
	for _, v := range seamans {
		if uint16(v.OrgId) == currentOrg.Id {
			v.UpToMaxLevel(t)
		}
	}

	resp := makeResponse(200)
	return formatResponse(&resp)
}

func enhanceShip(payload string) string {
	if currentOrg == nil {
		return makeErrorResponse("", 400)
	}

	ships := trainer.ListShip(t)
	for _, v := range ships {
		if v.Valid && v.Id < 100 {
			// 船只改为连射炮
			v.SetGun(t, trainer.SHIP_GUN_TYPE_05)
		}
	}

	resp := makeResponse(200)
	return formatResponse(&resp)
}

func turnOnAllPorts(payload string) string {
	if currentOrg == nil {
		return makeErrorResponse("", 400)
	}

	trainer.ToggleOnAllFeedPort(t)

	resp := makeResponse(200)
	return formatResponse(&resp)
}

func toggleLockFlag(payload string) string {
	data, _ := fastjson.Parse(payload)
	key := utils.GetJSONString(data, "key")
	value := data.GetBool("value")

	switch key {
	case "food":
		lockFoodFlag = value
	case "ship":
		lockShipFlag = value
	case "fatigue":
		lockFatigueFlag = value
	}

	fmt.Printf("%s: %v; lockFoodFlag: %v; lockShipFlag: %v; lockFatigueFlag: %v\n", key, value, lockFoodFlag, lockShipFlag, lockFatigueFlag)

	resp := makeResponse(200)
	return formatResponse(&resp)
}