package logger

import (
	"fmt"
	"time"

	config "github.com/dabi-ngin/discgo-bot/Config"
)

var startTime time.Time
var timerName string

func Debug(guildId string, logText string, a ...any) {
	if config.LoggingLevel <= config.LoggingLevelDebug {
		infoLine := fmt.Sprintf("%v | %v", time.Now().Format("02/01/06 15:04:05.000"), GetStack())
		if guildId != "" {
			infoLine += " | " + guildId
		}
		formattedLogText := logText
		if len(a) > 0 {
			formattedLogText = FormatInboundLogText(logText, a...)
		}
		SendToConsole(infoLine, formattedLogText, config.LoggingLevelDebug)
		if config.IsDev {
			infoLine += " | " + config.HostName
		}
		SendLogToDiscord(infoLine, formattedLogText, config.LoggingLevelDebug)
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
