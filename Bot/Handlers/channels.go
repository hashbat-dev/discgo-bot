package handlers

import (
	"time"

	"github.com/bwmarrin/discordgo"
	cache "github.com/dabi-ngin/discgo-bot/Cache"
	config "github.com/dabi-ngin/discgo-bot/Config"
	helpers "github.com/dabi-ngin/discgo-bot/Helpers"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
)

func init() {

	// Initialise the required Channel variables
	for i := 0; i <= config.LastPoolIota; i++ {
		Pools[i] = make(chan *ChannelRequest)
		PoolDurations[i] = PoolDuration{}
		PoolAvgDurations[i] = 0
		PoolProcessing[i] = 0
		PoolQueue[i] = 0
		PoolLastAdded[i] = helpers.GetNullDateTime()
		QueueLastAdded[i] = helpers.GetNullDateTime()
	}

	// Spin up the GoRoutines
	for i := 0; i <= config.LastPoolIota; i++ {
		for i := 0; i < config.ProcessPools[i].MaxWorkers; i++ {
			go commandWorker(i, Pools[i])
		}
	}
}

type ChannelRequest struct {
	Message     *discordgo.MessageCreate
	CommandName string
	Command     Command
}

type DashboardChannelInfo struct {
	Name                string
	ProcessingCount     int
	ProcessingLastAdded time.Time
	QueueCount          int
	QueueLastAdded      time.Time
	AverageDuration     time.Duration
}

type PoolDuration struct {
	Durations   []time.Duration
	AvgDuration time.Duration
}

var (
	Pools            = make(map[int]chan *ChannelRequest)
	PoolDurations    = map[int]PoolDuration{}
	PoolAvgDurations = map[int]time.Duration{}
	PoolProcessing   = map[int]int{}
	PoolQueue        = map[int]int{}
	PoolLastAdded    = map[int]time.Time{}
	QueueLastAdded   = map[int]time.Time{}
)

func commandWorker(id int, ch <-chan *ChannelRequest) {
	for {
		select {
		case chanMessage, ok := <-ch:
			if !ok {
				// Channel is closed, exit goroutine
				logger.Info(chanMessage.Message.GuildID, "CommandWorker %d: Channel closed, exiting...\n", id)
				return
			}

			logger.Info(chanMessage.Message.GuildID, "CommandWorker :: [%v] Starting", chanMessage.CommandName)
			timeStart := time.Now()

			PoolQueue[chanMessage.Command.ProcessPool().ProcessPoolIota]--
			PoolProcessing[chanMessage.Command.ProcessPool().ProcessPoolIota]++

			err := chanMessage.Command.Execute(chanMessage.Message, chanMessage.CommandName)

			PoolProcessing[chanMessage.Command.ProcessPool().ProcessPoolIota]--

			callDuration := time.Since(timeStart)
			cache.AddToCommandCache(config.CommandTypeBang, chanMessage.CommandName, chanMessage.Message.GuildID, chanMessage.Message.Author.ID, chanMessage.Message.Author.Username, timeStart, callDuration)
			AddChannelTimings(chanMessage.Command.ProcessPool().ProcessPoolIota, callDuration)

			if err != nil {
				logger.ErrorText(chanMessage.Message.GuildID, "CommandWorker :: [%v] Ended in error after %v :: %v", chanMessage.CommandName, callDuration, err.Error())
			} else {
				logger.Info(chanMessage.Message.GuildID, "CommandWorker :: [%v] Ended successfully after %v", chanMessage.CommandName, callDuration)
			}
		}
	}
}

func AddChannelTimings(channelIota int, callDuration time.Duration) {

	cacheDuration := PoolDurations[channelIota]
	cacheDuration.Durations = append(PoolDurations[channelIota].Durations, callDuration)
	cacheDuration.AvgDuration = helpers.AverageDuration(cacheDuration.Durations)

	if len(cacheDuration.Durations) > config.CommandAveragePool {
		cacheDuration.Durations = cacheDuration.Durations[1:]
	}

	PoolDurations[channelIota] = cacheDuration

}
