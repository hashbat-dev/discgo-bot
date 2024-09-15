package imgwork

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"

	"github.com/chai2010/webp"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
)

// DownloadImageToReader takes an imageUrl as found in a discord message, downloads from the CDN, and returns it as
// a single-use io.Reader object. Once read this will be empty. If intending to use multiple times, this should be
// consumed by being read into a bytes.Buffer object.
func DownloadImageToReader(guildId string, imageUrl string, isAnimated bool) (io.Reader, error) {
	// 1. Download the source image
	var downloadErr error
	var downloadedStaticImg *image.Image
	var downloadedGifImg *gif.GIF

	if isAnimated {
		downloadedGifImg, downloadErr = downloadGif(guildId, imageUrl)
	} else {
		downloadedStaticImg, downloadErr = downloadImage(guildId, imageUrl)
	}
	if downloadErr != nil {
		return nil, downloadErr
	}

	imageExtension := GetExtensionFromURL(imageUrl)
	if imageExtension == "" {
		err := errors.New("invalid image extension")
		logger.Error(guildId, err)
		return nil, err
	}

	var imageReader io.Reader
	var buf bytes.Buffer
	if isAnimated {
		encodeErr := gif.EncodeAll(&buf, downloadedGifImg)
		if encodeErr != nil {
			return nil, encodeErr
		}
		imageReader = &buf
	} else {
		switch imageExtension {
		case ".png":
			encodeErr := png.Encode(&buf, *downloadedStaticImg)
			if encodeErr != nil {
				return nil, encodeErr
			}
		case ".jpg":
			encodeErr := jpeg.Encode(&buf, *downloadedStaticImg, nil)
			if encodeErr != nil {
				return nil, encodeErr
			}
		case ".webp":
			webpBytes, encodeErr := webp.EncodeRGBA(*downloadedStaticImg, 100.0)
			if encodeErr != nil {
				return nil, encodeErr
			}
			buf = *bytes.NewBuffer(webpBytes)
		default:
			return nil, errors.New("unsupported file extension of " + imageExtension)
		}
		imageReader = &buf
	}
	return imageReader, nil
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
