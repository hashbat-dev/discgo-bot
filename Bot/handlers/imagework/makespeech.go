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
	"os"
	"path/filepath"
	"time"

	"github.com/dabi-ngin/discgo-bot/Bot/audit"
	"golang.org/x/image/webp"
)

func AddSpeechBubbleToImage(resizedImageReader io.Reader, buffer *bytes.Buffer, imageHeight int, isGif bool, imageExtension string) error {
	start := time.Now()
	// 1. Get our Speech Bubble overlay template
	// 1A. Get the filepath to the most appropriate overlay speech bubble image.
	//	   We used smaller speech bubbles for smaller heights to not obscure the image.
	rootDir, err := filepath.Abs(filepath.Dir("."))
	if err != nil {
		audit.Error(err)
		return err
	}

	overlayPath := "Bot/handlers/imagework/speechtemplates/template"
	if imageHeight < 100 {
		overlayPath += "S"
	} else if imageHeight < 200 {
		overlayPath += "M"
	} else {
		overlayPath += "L"
	}
	overlayPath += ".png"
	templateFilePath := filepath.Join(rootDir, overlayPath)

	audit.Log("Height is " + fmt.Sprint(imageHeight) + "px, using template file: " + overlayPath)

	// 2. Open the Speech Bubble overlay template
	overlayFile, err := os.Open(templateFilePath)
	if err != nil {
		audit.Error(err)
		return err
	}
	defer overlayFile.Close()

	overlayImage, err := png.Decode(overlayFile)
	if err != nil {
		audit.Error(err)
		return err
	}

	// Define the Transparent colour, GIFs are wacky so we can't always add a new
	// colour to set as transparent. This colour is Discord's BG, a good backup.
	transparentColor := color.RGBA{49, 51, 56, 0}

	if isGif {
		// GIFs ----------------
		// Decode the file
		gifImage, err := gif.DecodeAll(resizedImageReader)
		if err != nil {
			audit.Error(err)
			return err
		}

		for _, frame := range gifImage.Image {
			// Loop over each pixel in the frame
			for y := frame.Bounds().Min.Y; y < frame.Bounds().Max.Y; y++ {
				for x := frame.Bounds().Min.X; x < frame.Bounds().Max.X; x++ {
					// If the pixel in the PNG is black, set the corresponding pixel in the frame to transparent
					r, g, b, _ := overlayImage.At(x, y).RGBA()
					if r == 0 && g == 0 && b == 0 {
						frame.Set(x, y, transparentColor) // Set pixel to transparent color
					}
				}
			}
		}

		backgroundIndex := ClosestColorIndex(gifImage.Image[0].Palette, transparentColor)
		gifImage.BackgroundIndex = byte(backgroundIndex)

		// Encode the modified GIF to a buffer
		err = gif.EncodeAll(buffer, gifImage)
		if err != nil {
			audit.Error(err)
			return err
		}

	} else {
		// Static Images --------
		// Read and decode the input PNG image
		var inputImage image.Image
		if imageExtension == ".png" || imageExtension == ".jpg" {
			// .jpg images are resized to a .png before this function is called
			inputImage, err = png.Decode(resizedImageReader)
		} else if imageExtension == ".webp" {
			inputImage, err = webp.Decode(resizedImageReader)
		}
		if err != nil {
			audit.Error(err)
			return err
		}
		// Ensure both images have the same dimensions
		templateBounds := overlayImage.Bounds()
		inputBounds := inputImage.Bounds()
		if templateBounds.Dx() != inputBounds.Dx() {
			err = errors.New("template and input dimensions do not match")
			audit.Error(err)
			return err
		}

		// Create a new image with the same dimensions as the input image
		outputImg := image.NewRGBA(inputBounds)
		draw.Draw(outputImg, inputBounds, inputImage, image.Point{}, draw.Src)

		// Loop over each pixel in the input image and set corresponding pixel in the output image to transparent if condition is met
		for y := inputBounds.Min.Y; y < inputBounds.Max.Y; y++ {
			for x := inputBounds.Min.X; x < inputBounds.Max.X; x++ {
				// If the pixel in the template PNG is black, set the corresponding pixel in the output image to transparent
				r, g, b, _ := overlayImage.At(x, y).RGBA()
				if r == 0 && g == 0 && b == 0 {
					outputImg.Set(x, y, transparentColor) // Set pixel to transparent color
				}
			}
		}

		// Encode the modified PNG to buffer
		err = png.Encode(buffer, outputImg)
		if err != nil {
			err = fmt.Errorf("error encoding modified PNG: %w", err)
			audit.Error(err)
			return err
		}
	}

	elapsed := time.Since(start)
	audit.Log(fmt.Sprintf("AddSpeechBubbleToImage - Elapsed %dms", elapsed.Milliseconds()))
	return nil
}
