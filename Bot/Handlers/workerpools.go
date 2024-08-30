package handlers

import (
	"runtime"
	"time"

	"github.com/bwmarrin/discordgo"
	commands "github.com/dabi-ngin/discgo-bot/Bot/Commands"
	cache "github.com/dabi-ngin/discgo-bot/Cache"
	config "github.com/dabi-ngin/discgo-bot/Config"
	dashboard "github.com/dabi-ngin/discgo-bot/Dashboard"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
	"github.com/google/uuid"
)

type CommandTask struct {
	Message       *discordgo.MessageCreate
	Command       commands.Command
	CorrelationId uuid.UUID
}

var (
	IO_TASKS      = make(chan *CommandTask)
	CPU_TASKS     = make(chan *CommandTask)
	TRIVIAL_TASKS = make(chan *CommandTask)
)

func init() {
	for i := 0; i < config.N_TRIVIAL_WORKERS; i++ {
		go commandWorker(i, TRIVIAL_TASKS)
	}
	for i := 0; i < config.N_IO_WORKERS; i++ {
		go commandWorker(i, IO_TASKS)
	}
	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		go commandWorker(i, CPU_TASKS)
	}

}

func commandWorker(id int, ch <-chan *CommandTask) {
	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				// Channel is closed, exit goroutine
				logger.Info(msg.Message.GuildID, "commandWorker %d: Channel closed, exiting...\n", id)
				return
			}
			start := time.Now()
			logger.Info(msg.Message.GuildID, "commandWorker[%d] :: processing command [%v] correlation-id :: %v", id, msg.Command.Name(), msg.CorrelationId)

			execErr := msg.Command.Execute(msg.Message, msg.Command.Name())
			if execErr != nil {
				logger.ErrorText(msg.Message.GuildID, "commandWorker[%d] :: [%v] error :: %v :: correlation-id :: %v", id, msg.Command.Name(), execErr.Error(), msg.CorrelationId)
				continue // Failed to execute, skip loop iteration
			}

			duration := time.Since(start)
			cache.AddToCommandCache(config.CommandTypeBang, msg.Command.Name(), msg.Message.GuildID, msg.Message.Author.ID, msg.Message.Author.Username, start, duration)
			latencyMsg := dashboard.LatencyMessage{CommandName: msg.Command.Name(), TimeTaken: duration}
			dashboard.LatencyCh <- latencyMsg
			logger.Info(msg.Message.GuildID, "CommandWorker :: [%v] Ended successfully after %v :: correlation-id :: %v", msg.Command.Name(), duration, msg.CorrelationId)
		}
	}
}
