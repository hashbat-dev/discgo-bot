package helpers

import (
	"github.com/ZestHusky/femboy-control/Bot/audit"
	"github.com/ZestHusky/femboy-control/Bot/config"
	"github.com/bwmarrin/discordgo"
)

func ReplyToMessageWithText(message *discordgo.MessageCreate, textToReplyWith string) error {

	_, err := config.Session.ChannelMessageSendReply(message.ChannelID, textToReplyWith, message.Reference())
	if err != nil {
		audit.Error(err)
	}

	return err

}
