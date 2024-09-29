package cache

import (
	"time"

	triggers "github.com/hashbat-dev/discgo-bot/Bot/Commands/Triggers"
	data "github.com/hashbat-dev/discgo-bot/Data"
	helpers "github.com/hashbat-dev/discgo-bot/Helpers"
)

var ActiveGuilds map[string]Guild = make(map[string]Guild)
var ActiveAdminChannels map[string]time.Time = make(map[string]time.Time)

func AddToActiveGuildCache(dbId int, guildId string, isDev bool, guildName string, triggers []triggers.Phrase, starUpChannel string,
	starDownChannel string, ownerId string, adminRoleId string, reactionEmojis []data.GuildEmoji) {
	ActiveGuilds[guildId] = Guild{
		DbID:            dbId,
		DiscordID:       guildId,
		IsDev:           isDev,
		Name:            guildName,
		LastCommand:     helpers.GetNullDateTime(),
		Triggers:        triggers,
		StarUpChannel:   starUpChannel,
		StarDownChannel: starDownChannel,
		ServerOwner:     ownerId,
		BotAdminRole:    adminRoleId,
		ReactionEmojis:  reactionEmojis,
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

func UpdateBotAdminRole(guildId string, newRoleId string) {
	guildInfo := ActiveGuilds[guildId]
	guildInfo.BotAdminRole = newRoleId
	ActiveGuilds[guildId] = guildInfo
}
