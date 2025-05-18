package logger

import (
	config "github.com/hashbat-dev/discgo-bot/Config"
)

func Dev(guildId string, logText string, a ...any) {
	if !config.ServiceSettings.ISDEV {
		return
	}
	if config.ServiceSettings.LOGGINGLEVEL <= config.LoggingLevelDebug {
		infoLine, formattedLogText := ParseLoggingText(guildId, logText, a...)
		SendLogs(infoLine, formattedLogText, config.LoggingLevelDebug, true)
	}
}

func Dev_IgnoreDiscord(guildId string, logText string, a ...any) {
	if !config.ServiceSettings.ISDEV {
		return
	}
	if config.ServiceSettings.LOGGINGLEVEL <= config.LoggingLevelDebug {
		infoLine, formattedLogText := ParseLoggingText(guildId, logText, a...)
		SendLogs(infoLine, formattedLogText, config.LoggingLevelDebug, false)
	}
}
