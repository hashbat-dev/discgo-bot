package discord

import (
	config "github.com/hashbat-dev/discgo-bot/Config"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
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
