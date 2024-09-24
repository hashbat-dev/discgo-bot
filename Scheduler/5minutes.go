package scheduler

import (
	cache "github.com/dabi-ngin/discgo-bot/Cache"
)

func RunEvery5Minutes() {
	go cache.PruneCache()
}
