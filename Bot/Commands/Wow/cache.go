package wow

import (
	"time"

	config "github.com/hashbat-dev/discgo-bot/Config"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

var cacheBank map[string]CacheItem = make(map[string]CacheItem)

type CacheItem struct {
	Generation Generation
	Added      time.Time
}

func addToCache(wow *Generation) {
	cacheBank[wow.WowMessageID] = CacheItem{
		Generation: *wow,
		Added:      time.Now(),
	}
}

func CleanCache() {
	cutoff := time.Now().Add(-time.Duration(config.ServiceSettings.WOWRETENTIONMINS) * time.Minute)
	before := len(cacheBank)
	for key, item := range cacheBank {
		if item.Added.Before(cutoff) {
			delete(cacheBank, key)
		}
	}
	logger.Info("WOW", "CleanCache completed, Cache count went from %d to %d", before, len(cacheBank))
}
