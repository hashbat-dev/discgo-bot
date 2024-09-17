package tempfiles

import (
	"os"
	"path/filepath"
	"time"

	config "github.com/dabi-ngin/discgo-bot/Config"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
)

const (
	expiryPeriod = 5 * time.Minute
)

func CleanUpTempFiles() error {
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

	logger.Info("TEMPFILES", "CleanUpTimeFiles() completes, deleted %v files", deletedCount)
	return nil
}
