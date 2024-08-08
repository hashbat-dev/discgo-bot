package main

import (
	"sync"

	bot "github.com/dabi-ngin/discgo-bot/Bot"
	ping "github.com/dabi-ngin/discgo-bot/Ping"
	scheduler "github.com/dabi-ngin/discgo-bot/Scheduler"
)

func main() {

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		bot.Init()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		scheduler.Init()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ping.Run()
	}()

	wg.Wait()
}
