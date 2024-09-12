package bang

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
	"os"
	"path/filepath"
	"time"

	"github.com/bwmarrin/discordgo"
	config "github.com/dabi-ngin/discgo-bot/Config"
	discord "github.com/dabi-ngin/discgo-bot/Discord"
	helpers "github.com/dabi-ngin/discgo-bot/Helpers"
	imgwork "github.com/dabi-ngin/discgo-bot/ImgWork"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
	"github.com/google/uuid"
	"golang.org/x/image/webp"
)

type MakeSpeech struct{}

func (s MakeSpeech) Name() string {
	return "makespeech"
}

func (s MakeSpeech) PermissionRequirement() int {
	return config.CommandLevelUser
}

func (s MakeSpeech) Complexity() int {
	return config.CPU_BOUND_TASK
}

func (s MakeSpeech) Execute(message *discordgo.MessageCreate, command string) error {
	// 1. Check we have a valid Image and Extension
	imgUrl := helpers.GetImageFromMessage(message.Message, "")
	if imgUrl == "" {
		return errors.New("no image found")
	}

	imgExtension := imgwork.GetExtensionFromURL(imgUrl)
	if imgExtension == "" {
		return errors.New("invalid extension")
	}

	// 2. Generate the Output Image name
	//	  This will always be either a .gif (for animated) or .png (for static)
	outputImageName := uuid.New().String()
	isAnimated := false
	if imgExtension == ".gif" {
		outputImageName += ".gif"
		isAnimated = true
	} else {
		outputImageName += ".png"
	}

	// 3. Get the image as an io.Reader object
	imageReader, newHeight, err := imgwork.DownloadImageToReader(message.GuildID, imgUrl, isAnimated, 300)
	if err != nil {
		return err
	}

	// 4. Write the new Image to a Bytes Buffer
	var newImageBuffer bytes.Buffer
	err = addSpeechBubbleToImage(message.GuildID, imageReader, &newImageBuffer, newHeight, isAnimated, imgExtension)
	if err != nil {
		return err
	}

	// 5. Send the new Image back to the User
	return discord.ReplyToMessageWithImageBuffer(message, true, outputImageName, &newImageBuffer)
}

func addSpeechBubbleToImage(guildId string, imageReader io.Reader, newImgBuffer *bytes.Buffer, imageHeight int, isAnimated bool, imgExtension string) error {
	start := time.Now()
	// 1. Get our Speech Bubble overlay template
	// 1A. Get the filepath to the most appropriate overlay speech bubble image.
	//	   We used smaller speech bubbles for smaller heights to not obscure the image.
	rootDir, err := filepath.Abs(filepath.Dir("."))
	if err != nil {
		logger.Error(guildId, err)
		return err
	}

	overlayPath := "Resources/SpeechTemplates/"
	if imageHeight < 100 {
		overlayPath += "S"
	} else if imageHeight < 200 {
		overlayPath += "M"
	} else {
		overlayPath += "L"
	}
	overlayPath += ".png"
	templateFilePath := filepath.Join(rootDir, overlayPath)

	// 2. Open the Speech Bubble overlay template
	overlayFile, err := os.Open(templateFilePath)
	if err != nil {
		logger.Error(guildId, err)
		return err
	}
	defer overlayFile.Close()

	overlayImage, err := png.Decode(overlayFile)
	if err != nil {
		logger.Error(guildId, err)
		return err
	}

	// Define the Transparent colour, GIFs are wacky so we can't always add a new
	// colour to set as transparent. This colour is Discord's BG, a good backup.
	transparentColor := color.RGBA{49, 51, 56, 0}

	if isAnimated {
		// GIFs ----------------
		// Decode the file
		gifImage, err := gif.DecodeAll(imageReader)
		if err != nil {
			logger.Error(guildId, err)
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

		backgroundIndex := closestColorIndex(gifImage.Image[0].Palette, transparentColor)
		gifImage.BackgroundIndex = byte(backgroundIndex)

		// Encode the modified GIF to a buffer
		err = gif.EncodeAll(newImgBuffer, gifImage)
		if err != nil {
			logger.Error(guildId, err)
			return err
		}

	} else {
		// Static Images --------
		// Read and decode the input PNG image
		var inputImage image.Image
		if imgExtension == ".png" || imgExtension == ".jpg" {
			// .jpg images are resized to a .png before this function is called
			inputImage, err = png.Decode(imageReader)
		} else if imgExtension == ".webp" {
			inputImage, err = webp.Decode(imageReader)
		}
		if err != nil {
			logger.Error(guildId, err)
			return err
		}
		// Ensure both images have the same dimensions
		templateBounds := overlayImage.Bounds()
		inputBounds := inputImage.Bounds()
		if templateBounds.Dx() != inputBounds.Dx() {
			err = errors.New("template and input dimensions do not match")
			logger.Error(guildId, err)
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
		err = png.Encode(newImgBuffer, outputImg)
		if err != nil {
			err = fmt.Errorf("error encoding modified PNG: %w", err)
			logger.Error(guildId, err)
			return err
		}
	}

	logger.Info(guildId, fmt.Sprintf("addSpeechBubbleToImage completed (%v) after %v", imgExtension, time.Since(start)))
	return nil
}

func closestColorIndex(p color.Palette, target color.Color) int {
	minDistance := math.MaxFloat64
	index := 0

	for i, c := range p {
		distance := colourDistance(c, target)
		if distance < minDistance {
			minDistance = distance
			index = i
		}
	}

	return index
}

func colourDistance(c1, c2 color.Color) float64 {
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()

	dr := float64(r1>>8) - float64(r2>>8)
	dg := float64(g1>>8) - float64(g2>>8)
	db := float64(b1>>8) - float64(b2>>8)

	return math.Sqrt(dr*dr + dg*dg + db*db)
}
