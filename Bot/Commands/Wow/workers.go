package wow

import (
	"github.com/bwmarrin/discordgo"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

var (
	QueueDetect   chan (*discordgo.MessageCreate) = make(chan *discordgo.MessageCreate)
	queueGenerate chan (*discordgo.MessageCreate) = make(chan *discordgo.MessageCreate)
	queueRespond  chan (*Generation)              = make(chan *Generation)
	queueDatabase chan (*Generation)              = make(chan *Generation)
)

type Response struct {
	GuildID   string
	ChannelID string
	ReplyRef  *discordgo.MessageReference
	Wow       Generation
}

func Start() {
	go workerDetect()
	go workerGenerate()
	go workerRespond()
	go workerDatabase()

	// Pokemon Data
	go getAllPokemon()
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
		go func(i Generation) {
			respond(&i)
		}(*item)
	}
}

func workerDatabase() {
	logger.Info("WOW", "Database Queue starting...")
	for item := range queueDatabase {
		go func(i Generation) {
			postToDatabase(&i)
		}(*item)
	}
}
