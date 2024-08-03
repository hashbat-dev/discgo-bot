package scheduler

import (
	"time"
)

func Init() {

	nextRun2Seconds = time.Now()
	nextRun5Seconds = time.Now()
	nextRun1Minute = time.Now()
	nextRun5Minutes = time.Now()
	nextRun12Hours = time.Now()

	stopLoop := make(chan bool)
	Loop(stopLoop)
}

var nextRun2Seconds time.Time
var nextRun5Seconds time.Time
var nextRun1Minute time.Time
var nextRun5Minutes time.Time
var nextRun12Hours time.Time

var pollingSeconds int = 2

func Loop(stop chan bool) {
	for {
		select {
		case <-stop:
			// Should this break out of the loop??
		case <-time.After(time.Duration(pollingSeconds) * time.Second):
			if time.Now().After(nextRun2Seconds) {
				RunEvery2Seconds()
				nextRun2Seconds = time.Now().Add(time.Second * 2)
			}

			if time.Now().After(nextRun5Seconds) {
				RunEvery5Seconds()
				nextRun5Seconds = time.Now().Add(time.Second * 5)
			}

			if time.Now().After(nextRun5Seconds) {
				RunEvery5Seconds()
				nextRun5Seconds = time.Now().Add(time.Second * 5)
			}

			if time.Now().After(nextRun1Minute) {
				RunEvery1Minute()
				nextRun1Minute = time.Now().Add(time.Minute * 1)
			}

			if time.Now().After(nextRun5Minutes) {
				RunEvery5Minutes()
				nextRun5Minutes = time.Now().Add(time.Minute * 5)
			}

			if time.Now().After(nextRun12Hours) {
				RunEvery12Hours()
				nextRun12Hours = time.Now().Add(time.Hour * 12)
			}

		}
	}
}
