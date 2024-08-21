package handlers

import (
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	bangCommands "github.com/dabi-ngin/discgo-bot/Bot/BangCommands"
	triggerCommands "github.com/dabi-ngin/discgo-bot/Bot/TriggerCommands"
	cache "github.com/dabi-ngin/discgo-bot/Cache"
	config "github.com/dabi-ngin/discgo-bot/Config"
	database "github.com/dabi-ngin/discgo-bot/Database"
	helpers "github.com/dabi-ngin/discgo-bot/Helpers"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
)

var chBang chan *DispatchInfo = make(chan *DispatchInfo)
var chPhrase chan *DispatchInfo = make(chan *DispatchInfo)

// HandleNewMessage checks for Bot actions whenever a new Message is posted in a Server
func HandleNewMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	// 1. Do we want to skip this message?
	if SkipMessageCheck(session, message) {
		return
	}

	// 2. Decode the message to determine how to handle it
	// 	  => Do we have an bang !command?
	bangCommand := helpers.CheckForBangCommand(message.Content)

	// 3. Determine permissions of the sending user
	if bangCommand != "" {
		if !helpers.DoesUserHavePermissionToUseCommand(message) {
			return
		}
	}

	// => If not, do we have a trigger phrase command?
	triggerPhrase := ""
	if bangCommand == "" {
		triggerPhrase = triggerCommands.CheckForTriggerPhrase(message.Content)
	}

	// 4. Create Channels

	go func() {
		for cmd := range chBang {
			DispatchBangCommand(cmd.Message, cmd.Command)
		}
	}()

	go func() {
		for cmd := range chPhrase {
			DispatchTriggerCommand(cmd.Message, cmd.Command)
		}
	}()

	// 5. Add command to the relevant Channel
	if bangCommand != "" {
		chBang <- &DispatchInfo{Message: message, Command: bangCommand}
	}

	if triggerPhrase != "" {
		chPhrase <- &DispatchInfo{Message: message, Command: triggerPhrase}
	}

}

// Determines whether we should ignore the inbound Message
func SkipMessageCheck(session *discordgo.Session, message *discordgo.MessageCreate) bool {
	if message.Author == nil {
		return true
	}

	if message.Author.ID == session.State.User.ID {
		return true
	}

	return false
}

// CheckForTriggerPhrase Determines whether a string contains a trigger phrase for bot action
func CheckForTriggerPhrase(trigger string) string {
	return triggerCommands.CheckForTriggerPhrase(trigger)
}

type DispatchInfo struct {
	Message *discordgo.MessageCreate
	Command string
}

// DispatchBangCommand sends !commands to the relevant handler
func DispatchBangCommand(message *discordgo.MessageCreate, command string) bool {

	// Setup the Command
	command = strings.ToLower(command)
	logger.Event(message.GuildID, "User: [%v] has requested [!%v]", message.Author.Username, command)
	commandType := config.CommandTypeBang
	timeStart := time.Now()

	// Check for a command
	if foundCommand, ok := bangCommands.CommandTable[command]; ok {
		err := foundCommand.Execute(message, foundCommand)
		if err != nil {
			// Error during Processing - The error logging / reporting to users is done within the functions to ensure
			// we can deliver relevant error messages where needed.
			return false
		}

		// Log the Command
		timeFinish := time.Now()
		database.LogCommandUsage(message.GuildID, message.Author.ID, commandType, command)
		cache.AddToCommandCache(commandType, command, message.GuildID, message.Author.ID, message.Author.Username, timeStart, timeFinish)

		return true
	} else {
		logger.Info(message.GuildID, "User [%s] tried to use unknown command [!%s]", message.Author.Username, command)
		return false
	}
}

// DispatchTriggerCommand sends trigger commands to the relevant handler
func DispatchTriggerCommand(message *discordgo.MessageCreate, command string) bool {
	logger.Event(message.GuildID, "User: [%v] has triggered [%v]", message.Author.Username, command)
	timeStart := time.Now()
	commandType := config.CommandTypePhrase

	err := triggerCommands.RunTriggerCommand(command, message)
	if err != nil {
		logger.Error(message.GuildID, err)
		return false
	} else {
		database.LogCommandUsage(message.GuildID, message.Author.ID, 2, command)
	}

	timeFinish := time.Now()
	database.LogCommandUsage(message.GuildID, message.Author.ID, commandType, command)
	cache.AddToCommandCache(commandType, command, message.GuildID, message.Author.ID, message.Author.Username, timeStart, timeFinish)
	return true
}
