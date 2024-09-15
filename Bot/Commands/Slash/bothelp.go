package slash

import (
	"github.com/bwmarrin/discordgo"
	config "github.com/dabi-ngin/discgo-bot/Config"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
)

func SendHelp(i *discordgo.InteractionCreate, correlationId string) {
	userText := "## Chat Commands\n"
	userText += config.UserBangHelpText + "\n"
	userText += "## Slash Commands\n"
	userText += config.UserSlashHelpText
	err := config.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: userText,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		logger.Error(i.GuildID, err)
	}
}
