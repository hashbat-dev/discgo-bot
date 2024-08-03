package animegif

import (
	"fmt"

	"github.com/dabi-ngin/discgo-bot/Bot/audit"

	nb "github.com/Yakiyo/nekos_best.go"
	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	"github.com/dabi-ngin/discgo-bot/Bot/helpers"
	logger "github.com/dabi-ngin/discgo-bot/Bot/logging"
)

// AnimeGif https://docs.nekos.best/examples/unofficial/go.html
func AnimeGif(discord *discordgo.Session, message *discordgo.MessageCreate, command string) {

	// Get the receiver
	userFrom, userTo := helpers.GetReactionSourceUsers(message)

	// Get the GIF
	res, err := nb.Fetch(command)
	if err != nil {
		fmt.Println("[ANIMEGIF] ERROR GETTING GIF FROM API")
		fmt.Println("[ANIMEGIF] " + err.Error())
		logger.SendError(message)
		return
	} else {

		// Set Text
		showText := helpers.GetText(command, true, userFrom, userTo)

		fmt.Println("[ANIMEGIF] showText: " + showText + ", animeUrl: " + res.Url)
		e := embed.NewEmbed()
		e.SetDescription(showText)
		e.SetImage(res.Url)

		if message.ReferencedMessage != nil {
			_, err = discord.ChannelMessageSendEmbedReply(message.ChannelID, e.MessageEmbed, message.MessageReference)
		} else {
			_, err = discord.ChannelMessageSendEmbed(message.ChannelID, e.MessageEmbed)
		}

		if err == nil {
			fmt.Println("[ANIMEGIF] Sent AnimeGif for !" + command)
		} else {
			fmt.Println("[ANIMEGIF] ERROR SENDING ANIMEGIF")
			fmt.Println("[ANIMEGIF] " + err.Error())
			logger.SendError(message)
			return
		}
	}

}

func GiveListHandler(discord *discordgo.Session, interaction *discordgo.InteractionCreate) {

	fmt.Println("[ANIMEGIF] Received GiveList Handler command")

	var embedContent []*discordgo.MessageEmbed
	embedContent = append(embedContent, helpers.ReactList())
	if embedContent[0].Title == "Error" {
		return
	}

	err := discord.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: embedContent,
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		audit.Error(fmt.Errorf("animegif :: GiveListHandler :: error message=%s\n", err))
	}
}
