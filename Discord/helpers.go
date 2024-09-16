package discord

import (
	"bytes"

	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	config "github.com/dabi-ngin/discgo-bot/Config"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
)

func SendUserError(message *discordgo.MessageCreate, errorText string) {
	sendText := "An Error occured: " + errorText
	_, err := config.Session.ChannelMessageSendReply(message.ChannelID, sendText, message.Reference())
	if err != nil {
		logger.Error(message.GuildID, err)
	}
}

func SendUserMessage(message *discordgo.MessageCreate, messageText string) {
	_, err := config.Session.ChannelMessageSendReply(message.ChannelID, messageText, message.Reference())
	if err != nil {
		logger.Error(message.GuildID, err)
	}
}

func SendUserMessageReply(message *discordgo.MessageCreate, replyToQuoted bool, messageText string) *discordgo.Message {
	replyTo := message.Reference()
	if replyToQuoted && message.ReferencedMessage != nil {
		replyTo = message.ReferencedMessage.Reference()
	}
	msg, err := config.Session.ChannelMessageSendReply(message.ChannelID, messageText, replyTo)
	if err != nil {
		logger.Error(message.GuildID, err)
		return nil
	}

	return msg
}

func EditMessage(message *discordgo.Message, messageText string) {
	_, err := config.Session.ChannelMessageEdit(message.ChannelID, message.ID, messageText)
	if err != nil {
		logger.Error(message.GuildID, err)
	}
}

func GenericErrorEmbed() *embed.Embed {
	errEmbed := embed.NewEmbed()
	errEmbed.SetTitle("Error")
	errEmbed.SetDescription("An Error occured processing your request. Please try again or contact us through /support if this continues.")
	return errEmbed
}

func SendGenericErrorFromInteraction(i *discordgo.InteractionCreate) {
	errEmbed := GenericErrorEmbed()
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

func SendEmbedFromInteraction(i *discordgo.InteractionCreate, embedTitle string, embedText string) {
	errEmbed := embed.NewEmbed()
	errEmbed.SetTitle(embedTitle)
	errEmbed.SetDescription(embedText)
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

func UpdateInteractionResponse(i *discordgo.InteractionCreate, embedTitle string, embedText string) {
	newEmbed := embed.NewEmbed()
	newEmbed.SetTitle(embedTitle)
	newEmbed.SetDescription(embedText)

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

func UpdateInteractionResponseWithGenericError(i *discordgo.InteractionCreate) {
	errEmbed := GenericErrorEmbed()
	embeds := []*discordgo.MessageEmbed{errEmbed.MessageEmbed}
	_, err := config.Session.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &embeds,
	})

	if err != nil {
		logger.Error(i.GuildID, err)
	}
}

func ReplyToMessageWithImageBuffer(message *discordgo.MessageCreate, replyToQuotedMessage bool, imageName string, imageBuffer *bytes.Buffer) error {
	fileObj := &discordgo.File{
		Name:   imageName,
		Reader: imageBuffer,
	}

	replyToMsg := message.Reference()
	if replyToQuotedMessage {
		replyToMsg = message.ReferencedMessage.MessageReference
	}

	_, err := config.Session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
		Files:     []*discordgo.File{fileObj},
		Reference: replyToMsg,
	})

	if err != nil {
		logger.Error(message.GuildID, err)
	}

	return err
}

func ReplyToInteractionWithEmbed(interaction *discordgo.InteractionCreate, embed *embed.Embed, private bool) {
	var err error
	if private {
		err = config.Session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
				Flags:  discordgo.MessageFlagsEphemeral,
			},
		})
	} else {
		err = config.Session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
			},
		})
	}

	if err != nil {
		logger.Error(interaction.GuildID, err)
	}
}

func DeleteMessage(message *discordgo.MessageCreate) {
	err := config.Session.ChannelMessageDelete(message.ChannelID, message.ID)
	if err != nil {
		logger.Error(message.GuildID, err)
	}
}

func DeleteMessageObject(message *discordgo.Message) {
	err := config.Session.ChannelMessageDelete(message.ChannelID, message.ID)
	if err != nil {
		logger.Error(message.GuildID, err)
	}
}
