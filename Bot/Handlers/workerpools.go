package handlers

import (
	"runtime"
	"time"

	"github.com/bwmarrin/discordgo"
	commands "github.com/dabi-ngin/discgo-bot/Bot/Commands"
	config "github.com/dabi-ngin/discgo-bot/Config"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
	reporting "github.com/dabi-ngin/discgo-bot/Reporting"
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

			logger.Info(msg.Message.GuildID, "commandWorker[%d] :: processing command [%v] correlation-id :: %v", id, msg.Command.Name(), msg.CorrelationId)
			timeStart := time.Now()

			execErr := msg.Command.Execute(msg.Message, msg.Command.Name())
			if execErr != nil {
				logger.ErrorText(msg.Message.GuildID, "commandWorker[%d] :: [%v] error :: %v :: correlation-id :: %v", id, msg.Command.Name(), execErr.Error(), msg.CorrelationId)
				continue // Failed to execute, skip loop iteration
			}

			reporting.Command(config.CommandTypeBang, msg.Message, msg.Command, msg.CorrelationId, timeStart)
		}
	}
}
