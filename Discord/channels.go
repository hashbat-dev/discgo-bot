package discord

import (
	"time"

	"github.com/bwmarrin/discordgo"
	cache "github.com/hashbat-dev/discgo-bot/Cache"
	config "github.com/hashbat-dev/discgo-bot/Config"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

func CreateAdminChannel(guildId string, channelName string) (*discordgo.Channel, error) {
	adminRoleId := cache.ActiveGuilds[guildId].BotAdminRole
	guildOwnerId := ""
	if adminRoleId == "" {
		guildOwnerId = cache.ActiveGuilds[guildId].ServerOwner
	}

	var overwrites []*discordgo.PermissionOverwrite

	if guildOwnerId != "" {
		// Server Owner
		overwrites = []*discordgo.PermissionOverwrite{
			{
				ID:    guildOwnerId,
				Type:  discordgo.PermissionOverwriteTypeMember,
				Allow: discordgo.PermissionViewChannel,
			},
			{
				ID:   guildId,
				Type: discordgo.PermissionOverwriteTypeRole,
				Deny: discordgo.PermissionViewChannel,
			},
		}
	} else {
		// Bot Admin Role
		overwrites = []*discordgo.PermissionOverwrite{
			{
				ID:    adminRoleId,
				Type:  discordgo.PermissionOverwriteTypeRole,
				Allow: discordgo.PermissionViewChannel,
			},
			{
				ID:   guildId,
				Type: discordgo.PermissionOverwriteTypeRole,
				Deny: discordgo.PermissionViewChannel,
			},
		}
	}

	// Create the private channel
	channel, err := config.Session.GuildChannelCreateComplex(guildId, discordgo.GuildChannelCreateData{
		Name:                 channelName,
		Type:                 discordgo.ChannelTypeGuildText,
		PermissionOverwrites: overwrites,
	})
	if err != nil {
		logger.Error(guildId, err)
	}

	cache.ActiveAdminChannels[channel.ID] = time.Now()
	return channel, err
}

func DeleteAdminChannel(guildId string, channelId string) {
	_, err := config.Session.ChannelDelete(channelId)
	if err != nil {
		logger.Error(guildId, err)
	}
	delete(cache.ActiveAdminChannels, channelId)
}
