package main

import (
	"dk4/trainer"
	"dk4/winapi"
)

func main() {
  defer winapi.Unload()

  trainer := trainer.Trainer{}
  trainer.Init()

  // ui.Run()
}
