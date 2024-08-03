package handlers

import (
	"errors"
	"strings"

	"github.com/ZestHusky/femboy-control/Bot/audit"
	"github.com/ZestHusky/femboy-control/Bot/commands"
	"github.com/ZestHusky/femboy-control/Bot/handlers/meme"
	"github.com/bwmarrin/discordgo"
)

func AddInteractions(session *discordgo.Session) {
	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {

		if i.Type == discordgo.InteractionMessageComponent {
			// Handle button interactions
			customID := i.MessageComponentData().CustomID

			passedAction := ""
			var splitVal []string
			if strings.Contains(customID, "-") {
				splitVal = strings.Split(customID, "-")
				passedAction = splitVal[0]
			} else {
				passedAction = customID
			}

			switch passedAction {
			case "regenerate":
				meme.HandleRegenerate(i, splitVal[1], splitVal[2], splitVal[3], splitVal[4])
			case "delete":
				meme.HandleDelete(i, splitVal[1], splitVal[2])
			default:
				audit.Error(errors.New("unknown command"))
			}

		} else if h, ok := commands.CommandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}
