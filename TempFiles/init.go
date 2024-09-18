package tempfiles

import (
	"os"

	config "github.com/dabi-ngin/discgo-bot/Config"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
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
