package logger

import (
	"database/sql"

	config "github.com/hashbat-dev/discgo-bot/Config"
)

var Db *sql.DB

// Error - logs include stackTrace, requires LoggingLevel 1 or lower in config
func Error(guildId string, err error, a ...any) {
	if config.ServiceSettings.LOGGINGLEVEL <= config.LoggingLevelDebug {
		infoLine, formattedLogText := ParseLoggingText(guildId, err.Error(), a...)
		SendLogs(infoLine, formattedLogText, config.LoggingLevelError, true)
		InsertIntoDB(infoLine.CodeSource, err.Error(), guildId)
	}
}

func ErrorText(guildId string, message string, a ...any) {
	if config.ServiceSettings.LOGGINGLEVEL <= config.LoggingLevelDebug {
		infoLine, formattedLogText := ParseLoggingText(guildId, message, a...)
		SendLogs(infoLine, formattedLogText, config.LoggingLevelError, true)
		InsertIntoDB(infoLine.CodeSource, message, guildId)
	}
}

func ErrorWithoutDB(guildId string, err error, a ...any) {
	if config.ServiceSettings.LOGGINGLEVEL <= config.LoggingLevelDebug {
		infoLine, formattedLogText := ParseLoggingText(guildId, err.Error(), a...)
		SendLogs(infoLine, formattedLogText, config.LoggingLevelError, true)
	}
}

func Error_IgnoreDiscord(guildId string, err error, a ...any) {
	if config.ServiceSettings.LOGGINGLEVEL <= config.LoggingLevelDebug {
		infoLine, formattedLogText := ParseLoggingText(guildId, err.Error(), a...)
		SendLogs(infoLine, formattedLogText, config.LoggingLevelError, false)
		InsertIntoDB(infoLine.CodeSource, err.Error(), guildId)
	}
}

func InsertIntoDB(codeSource string, errorText string, guildId string) {
	stmt, err := Db.Prepare("INSERT INTO Errors (IsDev, CodeSource, ErrorText, GuildID) VALUES (?, ?, ?, ?)")
	if err != nil {
		ErrorWithoutDB(guildId, err)
		return
	}
	defer func(g string) {
		err := stmt.Close()
		if err != nil {
			ErrorWithoutDB(g, err)
		}
	}(guildId)

	isDev := 0
	if config.ServiceSettings.ISDEV {
		isDev = 1
	}
	_, err = stmt.Exec(isDev, codeSource, errorText, guildId)
	if err != nil {
		ErrorWithoutDB(guildId, err)
	}
}
