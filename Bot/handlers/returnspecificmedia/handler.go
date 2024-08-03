package returnspecificmedia

import (
	"github.com/ZestHusky/femboy-control/Bot/audit"
	"github.com/ZestHusky/femboy-control/Bot/config"
	"github.com/ZestHusky/femboy-control/Bot/logging"
	"github.com/bwmarrin/discordgo"
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
