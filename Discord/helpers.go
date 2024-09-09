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
	if private {
		config.Session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
				Flags:  discordgo.MessageFlagsEphemeral,
			},
		})
	} else {
		config.Session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
			},
		})
	}
}
