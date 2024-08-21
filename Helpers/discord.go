package helpers

import (
	"github.com/bwmarrin/discordgo"
	config "github.com/dabi-ngin/discgo-bot/Config"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
)

func SendUserError(message *discordgo.MessageCreate, errorText string) {

	//TODO - Flesh this out
	sendText := "An Error occured: " + errorText
	_, err := config.Session.ChannelMessageSendReply(message.ChannelID, sendText, message.Reference())
	if err != nil {
		logger.Error(message.GuildID, err)
	}

}

func SendUserMessage(message *discordgo.MessageCreate, messageText string) {

	//TODO - Flesh this out
	_, err := config.Session.ChannelMessageSendReply(message.ChannelID, messageText, message.Reference())
	if err != nil {
		logger.Error(message.GuildID, err)
	}

}
