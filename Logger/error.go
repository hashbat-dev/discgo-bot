package logger

import (
	config "github.com/dabi-ngin/discgo-bot/Config"
)

// Error - logs include stackTrace, requires LoggingLevel 1 or lower in config
func Error(guildId string, err error, a ...any) {
	if config.LoggingLevel <= config.LoggingLevelDebug {
		infoLine, formattedLogText := ParseLoggingText(guildId, err.Error(), a...)
		SendLogs(infoLine, formattedLogText, config.LoggingLevelError, true)
	}
}

func Error_IgnoreDiscord(guildId string, err error, a ...any) {
	if config.LoggingLevel <= config.LoggingLevelDebug {
		infoLine, formattedLogText := ParseLoggingText(guildId, err.Error(), a...)
		SendLogs(infoLine, formattedLogText, config.LoggingLevelError, false)
	}
}
