package handlers

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	commands "github.com/dabi-ngin/discgo-bot/Bot/Commands"
	triggers "github.com/dabi-ngin/discgo-bot/Bot/Commands/Triggers"
	cache "github.com/dabi-ngin/discgo-bot/Cache"
	config "github.com/dabi-ngin/discgo-bot/Config"
	database "github.com/dabi-ngin/discgo-bot/Database"
	discord "github.com/dabi-ngin/discgo-bot/Discord"
	helpers "github.com/dabi-ngin/discgo-bot/Helpers"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
	reporting "github.com/dabi-ngin/discgo-bot/Reporting"
	"github.com/google/uuid"
)

type Task struct {
	CommandType          int
	Complexity           int
	BangDetails          *BangTaskDetails
	SlashDetails         *SlashTaskDetails
	SlashResponseDetails *SlashResponseDetails
	PhraseDetails        *PhraseTaskDetails
}

type BangTaskDetails struct {
	Message       *discordgo.MessageCreate
	Command       commands.Command
	CorrelationId uuid.UUID
}

type SlashTaskDetails struct {
	Interaction  *discordgo.InteractionCreate
	SlashCommand SlashCommand
}

type SlashResponseDetails struct {
	Interaction   *discordgo.InteractionCreate
	ObjectID      string
	CorrelationID string
}

type PhraseTaskDetails struct {
	Message        *discordgo.MessageCreate
	TriggerPhrases []triggers.Phrase
}

var (
	MAX_QUEUE_IO_TASKS      = 100
	MAX_QUEUE_CPU_TASKS     = 100
	MAX_QUEUE_TRIVIAL_TASKS = 1000
	IO_TASKS                = make(chan *Task, MAX_QUEUE_IO_TASKS)
	CPU_TASKS               = make(chan *Task, MAX_QUEUE_CPU_TASKS)
	TRIVIAL_TASKS           = make(chan *Task, MAX_QUEUE_TRIVIAL_TASKS)
)

func init() {
	reporting.CreateWorkerChannel(config.IO_BOUND_TASK, "IO Bound", MAX_QUEUE_IO_TASKS, config.N_IO_WORKERS)
	reporting.CreateWorkerChannel(config.CPU_BOUND_TASK, "CPU Bound", MAX_QUEUE_CPU_TASKS, runtime.GOMAXPROCS(0))
	reporting.CreateWorkerChannel(config.TRIVIAL_TASK, "Trivial", MAX_QUEUE_TRIVIAL_TASKS, config.N_TRIVIAL_WORKERS)
	for i := 0; i < config.N_TRIVIAL_WORKERS; i++ {
		go worker(i, config.TRIVIAL_TASK, TRIVIAL_TASKS)
	}
	for i := 0; i < config.N_IO_WORKERS; i++ {
		go worker(i, config.IO_BOUND_TASK, IO_TASKS)
	}
	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		go worker(i, config.CPU_BOUND_TASK, CPU_TASKS)
	}

}

func worker(id int, taskId int, ch <-chan *Task) {
	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				// Channel is closed, exit goroutine
				logger.Info("WORKER", "worker %d: Channel closed, exiting...", id)
				return
			}

			reporting.WorkerProcessingStart(taskId)
			switch msg.CommandType {
			case config.CommandTypeBang:
				workerBang(msg.BangDetails)
			case config.CommandTypeSlash:
				workerSlash(msg.SlashDetails)
			case config.CommandTypeSlashResponse:
				workerSlashResponse(msg.SlashResponseDetails)
			case config.CommandTypePhrase:
				workerPhrase(msg.PhraseDetails)
			default:
				logger.ErrorText("WORKER", "Unknown CommandType value [%v]", msg.CommandType)
			}

			reporting.WorkerProcessingFinish(taskId)
		}
	}
}

func workerBang(msg *BangTaskDetails) {
	logger.Info(msg.Message.GuildID, "worker :: processing command [%v] correlation-id :: %v", msg.Command.Name(), msg.CorrelationId)
	timeStart := time.Now()

	execErr := msg.Command.Execute(msg.Message, msg.Command.Name())
	if execErr != nil {
		logger.ErrorText(msg.Message.GuildID, "worker :: [%v] error :: %v :: correlation-id :: %v", msg.Command.Name(), execErr.Error(), msg.CorrelationId)
		return // Failed to execute, skip loop iteration
	}

	reporting.Command(config.CommandTypeBang, msg.Message.GuildID, msg.Message.Author.ID, msg.Message.Author.Username, msg.Command.Name(), msg.CorrelationId.String(), timeStart)
}

func workerSlash(msg *SlashTaskDetails) {
	correlationId := cache.AddInteraction(msg.Interaction, msg.SlashCommand.Command.Name)
	timeStarted := time.Now()
	msg.SlashCommand.Handler(msg.Interaction, correlationId)
	reporting.Command(config.CommandTypeSlash, msg.Interaction.GuildID, msg.Interaction.Member.User.ID, msg.Interaction.Member.User.Username, msg.SlashCommand.Command.Name, correlationId, timeStarted)
}

func workerSlashResponse(msg *SlashResponseDetails) {
	logger.Info(msg.Interaction.GuildID, "Interaction ID: [%v] Processing Response, Object: [%v]", msg.CorrelationID, msg.ObjectID)
	timeStarted := time.Now()

	// Update the Interaction Cache with the provided options
	cache.UpdateInteraction(msg.CorrelationID, msg.Interaction)

	// Execute the Command Response
	discord.InteractionResponseHandlers[msg.ObjectID].Execute(msg.Interaction, msg.CorrelationID)
	reporting.Command(config.CommandTypeSlash, msg.Interaction.GuildID, msg.Interaction.Member.User.ID, msg.Interaction.Member.User.Username, msg.ObjectID, msg.CorrelationID, timeStarted)
}

func workerPhrase(msg *PhraseTaskDetails) {
	var notifyPhrases []string
	timeStarted := time.Now()

	// Process all generic Trigger Phrases first
	for _, phrase := range msg.TriggerPhrases {
		if phrase.IsSpecial {
			continue
		}
		reporting.Command(config.CommandTypePhrase, msg.Message.GuildID, msg.Message.Author.ID, msg.Message.Author.Username, phrase.Phrase, uuid.New().String(), timeStarted)
		if phrase.NotifyOnDetection {
			notifyPhrases = append(notifyPhrases, phrase.Phrase)
		}
	}

	if len(notifyPhrases) > 0 {
		showText := strings.ToUpper(helpers.ConcatStringWithAnd(notifyPhrases)) + " MENTIONED"
		_, err := config.Session.ChannelMessageSend(msg.Message.ChannelID, showText)
		if err != nil {
			logger.Error(msg.Message.GuildID, err)
		}
	}

	// Now process any Special ones
	for _, phrase := range msg.TriggerPhrases {
		if !phrase.IsSpecial {
			continue
		}
		switch phrase.Phrase {
		case "jason statham":
			webm, err := database.GetRandomResource(msg.Message.GuildID, 1)
			if err != nil {
				continue
			}
			discord.SendUserMessageReply(msg.Message, false, fmt.Sprintf("[%s](%s)", "Jason Statham?", webm))
		default:
			logger.ErrorText(msg.Message.GuildID, "Unhandled Special Phrase [%v]", phrase.Phrase)
		}
	}
}
