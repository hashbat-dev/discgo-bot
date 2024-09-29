package logger

import (
	config "github.com/hashbat-dev/discgo-bot/Config"
)

// Logs info, Bozo... requires LoggingLevel 2 or lower in config
func Info(guildId string, logText string, a ...any) {
	if config.ServiceSettings.LOGGINGLEVEL <= config.LoggingLevelEvent {
		infoLine, formattedLogText := ParseLoggingText(guildId, logText, a...)
		SendLogs(infoLine, formattedLogText, config.LoggingLevelInfo, true)
	}
}

func Info_IgnoreDiscord(guildId string, logText string, a ...any) {
	if config.ServiceSettings.LOGGINGLEVEL <= config.LoggingLevelEvent {
		infoLine, formattedLogText := ParseLoggingText(guildId, logText, a...)
		SendLogs(infoLine, formattedLogText, config.LoggingLevelInfo, false)
	}
}
