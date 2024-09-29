package scheduler

import (
	fakeyou "github.com/hashbat-dev/discgo-bot/External/FakeYou"
	imgur "github.com/hashbat-dev/discgo-bot/External/Imgur"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

func RunEvery12Hours() {
	logger.Debug("SCHEDULER", "Run every 12 hours started")
	go fakeyou.UpdateModels()
	go imgur.TidySubmissions()
}
