package reactions

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/dabi-ngin/discgo-bot/Bot/audit"
	"github.com/dabi-ngin/discgo-bot/Bot/config"
)

func CheckMessage(message *discordgo.MessageCreate) {
	var icon []string

	if strings.Contains(strings.ToLower(message.Content), "pug") {
		icon = append(icon, "pugface:1224510500260417637")
	}
	if strings.Contains(strings.ToLower(message.Content), "milk") {
		icon = append(icon, "slurp:1209269997487259668")
	}
	if strings.Contains(strings.ToLower(message.Content), "one piece") {
		icon = append(icon, "luffy:1235700349327900702")
	}

	if len(icon) > 0 {
		for _, i := range icon {
			err := config.Session.MessageReactionAdd(message.ChannelID, message.ID, i)
			if err != nil {
				audit.Error(err)
			}
		}

	}

}
