package dashboard

import "time"

type LatencyMessage struct {
	CommandName string
	TimeTaken   time.Duration
}

var LatencyCh chan LatencyMessage

func init() {
	LatencyCh = make(chan LatencyMessage, 50)
	go latencyChWorker()
}

// latencyChWorker is spun up as a goroutine to read from LatencyCh, a buffered channel,
// which has its contents written to from the CommandWorkers in Handlers, tracking
// how long operations have run for
func latencyChWorker() {
	for {
		select {
		case msg, ok := <-LatencyCh:
			if !ok {
				// channel closed
				return
			}
			latencies, latenciesExist := CommandLatencies[msg.CommandName]
			if !latenciesExist {
				CommandLatencies[msg.CommandName] = []time.Duration{msg.TimeTaken}
			} else {
				latencies = append(latencies, msg.TimeTaken)
			}
		}
	}
}
