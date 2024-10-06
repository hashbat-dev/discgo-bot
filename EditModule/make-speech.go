package editmodule

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	cache "github.com/hashbat-dev/discgo-bot/Cache"
	config "github.com/hashbat-dev/discgo-bot/Config"
	discord "github.com/hashbat-dev/discgo-bot/Discord"
	helpers "github.com/hashbat-dev/discgo-bot/Helpers"
	imgwork "github.com/hashbat-dev/discgo-bot/ImgWork"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
	"golang.org/x/image/webp"
)

type MakeSpeech struct{}

func (s MakeSpeech) SelectName() string {
	return "Add Speech Bubble"
}

func (s MakeSpeech) Emoji() *discordgo.ComponentEmoji {
	return &discordgo.ComponentEmoji{Name: "💬"}
}

func (s MakeSpeech) PermissionRequirement() int {
	return config.CommandLevelUser
}

func (s MakeSpeech) Complexity() int {
	return config.TRIVIAL_TASK
}

func (s MakeSpeech) Execute(i *discordgo.InteractionCreate, correlationId string) {
	discord.Interactions_SendMessage(i, "Add Speech Bubble", "Finding image...")

	// 1. Get the Message object associated with the Interaction request
	_, message := discord.GetAssociatedMessageFromInteraction(i)

	// 2. Check there's an associated Image
	imgUrl := helpers.GetImageFromMessage(message, "")
	if imgUrl == "" {
		discord.Interactions_EditIntoError(i, "No image found in message")
		cache.InteractionComplete(correlationId)
		return
	}

	imgExtension := imgwork.GetExtensionFromURL(imgUrl)
	if imgExtension == "" {
		discord.Interactions_EditIntoError(i, fmt.Sprintf("Invalid image extension (%s)", imgExtension))
		cache.InteractionComplete(correlationId)
		return
	}

	// 3. Generate the Output Image name
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
	discord.Interactions_EditText(i, "Add Speech Bubble", "Downloading image...")
	imageReader, downloadErr := imgwork.DownloadImageToReader(message.GuildID, imgUrl, isAnimated)
	if downloadErr != nil {
		logger.ErrorText(i.GuildID, "Couldn't download image")
		discord.Interactions_EditIntoError(i, "")
		cache.InteractionComplete(correlationId)
		return
	}

	// 4. Write the new Image to a Bytes Buffer
	var newImageBuffer bytes.Buffer
	discord.Interactions_EditText(i, "Add Speech Bubble", "Adding speech bubble...")
	addBubbleErr := addSpeechBubbleToImage(message.GuildID, imageReader, &newImageBuffer, isAnimated, imgExtension)
	if addBubbleErr != nil {
		logger.ErrorText(i.GuildID, "Error creating new image")
		discord.Interactions_EditIntoError(i, "")
		cache.InteractionComplete(correlationId)
		return
	}

	// 5. Send the new Image back to the User
	replyErr := discord.Message_ReplyWithImage(message, false, outputImageName, &newImageBuffer)
	if replyErr != nil {
		logger.ErrorText(i.GuildID, "Error sending new image")
		discord.Interactions_EditIntoError(i, "")
		cache.InteractionComplete(correlationId)
		return
	}

	discord.Interactions_EditText(i, "Add Speech Bubble Completed", "")
	cache.InteractionComplete(correlationId)
}

