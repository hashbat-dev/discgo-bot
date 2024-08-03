package imagework

import (
	"bytes"
	"fmt"
	"image/gif"
	"io"
	"time"

	"github.com/dabi-ngin/discgo-bot/Bot/audit"
)

func ReverseGif(resizedImageReader io.Reader, buffer *bytes.Buffer) error {
	start := time.Now()
	// 1. Decode the file into a GIF object
	gifImage, err := gif.DecodeAll(resizedImageReader)
	if err != nil {
		audit.Error(err)
		return err
	}

	// 2. Add all the frames to a new GIF object in reverse order
	outGif := &gif.GIF{}
	for i := len(gifImage.Image) - 1; i >= 0; i-- {
		outGif.Image = append(outGif.Image, gifImage.Image[i])
		outGif.Delay = append(outGif.Delay, gifImage.Delay[i])
	}
	outGif.BackgroundIndex = gifImage.BackgroundIndex
	outGif.Config = gifImage.Config
	outGif.Disposal = gifImage.Disposal
	outGif.LoopCount = gifImage.LoopCount

	// 3. Write the new GIF to the Return buffer
	err = gif.EncodeAll(buffer, outGif)
	if err != nil {
		audit.Error(err)
		return err
	}

	elapsed := time.Since(start)
	audit.Log(fmt.Sprintf("ReverseGif - elapsed %dms", elapsed.Milliseconds()))
	return nil
}
