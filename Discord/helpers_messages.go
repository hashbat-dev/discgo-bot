package discord

import (
	"bytes"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	config "github.com/hashbat-dev/discgo-bot/Config"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

// Replies to the Message Object with an Error, if Error Text is provided it will show this to the user.
func Message_ReplyWithError(message *discordgo.Message, replyToQuoted bool, errorText string) *discordgo.Message {
	sendText := "An Error occured, please try again."
	if errorText != "" {
		sendText = "An Error occured: " + errorText
	}

	return Message_ReplyWithMessage(message, replyToQuoted, sendText)
}

func Message_ReplyWithMessage(message *discordgo.Message, replyToQuoted bool, messageText string) *discordgo.Message {
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

func Message_EditText(message *discordgo.Message, messageText string) {
	_, err := config.Session.ChannelMessageEdit(message.ChannelID, message.ID, messageText)
	if err != nil {
		logger.Error(message.GuildID, err)
	}
}

func Message_ReplyWithImage(message *discordgo.Message, replyToQuotedMessage bool, imageName string, imageBuffer *bytes.Buffer) error {
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

func Message_SendImage(channelId string, guildId string, imgExtension string, imageBuffer *bytes.Buffer) (string, error) {
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

func Message_Delete(message *discordgo.Message) {
	err := config.Session.ChannelMessageDelete(message.ChannelID, message.ID)
	if err != nil {
		logger.Error(message.GuildID, err)
	}
}

func Message_GetObject(guildId string, channelId string, messageId string) (*discordgo.Message, error) {
	message, err := config.Session.ChannelMessage(channelId, messageId)
	if err != nil {
		logger.Error(guildId, err)
		return nil, err
	}

	return message, nil
}