func addSpeechBubbleToImage(
	guildId string,
	imageReader io.Reader,
	newImgBuffer *bytes.Buffer,
	isAnimated bool,
	imgExtension string,
) error {
	start := time.Now()
	// setup values
	var gifImage *gif.GIF
	var imageHeight int
	var imageWidth int
	var inputImage image.Image
	var decodeErr error

	// Get the image height
	if isAnimated {
		gifImage, decodeErr = gif.DecodeAll(imageReader)
		if decodeErr != nil {
			return decodeErr
		}
		imageHeight = gifImage.Config.Height
		imageWidth = gifImage.Config.Width
	} else {
		switch imgExtension {
		case ".png":
			inputImage, decodeErr = png.Decode(imageReader)
		case ".jpg", ".jpeg":
			inputImage, decodeErr = jpeg.Decode(imageReader)
		case ".webp":
			inputImage, decodeErr = webp.Decode(imageReader)
		default:
			decodeErr = fmt.Errorf("unsupported extension: %s", imgExtension)
		}
		if decodeErr != nil {
			logger.Error(guildId, decodeErr)
			return decodeErr
		}
		imageHeight = inputImage.Bounds().Max.Y
		imageWidth = inputImage.Bounds().Max.X
	}

	// Get corresponding overlay image
	overlayImage, overlayErr := getOverlayImage(imageHeight, imageWidth)
	if overlayErr != nil {
		logger.Error(guildId, overlayErr)
		return overlayErr
	}

	// Define the Transparent colour, GIFs are wacky so we can't always add a new
	// colour to set as transparent. This colour is Discord's BG, a good backup.
	transparentColor := color.RGBA{49, 51, 56, 0}
	if isAnimated {
		// GIFs ----------------
		// Resize our overlay based on gif dimensions
		resizedOverlayReader, resizeErr := imgwork.ResizeImage(guildId, overlayImage, uint(gifImage.Config.Width))
		if resizeErr != nil {
			logger.Error(guildId, resizeErr)
			return resizeErr
		}
		resizedOverlay, overlayDecodeErr := png.Decode(resizedOverlayReader)
		if overlayDecodeErr != nil {
			logger.Error(guildId, overlayDecodeErr)
		}
		logger.Debug(guildId, "resized overlay is %dx%d", resizedOverlay.Bounds().Dy(), resizedOverlay.Bounds().Dx())
		logger.Debug(guildId, "gif is %dx%d", gifImage.Image[0].Bounds().Dy(), gifImage.Image[0].Bounds().Dx())

		for _, frame := range gifImage.Image {
			// Loop over each pixel in the frame
			for y := frame.Bounds().Min.Y; y < frame.Bounds().Max.Y; y++ {
				for x := frame.Bounds().Min.X; x < frame.Bounds().Max.X; x++ {
					// If the pixel in the PNG is black, set the corresponding pixel in the frame to transparent
					r, g, b, _ := resizedOverlay.At(x, y).RGBA()
					if r == 0 && g == 0 && b == 0 {
						frame.Set(x, y, transparentColor) // Set pixel to transparent color
					}
				}
			}
		}

		backgroundIndex := closestColorIndex(gifImage.Image[0].Palette, transparentColor)
		gifImage.BackgroundIndex = byte(backgroundIndex)

		// Encode the modified GIF to a buffer
		encodeErr := gif.EncodeAll(newImgBuffer, gifImage)
		if encodeErr != nil {
			logger.Error(guildId, encodeErr)
			return encodeErr
		}

	} else {
		// Static Images --------
		// Read and decode the input PNG image
		resizedOverlayReader, resizeErr := imgwork.ResizeImage(guildId, overlayImage, uint(inputImage.Bounds().Dx()))
		if resizeErr != nil {
			logger.Error(guildId, resizeErr)
			return resizeErr
		}
		resizedOverlay, overlayDecodeErr := png.Decode(resizedOverlayReader)
		if overlayDecodeErr != nil {
			logger.Error(guildId, overlayDecodeErr)
		}

		// Create a new image with the same dimensions as the input image
		outputImg := image.NewRGBA(inputImage.Bounds())
		draw.Draw(outputImg, inputImage.Bounds(), inputImage, image.Point{}, draw.Src)

		// Loop over each pixel in the input image and set corresponding pixel in the output image to transparent if condition is met
		for y := inputImage.Bounds().Min.Y; y < inputImage.Bounds().Max.Y; y++ {
			for x := inputImage.Bounds().Min.X; x < inputImage.Bounds().Max.X; x++ {
				// If the pixel in the template PNG is black, set the corresponding pixel in the output image to transparent
				r, g, b, _ := resizedOverlay.At(x, y).RGBA()
				if r == 0 && g == 0 && b == 0 {
					outputImg.Set(x, y, transparentColor) // Set pixel to transparent color
				}
			}
		}

		// Encode the modified PNG to buffer
		encodeErr := png.Encode(newImgBuffer, outputImg)
		if encodeErr != nil {
			encodeErr = fmt.Errorf("error encoding modified PNG: %w", encodeErr)
			logger.Error(guildId, encodeErr)
			return encodeErr
		}
	}

	logger.Info(guildId, "addSpeechBubbleToImage completed (%s) after %v", imgExtension, time.Since(start))
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

// determines which image file to use for the operation based on the height of the input image
func getOverlayImage(height int, width int) (image.Image, error) {
	rootDir, fpErr := filepath.Abs(filepath.Dir("."))
	if fpErr != nil {
		return nil, fpErr
	}
	overlayPath := "Resources/SpeechTemplates/"
	dimension := float64(height) / float64(width)
	if dimension >= 2.5 {
		overlayPath += "S"
	} else if dimension >= 1.5 {
		overlayPath += "M"
	} else {
		overlayPath += "L"
	}
	overlayPath += ".png"

	templateFilePath := filepath.Join(rootDir, overlayPath)
	overlayFile, fopenErr := os.Open(templateFilePath)
	if fopenErr != nil {
		return nil, fopenErr
	}
	defer overlayFile.Close()

	overlayImage, overlayDecodeErr := png.Decode(overlayFile)
	if overlayDecodeErr != nil {
		return nil, overlayDecodeErr
	}
	return overlayImage, nil
}
