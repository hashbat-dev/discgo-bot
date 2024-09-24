package scheduler

import reporting "github.com/dabi-ngin/discgo-bot/Reporting"

func RunEvery5Seconds() {
	go reporting.Guilds()
}
