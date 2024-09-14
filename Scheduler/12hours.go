package scheduler

import (
	fakeyou "github.com/dabi-ngin/discgo-bot/External/FakeYou"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
)

func RunEvery12Hours() {
	logger.Debug("SCHEDULER", "Run every 12 hours started")
	go fakeyou.UpdateModels()
}
