package cache

import (
	"time"

	"github.com/bwmarrin/discordgo"
	triggers "github.com/dabi-ngin/discgo-bot/Bot/Commands/Triggers"
	helpers "github.com/dabi-ngin/discgo-bot/Helpers"
)

var ActiveGuilds map[string]Guild = make(map[string]Guild)

func AddToActiveGuildCache(guild *discordgo.GuildCreate, dbId int, isDev bool, triggers []triggers.Phrase, starUpId string, starDownId string) {
	ActiveGuilds[guild.ID] = Guild{
		DbID:            dbId,
		DiscordID:       guild.ID,
		IsDev:           isDev,
		Name:            guild.Name,
		LastCommand:     helpers.GetNullDateTime(),
		Triggers:        triggers,
		StarUpChannel:   starUpId,
		StarDownChannel: starDownId,
	}
}

func UpdateLastGuildCommand(guildId string) {
	guildInfo := ActiveGuilds[guildId]
	guildInfo.LastCommand = time.Now()
	guildInfo.CommandCount++
	ActiveGuilds[guildId] = guildInfo
}

func UpdateStarboardChannel(guildId string, channelId string, isUp bool) {
	guildInfo := ActiveGuilds[guildId]
	if isUp {
		guildInfo.StarUpChannel = channelId
	} else {
		guildInfo.StarDownChannel = channelId
	}
	ActiveGuilds[guildId] = guildInfo
}
