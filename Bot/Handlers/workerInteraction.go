package handlers

import (
	"time"

	"github.com/bwmarrin/discordgo"
	config "github.com/hashbat-dev/discgo-bot/Config"
	interactions "github.com/hashbat-dev/discgo-bot/Interactions"
	reporting "github.com/hashbat-dev/discgo-bot/Reporting"
)

func WorkerInteractionCall(i *discordgo.InteractionCreate, call InteractionCall) {
	timeStarted := time.Now()
	interactions.SendUserText(i, "Processing...")
	call.Execute(i, call.CorrelationID)
	reporting.Command(config.CommandTypeSlash, i.GuildID, i.Member.User.ID, i.Member.User.Username, call.ObjectID, call.CorrelationID, timeStarted)
}
