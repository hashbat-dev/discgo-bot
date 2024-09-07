package dashboard

import (
	"time"

	config "github.com/dabi-ngin/discgo-bot/Config"
	helpers "github.com/dabi-ngin/discgo-bot/Helpers"
)

// TODO - Implement a channel which will receive task information from the handlers package

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
	PoolDurations    = map[int]PoolDuration{}
	PoolAvgDurations = map[int]time.Duration{}
	PoolProcessing   = map[int]int{}
	PoolQueue        = map[int]int{}
	PoolLastAdded    = map[int]time.Time{}
	QueueLastAdded   = map[int]time.Time{}
)

func init() {
	// Initialise the required Channel variables
	for i := 0; i <= config.LastPoolIota; i++ {
		// Pools[i] = make(chan *handlers.CommandTask)
		PoolDurations[i] = PoolDuration{}
		PoolAvgDurations[i] = 0
		PoolProcessing[i] = 0
		PoolQueue[i] = 0
		PoolLastAdded[i] = helpers.GetNullDateTime()
		QueueLastAdded[i] = helpers.GetNullDateTime()
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
