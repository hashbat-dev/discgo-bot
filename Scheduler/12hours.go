package scheduler

import (
	wow "github.com/hashbat-dev/discgo-bot/Bot/Commands/Wow"
	database "github.com/hashbat-dev/discgo-bot/Database"
	imgur "github.com/hashbat-dev/discgo-bot/External/Imgur"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

func RunEvery12Hours() {
	logger.Debug("SCHEDULER", "Run every 12 hours started")
	go imgur.TidySubmissions()
	go ImageBankCheck()
	go wow.UpdatePokemonDatabase()
}

func ImageBankCheck() {
	err := database.TidyImgStorage("SCHEDULER")
	if err != nil {
		logger.ErrorText("SCHEDULER", "Error tidying image storage")
	}
	logger.Info("SCHEDULER", "Image Bank Check completed")
}
