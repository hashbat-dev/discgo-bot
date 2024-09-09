package handlers

import (
	"time"

	"github.com/bwmarrin/discordgo"
	slash "github.com/dabi-ngin/discgo-bot/Bot/Commands/Slash"
	config "github.com/dabi-ngin/discgo-bot/Config"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
	reporting "github.com/dabi-ngin/discgo-bot/Reporting"
	"github.com/google/uuid"
)

// ===[Add Slash Commands]=============================================
var slashCommands = []SlashCommand{
	{
		Command: &discordgo.ApplicationCommand{
			Name:        "support",
			Description: "How to get Help & Support for Discgo Bot",
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate, correlationId string) {
			slash.SupportInfo(s, i, correlationId)
		},
	},
}

type SlashCommand struct {
	Command *discordgo.ApplicationCommand
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate, correlationId string)
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
	cmdName := i.ApplicationCommandData().Name

	for _, cmd := range slashCommands {
		if cmd.Command.Name == cmdName {
			correlationId := uuid.New()
			timeStarted := time.Now()

			logger.Info(i.GuildID, "Processing slash command [%v] %v", cmdName, correlationId.String())
			cmd.Handler(s, i, correlationId.String())

			reporting.Command(config.CommandTypeSlash, i.GuildID, i.Member.User.ID, i.Member.User.Username, cmdName, correlationId, timeStarted)
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
