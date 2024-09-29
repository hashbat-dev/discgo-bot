package bang

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/png"
	"io"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	config "github.com/hashbat-dev/discgo-bot/Config"
	discord "github.com/hashbat-dev/discgo-bot/Discord"
	helpers "github.com/hashbat-dev/discgo-bot/Helpers"
	imgwork "github.com/hashbat-dev/discgo-bot/ImgWork"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
	"golang.org/x/image/webp"
)

type Wide struct{}

func (s Wide) Name() string {
	return "wide"
}

func (s Wide) PermissionRequirement() int {
	return config.CommandLevelUser
}

func (s Wide) Complexity() int {
	return config.CPU_BOUND_TASK
}

func (s Wide) Execute(message *discordgo.MessageCreate, command string) error {
	progressMessage := discord.SendUserMessageReply(message, false, "Wide mode: Finding image...")

	// 1. Check we have a valid Image and Extension
	imgUrl := helpers.GetImageFromMessage(message.Message, "")
	if imgUrl == "" {
		discord.EditMessage(progressMessage, "Wide mode: Invalid image")
		return errors.New("no image found")
	}

	imgExtension := imgwork.GetExtensionFromURL(imgUrl)
	if imgExtension == "" {
		discord.EditMessage(progressMessage, "Wide mode: Invalid image")
		return errors.New("invalid extension")
	}

	outputImageName := uuid.New().String()
	isAnimated := false
	if imgExtension == ".gif" {
		outputImageName += ".gif"
		isAnimated = true
	} else {
		outputImageName += ".png"
	}

	// 3. Get the image as an io.Reader object
	discord.EditMessage(progressMessage, "Wide mode: Downloading Image...")
	imageReader, downloadErr := imgwork.DownloadImageToReader(message.GuildID, imgUrl, isAnimated)
	if downloadErr != nil {
		discord.SendUserMessageReply(message, false, "Error creating Image")
		return downloadErr
	}

	// 4. Write the new Image to a Bytes Buffer
	var newImageBuffer bytes.Buffer
	discord.EditMessage(progressMessage, "Wide mode: Widening...")
	widenImageErr := widenImage(message.GuildID, imageReader, &newImageBuffer, isAnimated, imgExtension)
	if widenImageErr != nil {
		discord.SendUserMessageReply(message, false, "Error creating Image")
		return widenImageErr
	}

	// 5. Send the new Image back to the User
	replyErr := discord.ReplyToMessageWithImageBuffer(message, true, outputImageName, &newImageBuffer)
	if replyErr != nil {
		logger.Error(message.GuildID, replyErr)
		return replyErr
	}

	discord.DeleteMessageObject(progressMessage)
	discord.DeleteMessage(message)
	return nil
}

func widenImage(
	guildId string,
	imageReader io.Reader,
	newImgBuffer *bytes.Buffer,
	isAnimated bool,
	imgExtension string,
) error {
	start := time.Now()
	// setup values
	var gifImage *gif.GIF
	var inputImage image.Image
	var decodeErr error

	if isAnimated {
		gifImage, decodeErr = gif.DecodeAll(imageReader)
		if decodeErr != nil {
			return decodeErr
		}
	} else {
		if imgExtension == ".png" || imgExtension == ".jpg" {
			// .jpg images are resized to a .png before this function is called
			inputImage, decodeErr = png.Decode(imageReader)
		} else if imgExtension == ".webp" {
			inputImage, decodeErr = webp.Decode(imageReader)
		}
		if decodeErr != nil {
			return decodeErr
		}
	}

	if isAnimated {
		// GIFs ----------------
		newHeight := uint(float64(gifImage.Config.Height) * 0.6)
		newWidth := uint(float64(gifImage.Config.Width) * 2)
		resizedGifReader, resizeErr := imgwork.ResizeGif(guildId, gifImage, newWidth, newHeight)
		if resizeErr != nil {
			logger.Error(guildId, resizeErr)
			return resizeErr
		}
		resizedGif, err := gif.DecodeAll(resizedGifReader)
		if err != nil {
			logger.Error(guildId, err)
			return err
		}

		// Encode the modified GIF to a buffer
		encodeErr := gif.EncodeAll(newImgBuffer, resizedGif)
		if encodeErr != nil {
			logger.Error(guildId, encodeErr)
			return encodeErr
		}

	} else {
		// Static Images --------
		// Read and decode the input PNG image
		stretchedImageReader, stretchErr := imgwork.StretchImage(guildId, inputImage, uint(inputImage.Bounds().Dx()))
		if stretchErr != nil {
			logger.Error(guildId, stretchErr)
			return stretchErr
		}
		stretchedImg, err := png.Decode(stretchedImageReader)
		if err != nil {
			logger.Error(guildId, err)
			return err
		}

		// Encode the modified PNG to buffer
		encodeErr := png.Encode(newImgBuffer, stretchedImg)
		if encodeErr != nil {
			encodeErr = fmt.Errorf("error encoding modified PNG: %w", encodeErr)
			logger.Error(guildId, encodeErr)
			return encodeErr
		}
	}

	logger.Info(guildId, "wide mode completed (%s) after %v", imgExtension, time.Since(start))
	return nil
}
