package logger

import (
	"fmt"

	config "github.com/hashbat-dev/discgo-bot/Config"
)

var LogsForDiscord []string
var LogsForDashboard []DashboardLog

type DashboardLog struct {
	LogInfo  LogInfo
	LogText  string
	LogLevel int
}

// Adds logs to the queue to send to the Discord Channel associated with the instance
func SendLogToDiscord(logInfo LogInfo, logText string, logLevel int) {
	addLog := ""
	if logLevel == 0 {
		addLog += "\n" + "**ERROR**"
	}
	infoLine := fmt.Sprintf("%v | %v", logInfo.DateTime.Format("02/01/06 15:04:05.000"), logInfo.CodeSource)
	if logInfo.GuildID != "" {
		infoLine += " | " + logInfo.GuildID
	}
	addLog += "```asciidoc\n[" + infoLine + "]\n" + logText + "```"

	LogsForDiscord = append(LogsForDiscord, addLog)
}

// Adds logs to the queue to send to the Dashboard
func SendLogsToDashboard(logInfo LogInfo, logText string, logLevel int) {
	newLog := DashboardLog{
		LogInfo:  logInfo,
		LogText:  logText,
		LogLevel: logLevel,
	}
	NewLogsForDashboard := append([]DashboardLog{newLog}, LogsForDashboard...)

	if len(NewLogsForDashboard) > config.ServiceSettings.DASHBOARDMAXLOGS {
		NewLogsForDashboard = NewLogsForDashboard[:len(NewLogsForDashboard)-1]
	}

	LogsForDashboard = NewLogsForDashboard
}
