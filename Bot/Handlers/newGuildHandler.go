package handlers

import (
	"sync"

	"github.com/bwmarrin/discordgo"
	triggers "github.com/hashbat-dev/discgo-bot/Bot/Commands/Triggers"
	cache "github.com/hashbat-dev/discgo-bot/Cache"
	database "github.com/hashbat-dev/discgo-bot/Database"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
	reactions "github.com/hashbat-dev/discgo-bot/Reactions"
	reporting "github.com/hashbat-dev/discgo-bot/Reporting"
)

var guildMutex sync.Mutex

// Calls whenever a new Guild connects to the bot. This also runs for all active Guilds on startup.
func HandleNewGuild(session *discordgo.Session, newGuild *discordgo.GuildCreate) {
	guildMutex.Lock()
	defer guildMutex.Unlock()

	// 1. Do we have any existing records for the Guild?
	guild, err := database.Get(newGuild.ID)
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

	// 3. Get the Guild Starboard Emojis
	guildEmojis, err := database.GetAllGuildEmojis(guild.GuildID)
	if err == nil && len(guildEmojis) == 0 {
		// None assigned, provide the bare basics
		err = reactions.AddGuildEmoji(guild.GuildID, "", reactions.StandardUp, reactions.EmojiCategoryUp)
		if err != nil {
			logger.ErrorText(guild.GuildID, "Failed to add Standard 'Up' Emoji")
		}
		err = reactions.AddGuildEmoji(guild.GuildID, "", reactions.StandardDown, reactions.EmojiCategoryDown)
		if err != nil {
			logger.ErrorText(guild.GuildID, "Failed to add Standard 'Down' Emoji")
		}
		guildEmojis, err = database.GetAllGuildEmojis(guild.GuildID)
		if err != nil {
			logger.ErrorText(guild.GuildID, "Failed to get newly inserted Guild Emojis")
		}
	}

	// 4. Update our Guild Information
	// => Add Global Phrases
	triggerList = append(triggerList, triggers.GlobalPhrases...)

	// => Update the Guild Info
	guild.GuildName = newGuild.Name
	guild.GuildMemberCount = newGuild.MemberCount
	guild.GuildOwnerID = newGuild.OwnerID

	// 5. Update the Database with this information
	newG, err := database.Upsert(guild)
	if err != nil {
		logger.ErrorText(guild.GuildID, "Error updating Database")
	} else {
		guild = newG
	}

	// 6. Add to the Active Cache
	cache.AddToActiveGuildCache(guild.ID, guild.GuildID, guild.IsDevServer, guild.GuildName, triggerList, guild.StarUpChannel,
		guild.StarDownChannel, guild.GuildOwnerID, guild.GuildAdminRole, guildEmojis)
	reporting.Guilds()
}
