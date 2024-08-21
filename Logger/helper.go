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

var (
	ColourReset   = "\033[0m"
	ColourRed     = "\033[31m"
	ColourGreen   = "\033[32m"
	ColourYellow  = "\033[33m"
	ColourWhite   = "\033[97m"
	ColourMagenta = "\033[35m"
)

type LogInfo struct {
	DateTime   time.Time
	CodeSource string
	GuildID    string
}

func SendLogs(infoLine LogInfo, logText string, logLevel int, sendToDiscord bool) {
	SendToConsole(infoLine, logText, logLevel)
	SendLogsToDashboard(infoLine, logText, logLevel)
	if !config.LogToDiscord || !sendToDiscord {
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
	var useColour string
	var logType string
	switch logLevel {
	case config.LoggingLevelAdmin:
		useColour = ColourMagenta
		logType = "[ADMIN]"
	case config.LoggingLevelError:
		useColour = ColourRed
		logType = "[ERROR]"
	case config.LoggingLevelWarn:
		useColour = ColourYellow
		logType = "[WARN]"
	case config.LoggingLevelEvent:
		useColour = ColourGreen
		logType = "[EVENT]"
	case config.LoggingLevelInfo:
		useColour = ColourWhite
		logType = "[INFO]"
	default:
		useColour = ColourWhite
	}

	infoLine := fmt.Sprintf("%v | %v", logInfo.DateTime.Format("02/01/06 15:04:05.000"), logInfo.CodeSource)
	if logInfo.GuildID != "" {
		infoLine += " | " + logInfo.GuildID
	}
	if config.IsDev {
		infoLine += " | " + config.HostName
	}

	fmt.Printf("%v%v %v :: %v %v \n", useColour, logType, infoLine, logText, ColourReset)
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

		botIndex := strings.Index(line, config.BOT_SUB_FOLDER)
		botIndexLength := len(config.BOT_SUB_FOLDER)
		if botIndex == -1 {
			botIndex = strings.Index(line, config.ROOT_FOLDER)
			botIndexLength = len(config.ROOT_FOLDER)
		}
		lastIndex := 0

		isFileLine := strings.Contains(line, " +")
		if isFileLine {
			lastIndex = strings.LastIndex(line, " +") - 1
		} else {
			lastIndex = strings.LastIndex(line, ")")
		}

		if !isFileLine && !config.LoggingLogFunctions {
			continue
		}

		if botIndex != -1 && lastIndex != -1 {

			appended := false
			if isFirst {
				isFirst = false
			} else {
				if !config.LoggingVerboseStack {
					break
				}
				retVal += " <= "
				appended = true
			}

			retVal += line[botIndex+botIndexLength : lastIndex+1]

			if appended && !config.LoggingVerboseStack {
				break
			}

		}
	}

	return RemoveTextInParentheses(retVal)
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
