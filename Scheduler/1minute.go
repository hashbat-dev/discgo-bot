package scheduler

import tempfiles "github.com/dabi-ngin/discgo-bot/TempFiles"

func RunEvery1Minute() {
	go tempfiles.DeleteAllExpired()
}
