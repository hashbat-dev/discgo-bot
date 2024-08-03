package gifbank

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	"github.com/dabi-ngin/discgo-bot/Bot/audit"
	"github.com/dabi-ngin/discgo-bot/Bot/config"
	dbhelpers "github.com/dabi-ngin/discgo-bot/Bot/dbhelpers"
	"github.com/dabi-ngin/discgo-bot/Bot/helpers"
	logger "github.com/dabi-ngin/discgo-bot/Bot/logging"
)

func Post(session *discordgo.Session, message *discordgo.MessageCreate, command string) {
	userFrom, userTo := helpers.GetReactionSourceUsers(message)

	// Send Message
	gif, err := dbhelpers.GetRandGifURL(command)
	if err != nil {
		audit.ErrorWithText("Command: "+command, err)
		logger.SendError(message)
		return
	}

	// See if we can get a caption for the gif.
	// If we CAN, we need to send it as an embed with the caption.
	// If NOT, we want to send it as a regular message.
	showText := helpers.GetText(command, false, userFrom, userTo)
	if showText == "" {
		if message.ReferencedMessage != nil {
			_, err = config.Session.ChannelMessageSendReply(message.ChannelID, gif, message.MessageReference)
		} else {
			_, err = config.Session.ChannelMessageSend(message.ChannelID, gif)
		}
	} else {
		e := embed.NewEmbed()
		e.SetDescription(showText)
		e.SetImage(gif)

		if message.ReferencedMessage != nil {
			_, err = config.Session.ChannelMessageSendEmbedReply(message.ChannelID, e.MessageEmbed, message.MessageReference)
		} else {
			_, err = config.Session.ChannelMessageSendEmbed(message.ChannelID, e.MessageEmbed)
		}
	}

	if err == nil {
		audit.Log("Sent Gif for category: " + command)
	} else {
		audit.Error(err)
		logger.SendError(message)
		return
	}
}

func PostFromEdit(session *discordgo.Session, message *discordgo.MessageUpdate, command string) {
	userFrom, userTo := helpers.GetReactionSourceUsers_FromEdit(message)

	// Send Message
	gif, err := dbhelpers.GetRandGifURL(command)
	if err != nil {
		audit.Error(err)
		logger.SendErrorFromEdit(message)
		return
	}

	// Set Text
	showText := helpers.GetText(command, true, userFrom, userTo)
	e := embed.NewEmbed()
	e.SetDescription(showText)
	e.SetImage(gif)

	if message.ReferencedMessage != nil {
		_, err = config.Session.ChannelMessageSendEmbedReply(message.ChannelID, e.MessageEmbed, message.MessageReference)
	} else {
		_, err = config.Session.ChannelMessageSendEmbed(message.ChannelID, e.MessageEmbed)
	}

	if err == nil {
		audit.Log("Sent for Category: " + command)
	} else {
		audit.Error(err)
		logger.SendErrorFromEdit(message)
		return
	}
}

func AddGIF(message *discordgo.MessageCreate, category string) {

	if message.ReferencedMessage == nil {
		logger.SendErrorMsg(message, "You've not replied to an Image, dummy!")
		return
	}

	gifLink, err := helpers.DoesMessageHaveImage(message.ReferencedMessage, "")
	if err != nil {
		logger.SendErrorMsg(message, err.Error())
		return
	}

	id, err := dbhelpers.InsertGIF(category, gifLink, message.Author.ID)
	if err != nil {
		audit.Error(err)
		logger.SendError(message)
	} else {
		e := embed.NewEmbed()
		e.SetTitle("Gif Added to the " + category + " bank")
		e.SetDescription("I will be sure to use this in future! " + helpers.GetEmote("yep", true))
		msgRef := logger.MessageRefObj(message.Message)
		config.Session.ChannelMessageSendEmbedReply(message.ChannelID, e.MessageEmbed, &msgRef)
		audit.Log("Added ID: " + fmt.Sprint(id) + ", Category: " + category)
	}

}

func Delete(message *discordgo.MessageCreate, category string) {
	// Can we find a gif in the message?
	if message.ReferencedMessage == nil {
		logger.SendErrorMsg(message, "You've not replied to a gif, dummy!")
		return
	}

	// Get the Gif
	msgContent := strings.Trim(message.ReferencedMessage.Content, " ")
	msgContentLower := strings.ToLower(msgContent)

	gifLink := ""
	if strings.Contains(msgContentLower, ".gif") || strings.Contains(msgContentLower, ".png") ||
		strings.Contains(msgContentLower, ".jpg") || strings.Contains(msgContentLower, ".webp") || strings.Contains(msgContentLower, "tenor.com/") {
		gifLink = msgContent
	}

	if gifLink == "" {
		if len(message.ReferencedMessage.Attachments) > 0 {
			attachUrl := message.ReferencedMessage.Attachments[0].URL
			if strings.Contains(attachUrl, ".gif") || strings.Contains(attachUrl, ".png") ||
				strings.Contains(attachUrl, ".jpg") || strings.Contains(attachUrl, ".webp") || strings.Contains(msgContentLower, "tenor.com/") {
				gifLink = attachUrl
			}
		}
	}

	if gifLink == "" {
		errText := "That's not an image!"
		logger.SendErrorMsg(message, errText)
		return
	}

	err := dbhelpers.DeleteGIF(category, gifLink)
	if err != nil {
		audit.Error(err)
		logger.SendError(message)
	} else {
		e := embed.NewEmbed()
		e.SetTitle("Gif Deleted")
		e.SetDescription("I've deleted this gif, thanks! " + helpers.GetEmote("yep", true))
		config.Session.ChannelMessageSendEmbed(message.ChannelID, e.MessageEmbed)
	}
}
