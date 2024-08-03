package toothjack

import (
	"math/rand"

	"github.com/bwmarrin/discordgo"
	"github.com/dabi-ngin/discgo-bot/Bot/audit"
	"github.com/dabi-ngin/discgo-bot/Bot/config"
	dbhelper "github.com/dabi-ngin/discgo-bot/Bot/dbhelpers"
	"github.com/dabi-ngin/discgo-bot/Bot/gifbank"
	"github.com/dabi-ngin/discgo-bot/Bot/helpers"
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
