package handlers

import (
	"github.com/ZestHusky/femboy-control/Bot/audit"
	"github.com/ZestHusky/femboy-control/Bot/config"
	"github.com/ZestHusky/femboy-control/Bot/constants"
	dbhelpers "github.com/ZestHusky/femboy-control/Bot/dbhelpers"
	"github.com/ZestHusky/femboy-control/Bot/gifbank"
	"github.com/ZestHusky/femboy-control/Bot/handlers/pogcorrection"
	"github.com/ZestHusky/femboy-control/Bot/helpers"
	"github.com/bwmarrin/discordgo"
)

// Called when a Message is Edited, shifty shenanigans afoot
func EditMessageHandler(session *discordgo.Session, message *discordgo.MessageUpdate) {
	if config.IsDev {
		if message.ChannelID != constants.CHANNEL_BOT_TEST {
			return
		}
	} else {
		if message.ChannelID == constants.CHANNEL_BOT_TEST {
			return
		}
	}
	go ProcessEdit(session, message)
}

func ProcessEdit(session *discordgo.Session, message *discordgo.MessageUpdate) {
	var err error
	var skipped bool

	if skipped, err = helpers.SkipProcessing(session, nil, message); skipped {
		return
	}
	if err != nil {
		audit.Error(err)
	}

	// POG CHECKS
	pogCorr := pogcorrection.Detect(message.Content, message.Author.ID, message.ReferencedMessage)
	switch pogCorr {
	case "tooth":
		dbhelpers.CountCommand("pog-correction", message.Author.ID)
		gifbank.PostFromEdit(session, message, "tooth")
	case "correct":
		pogcorrection.CorrectionFromEdit(message)
		dbhelpers.CountCommand("pog-correction", message.Author.ID)
	}

}
