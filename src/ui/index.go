package ui

import (
	"dk4/config"
	"fmt"

	"github.com/jchv/go-webview2"
)

func Run() {
  w := webview2.NewWithOptions(webview2.WebViewOptions{
		Debug:     config.DEBUG,
		AutoFocus: true,
		WindowOptions: webview2.WindowOptions{
			Title:  "大航海时代4HD修改器",
			Width:  800,
			Height: 600,
			IconId: 2, // icon resource id
			Center: true,
		},
	})
  defer w.Destroy()

  // w.SetTitle("大航海时代4HD修改器")
  // w.SetSize(800, 600, webview2.HintFixed)
  w.Bind("dk4", func(arg1 string) string {
    fmt.Println("camm dk4 func")
    return "debug"
  })

  w.Navigate("https://www.baidu.com")
  w.Run()
}