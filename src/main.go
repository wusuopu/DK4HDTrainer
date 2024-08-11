package main

import (
	"dk4/ui"
	"dk4/winapi"
	"embed"
)

//go:embed views/*
var embededFiles embed.FS

func main() {
  defer winapi.Unload()

  ui.Run(embededFiles)
}
