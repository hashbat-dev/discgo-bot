package bot

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	handlers "github.com/dabi-ngin/discgo-bot/Bot/Handlers"
	config "github.com/dabi-ngin/discgo-bot/Config"
	database "github.com/dabi-ngin/discgo-bot/Database"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
)

func Init() {
	// 1. Database Init
	if !database.Init() {
		logger.Error("", errors.New("Failed to initialise database"))
		return
	}

	// 2. Discord Session Init
	if !sessionInit() {
		logger.Error("", errors.New("Failed to initialise session"))
		return
	}

	// 3. Add Handlers to the Session
	if !addHandlers() {
		logger.Error("", errors.New("Failed to add handlers"))
		return
	}

	// 4. Open the Discord session
	if !sessionOpen() {
		logger.Error("", errors.New("Failed to open session"))
		return
	}

	// 5. Log Init
	if !logger.Init() {
		logger.Error("", errors.New("Failed to initialise logging"))
		return
	}

	var initSuffix string
	if config.ServiceSettings.DASHBOARDURL != "" {
		initSuffix = ", Dashboard open at: " + config.ServiceSettings.DASHBOARDURL
	}
	logger.Info("", "Bot intialisation successful%s", initSuffix)

	// 6. Register Discord /commands
	if !registerCommands() {
		logger.Error("", errors.New("Failed to register commands"))
		return
	}

	// 7. Reset Global Discord /commands
	handlers.RefreshSlashCommands("")

}

func sessionInit() bool {
	session, err := discordgo.New("Bot " + config.ServiceSettings.BOTTOKEN)
	if err != nil {
		logger.Error("", err)
		return false
	} else if session == nil {
		logger.Error("", err)
		return false
	}

	config.Session = session
	return true
}

func sessionOpen() bool {
	err := config.Session.Open()
	if err != nil {
		logger.Error("", err)
		return false
	}

	if config.Session == nil {
		logger.Error("FUCK", err)
	}
	return true
}

func addHandlers() bool {
	config.Session.AddHandler(handlers.HandleNewMessage)          //    New Messages
	config.Session.AddHandler(handlers.HandleNewGuild)            //	Server connected to the bot
	config.Session.AddHandler(handlers.HandleInteractionResponse) //	Responses from Interaction objects
	config.Session.AddHandler(handlers.HandleReactionAdd)         //	Message Reactions: Add
	config.Session.AddHandler(handlers.HandleReactionRemove)      //	Message Reactions: Remove
	return true
}

func registerCommands() bool {
	return true
}
