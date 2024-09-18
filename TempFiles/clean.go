package tempfiles

import (
	"os"
	"path/filepath"
	"time"

	config "github.com/dabi-ngin/discgo-bot/Config"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
)

func DeleteAllExpired() {
	expiryPeriod := time.Duration(config.ServiceSettings.TEMPFILEEXPIRYMINS) * time.Minute
	now := time.Now()
	deletedCount := 0
	err := filepath.Walk(config.TEMP_FOLDER, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Error("TEMPFILES", err)
		} else {
			if !info.IsDir() {
				modTime := info.ModTime()
				if now.Sub(modTime) > expiryPeriod {
					err := os.Remove(path)
					if err != nil {
						logger.Error("TEMPFILES", err)
					} else {
						logger.Info("TEMPFILES", "Deleted TempFile: %s", path)
						deletedCount++
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		logger.Error("TEMPFILES", err)
	}

	logger.Debug("TEMPFILES", "CleanUpTimeFiles() completes, deleted %v files", deletedCount)
}
