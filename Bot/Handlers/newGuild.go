package handlers

import (
	"github.com/bwmarrin/discordgo"
	database "github.com/dabi-ngin/discgo-bot/Database"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
)

// Calls whenever a new Guild connects to the bot. This also runs for all active Guilds on startup.
func HandleNewGuild(session *discordgo.Session, newGuild *discordgo.GuildCreate) {

	// 1. Do we have any existing records for the Guild?
	guildExists, err := database.Guild_DoesGuildExist(newGuild.ID)
	if err != nil {
		logger.Error(newGuild.ID, err)
	}

	// => Guild already exists, update the Member Count and exit
	if guildExists {
		err = database.Guild_UpdateMemberCount(newGuild.ID, newGuild.MemberCount)
		if err != nil {
			logger.Error(newGuild.ID, err)
		}

		logger.Event(newGuild.ID, "Existing Guild connected: %v", newGuild.Name)
		return
	}

	// 2. This is a new Guild, perform our First Time setup
	err = database.Guild_InsertNewEntry(newGuild.ID, newGuild.Name, newGuild.MemberCount, newGuild.OwnerID)
	if err != nil {
		logger.Error(newGuild.ID, err)
		return
	}

}
