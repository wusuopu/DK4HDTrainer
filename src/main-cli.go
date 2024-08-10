package main

import (
	"bufio"
	"dk4/cli"
	"dk4/winapi"
	"os"
	"strings"
	"time"
)

func main() {
  defer winapi.Unload()

	handle, _ := cli.Actions["r"]
	handle()

	handle, _ = cli.Actions["p"]
	handle()
	cli.PrintHelp()

	// 定时修改游戏数据
	ticker := time.NewTicker(2 * time.Second)
	go func () {
		for {
			select {
			case <-ticker.C:
				cli.LockValueTick()
			}
		}
	}()



  scanner := bufio.NewScanner(os.Stdin)
  // Iterate over each line in the file
  for scanner.Scan() {
    input := scanner.Text()
		input = strings.TrimSpace(strings.ToLower(input))

		if input == "q" {
			ticker.Stop()
			break
		}

		if handle, ok := cli.Actions[input]; ok {
			handle()
			continue
		}
  }
}