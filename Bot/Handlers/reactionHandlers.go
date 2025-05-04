package handlers

import (
	"github.com/bwmarrin/discordgo"
	cache "github.com/hashbat-dev/discgo-bot/Cache"
	config "github.com/hashbat-dev/discgo-bot/Config"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

func HandleReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	handleMessageReactionAdd <- r
}

func HandleReactionRemove(s *discordgo.Session, r *discordgo.MessageReactionRemove) {
	handleMessageReactionRemove <- r
}

func ProcessReactionAdd(r *discordgo.MessageReactionAdd) {
	if skipReactionCheck(r.UserID, r.GuildID, r.ChannelID) {
		return
	}
	startReactionCheck(r.GuildID, r.ChannelID, r.MessageID)
}

func ProcessReactionRemove(r *discordgo.MessageReactionRemove) {
	if skipReactionCheck(r.UserID, r.GuildID, r.ChannelID) {
		return
	}
	startReactionCheck(r.GuildID, r.ChannelID, r.MessageID)
}

func startReactionCheck(guildId string, channelId string, messageId string) {
	message, err := config.Session.ChannelMessage(channelId, messageId)
	if err != nil {
		logger.Error(guildId, err)
	}
	message.GuildID = guildId
	DispatchTask(&Task{
		CommandType: config.CommandTypeReactionCheck,
		Complexity:  config.IO_BOUND_TASK,
		MessageObj:  message,
	})
}

func skipReactionCheck(userId string, guildId string, channelId string) bool {
	if userId == config.Session.State.User.ID {
		return true
	}

	if config.ServiceSettings.ISDEV != cache.ActiveGuilds[guildId].IsDev {
		return true
	}

	if _, exists := cache.ActiveAdminChannels[channelId]; exists {
		logger.Debug(guildId, "Skipping Reaction check as sourced from an Admin Channel")
		return true
	}
	return false
}
