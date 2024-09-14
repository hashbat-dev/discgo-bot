package logger

import (
	config "github.com/dabi-ngin/discgo-bot/Config"
)

// Logs important Events, ie. The requesting/delivery of user input
func Event(guildId string, logText string, a ...any) {
	if config.ServiceSettings.LOGGINGLEVEL <= config.LoggingLevelEvent {
		infoLine, formattedLogText := ParseLoggingText(guildId, logText, a...)
		SendLogs(infoLine, formattedLogText, config.LoggingLevelEvent, true)
	}
}

func Event_IgnoreDiscord(guildId string, logText string, a ...any) {
	if config.ServiceSettings.LOGGINGLEVEL <= config.LoggingLevelEvent {
		infoLine, formattedLogText := ParseLoggingText(guildId, logText, a...)
		SendLogs(infoLine, formattedLogText, config.LoggingLevelEvent, false)
	}
}
