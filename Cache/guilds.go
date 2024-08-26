package cache

import (
	"github.com/bwmarrin/discordgo"
	triggers "github.com/dabi-ngin/discgo-bot/Bot/Commands/Triggers"
	helpers "github.com/dabi-ngin/discgo-bot/Helpers"
)

var ActiveGuilds []Guild

func AddToActiveGuildCache(guild *discordgo.GuildCreate, dbId int, triggers []triggers.Phrase) {
	if ActiveGuilds == nil {
		ActiveGuilds = []Guild{}
	}

	ActiveGuilds = append(ActiveGuilds, Guild{
		DbID:        dbId,
		DiscordID:   guild.ID,
		Name:        guild.Name,
		LastCommand: helpers.GetNullDateTime(),
		Triggers:    triggers,
	})

}

func GetGuildIndex(guildId string) int {
	guildIndex := -1
	for i, guild := range ActiveGuilds {
		if guild.DiscordID == guildId {
			guildIndex = i
			break
		}
	}

	return guildIndex
}
