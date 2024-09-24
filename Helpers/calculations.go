package helpers

import "time"

func AverageDuration(durations []time.Duration) time.Duration {
	var total time.Duration

	for _, duration := range durations {
		total += duration
	}

	if len(durations) == 0 {
		return 0
	}

	average := total / time.Duration(len(durations))
	return average
}
