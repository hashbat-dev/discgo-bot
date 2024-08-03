package returnspecificmedia

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dabi-ngin/discgo-bot/Bot/audit"
	"github.com/dabi-ngin/discgo-bot/Bot/config"
	"github.com/dabi-ngin/discgo-bot/Bot/logging"
)

func HandleMessage(message *discordgo.MessageCreate, requestType string) {
	// grab media from requst
	returnMedia := ReturnStuff(requestType)

	_, err := config.Session.ChannelMessageSendReply(message.ChannelID, returnMedia, message.MessageReference)

	if err != nil {
		logging.SendErrorMsg(message, err.Error())
		audit.Error(err)
		return
	}

}
