package bangCommands

import (
	"github.com/bwmarrin/discordgo"
	testhandler "github.com/dabi-ngin/discgo-bot/Bot/Handlers/TestHandler"
)

var (
	commandTable = make(map[string]func(message *discordgo.MessageCreate) error)
)

func Init() bool {
	commandTable["test"] = testhandler.HandleNewMessage
	return true
}

func GetCommand(query string) func(message *discordgo.MessageCreate) error {
	if val, ok := commandTable[query]; ok {
		return val
	} else {
		return nil
	}
}

func RunCommand(command string, message *discordgo.MessageCreate) error {
	err := commandTable[command](message)
	return err
}
