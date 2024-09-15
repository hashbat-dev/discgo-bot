package handlers

import (
	"fmt"
	"sort"

	"github.com/bwmarrin/discordgo"
	slash "github.com/dabi-ngin/discgo-bot/Bot/Commands/Slash"
	cache "github.com/dabi-ngin/discgo-bot/Cache"
	config "github.com/dabi-ngin/discgo-bot/Config"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
)

// ===[Add Slash Commands]=============================================
var slashCommands = []SlashCommand{
	//	/support
	{
		Command: &discordgo.ApplicationCommand{
			Name:        "support",
			Description: "Get Support for Discgo Bot and its features",
		},
		Handler: func(i *discordgo.InteractionCreate, correlationId string) {
			slash.SupportInfo(i, correlationId)
		},
		Complexity: config.TRIVIAL_TASK,
	},
	//	/tts-play
	{
		Command: &discordgo.ApplicationCommand{
			Name:        "tts-play",
			Description: "Convert Text to Speech through one of FakeYou.com's thousands of Voice Models",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "voice",
					Description: "The Voice model to use, you can choose from a list of results once submitted.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "text",
					Description: "The Text to convert to speech.",
					Required:    true,
				},
			},
		},
		Handler: func(i *discordgo.InteractionCreate, correlationId string) {
			slash.TtsPlay(i, correlationId)
		},
		Complexity: config.IO_BOUND_TASK,
	},
	//	/help
	{
		Command: &discordgo.ApplicationCommand{
			Name:        "help",
			Description: "View all the Bot's commands and their descriptions",
		},
		Handler: func(i *discordgo.InteractionCreate, correlationId string) {
			slash.SendHelp(i, correlationId)
		},
		Complexity: config.TRIVIAL_TASK,
	},
}

type SlashCommand struct {
	Command    *discordgo.ApplicationCommand
	Handler    func(i *discordgo.InteractionCreate, correlationId string)
	Complexity int
}

var SlashCommands map[string]SlashCommand = make(map[string]SlashCommand)

func InitSlashCommands() {
	// Handle adding Slash Commands (these are added below this function)
	for _, cmd := range slashCommands {
		SlashCommands[cmd.Command.Name] = cmd
	}

	writeSlashHelpText()

	config.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		SlashCommandHandler(s, i)
	})
}

func SlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionMessageComponent {
		return
	}

	if config.ServiceSettings.ISDEV != cache.ActiveGuilds[i.GuildID].IsDev {
		return
	}

	cmdName := i.ApplicationCommandData().Name
	if cmd, exists := SlashCommands[i.ApplicationCommandData().Name]; exists {
		DispatchTask(&Task{
			CommandType: config.CommandTypeSlash,
			Complexity:  cmd.Complexity,
			SlashDetails: &SlashTaskDetails{
				Interaction:  i,
				SlashCommand: cmd,
			},
		})
		return
	}

	logger.ErrorText(i.GuildID, "No Handler found for Slash Command: %v", cmdName)
}

func RefreshSlashCommands(guildId string) {
	if len(SlashCommands) == 0 {
		InitSlashCommands()
	}

	// Fetch all registered commands assigned from the Bot -> the Guild
	registeredCommands, err := config.Session.ApplicationCommands(config.Session.State.User.ID, guildId)
	if err != nil {
		logger.Error(guildId, err)
		return
	}

	// Check Existing Commands
	var validated map[string]interface{} = make(map[string]interface{})
	for _, cmd := range registeredCommands {
		delete := false
		if local, exists := SlashCommands[cmd.Name]; exists {
			// Slash Command already registered, are the options the same?
			var liveOpts map[string]interface{} = make(map[string]interface{})
			for _, opt := range cmd.Options {
				liveOpts[fmt.Sprintf("%s|%s", opt.Name, opt.Type.String())] = struct{}{}
			}
			var currOpts map[string]interface{} = make(map[string]interface{})
			for _, opt := range local.Command.Options {
				currOpts[fmt.Sprintf("%s|%s", opt.Name, opt.Type.String())] = struct{}{}
			}

			// Different option lengths?
			if len(liveOpts) != len(currOpts) {
				delete = true
			}

			// Matching options?
			if !delete {
				for item := range liveOpts {
					if _, found := currOpts[item]; !found {
						delete = true
					}
				}
			}

			if !delete {
				validated[cmd.Name] = struct{}{}
			}
		} else {
			// Slash Command exists externally but not in our map, delete it
			delete = true
		}

		if delete {
			err := config.Session.ApplicationCommandDelete(config.Session.State.User.ID, guildId, cmd.ID)
			if err != nil {
				logger.ErrorText(guildId, "Error deleting Command: %s, Error: %e", cmd.Name, err)
			}
		} else {
			validated[cmd.Name] = struct{}{}
		}
	}

	// Register the Commands again
	for _, cmd := range SlashCommands {
		if _, validated := validated[cmd.Command.Name]; validated {
			// Already registered
			logger.Debug(guildId, "Command already registered: %v", cmd.Command.Name)
			continue
		}

		_, err := config.Session.ApplicationCommandCreate(config.Session.State.User.ID, guildId, cmd.Command)
		if err != nil {
			logger.ErrorText(guildId, "Error registering Command: %v, Error: %v", cmd.Command.Name, err)
		} else {
			logger.Debug(guildId, "Successfully registered command: %v", cmd.Command.Name)
		}
	}
}

func writeSlashHelpText() {
	// 1. Sort the slice Alphabetically
	slashCmds := slashCommands
	sort.Slice(slashCmds, func(i, j int) bool {
		return slashCmds[i].Command.Name < slashCmds[j].Command.Name
	})

	// 2. Generate the Help Text
	text := ""
	for i, cmd := range slashCmds {
		cmdText := "* **/" + cmd.Command.Name + "**: " + cmd.Command.Description
		if len(cmd.Command.Options) > 0 {
			cmdText += " It accepts these parameters:"
		}
		for _, opt := range cmd.Command.Options {
			cmdText += "\n * **" + opt.Name + "**: " + opt.Description
		}
		if i > 0 {
			text += "\n"
		}
		text += cmdText
	}

	// 3. Set it in the Config
	config.UserSlashHelpText = text
}
