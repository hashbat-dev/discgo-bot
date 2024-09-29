package reactions

import (
	"github.com/bwmarrin/discordgo"
	cache "github.com/hashbat-dev/discgo-bot/Cache"
	config "github.com/hashbat-dev/discgo-bot/Config"
	database "github.com/hashbat-dev/discgo-bot/Database"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

func DeleteIfExists(message *discordgo.Message) {
	dbEntry := database.Starboard_Get(message.GuildID, message.ID)
	if dbEntry.ID > 0 {
		// 1. Delete the Starboard Message first
		if dbEntry.StarboardMessageID != "" {
			channelId := ""
			if dbEntry.IsUpChannel {
				channelId = cache.ActiveGuilds[message.GuildID].StarUpChannel
			} else {
				channelId = cache.ActiveGuilds[message.GuildID].StarDownChannel
			}

			if channelId != "" {
				err := config.Session.ChannelMessageDelete(channelId, dbEntry.StarboardMessageID)
				if err != nil {
					logger.Error(message.GuildID, err)
				}
			}
		}

		// 2. Then delete from the Database
		database.Starboard_Delete(message.GuildID, dbEntry.ID)
	}
}
