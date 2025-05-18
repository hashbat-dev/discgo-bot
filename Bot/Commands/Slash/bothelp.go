package slash

import (
	"github.com/bwmarrin/discordgo"
	config "github.com/hashbat-dev/discgo-bot/Config"
	discord "github.com/hashbat-dev/discgo-bot/Discord"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

func SendHelp(i *discordgo.InteractionCreate, correlationId string) {

	discord.SendEmbedFromInteraction(i, "", "Sending Help information to your DM's", 0)

	if len(config.UserBangHelpText) > 0 {
		userText := "## Chat Commands\n"
		userText += config.UserBangHelpText + "\n"
		err := discord.SendDM(i.GuildID, i.Member.User.ID, userText)
		if err != nil {
			logger.ErrorText(i.GuildID, "Error sending Help DM")
		}
	}

	if len(config.UserSlashHelpText) > 0 {
		userText := "## Slash Commands\n"
		userText += config.UserSlashHelpText
		err := discord.SendDM(i.GuildID, i.Member.User.ID, userText)
		if err != nil {
			logger.ErrorText(i.GuildID, "Error sending Help DM")
		}
	}

	if len(config.UserMsgCmdHelpText) > 0 {
		userText := "## Message Commands\n"
		userText += "*Right click a Message these options will be available in the 'Apps' sub-menu...*\n"
		userText += config.UserMsgCmdHelpText
		err := discord.SendDM(i.GuildID, i.Member.User.ID, userText)
		if err != nil {
			logger.ErrorText(i.GuildID, "Error sending Help DM")
		}
	}

}
