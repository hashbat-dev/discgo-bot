package imagework

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"io"
	"math"
	"time"

	"github.com/dabi-ngin/discgo-bot/Bot/audit"
)

func ChangeSpeedGif(resizedImageReader io.Reader, buffer *bytes.Buffer, speedUp bool) error {
	start := time.Now()
	// 1. Decode the file into a GIF object
	gifImage, err := gif.DecodeAll(resizedImageReader)
	if err != nil {
		audit.Error(err)
		return err
	}

	// 2. Lower the Delay of the frames
	outGif := &gif.GIF{}
	var snippedDelay []int
	alreadySlowestCount := 0
	for i := range gifImage.Image {
		frameDelay := gifImage.Delay[i]
		if frameDelay <= 2 {
			alreadySlowestCount++
		}
		newValue := 1
		if speedUp {
			newValue = int(math.Round(float64(frameDelay) / 4.0))
		} else {
			newValue = int(math.Round(float64(frameDelay) * 2.0))
		}
		if newValue < 2 {
			newValue = 2
		}
		snippedDelay = append(snippedDelay, newValue)
	}

	// 3. If it's already as slow as possible due to the delay, cut frames out
	var newPaletted []*image.Paletted
	var newDelay []int
	var newDisposal []byte
	if alreadySlowestCount == len(gifImage.Delay) {
		keptLast := false
		for i, frame := range gifImage.Image {
			if !keptLast {
				newPaletted = append(newPaletted, frame)
				newDelay = append(newDelay, snippedDelay[i])
				newDisposal = append(newDisposal, gifImage.Disposal[i])
				keptLast = true
			} else {
				keptLast = false
			}
		}
	} else {
		newPaletted = gifImage.Image
		newDelay = snippedDelay
		newDisposal = gifImage.Disposal
	}

	outGif.BackgroundIndex = gifImage.BackgroundIndex
	outGif.Config = gifImage.Config
	outGif.Disposal = newDisposal
	outGif.LoopCount = gifImage.LoopCount
	outGif.Image = newPaletted
	outGif.Delay = newDelay

	// 3. Write the new GIF to buffer arg
	err = gif.EncodeAll(buffer, outGif)
	if err != nil {
		audit.Error(err)
		return err
	}

	elapsed := time.Since(start)
	audit.Log(fmt.Sprintf("ChangeSpeedGif - elapsed: %dms", elapsed.Milliseconds()))
	return nil
}
