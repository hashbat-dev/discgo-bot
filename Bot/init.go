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
	// 1. Config Init
	if !config.Init() {
		logger.Error("", errors.New("Failed to load configs"))
		return
	}

	// 2. Database Init
	if !database.Init() {
		logger.Error("", errors.New("Failed to initialise database"))
		return
	}

	// 3. Discord Session Init
	if !sessionInit() {
		logger.Error("", errors.New("Failed to initialise session"))
		return
	}

	// 4. Add Handlers to the Session
	if !addHandlers() {
		logger.Error("", errors.New("Failed to add handlers"))
		return
	}

	// 5. Open the Discord session
	if !sessionOpen() {
		logger.Error("", errors.New("Failed to open session"))
		return
	}

	// 6. Log Init
	if !logger.Init() {
		logger.Error("", errors.New("Failed to initialise logging"))
		return
	}

	initText := "Bot initalisation successful"
	if config.DashboardUrl != "" {
		initText += ", Dashboard open at: " + config.DashboardUrl
	}
	logger.Info("", initText)

	// 7. Register Discord /commands
	if !registerCommands() {
		logger.Error("", errors.New("Failed to register commands"))
		return
	}

}

func sessionInit() bool {
	session, err := discordgo.New("Bot " + config.BotToken)
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

	config.Session.AddHandler(handlers.HandleNewMessage) // New Messages
	config.Session.AddHandler(handlers.HandleNewGuild)   //	Added to a new Server
	return true

}

func registerCommands() bool {
	return true
}
