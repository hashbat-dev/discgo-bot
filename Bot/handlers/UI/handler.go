package UI

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dabi-ngin/discgo-bot/Bot/audit"
)

var (
	componentHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"fd_yes": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		"fd_no": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Huh. I see, maybe some of these resources might help you?",
					Flags:   discordgo.MessageFlagsEphemeral,
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									Emoji: discordgo.ComponentEmoji{
										Name: "📜",
									},
									Label: "Documentation",
									Style: discordgo.LinkButton,
									URL:   "https://discord.com/developers/docs/interactions/message-components#buttons",
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
				},
			})
			if err != nil {
				audit.Error(err)
			}
		},
	}
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"button": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Select a button option :)",
					Flags:   discordgo.MessageFlagsEphemeral,
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									Emoji: discordgo.ComponentEmoji{
										Name: "✅",
									},
									Label:    "Yes",
									Style:    discordgo.SuccessButton,
									Disabled: false,
									CustomID: "fd_yes",
								},
								discordgo.Button{
									Emoji: discordgo.ComponentEmoji{
										Name: "❌",
									},
									Label:    "No",
									Style:    discordgo.DangerButton,
									Disabled: false,
									CustomID: "fd_no",
								},
								discordgo.Button{
									Label:    "I don't know",
									Style:    discordgo.LinkButton,
									Disabled: false,
									URL:      "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
									Emoji: discordgo.ComponentEmoji{
										Name: "🤷‍♂️",
									},
								},
							},
						},
					},
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
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case discordgo.InteractionMessageComponent:
			if h, ok := componentHandlers[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			}
		}
	})
}
