package handlers

import (
	"fmt"
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
	structs "github.com/dabi-ngin/discgo-bot/Structs"
)

var chBang chan *bangChannelMessage = make(chan *bangChannelMessage)
var chPhrase chan *DispatchInfo = make(chan *DispatchInfo)

// init spins up some workers which will process messages passed into the
// chBang and chPhrase channels
func init() {
	for i := 0; i < 20; i++ {
		go triggerCommandWorker(i, chPhrase)
	}
	for i := 0; i < 3; i++ {
		go bangCommandWorker(i, chBang)
	}
}

// HandleNewMessage checks for Bot actions whenever a new Message is posted in a Server
func HandleNewMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	// 1. Do we want to skip this message?
	if SkipMessageCheck(session, message) {
		return
	}

	if command, ok := bangCommands.CommandTable[]
	if !helpers.DoesUserHavePermissionToUseCommand(message) {
		return
	}

	// 2. Decode the message to determine how to handle it
	// 	  => Do we have an bang !command?
	bangCommand := helpers.CheckForBangCommand(message.Content)

	// 3. Determine permissions of the sending user
	if bangCommand != "" {
		go DispatchBangCommand(message, bangCommand)
		return
	}

	// => If not, do we have a trigger phrase command?
	triggerPhrase := ""
	if bangCommand == "" {
		triggerPhrase = triggerCommands.CheckForTriggerPhrase(message.Content)
	}
	if triggerPhrase != "" {
		go DispatchTriggerCommand(message, triggerPhrase)
		return
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
	Message     *discordgo.MessageCreate
	CommandName string
}

type bangChannelMessage struct {
	Message     *discordgo.MessageCreate
	CommandName string
	Command     structs.BangCommand
}

// DispatchBangCommand sends !commands to the relevant handler
func DispatchBangCommand(message *discordgo.MessageCreate, commandName string) bool {
	// Setup the Command
	commandName = strings.ToLower(commandName)
	logger.Event(message.GuildID, "User: [%v] has requested [!%v]", message.Author.Username, commandName)
	commandType := config.CommandTypeBang
	// Check for a command
	command, ok := bangCommands.CommandTable[commandName]
	if !ok {
		errMsg := fmt.Sprintf("User '%s' tried to use unknown command '!%s'", message.Author.Username, commandName)
		logger.Info(message.GuildID, errMsg)
		return false
	}
	chBang <- &bangChannelMessage{message, commandName, command}

	database.LogCommandUsage(message.GuildID, message.Author.ID, commandType, commandName)
	// cache.AddToCommandCache(commandType, commandName, message.GuildID, message.Author.ID, message.Author.Username, timeStart, timeFinish)
	return true
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

// bangCommandWorker is intended to be used as a goroutine which loops over the messages in a channel
// of type *DispatchAction, reading out the message and actioning any functions necessary.
func bangCommandWorker(id int, ch <-chan *bangChannelMessage) {
	for {
		select {
		case bangChanMessage, ok := <-ch:
			if !ok {
				// Channel is closed, exit goroutine
				logger.Info(bangChanMessage.Message.GuildID, "BangCommandWorker %d: Channel closed, exiting...\n", id)
				return
			}
			fmt.Printf("BangCommandWorker %d: Processing command '%s'\n", id, bangChanMessage.CommandName)
			err := bangChanMessage.Command.Begin(bangChanMessage.Message, bangChanMessage.Command)
			if err != nil {
				logger.Error(bangChanMessage.Message.GuildID, err)
			}
		}
	}
}

// triggerCommandWorker is intended to be used as a goroutine which loops over the messages in a channel
// of type *DispatchAction, reading out the message and actioning any functions necessary.
func triggerCommandWorker(id int, ch <-chan *DispatchInfo) {
	for {
		select {
		case command, ok := <-ch:
			if !ok {
				logger.Info(command.CommandName, "TriggerCommandWorker %d: Channel closing...", id)
				return
			}
			logger.Info(command.Message.GuildID, "TriggerCommandWorker %d: Processing command %s", id, command.CommandName)
			err := triggerCommands.RunTriggerCommand(command.CommandName, command.Message)
			if err != nil {
				logger.Error(command.Message.GuildID, err)
			}
		}
	}
}
