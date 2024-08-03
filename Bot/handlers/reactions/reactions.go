package reactions

import (
	"strings"

	"github.com/ZestHusky/femboy-control/Bot/audit"
	"github.com/ZestHusky/femboy-control/Bot/config"
	"github.com/ZestHusky/femboy-control/Bot/constants"
	"github.com/ZestHusky/femboy-control/Bot/friday"
	"github.com/bwmarrin/discordgo"
)

func CheckMessage(message *discordgo.MessageCreate) {
	var icon []string

	if strings.Contains(strings.ToLower(message.Content), "balls") {
		icon = append(icon, "pognuts:1209278172626026496")
	}
	if strings.Contains(strings.ToLower(message.Content), "pug") {
		icon = append(icon, "pugface:1224510500260417637")
	}
	if strings.Contains(strings.ToLower(message.Content), "milk") {
		icon = append(icon, "slurp:1209269997487259668")
	}
	if strings.Contains(strings.ToLower(message.Content), "one piece") {
		icon = append(icon, "luffy:1235700349327900702")
	}
	if friday.IsItFwiday() && message.Author.ID == constants.USER_ID_POG {
		icon = append(icon, "nicedick:1235701740926402600")
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
