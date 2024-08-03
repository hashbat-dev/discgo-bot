package imagework

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/png"
	"io"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dabi-ngin/discgo-bot/Bot/helpers"

	"github.com/dabi-ngin/discgo-bot/Bot/audit"
	"github.com/nfnt/resize"
)

func DownloadImage(url string) (*image.Image, error) {
	start := time.Now()
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	img, _, err := image.Decode(response.Body)
	if err != nil {
		return nil, err
	}
	elapsed := time.Since(start)
	audit.Log(fmt.Sprintf("DownloadImage - Downloading %s took %dms\n", url, elapsed.Milliseconds()))

	return &img, nil
}

func ResizeImageGif(gifImg *gif.GIF, width uint) (io.Reader, error) {
	start := time.Now()
	// Resize each frame
	var aspectRatio float64
	var height uint = 0
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

	audit.Log("New Size: " + fmt.Sprint(gifImg.Config.Width) + "w x " + fmt.Sprint(gifImg.Config.Height) + "h")
	// Encode the resized frames back to GIF format
	var buf bytes.Buffer
	err := gif.EncodeAll(&buf, gifImg)
	if err != nil {
		return nil, err
	}

	end := time.Since(start)
	audit.Log(fmt.Sprintf("ResizeImageGif - Time Elapsed: %dms\n", end.Milliseconds()))
	return &buf, nil
}

func ResizeImageStatic(img image.Image, width uint) (io.Reader, error) {
	start := time.Now()
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
		return nil, err
	}

	end := time.Since(start)
	audit.Log(fmt.Sprintf("ResizeImageStatic - Time Elapsed: %dms\n", end.Milliseconds()))
	return &buf, nil
}

func ClosestColorIndex(p color.Palette, target color.Color) int {
	minDistance := math.MaxFloat64
	index := 0

	for i, c := range p {
		distance := ColourDistance(c, target)
		if distance < minDistance {
			minDistance = distance
			index = i
		}
	}

	return index
}

func ColourDistance(c1, c2 color.Color) float64 {
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()

	dr := float64(r1>>8) - float64(r2>>8)
	dg := float64(g1>>8) - float64(g2>>8)
	db := float64(b1>>8) - float64(b2>>8)

	return math.Sqrt(dr*dr + dg*dg + db*db)
}

func GetImageURLAndExtension(message *discordgo.MessageCreate) (string, string, error) {
	// Have they replied to a Message?
	if message.ReferencedMessage == nil {
		err := errors.New("you didn't reply to a message")
		return "", "", err
	}

	imgLink := ""

	// Is this a Tenor Link?
	if strings.Contains(message.ReferencedMessage.Content, "tenor.com/") {
		tenorLink, err := helpers.GetGifURLFromTenorLink(message.ReferencedMessage.Content)
		if err != nil {
			audit.Error(err)
		} else if tenorLink != "" {
			imgLink = tenorLink
		}
	}

	// Did that message contain an embed?
	if imgLink == "" {
		if len(message.ReferencedMessage.Embeds) > 0 {
			imgLink = message.ReferencedMessage.Embeds[0].Thumbnail.ProxyURL
		} else if len(message.ReferencedMessage.Attachments) > 0 {
			imgLink = message.ReferencedMessage.Attachments[0].ProxyURL
		}
	}

	if imgLink == "" {
		err := errors.New("that message didn't have any embeds or attachments")
		return "", "", err
	}

	// Return valid extensions
	if strings.Contains(imgLink, ".gif") {
		return imgLink, ".gif", nil
	} else if strings.Contains(imgLink, ".png") {
		return imgLink, ".png", nil
	} else if strings.Contains(imgLink, ".jpg") {
		return imgLink, ".jpg", nil
	} else if strings.Contains(imgLink, ".webp") {
		return imgLink, ".webp", nil
	}

	err := errors.New("image wasn't a supported filetype")
	return "", "", err
}
