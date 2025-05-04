package handlers

import (
	"github.com/bwmarrin/discordgo"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

var (
	handleNewMessageQueue          chan (*discordgo.MessageCreate)         = make(chan *discordgo.MessageCreate)
	handleNewGuildQueue            chan (*discordgo.GuildCreate)           = make(chan *discordgo.GuildCreate)
	handleInteractionResponseQueue chan (*discordgo.InteractionCreate)     = make(chan *discordgo.InteractionCreate)
	handleMessageReactionAdd       chan (*discordgo.MessageReactionAdd)    = make(chan *discordgo.MessageReactionAdd)
	handleMessageReactionRemove    chan (*discordgo.MessageReactionRemove) = make(chan *discordgo.MessageReactionRemove)
)

func Start() {
	go workerNewMessageQueue()
	go workerNewGuildQueue()
	go workerInteractionResponseQueue()
	go workerMessageReactionAdd()
	go workerMessageReactionRemove()
}

func workerNewMessageQueue() {
	logger.Info("HANDLERS", "New Message Queue starting...")
	for item := range handleNewMessageQueue {
		go func(i *discordgo.MessageCreate) {
			ProcessNewMessage(i)
		}(item)
	}
}

func workerNewGuildQueue() {
	logger.Info("HANDLERS", "New Guild Queue starting...")
	for item := range handleNewGuildQueue {
		go func(i *discordgo.GuildCreate) {
			ProcessNewGuild(i)
		}(item)
	}
}

func workerInteractionResponseQueue() {
	logger.Info("HANDLERS", "Interaction Response Queue starting...")
	for item := range handleInteractionResponseQueue {
		go func(i *discordgo.InteractionCreate) {
			ProcessInteractionResponse(i)
		}(item)
	}
}

func workerMessageReactionAdd() {
	logger.Info("HANDLERS", "Message Reaction Add Queue starting...")
	for item := range handleMessageReactionAdd {
		go func(i *discordgo.MessageReactionAdd) {
			ProcessReactionAdd(i)
		}(item)
	}
}

func workerMessageReactionRemove() {
	logger.Info("HANDLERS", "Message Reaction Remove Queue starting...")
	for item := range handleMessageReactionRemove {
		go func(i *discordgo.MessageReactionRemove) {
			ProcessReactionRemove(i)
		}(item)
	}
}
