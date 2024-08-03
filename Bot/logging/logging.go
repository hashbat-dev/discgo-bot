package logging

import (
	"math/rand"
	"runtime"

	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	"github.com/dabi-ngin/discgo-bot/Bot/config"
	"github.com/dabi-ngin/discgo-bot/Bot/constants"
)

// Send a generic "Bot has had an error" message back to the channel which created the message
func SendError(message *discordgo.MessageCreate) {
	e := embed.NewEmbed()
	e.SetTitle("Error")
	e.SetDescription(getRandomText(constants.ERROR_RAND_TEXT))
	e.SetImage(constants.GIF_AE_CRY)
	go config.Session.ChannelMessageSendEmbed(message.ChannelID, e.MessageEmbed)
}

// Same function as above, but takes a MessageUpdate object for edits
func SendErrorFromEdit(message *discordgo.MessageUpdate) {
	e := embed.NewEmbed()
	e.SetTitle("Error")
	e.SetDescription(getRandomText(constants.ERROR_RAND_TEXT))
	e.SetImage(constants.GIF_AE_CRY)
	go config.Session.ChannelMessageSendEmbed(message.ChannelID, e.MessageEmbed)
}

// Sends a "Bot has had an error" message but with a custom provided message
func SendErrorMsg(message *discordgo.MessageCreate, customMsg string) {
	e := embed.NewEmbed()
	e.SetTitle("Error")
	e.SetDescription(customMsg)
	e.SetImage(constants.GIF_AE_CRY)
	go config.Session.ChannelMessageSendEmbed(message.ChannelID, e.MessageEmbed)
}

// Sends a "Bot has had an error" message but with a custom provided message, replying to a message
func SendErrorMsgReply(message *discordgo.MessageCreate, customMsg string) {
	e := embed.NewEmbed()
	e.SetTitle("Error")
	e.SetDescription(customMsg)
	e.SetImage(constants.GIF_AE_CRY)
	go config.Session.ChannelMessageSendEmbedReply(message.ChannelID, e.MessageEmbed, message.ReferencedMessage.Reference())
}

// Sends an error from a /slash command input
func SendErrorInteraction(interaction *discordgo.InteractionCreate) {
	e := embed.NewEmbed()
	e.SetTitle("Error")
	e.SetDescription(getRandomText(constants.ERROR_RAND_TEXT))
	e.SetImage(constants.GIF_AE_CRY)
	go config.Session.ChannelMessageSendEmbed(interaction.ChannelID, e.MessageEmbed)
}

// Sends an error to a private /slash input that is only seen by the requesting user
func SendErrorMsgInteraction(interaction *discordgo.InteractionCreate, title string, message string, private bool) {
	embedNotAdmin := embed.NewEmbed()
	embedNotAdmin.SetTitle(title)
	embedNotAdmin.SetDescription(message)
	embedNotAdmin.SetImage(constants.GIF_AE_CRY)

	var embedContent []*discordgo.MessageEmbed
	embedContent = append(embedContent, embedNotAdmin.MessageEmbed)
	if embedContent[0].Title == "Error" {
		return
	}

	if private {
		go config.Session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: embedContent,
				Flags:  discordgo.MessageFlagsEphemeral,
			},
		})
	} else {
		go config.Session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: embedContent,
			},
		})
	}
}

func SendMessageInteraction(interaction *discordgo.InteractionCreate, title string, message string, imgUrl string, footer string, isPrivate bool) {
	e := embed.NewEmbed()
	e.SetTitle(title)
	e.SetDescription(message)
	if imgUrl != "" {
		e.SetImage(imgUrl)
	}
	if footer != "" {
		e.SetFooter(footer)
	}

	var embedContent []*discordgo.MessageEmbed
	embedContent = append(embedContent, e.MessageEmbed)
	if embedContent[0].Title == "Error" {
		return
	}

	if isPrivate {
		go config.Session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:  discordgo.MessageFlagsEphemeral,
				Embeds: embedContent,
			},
		})
	} else {
		go config.Session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: embedContent,
			},
		})
	}

}

func MessageRefObj(msg *discordgo.Message) discordgo.MessageReference {
	return discordgo.MessageReference{
		MessageID: msg.ID,
		ChannelID: msg.ChannelID,
		GuildID:   msg.GuildID,
	}
}

func getRandomText(inputSlice []string) string {
	return inputSlice[0+rand.Intn(len(inputSlice))]
}

func CurrentFunctionName() string {
	pc, _, _, ok := runtime.Caller(3)
	if !ok {
		return ""
	}
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return ""
	}
	return fn.Name()
}
