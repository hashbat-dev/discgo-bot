package cache

import (
	"time"

	"github.com/bwmarrin/discordgo"
	helpers "github.com/dabi-ngin/discgo-bot/Helpers"
)

var ActiveGuilds []Guild

type Guild struct {
	DbID         int
	DiscordID    string
	Name         string
	CommandCount int
	LastCommand  time.Time
}

type GuildPermissions struct {
	CommandType  int
	RequiredRole string
}

const (
	CommandTypeAdmin   = iota
	CommandTypeBang    = iota
	CommandTypeTrigger = iota
)

func AddToActiveGuildCache(guild *discordgo.GuildCreate, dbId int) {

	if ActiveGuilds == nil {
		ActiveGuilds = []Guild{}
	}

	ActiveGuilds = append(ActiveGuilds, Guild{
		DbID:        dbId,
		DiscordID:   guild.ID,
		Name:        guild.Name,
		LastCommand: helpers.GetNullDateTime(),
	})

}

func UpdateGuildLastCommand(guildId string) {

	guildIndex := -1
	for i, guild := range ActiveGuilds {
		if guild.DiscordID == guildId {
			guildIndex = i
			break
		}
	}

	if guildIndex > -1 {
		ActiveGuilds[guildIndex].CommandCount++
		ActiveGuilds[guildIndex].LastCommand = time.Now()
	}

}
