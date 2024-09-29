package scheduler

import reporting "github.com/hashbat-dev/discgo-bot/Reporting"

func RunEvery5Seconds() {
	go reporting.Guilds()
}
