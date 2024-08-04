package main

import (
	"fmt"
	"sync"

	bot "github.com/dabi-ngin/discgo-bot/Bot"
	"github.com/dabi-ngin/discgo-bot/Bot/config"
	ping "github.com/dabi-ngin/discgo-bot/Ping"
)

func main() {

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Setting up Bot...")
		config.BotToken = "MTI2OTYyOTY4NjIzMzgyNTI5MA.GmMv1T.vDmZJy5inKcVCxnx4ipe429o97Vg9PfaenSWDc"
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
