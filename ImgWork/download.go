package imgwork

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"time"

	"github.com/chai2010/webp"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
	"github.com/nfnt/resize"
)

func DownloadImageToReader(guildId string, imageUrl string, isAnimated bool, resizeToWidth int) (io.Reader, int, error) {

	// 1. Download the source image
	var err error
	var downloadedStaticImg *image.Image
	var downloadedGifImg *gif.GIF

	if isAnimated {
		downloadedGifImg, err = downloadGif(guildId, imageUrl)
	} else {
		downloadedStaticImg, err = downloadImage(guildId, imageUrl)
	}
	if err != nil {
		return nil, 0, err
	}

	imageExtension := GetExtensionFromURL(imageUrl)
	if imageExtension == "" {
		err = errors.New("invalid image extension")
		logger.Error(guildId, err)
		return nil, 0, err
	}

	var newHeight int = 0
	var imageReader io.Reader
	if resizeToWidth > 0 {
		// Resizing is needed
		if isAnimated {
			imageReader, newHeight, err = resizeGif(guildId, downloadedGifImg, uint(resizeToWidth))
		} else {
			imageReader, newHeight, err = resizeImage(guildId, *downloadedStaticImg, uint(resizeToWidth))
		}
	} else {
		// Write the downloaded file as is to the io.Reader
		var buf bytes.Buffer
		if isAnimated {
			err = gif.EncodeAll(&buf, downloadedGifImg)
			if err == nil {
				imageReader = &buf
			}
		} else {
			switch imageExtension {
			case ".png":
				err = png.Encode(&buf, *downloadedStaticImg)
			case ".jpg":
				err = jpeg.Encode(&buf, *downloadedStaticImg, nil)
			case ".webp":
				webpBytes, errWebp := webp.EncodeRGBA(*downloadedStaticImg, 100.0)
				if errWebp != nil {
					err = errWebp
				} else {
					buf = *bytes.NewBuffer(webpBytes)
				}
			default:
				err = errors.New("unsupported file extension of " + imageExtension)
			}
			if err == nil {
				imageReader = &buf
			}
		}
	}

	if err != nil {
		logger.Error(guildId, err)
		return nil, 0, err
	} else {
		return imageReader, newHeight, nil
	}
}

func downloadGif(guildId string, gifLink string) (*gif.GIF, error) {
	// Download the GIF
	response, err := http.Get(gifLink)
	if err != nil {
		logger.Error(guildId, err)
		return nil, err
	}
	defer response.Body.Close()

	// Check if the download was successful
	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("failed to download GIF: received status code %d", response.StatusCode)
		logger.Error(guildId, err)
		return nil, err
	}

	// Decode the GIF directly from the response body
	gifData, err := gif.DecodeAll(response.Body)
	if err != nil {
		logger.Error(guildId, err)
		return nil, fmt.Errorf("failed to decode GIF: %w", err)
	}

	logger.Debug(guildId, "Downloaded animated file: %v", gifLink)
	return gifData, nil
}

func downloadImage(guildId string, url string) (*image.Image, error) {
	response, err := http.Get(url)
	if err != nil {
		logger.Error(guildId, err)
		return nil, err
	}
	defer response.Body.Close()

	img, _, err := image.Decode(response.Body)
	if err != nil {
		logger.Error(guildId, err)
		return nil, err
	}

	logger.Debug(guildId, "Downloaded static file: %v", url)
	return &img, nil
}

func resizeGif(guildId string, gifImg *gif.GIF, width uint) (io.Reader, int, error) {
	// Resize each frame
	startTime := time.Now()
	var aspectRatio float64
	var height uint
	for i, frame := range gifImg.Image {
		bounds := frame.Bounds()
		aspectRatio = float64(bounds.Dy()) / float64(bounds.Dx())
		height = uint(float64(width) * aspectRatio)
		resizedFrame := resize.Resize(width, height, frame, resize.Lanczos3)

		// Convert the resized frame to *image.Paletted
		palettedFrame := image.NewPaletted(resizedFrame.Bounds(), gifImg.Image[i].Palette)

		// Draw the resized frame onto the paletted frame
		draw.FloydSteinberg.Draw(palettedFrame, palettedFrame.Rect, resizedFrame, image.Point{})

		// Update the frame in the GIF
		gifImg.Image[i] = palettedFrame
	}

	// Update the GIF configuration
	gifImg.Config.Width = int(width)
	gifImg.Config.Height = int(float64(width) * aspectRatio)

	// Encode the resized frames back to GIF format
	var buf bytes.Buffer
	err := gif.EncodeAll(&buf, gifImg)
	if err != nil {
		logger.Error(guildId, err)
		return nil, gifImg.Config.Height, err
	}

	logger.Info(guildId, "Resized GIF to %vpx width after %v", width, time.Since(startTime))
	return &buf, gifImg.Config.Height, nil
}

func resizeImage(guildId string, img image.Image, width uint) (io.Reader, int, error) {
	// Calculate the new height maintaining the aspect ratio
	bounds := img.Bounds()
	originalWidth := bounds.Dx()
	originalHeight := bounds.Dy()
	newHeight := uint(float64(originalHeight) * (float64(width) / float64(originalWidth)))

	// Resize the image using the resize package
	resizedImg := resize.Resize(width, newHeight, img, resize.Lanczos3)

	// Create a bytes buffer to write the PNG image to
	var buf bytes.Buffer
	err := png.Encode(&buf, resizedImg)
	if err != nil {
		logger.Error(guildId, err)
		return nil, int(newHeight), err
	}

	logger.Debug(guildId, "Resized static image to %vpx width", width)
	return &buf, int(newHeight), nil
}
