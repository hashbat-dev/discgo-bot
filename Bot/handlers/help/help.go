package help

import (
	"github.com/ZestHusky/femboy-control/Bot/logging"
	"github.com/bwmarrin/discordgo"
)

func GetHelpText(interaction *discordgo.InteractionCreate, commandString string) {

	cmdText := "Here are all commands available to the bot by using !<command>, some of these require you to reply to a message.\n" + commandString
	logging.SendMessageInteraction(interaction, "Bot Command List", cmdText, "", "", true)

}
