package logger

import config "github.com/dabi-ngin/discgo-bot/Config"

// Logs important Events, ie. The requesting/delivery of user input
func Admin(guildId string, logText string, a ...any) {
	if config.ServiceSettings.LOGGINGLEVEL <= config.LoggingLevelAdmin {
		infoLine, formattedLogText := ParseLoggingText(guildId, logText, a...)
		SendLogs(infoLine, formattedLogText, config.LoggingLevelAdmin, true)
	}
}

func Admin_IgnoreDiscord(guildId string, logText string, a ...any) {
	if config.ServiceSettings.LOGGINGLEVEL <= config.LoggingLevelAdmin {
		infoLine, formattedLogText := ParseLoggingText(guildId, logText, a...)
		SendLogs(infoLine, formattedLogText, config.LoggingLevelAdmin, false)
	}
}
