package handlers

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	triggers "github.com/dabi-ngin/discgo-bot/Bot/Commands/Triggers"
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

	// 2. Did we find any Bang Commands?
	bangCommand := ExtractCommand(message.Content)
	if bangCommand != "" {
		go DispatchCommand(message, bangCommand)
	}

	// 3. Check for and Process Triggers
	CheckForAndProcessTriggers(message)
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

func ExtractCommand(input string) string {
	if strings.HasPrefix(input, "!") {
		parts := strings.SplitN(input[1:], " ", 2)
		return strings.ToLower(parts[0])
	}
	return ""
}

func CheckForAndProcessTriggers(message *discordgo.MessageCreate) {
	var matchedPhrases []triggers.Phrase
	guildIndex := cache.GetGuildIndex(message.GuildID)
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

// Dispatches a Command to its appropriate channel. If SpecifiedChannel is -1 the channel will be determined from the Command object
func DispatchCommand(message *discordgo.MessageCreate, bangCommand string) {
	// See if can get match the request to a Command
	var cmd Command
	if value, exists := Commands[bangCommand]; !exists {
		logger.Info(message.GuildID, "User [%v] called command [!%v] which did not exist", message.Author.ID, bangCommand)
		return
	} else {
		cmd = value
	}

	// Build the Request
	chRequest := &ChannelRequest{
		Message:     message,
		CommandName: bangCommand,
		Command:     cmd,
	}

	// Check the Channel exists
	poolIota := cmd.ProcessPool().ProcessPoolIota
	if poolIota > config.LastPoolIota {
		logger.Error(message.GuildID, fmt.Errorf("processPoolIota of %v does not have a channel to handle it", cmd.ProcessPool().ProcessPoolIota))
		return
	}

	// Dispatch the Command
	PoolQueue[poolIota]++
	PoolLastAdded[poolIota] = time.Now()
	Pools[poolIota] <- chRequest

}
