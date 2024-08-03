package imagework

import (
	"bytes"
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"io"
	"sync"
	"time"

	"github.com/ZestHusky/femboy-control/Bot/audit"
	"github.com/disintegration/imaging"
)

func DeepFryThatShit(imageReader io.Reader, isGif bool, buffer *bytes.Buffer) error {
	var err error
	start := time.Now()
	audit.Log("")
	if isGif {
		*buffer, err = DeepFryGIF(imageReader)
	} else {
		*buffer, err = DeepFryIMG(imageReader)
	}
	if err != nil {
		audit.Error(err)
		return err
	}
	elapsed := time.Since(start)
	audit.Log(fmt.Sprintf("DeepFrythatShit - Elapsed:%dms", elapsed.Milliseconds()))
	return nil
}

func DeepFryIMG(imageReader io.Reader) (bytes.Buffer, error) {
	var contrastOffset float64 = 85
	var brightnessOffset float64 = -20
	var SIGMABALLS float64 = 5

	img, _, err := image.Decode(imageReader)
	if err != nil {
		return bytes.Buffer{}, err
	}

	audit.Log("DeepFry - Image Decoded")

	img = imaging.AdjustContrast(img, contrastOffset)
	img = imaging.AdjustBrightness(img, brightnessOffset)
	img = imaging.Sharpen(img, SIGMABALLS)

	var buf bytes.Buffer
	err = png.Encode(&buf, img)

	audit.Log("DeepFry - PNG Encoded")

	return buf, err
}

func DeepFryGIF(imageReader io.Reader) (bytes.Buffer, error) {
	var contrastOffset float64 = 55
	var brightnessOffset float64 = -5
	var SIGMABALLS float64 = 3.5

	audit.Log("DeepFry [Gif] - Started")

	img, err := gif.DecodeAll(imageReader)
	if err != nil {
		return bytes.Buffer{}, err
	}

	audit.Log("DeepFry [Gif] - Decoded Image")

	var wg sync.WaitGroup
	for i := range img.Image {
		frame := img.Image[i]

		wg.Add(1)
		go func(gi int, gframe image.Image, gcontrastOffset float64, gbrightnessOffset float64, gSigma float64) {

			defer wg.Done()

			audit.Log(fmt.Sprintf("DeepFry [Gif] - Start Frame [%v]", gi))
			gbounds := frame.Bounds()

			// Convert GIF frame to RGBA
			rgba := image.NewRGBA(gbounds)
			draw.Draw(rgba, gbounds, gframe, gbounds.Min, draw.Src)

			// Apply effects
			nrgba := imaging.AdjustContrast(rgba, gcontrastOffset)
			nrgba = imaging.AdjustBrightness(nrgba, gbrightnessOffset)
			nrgba = imaging.Sharpen(nrgba, gSigma)

			// Convert *image.NRGBA to *image.RGBA
			audit.Log(fmt.Sprintf("DeepFry [Gif] - Start Frame nrgbaToRGBA [%v]", gi))
			rgba = nrgbaToRGBA(nrgba)
			audit.Log(fmt.Sprintf("DeepFry [Gif] - Finish Frame nrgbaToRGBA [%v]", gi))

			// Convert back to GIF frame
			palettedImage := image.NewPaletted(gbounds, palette.Plan9)
			draw.FloydSteinberg.Draw(palettedImage, gbounds, rgba, image.Point{})
			img.Image[gi] = palettedImage

			audit.Log(fmt.Sprintf("DeepFry [Gif] - Finish Frame [%v]", gi))

		}(i, frame, contrastOffset, brightnessOffset, SIGMABALLS)

	}

	wg.Wait()
	audit.Log("DeepFry [Gif] - Frames Finished")

	// Encode the modified GIF
	var buf bytes.Buffer
	err = gif.EncodeAll(&buf, img)
	if err != nil {
		return bytes.Buffer{}, err
	}

	return buf, nil
}

func nrgbaToRGBA(nrgba *image.NRGBA) *image.RGBA {
	bounds := nrgba.Bounds()
	rgba := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba.Set(x, y, nrgba.At(x, y))
		}
	}
	return rgba
}
