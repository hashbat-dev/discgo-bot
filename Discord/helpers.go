package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

func GenericErrorEmbed(errorText string) *embed.Embed {
	errEmbed := embed.NewEmbed()
	errEmbed.SetTitle("Error")
	if errorText == "" {
		errEmbed.SetDescription("An Error occured processing your request. Please try again or contact us through /support if this continues.")
	} else {
		errEmbed.SetDescription(fmt.Sprintf("An Error occured processing your request: %s", errorText))
	}

	return errEmbed
}

func GetAssociatedMessageFromInteraction(i *discordgo.InteractionCreate) (string, *discordgo.Message) {
	messageID := i.ApplicationCommandData().TargetID
	if messageID == "" {
		logger.ErrorText(i.GuildID, "MakeMeme: No MessageID provided")
	}

	message := i.ApplicationCommandData().Resolved.Messages[messageID]

	return messageID, message
}
