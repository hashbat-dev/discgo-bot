package cache

import "time"

func PruneCache() {
	// 1. Delete any Cached Interactions started +6 hours ago
	for key, item := range ActiveInteractions {
		if time.Since(item.Started) > time.Duration(6*time.Hour) {
			delete(ActiveInteractions, key)
		}
	}
	// 2. Delete any Admin Channel created +3 hours ago
	for key, item := range ActiveAdminChannels {
		if time.Since(item) > time.Duration(3*time.Hour) {
			delete(ActiveAdminChannels, key)
		}
	}
}
