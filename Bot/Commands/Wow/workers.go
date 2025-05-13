package wow

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	config "github.com/hashbat-dev/discgo-bot/Config"
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
			msg := &discordgo.MessageSend{
				Content:   i.Output,
				Reference: i.Message.Reference(),
			}
			wowMsg, err := config.Session.ChannelMessageSendComplex(i.Message.ChannelID, msg)
			if err != nil {
				logger.Error(i.Message.GuildID, err)
			}
			i.WowMessageID = wowMsg.ID
			addToCache(&i)

			// New Record?
			if item.OCount > dataHighestWowInGuild[i.Message.GuildID] {
				emb := embed.NewEmbed()
				emb.SetTitle("New Wow Record")
				emb.SetDescription(fmt.Sprintf("Ladies and Gentlemen, <@%s> has just broken the all time Wow record!", i.Message.Author.ID))
				emb.SetThumbnail(config.TROPHY_IMG_URL)
				emb.SetFooter(fmt.Sprintf("%d level Wow", i.OCount))
				_, err = config.Session.ChannelMessageSendComplex(i.Message.ChannelID, &discordgo.MessageSend{
					Reference: wowMsg.Reference(),
					Embed:     emb.MessageEmbed,
				})

				if err != nil {
					logger.Error(i.Message.GuildID, err)
				}

				dataHighestWowInGuild[i.Message.GuildID] = item.OCount
			}
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
