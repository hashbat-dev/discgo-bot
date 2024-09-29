package logger

import (
	"time"

	config "github.com/hashbat-dev/discgo-bot/Config"
)

var startTime time.Time
var timerName string

func Debug(guildId string, logText string, a ...any) {
	if config.ServiceSettings.LOGGINGLEVEL <= config.LoggingLevelDebug {
		infoLine, formattedLogText := ParseLoggingText(guildId, logText, a...)
		SendLogs(infoLine, formattedLogText, config.LoggingLevelDebug, true)
	}
}

func Debug_IgnoreDiscord(guildId string, logText string, a ...any) {
	if config.ServiceSettings.LOGGINGLEVEL <= config.LoggingLevelDebug {
		infoLine, formattedLogText := ParseLoggingText(guildId, logText, a...)
		SendLogs(infoLine, formattedLogText, config.LoggingLevelDebug, false)
	}
}

// StartTimer is used to start a timer when you want to time a method
func StartTimer(tName string) {
	startTime = time.Now()
	timerName = tName
}

// EndTimer is called when you want to end the timer and log the results
func EndTimer() {
	timeNow := time.Now()

	elapsed := timeNow.Sub(startTime)

	Info("Timer %v took %v", timerName, elapsed)
}
