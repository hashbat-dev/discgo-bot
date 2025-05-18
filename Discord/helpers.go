package discord

import (
	"bytes"
	"strings"

	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	"github.com/google/uuid"
	config "github.com/hashbat-dev/discgo-bot/Config"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
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

func SendEmbedFromInteraction(i *discordgo.InteractionCreate, embedTitle string, embedText string, color int) {

	errEmbed := embed.NewEmbed()
	errEmbed.SetTitle(embedTitle)
	errEmbed.SetDescription(embedText)
	if color > 0 {
		errEmbed.SetColor(color)
	} else if strings.Contains(embedTitle, "Error") {
		errEmbed.SetColor(config.EmbedColourRed)
	}
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

func UpdateInteractionResponse(i *discordgo.InteractionCreate, embedTitle string, embedText string, color int) {
	newEmbed := embed.NewEmbed()
	newEmbed.SetTitle(embedTitle)
	newEmbed.SetDescription(embedText)
	if color > 0 {
		newEmbed.SetColor(color)
	}

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

func ReplyToMessage(message *discordgo.MessageCreate, text string) error {
	_, err := config.Session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
		Content:   text,
		Reference: message.Reference(),
	})

	if err != nil {
		logger.Error(message.GuildID, err)
	}

	return err
}

func ReplyToMessageWithEmbed(message *discordgo.Message, embed embed.Embed) {
	_, err := config.Session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
		Content: "",
		Embeds: []*discordgo.MessageEmbed{
			embed.MessageEmbed,
		},
		Reference: message.Reference(),
	})
	if err != nil {
		logger.Error(message.GuildID, err)
	}
}

func SendMessageWithImageBuffer(channelId string, guildId string, imgExtension string, imageBuffer *bytes.Buffer) (string, error) {
	imageName := uuid.New().String() + imgExtension
	fileObj := &discordgo.File{
		Name:   imageName,
		Reader: imageBuffer,
	}

	msg, err := config.Session.ChannelMessageSendComplex(channelId, &discordgo.MessageSend{
		Files: []*discordgo.File{fileObj},
	})

	if err != nil {
		logger.Error(guildId, err)
		return "", err
	}

	return msg.ID, err
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

func InteractionLoadingStart(i *discordgo.InteractionCreate) {
	err := config.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		logger.Error(i.GuildID, err)
	}
}

func InteractionLoadingFinish(i *discordgo.InteractionCreate) {
	err := config.Session.InteractionResponseDelete(i.Interaction)
	if err != nil {
		logger.Error(i.GuildID, err)
	}
}
