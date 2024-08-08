package logger

import (
	"fmt"

	config "github.com/dabi-ngin/discgo-bot/Config"
)

// Remind is used to remind us to fix things that are temporary, surely no workplace could possible operate like this. Right?
func Remind(logText string) {
	if config.LoggingLevel <= config.LoggingLevelInfo {
		fmt.Printf("%v%v %v%v\n", ColourMagenta, "[REMINDER]", logText, ColourReset)
	}
}
