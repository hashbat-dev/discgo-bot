package wow

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

func generate(message *discordgo.MessageCreate) {
	count := getRandomNumber(1, 10)

	for _, fn := range effectList {
		i, effect := fn(message)
		if effect != nil {
			logger.Info(message.GuildID, "Wow generation, MessageID: %s, Effect: %s, Count: %d", message.ID, effect.Name, i)
			count += i
		}
	}

	if count < 1 {
		count = 1
	}

	wow := fmt.Sprintf("w%sw", getOs(count))
	pushToRespondQueue(message, wow)
}

func pushToRespondQueue(message *discordgo.MessageCreate, wowText string) {
	logger.Info(message.GuildID, "Wow response queued, MessageID: %s", message.ID)
	queueRespond <- Response{
		GuildID:   message.GuildID,
		ChannelID: message.ChannelID,
		ReplyRef:  message.Reference(),
		WowText:   wowText,
	}
}
