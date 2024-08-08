package scheduler

import (
	"time"
)

func Init() {

	nextRun2Seconds = time.Now()
	nextRun1Minute = time.Now()
	nextRun5Minutes = time.Now()

	stopLoop := make(chan bool)
	Loop(stopLoop)
}

var nextRun2Seconds time.Time
var nextRun1Minute time.Time
var nextRun5Minutes time.Time
var pollingSeconds int = 2

func Loop(stop chan bool) {
	for {
		select {
		case <-stop:
			// Should this break out of the loop??
		case <-time.After(time.Duration(pollingSeconds) * time.Second):
			if time.Now().After(nextRun2Seconds) {
				RunEvery2Seconds()
			}

			if time.Now().After(nextRun1Minute) {
				RunEvery1Minute()
			}

			if time.Now().After(nextRun5Minutes) {
				RunEvery5Minutes()
			}

		}
	}
}
