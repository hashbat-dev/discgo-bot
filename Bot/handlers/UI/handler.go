package UI

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dabi-ngin/discgo-bot/Bot/audit"
)

var (
	commands = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"button": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "This is it. You've reached your destination. Your choice was monkey.\n" +
						"If you want to know more, check out the links below",
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									Emoji: discordgo.ComponentEmoji{
										Name: "📜",
									},
									Label: "Documentation",
									Style: discordgo.LinkButton,
									URL:   "https://discord.com/developers/docs/interactions/message-components#select-menus",
								},
								discordgo.Button{
									Emoji: discordgo.ComponentEmoji{
										Name: "🔧",
									},
									Label: "Discord developers",
									Style: discordgo.LinkButton,
									URL:   "https://discord.gg/discord-developers",
								},
								discordgo.Button{
									Emoji: discordgo.ComponentEmoji{
										Name: "🦫",
									},
									Label: "Discord Gophers",
									Style: discordgo.LinkButton,
									URL:   "https://discord.gg/7RuRrVHyXF",
								},
							},
						},
					},

					Flags: discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				audit.Error(err)
			}
		},
	}
)

func HandleMessage(session *discordgo.Session) {
	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := commands[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		}
	})

	//session.AddHandler(createInteraction)

	// var i *discordgo.InteractionCreate

	// var resp discordgo.InteractionResponse

	// resp.Type = discordgo.InteractionResponseChannelMessageWithSource

	// var respdata discordgo.InteractionResponseData

	// respdata.Content = "Test response"

	// respdata.Components = []discordgo.MessageComponent{
	// 	discordgo.Button{
	// 		Label: "testbutton",
	// 	},
	// }

	// resp.Data = &respdata

	// err := session.InteractionRespond(i.Interaction, &resp)

	// if err != nil {
	// 	audit.Error(err)
	// }

}

// func createInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
// 	resperr := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
// 		Type: discordgo.InteractionResponseChannelMessageWithSource,
// 		Data: &discordgo.InteractionResponseData{
// 			Content: "She took the kids!!!",
// 			Flags:   discordgo.MessageFlagsEphemeral,
// 			Components: []discordgo.MessageComponent{
// 				discordgo.Button{
// 					Emoji: discordgo.ComponentEmoji{
// 						Name: "🦫",
// 					},
// 					Label: ":fire: My Kids!!",
// 					Style: discordgo.DangerButton,
// 					URL:   "https://www.mills-reeve.com/services/family-and-children/divorce",
// 				},
// 			},
// 		},
// 	})
// 	if resperr != nil {
// 		audit.Error(resperr)
// 	}
// }
