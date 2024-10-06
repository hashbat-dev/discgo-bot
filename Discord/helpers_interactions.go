package discord

import (
	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	config "github.com/hashbat-dev/discgo-bot/Config"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

// Replies to an Interaction with an Error, if the errorText is blank it will provide a generic error message.
func Interactions_SendError(i *discordgo.InteractionCreate, errorText string) {
	errEmbed := GenericErrorEmbed(errorText)
	err := config.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{errEmbed.MessageEmbed},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		logger.Error(i.GuildID, err)
	}
}

func Interactions_SendMessage(i *discordgo.InteractionCreate, messageTitle string, messageText string) {
	errEmbed := embed.NewEmbed()
	errEmbed.SetTitle(messageTitle)
	errEmbed.SetDescription(messageText)
	err := config.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{errEmbed.MessageEmbed},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		logger.Error(i.GuildID, err)
	}
}

func Interactions_EditText(i *discordgo.InteractionCreate, messageTitle string, messageText string) {
	newEmbed := embed.NewEmbed()
	newEmbed.SetTitle(messageTitle)
	newEmbed.SetDescription(messageText)

	wipeContent := ""
	embeds := []*discordgo.MessageEmbed{newEmbed.MessageEmbed}
	_, err := config.Session.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content:    &wipeContent,
		Components: &[]discordgo.MessageComponent{},
		Embeds:     &embeds,
	})

	if err != nil {
		logger.Error(i.GuildID, err)
	}
}

func Interactions_EditIntoError(i *discordgo.InteractionCreate, errorText string) {
	errEmbed := GenericErrorEmbed(errorText)
	wipeContent := ""
	embeds := []*discordgo.MessageEmbed{errEmbed.MessageEmbed}
	_, err := config.Session.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content:    &wipeContent,
		Components: &[]discordgo.MessageComponent{},
		Embeds:     &embeds,
	})

	if err != nil {
		logger.Error(i.GuildID, err)
	}
}

func Interaction_SendPublicMessage(interaction *discordgo.InteractionCreate, title string, text string) {
	embed := embed.NewEmbed()
	embed.SetTitle(title)
	embed.SetDescription(text)
	err := config.Session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
		},
	})

	if err != nil {
		logger.Error(interaction.GuildID, err)
	}
}
