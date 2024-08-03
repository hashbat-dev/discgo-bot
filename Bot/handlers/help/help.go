package help

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dabi-ngin/discgo-bot/Bot/logging"
)

func GetHelpText(interaction *discordgo.InteractionCreate, commandString string) {
	cmdText := "Here are all commands available to the bot by using !<command>, some of these require you to reply to a message.\n" + commandString
	logging.SendMessageInteraction(interaction, "Bot Command List", cmdText, "", "", true)
}
