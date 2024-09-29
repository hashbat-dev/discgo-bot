package imgur

import (
	"time"

	database "github.com/hashbat-dev/discgo-bot/Database"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

func TidySubmissions() {
	allDatabaseEntries, err := database.GetAllImgurLogs("IMGUR")
	if err != nil {
		return
	}

	i := 0
	for _, entry := range allDatabaseEntries {
		if time.Since(entry.CreatedDateTime) >= time.Duration(12*time.Hour) {
			if DeleteImgurEntry("IMGUR", entry.ImgurDeleteHash) == nil {
				i++
			}
		}
	}

	logger.Info("IMGUR", "Tidy Imgur Submissions completed, deleted %v entries", i)
}
