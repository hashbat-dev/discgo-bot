package tempfiles

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"

	config "github.com/dabi-ngin/discgo-bot/Config"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
	"github.com/google/uuid"
)

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
	return fmt.Sprintf("%stemp/%s", config.ServiceSettings.DASHBOARDURL, url.PathEscape(fileName))
}
