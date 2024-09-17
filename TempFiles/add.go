package tempfiles

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	config "github.com/dabi-ngin/discgo-bot/Config"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
	"github.com/google/uuid"
)

func init() {
	if _, err := os.Stat(config.TEMP_FOLDER); os.IsNotExist(err) {
		err := os.Mkdir(config.TEMP_FOLDER, 0755)
		if err != nil {
			logger.Error("TEMPFILES", err)
			return
		}
		logger.Debug("TEMPFILES", "Created [/temp] Directory")
	}
}

func AddFile(file io.Reader, fileExtension string) string {
	fileName := fmt.Sprintf("%s%s", uuid.New().String(), fileExtension)
	filePath := filepath.Join(config.TEMP_FOLDER, fileName)

	// Create the file
	outFile, err := os.Create(filePath)
	if err != nil {
		logger.Error("TEMPFILES", err)
		return ""
	}
	defer outFile.Close()

	// Write to the file
	_, err = io.Copy(outFile, file)
	if err != nil {
		logger.Error("TEMPFILES", err)
		return ""
	}

	// Return the full URL
	return fmt.Sprintf("%s/temp/%s", config.ServiceSettings.DASHBOARDURL, url.PathEscape(fileName))
}

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
