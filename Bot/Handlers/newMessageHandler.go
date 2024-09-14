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
	logger "github.com/dabi-ngin/discgo-bot/Logger"
	reporting "github.com/dabi-ngin/discgo-bot/Reporting"
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
			DispatchTask(&WorkerItem{
				CommandType: config.CommandTypeBang,
				Complexity:  command.Complexity(),
				BangCommand: BangCommandWorker{
					Message:       message,
					Command:       command,
					CorrelationId: correlationId,
				},
			})
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
func DispatchTask(task *WorkerItem) {
	reporting.WorkerQueued(task.Complexity)
	switch task.Complexity {
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
	for _, trigger := range cache.ActiveGuilds[message.GuildID].Triggers {
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

	// Dispatch any matching Triggers
	if len(matchedPhrases) > 0 {
		DispatchTask(&WorkerItem{
			CommandType: config.CommandTypePhrase,
			Complexity:  config.TRIVIAL_TASK,
			Phrases: PhraseWorker{
				Message:        message,
				TriggerPhrases: matchedPhrases,
			},
		})
	}

}
