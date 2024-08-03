package audit

import (
	"fmt"

	"github.com/ZestHusky/femboy-control/Bot/config"
	"github.com/ZestHusky/femboy-control/Bot/constants"
	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
)

func SendNextLogBatch() {

	var newCache []string
	nextMsg := ""

	limitReached := false
	for _, s := range LogCache {
		if (len(nextMsg) + len(s)) >= constants.MAX_MESSAGE_LENGTH {
			limitReached = true
		}

		if limitReached {
			newCache = append(newCache, s)
		} else {
			nextMsg += s
		}
	}

	LogCache = newCache
	if nextMsg != "" {
		SendLogsToDiscordChannel(nextMsg)
	}
}

func SendLogsToDiscordChannel(logs string) {

	channelId := constants.CHANNEL_BOT_ERRORS_LIVE
	if config.IsDev {
		channelId = constants.CHANNEL_BOT_ERRORS_DEV
	}

	_, err := config.Session.ChannelMessageSend(channelId, logs)
	if err != nil {
		fmt.Println(err.Error())
	}

}

func SendInteractionResponse(interaction *discordgo.InteractionCreate, title string, message string, footer string, isError bool, private bool, imgUrl string) {
	e := embed.NewEmbed()
	e.SetTitle(title)
	e.SetDescription(message)
	if isError && imgUrl == "" {
		e.SetImage(constants.GIF_AE_CRY)
	} else if imgUrl != "" {
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

	var err error
	if private {
		err = config.Session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: embedContent,
				Flags:  discordgo.MessageFlagsEphemeral,
			},
		})
	} else {
		err = config.Session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: embedContent,
			},
		})
	}

	if err != nil {
		Error(err)
	}
}
