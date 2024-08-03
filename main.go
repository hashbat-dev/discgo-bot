package main

import (
	"fmt"
	"sync"

	bot "github.com/ZestHusky/femboy-control/Bot"
	"github.com/ZestHusky/femboy-control/Bot/config"
	ping "github.com/ZestHusky/femboy-control/Ping"
)

func main() {

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Setting up Bot...")
		config.BotToken = "MTIyMjEwMDUxOTg1MDQ3OTY0Nw.Gxij91.rED1SyMYxyoZyqER82p3r1p8Fmb7J2zwZWQd94"
		go bot.Run()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Setting up Ping...")
		ping.Run()
	}()

	wg.Wait()
}
