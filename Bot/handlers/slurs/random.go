package slurs

import (
	"math/rand"

	"github.com/ZestHusky/femboy-control/Bot/audit"
	"github.com/ZestHusky/femboy-control/Bot/config"
	dbhelper "github.com/ZestHusky/femboy-control/Bot/dbhelpers"
	"github.com/ZestHusky/femboy-control/Bot/helpers"
	"github.com/bwmarrin/discordgo"
)

func RandomSlur(message *discordgo.MessageCreate) bool {
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
		nBomb, err := GetRandomSlur()
		if err != nil {
			audit.Error(err)
			return false
		}

		_, err = config.Session.ChannelMessageSendReply(message.ChannelID, nBomb, message.Reference())
		if err != nil {
			audit.Error(err)
			return false
		}
		audit.Log("[!!!] Racial Slur sent to User: " + helpers.GetNicknameFromID(message.GuildID, message.Author.ID))
		dbhelper.CountCommand("randomslur", message.Message.Author.ID)
	} else if random == chance-10 {
		applied = true
		nBomb, err := GetRandomJobSlur()
		if err != nil {
			audit.Error(err)
			return false
		}

		_, err = config.Session.ChannelMessageSendReply(message.ChannelID, nBomb, message.Reference())
		if err != nil {
			audit.Error(err)
			return false
		}

		audit.Log("[!!!] Racial Job Slur sent to User: " + helpers.GetNicknameFromID(message.GuildID, message.Author.ID))
		dbhelper.CountCommand("randomjob", message.Message.Author.ID)
	}

	return applied
}
