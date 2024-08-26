package handlers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	commands "github.com/dabi-ngin/discgo-bot/Bot/Commands"
	triggers "github.com/dabi-ngin/discgo-bot/Bot/Commands/Triggers"
	cache "github.com/dabi-ngin/discgo-bot/Cache"
	config "github.com/dabi-ngin/discgo-bot/Config"
	database "github.com/dabi-ngin/discgo-bot/Database"
	helpers "github.com/dabi-ngin/discgo-bot/Helpers"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
	"github.com/google/uuid"
)

// HandleNewMessage checks for Bot actions whenever a new Message is posted in a Server
func HandleNewMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	// 1. Do we want to skip this message?
	if skipMessageCheck(session, message) {
		return
	}
	// 2. Generate a correlationID for the request
	correlationId, corrErr := uuid.NewUUID()
	if corrErr != nil {
		logger.Debug(message.GuildID, "failed to generate correlation-id for request. err=%s", corrErr.Error())
	}

	// 3. Extract the command name
	commandName := extractCommandName(message.Content)
	if commandName != "" {
		// 4. Retrieve command using the name
		command := getCommandByName(commandName)
		if command != nil {
			// 5. TODO - check user permissions
			task := &CommandTask{
				Message:       message,
				Command:       command,
				CorrelationId: correlationId,
			}
			dispatchTask(task)
		} else {
			logger.Debug(message.GuildID, "invalid message command attempt :: could not retrieve '%s' from jump table :: correlation-id :: %v", commandName, correlationId)
		}
	}

	// 6. Check for and Process Triggers
	checkForAndProcessTriggers(message)
}

// Determines whether we should ignore the inbound Message
func skipMessageCheck(session *discordgo.Session, message *discordgo.MessageCreate) bool {
	if message.Author == nil {
		return true
	}

	if message.Author.ID == session.State.User.ID {
		return true
	}

	return false
}

func extractCommandName(input string) string {
	if strings.HasPrefix(input, "!") {
		parts := strings.SplitN(input[1:], " ", 2)
		return strings.ToLower(parts[0])
	}
	return ""
}

func getCommandByName(commandName string) commands.Command {
	cmd, ok := commands.JumpTable[commandName]
	if !ok {
		return nil
	}
	return cmd
}

// Dispatches a Command to its appropriate channel.
func dispatchTask(task *CommandTask) {
	// TODO - add touch point to pass off queue info to the dashboard
	switch task.Command.Complexity() {
	case config.TRIVIAL_TASK:
		TRIVIAL_TASKS <- task
	case config.CPU_BOUND_TASK:
		CPU_TASKS <- task
	case config.IO_BOUND_TASK:
		IO_TASKS <- task
	default:
		TRIVIAL_TASKS <- task
	}
}

func checkForAndProcessTriggers(message *discordgo.MessageCreate) {
	var matchedPhrases []triggers.Phrase
	guildIndex := cache.GetGuildIndex(message.GuildID)
	if guildIndex < 0 {
		return
	}
	for _, trigger := range cache.ActiveGuilds[guildIndex].Triggers {
		var regexString string
		if trigger.WordOnlyMatch {
			regexString = `(?i)\b%s\b`
		} else {
			regexString = `(?i)%s`
		}

		check := regexp.MustCompile(fmt.Sprintf(regexString, regexp.QuoteMeta(trigger.Phrase)))
		if check.MatchString(message.Content) {
			matchedPhrases = append(matchedPhrases, trigger)
		}
	}

	// Process any matching Triggers
	var notifyPhrases []string
	for _, phrase := range matchedPhrases {
		database.LogCommandUsage(message.GuildID, message.Author.ID, config.CommandTypePhrase, phrase.Phrase)
		if phrase.NotifyOnDetection {
			notifyPhrases = append(notifyPhrases, phrase.Phrase)
		}
	}

	if len(notifyPhrases) > 0 {
		showText := strings.ToUpper(helpers.ConcatStringWithAnd(notifyPhrases)) + " MENTIONED"
		_, err := config.Session.ChannelMessageSend(message.ChannelID, showText)
		if err != nil {
			logger.Error(message.GuildID, err)
		}
	}
}
