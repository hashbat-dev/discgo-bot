package handlers

import (
	"errors"
	"strings"

	"github.com/bwmarrin/discordgo"
	triggers "github.com/dabi-ngin/discgo-bot/Bot/Commands/Triggers"
	cache "github.com/dabi-ngin/discgo-bot/Cache"
	database "github.com/dabi-ngin/discgo-bot/Database"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
	reporting "github.com/dabi-ngin/discgo-bot/Reporting"
)

// Calls whenever a new Guild connects to the bot. This also runs for all active Guilds on startup.
func HandleNewGuild(session *discordgo.Session, newGuild *discordgo.GuildCreate) {
	// 1. Do we have any existing records for the Guild?
	dbId, err := database.Guild_DoesGuildExist(newGuild.ID)
	if err != nil && !strings.Contains(err.Error(), "no rows") {
		logger.Error(newGuild.ID, err)
	}

	var triggerList []triggers.Phrase

	if dbId > 0 {
		// => Guild already exists, update the Member Count
		err = database.Guild_UpdateMemberCount(newGuild.ID, newGuild.MemberCount)
		if err != nil {
			logger.Error(newGuild.ID, err)
			return
		}

		// => Does the Guild have any Triggers? If so get for the Cache
		phraseLinks, err := database.GetAllGuildPhrases(newGuild.ID)
		if err != nil {
			logger.Error(newGuild.ID, err)
			return
		}

		for _, phrase := range phraseLinks {
			triggerList = append(triggerList, phrase.Phrase)
		}

		// => Add Global Phrases
		triggerList = append(triggerList, triggers.GlobalPhrases...)

		logger.Event(newGuild.ID, "Existing Guild connected: %v", newGuild.Name)
	} else {
		// 2. This is a new Guild, perform our First Time setup
		newId, err := database.Guild_InsertNewEntry(newGuild.ID, newGuild.Name, newGuild.MemberCount, newGuild.OwnerID)
		if err != nil {
			logger.Error(newGuild.ID, err)
			return
		}

		if newId > 0 {
			dbId = newId
		} else {
			logger.Error(newGuild.ID, errors.New("guild insert returned 0"))
			return
		}
	}

	// 3. Register Commands with the Guild
	RefreshSlashCommands(newGuild.ID)

	// 4. Add to the Active Cache
	cache.AddToActiveGuildCache(newGuild, dbId, triggerList)
	reporting.Guilds()
}
