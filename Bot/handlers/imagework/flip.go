package imagework

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/png"
	"io"
	"sync"
	"time"

	"github.com/dabi-ngin/discgo-bot/Bot/audit"
)

func FlipImage(imageReader io.Reader, isGif bool, buffer *bytes.Buffer, flipDirection string) error {
	start := time.Now()
	var err error

	flipDir := map[string]string{
		"flipleft":  "left",
		"flipup":    "up",
		"flipright": "right",
		"flipdown":  "down",
	}[flipDirection]
	if isGif {
		*buffer, err = FlipImageGif(imageReader, flipDir)
	} else {
		*buffer, err = FlipImageStatic(imageReader, flipDir)
	}
	if err != nil {
		audit.Error(err)
		return err
	}

	elapsed := time.Since(start)
	audit.Log(fmt.Sprintf("FlipImage - Time Elapsed: %dms", elapsed.Milliseconds()))
	return nil
}

func FlipImageGif(imageReader io.Reader, flipDirection string) (bytes.Buffer, error) {
	imageReaderGifObject, err := gif.DecodeAll(imageReader)
	if err != nil {
		return bytes.Buffer{}, err
	}

	mirroredGif := &gif.GIF{
		Image:           make([]*image.Paletted, len(imageReaderGifObject.Image)),
		Delay:           imageReaderGifObject.Delay,
		LoopCount:       imageReaderGifObject.LoopCount,
		Disposal:        imageReaderGifObject.Disposal,
		Config:          imageReaderGifObject.Config,
		BackgroundIndex: imageReaderGifObject.BackgroundIndex,
	}

	var wg sync.WaitGroup
	bounds := imageReaderGifObject.Image[0].Bounds()
	for i, frame := range imageReaderGifObject.Image {
		wg.Add(1)
		go func(flip string, x int, f *image.Paletted, b image.Rectangle) {
			defer wg.Done()
			mirroredGif.Image[x] = FlipGifFrame(f, flip, b)
		}(flipDirection, i, frame, bounds)
	}
	wg.Wait()

	for i := range mirroredGif.Disposal {
		mirroredGif.Disposal[i] = gif.DisposalNone
	}

	var buf bytes.Buffer
	err = gif.EncodeAll(&buf, mirroredGif)
	return buf, err
}

func FlipGifFrame(img *image.Paletted, flipDirection string, bounds image.Rectangle) *image.Paletted {
	width := bounds.Dx()
	height := bounds.Dy()

	newImg := image.NewPaletted(bounds, img.Palette)
	var wg sync.WaitGroup
	if flipDirection == "left" || flipDirection == "right" {
		for y := 0; y < height; y++ {
			wg.Add(1)
			if flipDirection == "left" {
				go func(y int) {
					defer wg.Done()
					for x := 0; x < width/2; x++ {
						leftColor := img.At(x, y)
						rightX := width - x - 1
						newImg.Set(x, y, leftColor)
						newImg.Set(rightX, y, leftColor)
					}
				}(y)
			} else {
				go func(y int) {
					defer wg.Done()
					for x := width / 2; x < width; x++ {
						rightColor := img.At(x, y)
						leftX := width - x - 1
						newImg.Set(leftX, y, rightColor)
						newImg.Set(x, y, rightColor)
					}
				}(y)
			}
		}
	} else {
		for x := 0; x < width; x++ {
			wg.Add(1)
			if flipDirection == "up" {
				go func(x int) {
					defer wg.Done()
					for y := 0; y < height/2; y++ {
						leftColor := img.At(x, y)
						downY := height - y - 1
						newImg.Set(x, y, leftColor)
						newImg.Set(x, downY, leftColor)
					}
				}(x)
			} else {
				go func(x int) {
					defer wg.Done()
					for y := height / 2; y < height; y++ {
						rightColor := img.At(x, y)
						upY := height - y - 1
						newImg.Set(x, upY, rightColor)
						newImg.Set(x, y, rightColor)
					}
				}(x)
			}
		}
	}
	wg.Wait()
	return newImg
}

func FlipImageStatic(imageReader io.Reader, flipDirection string) (bytes.Buffer, error) {
	img, _, err := image.Decode(imageReader)
	if err != nil {
		return bytes.Buffer{}, err
	}

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	newImg := image.NewRGBA(bounds)

	if flipDirection == "left" || flipDirection == "right" {
		if flipDirection == "left" {
			// Mirror left side to the right
			for y := 0; y < height; y++ {
				for x := 0; x <= width/2; x++ {
					newImg.Set(x, y, img.At(x, y))
					if x != width/2 || width%2 == 0 {
						newImg.Set(width-x-1, y, img.At(x, y))
					}
				}
			}

			// If the width is odd, shift left-hand pixels right by one pixel
			if width%2 != 0 {
				for y := 0; y < height; y++ {
					for x := width/2 - 1; x >= 0; x-- {
						newImg.Set(x+1, y, newImg.At(x, y))
					}
					newImg.Set(0, y, image.Transparent)
				}
			}
		} else {
			// Mirror right side to the left
			for y := 0; y < height; y++ {
				for x := width / 2; x < width; x++ {
					newImg.Set(x, y, img.At(x, y))
					newImg.Set(width-x-1, y, img.At(x, y))
				}
			}

			// If the width is odd, shift right-hand pixels left by one pixel
			if width%2 != 0 {
				for y := 0; y < height; y++ {
					for x := width/2 + 1; x < width; x++ {
						newImg.Set(x-1, y, newImg.At(x, y))
					}
					newImg.Set(width-1, y, image.Transparent)
				}
			}
		}
	} else {
		if flipDirection == "up" {
			// Mirror left side to the right
			for x := 0; x < width; x++ {
				for y := 0; y <= height/2; y++ {
					newImg.Set(x, y, img.At(x, y))
					if x != height/2 || height%2 == 0 {
						newImg.Set(x, height-y-1, img.At(x, y))
					}
				}
			}

			// If the height is odd, shift left-hand pixels right by one pixel
			if height%2 != 0 {
				for x := 0; x < width; x++ {
					for y := height/2 - 1; y >= 0; y-- {
						newImg.Set(x+1, y, newImg.At(x, y))
					}
					newImg.Set(x, 0, image.Transparent)
				}
			}
		} else {
			// Mirror right side to the left
			for x := 0; x < width; x++ {
				for y := height / 2; y < height; y++ {
					newImg.Set(x, y, img.At(x, y))
					newImg.Set(x, height-y-1, img.At(x, y))
				}
			}

			// If the width is odd, shift right-hand pixels left by one pixel
			if height%2 != 0 {
				for x := 0; x < width; x++ {
					for y := height/2 + 1; y < height; y++ {
						newImg.Set(x-1, y, newImg.At(x, y))
					}
					newImg.Set(x, height-1, image.Transparent)
				}
			}
		}
	}

	var buf bytes.Buffer
	err = png.Encode(&buf, newImg)
	return buf, err
}
