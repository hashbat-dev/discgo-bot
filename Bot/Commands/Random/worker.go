package random

import (
	"github.com/bwmarrin/discordgo"
	config "github.com/hashbat-dev/discgo-bot/Config"
	database "github.com/hashbat-dev/discgo-bot/Database"
	helpers "github.com/hashbat-dev/discgo-bot/Helpers"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

var (
	QueueRandom chan (*discordgo.MessageCreate) = make(chan *discordgo.MessageCreate)
)

func Start() {
	go workerRandom()
}

func workerRandom() {
	logger.Info("RANDOM", "Random Queue starting...")
	for item := range QueueRandom {
		go func(i *discordgo.MessageCreate) {
			execute(i)
		}(item)
	}
}

func execute(message *discordgo.MessageCreate) {
	// Speech?
	maxSpeech := 400
	if message.Author.ID == config.NON_SPECIFIC_USER {
		maxSpeech = 200
	}
	if helpers.GetRandomNumber(1, maxSpeech) == 150 {
		imgUrl, err := database.GetRandomImage(message.GuildID, 1)
		if err != nil {
			return
		}
		_, err = config.Session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
			Content:   imgUrl,
			Reference: message.Reference(),
		})
		if err != nil {
			logger.Error(message.GuildID, err)
			return
		}
		logger.Info(message.GuildID, "%s (%s) was randomly !speech'd", message.Author.ID, message.Author.Username)
	}
}
