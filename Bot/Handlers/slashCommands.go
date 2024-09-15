package handlers

import (
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
			Description: "How to get Help & Support for Discgo Bot",
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
					Description: "The Voice model. (Can enter /tts-search ID or search term)",
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
}

type SlashCommand struct {
	Command    *discordgo.ApplicationCommand
	Handler    func(i *discordgo.InteractionCreate, correlationId string)
	Complexity int
}

var SlashCommands []*discordgo.ApplicationCommand

func InitSlashCommands() {
	// Handle adding Slash Commands (these are added below this function)
	for _, cmd := range slashCommands {
		SlashCommands = append(SlashCommands, cmd.Command)
	}

	config.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		SlashCommandHandler(s, i)
	})
}

func SlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionMessageComponent {
		return
	}

	if !config.ServiceSettings.ISDEV != cache.ActiveGuilds[i.GuildID].IsDev {
		return
	}

	cmdName := i.ApplicationCommandData().Name

	for _, cmd := range slashCommands {
		if cmd.Command.Name == cmdName {
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
	}

	logger.ErrorText(i.GuildID, "No Handler found for Slash Command: %v", cmdName)
}

func RefreshSlashCommands(guildId string) {
	if SlashCommands == nil {
		InitSlashCommands()
	}

	// Fetch all registered commands assigned from the Bot -> the Guild
	registeredCommands, err := config.Session.ApplicationCommands(config.Session.State.User.ID, guildId)
	if err != nil {
		logger.Error(guildId, err)
		return
	}

	// Delete each command
	for _, cmd := range registeredCommands {
		err := config.Session.ApplicationCommandDelete(config.Session.State.User.ID, guildId, cmd.ID)
		if err != nil {
			logger.ErrorText(guildId, "Error deleting Command: %v, Error: %v", cmd.Name, err)
		}
	}

	// Register the Commands again
	for _, cmd := range SlashCommands {
		_, err := config.Session.ApplicationCommandCreate(config.Session.State.User.ID, guildId, cmd)
		if err != nil {
			logger.ErrorText(guildId, "Error registering Command: %v, Error: %v", cmd.Name, err)
		} else {
			logger.Debug(guildId, "Successfully registered command: %v", cmd.Name)
		}
	}
}
