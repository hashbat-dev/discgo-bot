package imagework

import (
	"errors"
	"sync"
)

const WorkerCount = 4

var (
	requestChan = make(chan FlipRequest, 100)
	wg          sync.WaitGroup
)

func init() {
	for i := 0; i < WorkerCount; i++ {
		wg.Add(1)
		go worker()
	}
}

func worker() {
	defer wg.Done()
	for req := range requestChan {
		err := processRequest(req)
		req.ResponseChan <- err
		close(req.ResponseChan)
	}
}

func EnqueueFlipRequest(req FlipRequest) {
	requestChan <- req
}

func CloseWorkerPool() {
	close(requestChan)
	wg.Wait()
}

func processRequest(req FlipRequest) error {
	var err error
	switch req.RequestType {
	case "makespeech":
		err = AddSpeechBubbleToImage(req.ImageReader, req.ResponseBuffer, 300, req.IsGif, ".png")
	case "reverse":
		err = ReverseGif(req.ImageReader, req.ResponseBuffer)
	case "flipleft", "flipright", "flipup", "flipdown":
		err = FlipImage(req.ImageReader, req.IsGif, req.ResponseBuffer, req.RequestType)
	case "speedup":
		err = ChangeSpeedGif(req.ImageReader, req.ResponseBuffer, true)
	case "slowdown":
		err = ChangeSpeedGif(req.ImageReader, req.ResponseBuffer, false)
	case "deepfry":
		err = DeepFryThatShit(req.ImageReader, req.IsGif, req.ResponseBuffer)
	default:
		err = errors.New("unknown request type: " + req.RequestType)
	}
	return err
}
