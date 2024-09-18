package discord

import (
	config "github.com/dabi-ngin/discgo-bot/Config"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
)

func SendDM(guildId string, userID string, message string) error {
	channel, err := config.Session.UserChannelCreate(userID)
	if err != nil {
		logger.Error(guildId, err)
		return err
	}

	_, err = config.Session.ChannelMessageSend(channel.ID, message)
	if err != nil {
		logger.Error(guildId, err)
		return err
	}

	return nil
}
