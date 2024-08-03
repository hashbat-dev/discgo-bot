package main

import (
	"sync"

	bot "github.com/dabi-ngin/discgo-bot/Bot"
	dashboard "github.com/dabi-ngin/discgo-bot/Dashboard"
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
		dashboard.Run()
	}()

	wg.Wait()
}
