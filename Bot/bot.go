package bot

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/dabi-ngin/discgo-bot/Bot/audit"
	"github.com/dabi-ngin/discgo-bot/Bot/commands"
	"github.com/dabi-ngin/discgo-bot/Bot/config"
	"github.com/dabi-ngin/discgo-bot/Bot/constants"
	dbhelpers "github.com/dabi-ngin/discgo-bot/Bot/dbhelpers"
	"github.com/dabi-ngin/discgo-bot/Bot/handlers"
	"github.com/dabi-ngin/discgo-bot/Bot/helpers"
)

var session *discordgo.Session

func Run() {

	// Setup our constant values first
	loadConfig()

	// Initialise the Discord session
	var err error
	session, err = discordgo.New("Bot " + config.BotToken)
	if err != nil {
		audit.Error(err)
		return
	} else if session == nil {
		audit.Error(errors.New("discord session object nil"))
		return
	}

	// Add Message event handlers to the Bot Session
	session.AddHandler(handlers.NewMessageHandler) // New message

	// Open the Discord Bot session
	err = session.Open()
	if err != nil {
		audit.Error(err)
		return
	} else {
		audit.Log("Discord Session opened successfully")
	}
	dbhelpers.LoadDB()

	config.Session = session

	// DELETE A /COMMAND HERE
	// session.ApplicationCommandDelete(session.State.Application.ID, session.State.Application.GuildID, "1236369287787057152")

	// Add the Bot's /slash commands
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands.Commands))
	commandCount := len(commands.Commands)
	completeCount := 0
	var commandWg sync.WaitGroup
	for i, v := range commands.Commands {
		commandWg.Add(1)
		go func(index int, cmd *discordgo.ApplicationCommand) {
			cmd, err := session.ApplicationCommandCreate(session.State.User.ID, session.State.Application.GuildID, cmd)
			if err != nil {
				audit.ErrorWithText("Command: "+cmd.Name, err)
			} else {
				completeCount++
				audit.Log("Added Command: " + cmd.Name + " (" + fmt.Sprint(completeCount) + "/" + fmt.Sprint(commandCount) + " added)")
				registeredCommands[index] = cmd

				if commandCount == completeCount {
					audit.Log("All Commands Added ===========================")
				}

			}
		}(i, v)
	}

	// Add the Handlers for the /slash commands
	handlers.AddInteractions(session)

	// Initialise the Database connection, cached for later use
	VersionCheck(session)

	config.Session = session
	// Keep the bot running until an OS Shutdown/Exit
	audit.Log("Setup complete, bot running on: " + config.ServerName)

	// Report the bot as running to the Server
	if !config.IsDev && !config.ReportedRunning {
		BotStartMsg := "Bot started on " + config.ServerName
		_, err := session.ChannelMessageSend(constants.CHANNEL_BOT_TESTING, BotStartMsg)
		if err != nil {
			audit.Error(err)
		}
		config.ReportedRunning = true
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// When shutting down, un-register the commands. If this isn't done
	// then commands need to be deleted manually by calling the DeleteCommand
	// function via Code, using the Command's ID. This can be found by going
	// to the Server Settings > Integrations > Webhooks
	for _, v := range registeredCommands {
		err := session.ApplicationCommandDelete(session.State.User.ID, session.State.Application.GuildID, v.ID)
		if err != nil {
			audit.Error(err)
		}
	}
}

func loadConfig() {
	config.SetVars()

	// Are we in Dev?
	hostname, err := os.Hostname()
	if err != nil {
		config.ServerName = "Unknown"
		audit.Error(err)
	} else {
		config.ServerName = hostname
	}

	// Make sure the Temp Directory exists on Startup
	if !helpers.CheckDirectoryExists("tmp") {
		audit.Error(errors.New("could not create temp directory"))
	}
}
