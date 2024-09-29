package scheduler

import tempfiles "github.com/hashbat-dev/discgo-bot/TempFiles"

func RunEvery1Minute() {
	go tempfiles.DeleteAllExpired()
}
