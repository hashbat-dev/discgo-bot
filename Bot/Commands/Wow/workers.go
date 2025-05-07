package wow

import (
	"github.com/bwmarrin/discordgo"
	config "github.com/hashbat-dev/discgo-bot/Config"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

var (
	QueueDetect   chan (*discordgo.MessageCreate) = make(chan *discordgo.MessageCreate)
	queueGenerate chan (*discordgo.MessageCreate) = make(chan *discordgo.MessageCreate)
	queueRespond  chan (Response)                 = make(chan Response)
)

type Response struct {
	GuildID   string
	ChannelID string
	ReplyRef  *discordgo.MessageReference
	WowText   string
}

func Start() {
	go workerDetect()
	go workerGenerate()
	go workerRespond()
}

func workerDetect() {
	logger.Info("WOW", "Detect Queue starting...")
	for item := range QueueDetect {
		go func(i *discordgo.MessageCreate) {
			Detect(i)
		}(item)
	}
}

func workerGenerate() {
	logger.Info("WOW", "Generate Queue starting...")
	for item := range queueGenerate {
		go func(i *discordgo.MessageCreate) {
			generate(i)
		}(item)
	}
}

func workerRespond() {
	logger.Info("WOW", "Respond Queue starting...")
	for item := range queueRespond {
		go func(i Response) {
			msg := &discordgo.MessageSend{
				Content:   i.WowText,
				Reference: i.ReplyRef,
			}
			_, err := config.Session.ChannelMessageSendComplex(i.ChannelID, msg)
			if err != nil {
				logger.Error(i.GuildID, err)
			}
		}(item)
	}
}
