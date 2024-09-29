package logger

import (
	"fmt"

	config "github.com/hashbat-dev/discgo-bot/Config"
)

// Remind is used to remind us to fix things that are temporary, surely no workplace could possible operate like this. Right?
func Remind(logText string) {
	if config.ServiceSettings.LOGGINGLEVEL <= config.LoggingLevelInfo {
		fmt.Printf("%v%v %v%v\n", config.Colours["magenta"].Terminal, "[REMINDER]", logText, config.Colours["default"].Terminal)
	}
}
