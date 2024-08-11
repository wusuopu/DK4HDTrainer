package ui

import (
	"dk4/config"
	"dk4/trainer"
	"dk4/utils"
	"embed"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/jchv/go-webview2"
)

var t = &trainer.Trainer{}
var leadSeaman *trainer.Seaman				// 当前玩家角色
var currentOrg *trainer.Organization	// 当前玩家势力
var lockFoodFlag bool = false
var lockFatigueFlag bool = false
var lockShipFlag bool = false

var viewDir string

func Run(embededFiles embed.FS) {
  w := webview2.NewWithOptions(webview2.WebViewOptions{
		Debug:     config.DEBUG,
		AutoFocus: true,
		WindowOptions: webview2.WindowOptions{
			Title:  "大航海时代4HD威力加强版修改器",
			Width:  1300,
			Height: 800,
			IconId: 2, // icon resource id
			Center: true,
		},
	})
  defer w.Destroy()

	// 定时修改游戏数据
	ticker := time.NewTicker(2 * time.Second)
	go func () {
		for {
			select {
			case <-ticker.C:
				LockValueTick()
			}
		}
	}()
	defer ticker.Stop()

  w.SetTitle("大航海时代4HD威力加强版修改器")
  w.SetSize(1300, 800, webview2.HintFixed)
	injectJSSDK(w)

	filename := "index.html"
	baseDir := ""
	if config.DEBUG {
		baseDir, _ = os.Getwd()
		filename = "index.debug.html"
	} else {
		tmpDir, err := expandEmbed(embededFiles)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer os.RemoveAll(tmpDir)
		baseDir = tmpDir
	}
	viewDir = path.Join(baseDir, "views")

	url := "file:///" + path.Join(viewDir, filename)
	fmt.Println("url:", url)
  w.Navigate(url)
  w.Run()
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
			if uint16(v.OrgId) == currentOrg.Id {
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