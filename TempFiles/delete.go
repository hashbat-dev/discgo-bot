package tempfiles

import (
	"os"
	"path/filepath"
	"strings"

	config "github.com/dabi-ngin/discgo-bot/Config"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
)

func DeleteFile(guildId string, inFile string) {
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
