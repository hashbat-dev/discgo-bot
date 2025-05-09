package scheduler

import (
	wow "github.com/hashbat-dev/discgo-bot/Bot/Commands/Wow"
	cache "github.com/hashbat-dev/discgo-bot/Cache"
)

func RunEvery5Minutes() {
	go cache.PruneCache()
	go wow.CleanCache()
}
