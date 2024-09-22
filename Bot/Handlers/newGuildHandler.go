package handlers

import (
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
	guild, err := database.Guild_Get(newGuild.ID)
	if err != nil {
		logger.ErrorText(newGuild.ID, "Failed to process Guild")
		return
	}

	// 2. Get the Guild Triggers
	var triggerList []triggers.Phrase
	if guild.ID > 0 {
		// => Does the Guild have any Triggers? If so get for the Cache
		phraseLinks, err := database.GetAllGuildPhrases(guild.GuildID)
		if err != nil {
			logger.Error(guild.GuildID, err)
			return
		}

		for _, phrase := range phraseLinks {
			triggerList = append(triggerList, phrase.Phrase)
		}
	}

	// => Add Global Phrases
	triggerList = append(triggerList, triggers.GlobalPhrases...)

	// => Update the Guild Info
	guild.GuildName = newGuild.Name
	guild.GuildMemberCount = newGuild.MemberCount
	guild.GuildOwnerID = newGuild.OwnerID

	// 3. Update the Database with this information
	newG, err := database.Guild_InsertUpdate(guild)
	if err != nil {
		logger.ErrorText(guild.GuildID, "Error updating Database")
	} else {
		guild = newG
	}

	// 4. Add to the Active Cache
	cache.AddToActiveGuildCache(guild.ID, guild.GuildID, guild.IsDevServer, guild.GuildName, triggerList, guild.StarUpChannel,
		guild.StarDownChannel, guild.GuildOwnerID, guild.GuildAdminRole)
	reporting.Guilds()
}
