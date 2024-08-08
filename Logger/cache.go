package logger

var LogsForDiscord []string

// Adds logs to the queue to send to the Discord Channel associated with the instance
func SendLogToDiscord(infoLine string, logText string, logLevel int) {
	addLog := ""
	if logLevel == 0 {
		addLog += "\n" + "**ERROR**"
	}
	addLog += "```asciidoc\n[" + infoLine + "]\n" + logText + "```"

	LogsForDiscord = append(LogsForDiscord, addLog)
}
