package handlers

import (
	"github.com/bwmarrin/discordgo"
	cache "github.com/dabi-ngin/discgo-bot/Cache"
	config "github.com/dabi-ngin/discgo-bot/Config"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
)

func HandleReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if skipReactionCheck(r.UserID, r.GuildID) {
		return
	}
	startReactionCheck(r.GuildID, r.ChannelID, r.MessageID)
}

func HandleReactionRemove(s *discordgo.Session, r *discordgo.MessageReactionRemove) {
	if skipReactionCheck(r.UserID, r.GuildID) {
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

func skipReactionCheck(userId string, guildId string) bool {
	if userId == config.Session.State.User.ID {
		return true
	}

	if config.ServiceSettings.ISDEV != cache.ActiveGuilds[guildId].IsDev {
		return true
	}

	return false
}
