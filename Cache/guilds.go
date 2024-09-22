package cache

import (
	"time"

	triggers "github.com/dabi-ngin/discgo-bot/Bot/Commands/Triggers"
	helpers "github.com/dabi-ngin/discgo-bot/Helpers"
)

var ActiveGuilds map[string]Guild = make(map[string]Guild)

func AddToActiveGuildCache(dbId int, guildId string, isDev bool, guildName string, triggers []triggers.Phrase, starUpChannel string,
	starDownChannel string, ownerId string, adminRoleId string) {
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
