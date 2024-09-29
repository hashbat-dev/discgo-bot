package reactions

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	cache "github.com/hashbat-dev/discgo-bot/Cache"
	config "github.com/hashbat-dev/discgo-bot/Config"
	database "github.com/hashbat-dev/discgo-bot/Database"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

var (
	upName   = "hall-of-fame"
	downName = "hall-of-shame"
)

func CreateChannel(guildId string, isUp bool) string {
	channelName := upName
	if !isUp {
		channelName = downName
	}

	permissionOverwrites := []*discordgo.PermissionOverwrite{
		{
			ID:    guildId,                                                            // @everyone role ID
			Type:  discordgo.PermissionOverwriteTypeRole,                              // Role permission
			Deny:  discordgo.PermissionViewChannel | discordgo.PermissionSendMessages, // Deny sending messages
			Allow: discordgo.PermissionViewChannel,                                    // Allow viewing
		},
		{
			ID:    config.Session.State.User.ID,                                       // The Bot's Use ID
			Type:  discordgo.PermissionOverwriteTypeMember,                            // Member specific permission
			Allow: discordgo.PermissionViewChannel | discordgo.PermissionSendMessages, // View and send messages
		},
	}

	channel, err := config.Session.GuildChannelCreateComplex(guildId, discordgo.GuildChannelCreateData{
		Name:                 channelName,
		Type:                 discordgo.ChannelTypeGuildText,
		PermissionOverwrites: permissionOverwrites,
	})
	if err != nil {
		logger.Error(guildId, err)
		return ""
	}

	if channel.ID == "0" || channel.ID == "" {
		err = errors.New("returned channel id = 0")
		logger.Error(guildId, err)
		return ""
	}

	guildObj, err := database.Get(guildId)
	if err != nil {
		return ""
	}

	if isUp {
		guildObj.StarUpChannel = channel.ID
	} else {
		guildObj.StarDownChannel = channel.ID
	}

	_, err = database.Upsert(guildObj)
	if err != nil {
		return ""
	}

	cache.UpdateStarboardChannel(guildId, channel.ID, isUp)
	return channel.ID
}
