package scheduler

import (
	cache "github.com/hashbat-dev/discgo-bot/Cache"
)

func RunEvery5Minutes() {
	go cache.PruneCache()
}
