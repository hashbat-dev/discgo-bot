package gifbank

import (
	"fmt"

	"github.com/ZestHusky/femboy-control/Bot/audit"

	dbhelper "github.com/ZestHusky/femboy-control/Bot/dbhelpers"

	"github.com/ZestHusky/femboy-control/Bot/config"
	"github.com/ZestHusky/femboy-control/Bot/helpers"
	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	"github.com/google/uuid"
)

func DumpCategory(interaction *discordgo.InteractionCreate) {

	optionMap := helpers.GetOptionMap(interaction)
	category := helpers.GetOptionStringValue(optionMap, "gif-category")

	// Create Embed for responses/errors
	embTitle := "GIF Dump"
	embedStart := embed.NewEmbed()
	embedStart.SetTitle(embTitle)
	embedStart.SetDescription("Getting your Dump, sit tight...")
	resp, err := config.Session.ChannelMessageSendEmbed(interaction.ChannelID, embedStart.MessageEmbed)
	if err != nil {
		audit.Error(err)
		return
	}

	// Get list of all GIFs in the Category
	gifs, err := dbhelper.GetAllGifs(category)
	if err != nil {
		audit.Error(err)
		config.Session.ChannelMessageEditEmbed(interaction.ChannelID, resp.ID, helpers.GenericErrorEmbed(embTitle, "Couldn't get GIF Collection"))
		return
	}

	if len(gifs) == 0 {
		config.Session.ChannelMessageEditEmbed(interaction.ChannelID, resp.ID, helpers.GenericEmbed(embTitle, "There are no GIFs in that collection"))
		return
	} else {
		config.Session.ChannelMessageEditEmbed(interaction.ChannelID, resp.ID, helpers.GenericEmbed(embTitle, "Outputting GIFs ("+fmt.Sprint(len(gifs))+" total)"))
	}

	random := uuid.New().String()[:5]
	threadName := "dump-" + category + "-" + random
	thread, err := config.Session.ThreadStart(interaction.ChannelID, threadName, discordgo.ChannelTypeGuildText, 60)
	if err != nil {
		audit.Error(err)
		config.Session.ChannelMessageEditEmbed(interaction.ChannelID, resp.ID, helpers.GenericErrorEmbed(embTitle, "Couldn't create Thread for the Dump"))
		return
	}

	// Begin Outputting GIFs
	failedItems := 0
	for _, gif := range gifs {
		_, err := config.Session.ChannelMessageSend(thread.ID, gif.GifURL)
		if err != nil {
			audit.ErrorWithText("GifID: "+fmt.Sprint(gif.ID), err)
			failedItems++
		}
	}

	// Conclude the Embed
	concludeText := "Output completed to thread: <#" + thread.ID + ">"
	if failedItems > 0 {
		concludeText += "\n\nFailed to post " + fmt.Sprint(failedItems) + "gif"
		if failedItems != 1 {
			concludeText += "s"
		}
	}
	config.Session.ChannelMessageEditEmbed(interaction.ChannelID, resp.ID, helpers.GenericEmbed(embTitle, concludeText))
}
