package helpers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dabi-ngin/discgo-bot/Bot/audit"
	"github.com/dabi-ngin/discgo-bot/Bot/config"
)

func ReplyToMessageWithText(message *discordgo.MessageCreate, textToReplyWith string) error {

	_, err := config.Session.ChannelMessageSendReply(message.ChannelID, textToReplyWith, message.Reference())
	if err != nil {
		audit.Error(err)
	}

	return err

}
