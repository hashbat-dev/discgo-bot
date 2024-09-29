package main

import (
	"sync"

	bot "github.com/hashbat-dev/discgo-bot/Bot"
	dashboard "github.com/hashbat-dev/discgo-bot/Dashboard"
	scheduler "github.com/hashbat-dev/discgo-bot/Scheduler"
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
