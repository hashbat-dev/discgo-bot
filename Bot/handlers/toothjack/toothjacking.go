package toothjack

import (
	"math/rand"

	"github.com/ZestHusky/femboy-control/Bot/audit"
	"github.com/ZestHusky/femboy-control/Bot/config"
	dbhelper "github.com/ZestHusky/femboy-control/Bot/dbhelpers"
	"github.com/ZestHusky/femboy-control/Bot/gifbank"
	"github.com/ZestHusky/femboy-control/Bot/helpers"
	"github.com/bwmarrin/discordgo"
)

func RandToothjack(message *discordgo.MessageCreate, discord *discordgo.Session) bool {
	chance := 500
	if message.Author.ID == config.BullyTarget {
		chance = 200
	} else {
		chance = 500
	}

	random := rand.Intn(chance + 1)
	applied := false
	if random == chance/2 {
		applied = true
		gifbank.Post(discord, message, "tooth")
		audit.Log("[!!!] Toothjack sent to User: " + helpers.GetNicknameFromID(message.GuildID, message.Author.ID))
		dbhelper.CountCommand("randomtooth", message.Message.Author.ID)
	}

	return applied
}
