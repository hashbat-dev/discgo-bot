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

// HandleNewMessage checks for Bot actions whenever a new Message is posted in a Server
func HandleNewMessage(session *discordgo.Session, message *discordgo.MessageCreate) {

	// 1. Do we want to skip this message?
	if SkipMessageCheck(session, message) {
		return
	}

	// 2. Decode the message to determine how to handle it
	// 	  => Do we have an exclamation !command?
	exclamationCommand := CheckForExclamationCommand(message.Content)

	// 3. Determine permissions of the sending user
	if exclamationCommand != "" {
		if !helpers.DoesUserHavePermissionToUseCommand(message) {
			return
		}
	}

	// => If not, do we have a trigger phrase command?
	triggerPhrase := ""
	if exclamationCommand == "" {
		triggerPhrase = CheckForTriggerPhrase(message.Content)
	}

	// 4. Send the message to the relevant handler
	if exclamationCommand != "" {
		DispatchExclamationCommand(message, exclamationCommand)
	}

	if triggerPhrase != "" {
		DispatchTriggerCommand(message, triggerPhrase)
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

// Checks for, and returns if exists a !command
func CheckForExclamationCommand(messageContent string) string {

	if len(messageContent) == 0 {
		return ""
	}

	if string([]rune(messageContent)[0]) == "!" {
		spaceIndex := strings.Index(messageContent, " ")
		if spaceIndex == -1 {
			// No spaces in the Content, we assume the whole message is the ! command
			return messageContent[1:]
		} else {
			return strings.Split(messageContent, " ")[0]
		}
	}
	return ""
}

// CheckForTriggerPhrase Determines whether a string contains a trigger phrase for bot action
func CheckForTriggerPhrase(trigger string) string {
	return triggerCommands.CheckForTriggerPhrase(trigger)
}

// DispatchExclamationCommand sends !commands to the relevant handler
func DispatchExclamationCommand(message *discordgo.MessageCreate, command string) {

	// Setup the Command
	logger.Event(message.GuildID, "User: [%v] has requested [!%v]", message.Author.Username, command)
	commandType := config.CommandTypeBang
	timeStart := time.Now()
	foundCommand := bangCommands.GetCommand(command)

	// Execute the Command if found
	if foundCommand != nil {
		err := bangCommands.RunCommand(command, message)
		if err != nil {
			logger.Error(message.GuildID, err)
			return
		}
	} else {
		logger.Info(message.GuildID, "User [%s] tried to use unknown command [!%s]", message.Author.Username, command)
		return
	}

	// Log the Command
	timeFinish := time.Now()
	database.LogCommandUsage(message.GuildID, message.Author.ID, commandType, command)
	cache.AddToCommandCache(commandType, command, message.GuildID, message.Author.ID, message.Author.Username, timeStart, timeFinish)
}

// DispatchTriggerCommand sends trigger commands to the relevant handler
func DispatchTriggerCommand(message *discordgo.MessageCreate, command string) {
	logger.Event(message.GuildID, "User: [%v] has triggered [%v]", message.Author.Username, command)
	timeStart := time.Now()
	commandType := config.CommandTypePhrase

	err := triggerCommands.RunTriggerCommand(command, message)
	if err != nil {
		logger.Error(message.GuildID, err)
	} else {
		database.LogCommandUsage(message.GuildID, message.Author.ID, 2, command)
	}

	timeFinish := time.Now()
	database.LogCommandUsage(message.GuildID, message.Author.ID, commandType, command)
	cache.AddToCommandCache(commandType, command, message.GuildID, message.Author.ID, message.Author.Username, timeStart, timeFinish)
}
