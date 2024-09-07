package scheduler

import (
	config "github.com/dabi-ngin/discgo-bot/Config"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
	reporting "github.com/dabi-ngin/discgo-bot/Reporting"
)

func RunEvery2Seconds() {
	pollHardwareStats()
	sendNextDiscordLogBatch()
	reporting.Logs()
}

func RunEvery5Seconds() {
	reporting.Guilds()
}

func sendNextDiscordLogBatch() {
	if !logger.InitComplete {
		return
	}

	var newCache []string
	nextMsg := ""

	limitReached := false

	for _, s := range logger.LogsForDiscord {
		if (len(nextMsg) + len(s)) >= config.MAX_MESSAGE_LENGTH {
			limitReached = true
		}

		if limitReached {
			newCache = append(newCache, s)
		} else {
			nextMsg += s
		}
	}

	logger.LogsForDiscord = newCache
	if nextMsg != "" {
		sendLogsToDiscordChannel(nextMsg)
	}
}

func sendLogsToDiscordChannel(logs string) {
	_, err := config.Session.ChannelMessageSend(config.LoggingChannelID, logs)
	if err != nil {
		logger.Error("", err)
	}

}

func pollHardwareStats() {
	reporting.Hardware()
}
