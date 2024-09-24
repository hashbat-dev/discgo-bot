package logger

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	config "github.com/dabi-ngin/discgo-bot/Config"
)

func FormatInboundLogText(logText string, a ...any) string {
	if len(a) > 0 {
		return fmt.Sprintf(logText, a...)
	} else {
		return logText
	}
}

type LogInfo struct {
	DateTime   time.Time
	CodeSource string
	GuildID    string
}

func SendLogs(infoLine LogInfo, logText string, logLevel int, sendToDiscord bool) {
	SendToConsole(infoLine, logText, logLevel)
	SendLogsToDashboard(infoLine, logText, logLevel)
	if !config.ServiceSettings.LOGTODISCORD || !sendToDiscord {
		return
	}
	SendLogToDiscord(infoLine, logText, logLevel)
}

func ParseLoggingText(guildId string, logText string, a ...any) (LogInfo, string) {
	logInfo := LogInfo{
		DateTime:   time.Now(),
		CodeSource: GetStack(),
		GuildID:    guildId,
	}
	formattedLogText := logText
	if len(a) > 0 {
		formattedLogText = FormatInboundLogText(logText, a...)
	}
	return logInfo, formattedLogText
}

func SendToConsole(logInfo LogInfo, logText string, logLevel int) {
	logType := config.LoggingLevels[logLevel]

	infoLine := fmt.Sprintf("%v | %v", logInfo.DateTime.Format("02/01/06 15:04:05.000"), logInfo.CodeSource)
	if logInfo.GuildID != "" {
		infoLine += " | " + logInfo.GuildID
	}
	if config.ServiceSettings.ISDEV {
		infoLine += " | " + config.ServiceSettings.HOSTNAME
	}

	fmt.Printf("%v[%v] %v :: %v %v \n", logType.Colour.Terminal, strings.ToUpper(logType.Name), infoLine, logText, config.Colours["default"].Terminal)
}

// GetStack gets that bread homie
func GetStack() string {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			return ParseStackTrace(string(buf[:n]))
		}
		buf = make([]byte, len(buf)*2)
	}
}

func ParseStackTrace(stack string) string {
	retVal := ""
	lines := strings.Split(stack, "\n")

	isFirst := true

	for i := 7; i < len(lines)-1; i++ {

		line := lines[i]
		if strings.Contains(line, "logger") || strings.Contains(line, "created by") || strings.Contains(line, "main.go") {
			continue
		}

		botIndex := strings.Index(line, config.ROOT_FOLDER)
		botIndexLength := len(config.ROOT_FOLDER)

		lastIndex := 0

		isFileLine := strings.Contains(line, " +")
		if isFileLine {
			lastIndex = strings.LastIndex(line, " +") - 1
		} else {
			lastIndex = strings.LastIndex(line, ")")
		}

		if !isFileLine && !config.ServiceSettings.LOGFUNCTIONS {
			continue
		}

		if botIndex != -1 && lastIndex != -1 {

			appended := false
			if isFirst {
				isFirst = false
			} else {
				if !config.ServiceSettings.LOGFUNCTIONS {
					break
				}
				retVal += " <= "
				appended = true
			}

			retVal += line[botIndex+botIndexLength : lastIndex+1]

			if appended && !config.ServiceSettings.VERBOSESTACK {
				break
			}

		}
	}

	return "./" + RemoveTextInParentheses(retVal)
}

func RemoveTextInParentheses(input string) string {
	var result strings.Builder
	inParens := false

	for _, char := range input {
		if char == '(' {
			inParens = true
			result.WriteRune(char)
		} else if char == ')' {
			inParens = false
			result.WriteRune(char)
		} else if !inParens {
			result.WriteRune(char)
		}
	}

	return result.String()
}
