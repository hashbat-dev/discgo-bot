package tempfiles

import (
	"os"
	"path/filepath"
	"strings"

	config "github.com/hashbat-dev/discgo-bot/Config"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

func DeleteFile(guildId string, inFile string) {
	// Skip due to TempFileGrace being enabled?
	if config.ServiceSettings.TEMPFILEGRACE {
		return
	}

	var fileName string

	// Check if we have a full filepath or not
	if strings.Contains(inFile, "/") {
		fileParts := strings.Split(inFile, "/")
		fileName = fileParts[len(fileParts)-1]
	} else {
		fileName = inFile
	}

	filePath := filepath.Join(config.TEMP_FOLDER, fileName)

	// Delete the file
	err := os.Remove(filePath)
	if err != nil {
		logger.Error(guildId, err)
		return
	}

	logger.Debug(guildId, "TempFile deleted: %s", fileName)
}
