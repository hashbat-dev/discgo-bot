package slurs

import (
	"github.com/ZestHusky/femboy-control/Bot/audit"
	dbhelper "github.com/ZestHusky/femboy-control/Bot/dbhelpers"
	"github.com/ZestHusky/femboy-control/Bot/helpers"
	"github.com/ZestHusky/femboy-control/Bot/logging"
	"github.com/bwmarrin/discordgo"
)

func DefineASlur(interaction *discordgo.InteractionCreate) {

	optionMap := helpers.GetOptionMap(interaction)
	inSlur := helpers.GetOptionStringValue(optionMap, "slur")
	target, description, err := dbhelper.GetSlurDefinition(inSlur)
	if err != nil {
		audit.Error(err)
	}

	if target == "" || description == "" {
		description = "That must be a super cool slur, I don't even know what it means!"
	}
	embedText := "**Target:** " + target + "\n"
	embedText += "**Definition:** " + description
	logging.SendMessageInteraction(interaction, inSlur, embedText, "", "", false)

}
