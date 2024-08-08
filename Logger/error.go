package logger

import (
	"fmt"
	"time"

	config "github.com/dabi-ngin/discgo-bot/Config"
)

// Error - logs include stackTrace, requires LoggingLevel 1 or lower in config
func Error(guildId string, err error, a ...any) {
	if config.LoggingLevel <= config.LoggingLevelError {
		infoLine := fmt.Sprintf("%v | %v", time.Now().Format("02/01/06 15:04:05.000"), GetStack())
		if guildId != "" {
			infoLine += " | " + guildId
		}
		formattedLogText := err.Error()
		if len(a) > 0 {
			formattedLogText = FormatInboundLogText(err.Error(), a...)
		}
		SendToConsole(infoLine, formattedLogText, config.LoggingLevelError)
		if config.IsDev {
			infoLine += " | " + config.HostName
		}
		SendLogToDiscord(infoLine, formattedLogText, config.LoggingLevelError)
	}

}
