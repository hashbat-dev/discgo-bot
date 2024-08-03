package audit

import (
	"fmt"
	"time"

	"github.com/ZestHusky/femboy-control/Bot/config"
)

var verboseStack bool = false
var logFunctions bool = false

func Log(logText string) {
	stack := GetStack()
	fmt.Println(stack + " | " + logText)
	if config.IsDev {
		stack += " | " + config.ServerName
	}
	go PostLog(time.Now().Format("02/01/06 15:04:05.000")+" | "+stack, logText, false)
}

func LogDevOnly(logText string) {
	if !config.IsDev {
		return
	}
	stack := GetStack()
	fmt.Println(stack + " | " + logText)
	if config.IsDev {
		stack += " | " + config.ServerName
	}
	go PostLog(time.Now().Format("02/01/06 15:04:05.000")+" | "+stack, logText, false)
}

func Error(err error) {
	stack := GetStack()
	fmt.Println(stack + " | " + err.Error())
	if config.IsDev {
		stack += " | " + config.ServerName
	}
	go PostLog(time.Now().Format("02/01/06 15:04:05.000")+" | "+stack, err.Error(), false)
}

func ErrorWithText(description string, err error) {
	stack := GetStack() + " | " + description
	fmt.Println(stack + " | " + err.Error())
	if config.IsDev {
		stack += " | " + config.ServerName
	}
	go PostLog(time.Now().Format("02/01/06 15:04:05.000")+" | "+stack, err.Error(), false)
}

var LogCache []string

func PostLog(infoLine string, logLine string, isError bool) {
	addLog := ""
	if isError {
		addLog += "\n" + "**ERROR**"
	}
	addLog += "```asciidoc\n[" + infoLine + "]\n" + logLine + "```"

	LogCache = append(LogCache, addLog)
}
