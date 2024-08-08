package triggerCommands

import (
	"github.com/bwmarrin/discordgo"
	testhandler "github.com/dabi-ngin/discgo-bot/Bot/Handlers/TestHandler"
)

var (
	commandTable = make(map[string]func(message *discordgo.MessageCreate, trigger string) error)
)

func Init() bool {
	commandTable["triggertest"] = testhandler.HandleNewTrigger
	return true
}

// CheckForTriggerPhrase returns the trigger if it exists and returns and empty string if it does not
func CheckForTriggerPhrase(triggerQuery string) string {
	for k := range commandTable {
		if k == triggerQuery {
			return triggerQuery
		}
	}

	return ""
}

// RunTriggerCommand runs the command associated with a trigger
func RunTriggerCommand(command string, message *discordgo.MessageCreate) error {
	err := commandTable[command](message, command)
	return err
}
