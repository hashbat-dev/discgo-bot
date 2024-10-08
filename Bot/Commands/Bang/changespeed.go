package bang

import (
	"bytes"
	"errors"
	"image"
	"image/gif"
	"io"
	"math"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	config "github.com/hashbat-dev/discgo-bot/Config"
	discord "github.com/hashbat-dev/discgo-bot/Discord"
	helpers "github.com/hashbat-dev/discgo-bot/Helpers"
	imgwork "github.com/hashbat-dev/discgo-bot/ImgWork"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

type ChangeSpeed struct {
	SpeedUp bool
}

func NewChangeSpeed(speedUp bool) *ChangeSpeed {
	return &ChangeSpeed{
		SpeedUp: speedUp,
	}
}

func (s ChangeSpeed) Name() string {
	return "changespeed"
}

func (s ChangeSpeed) PermissionRequirement() int {
	return config.CommandLevelUser
}

func (s ChangeSpeed) Complexity() int {
	return config.CPU_BOUND_TASK
}

func (s ChangeSpeed) Execute(message *discordgo.MessageCreate, command string) error {
	commandMessage := "Speed Up"
	if !s.SpeedUp {
		commandMessage = "Slow Down"
	}
	progressMessage := discord.SendUserMessageReply(message, false, commandMessage+": Finding GIF...")

	// 1. Check we have a valid Image and Extension
	imgUrl := helpers.GetImageFromMessage(message.Message, "")
	if imgUrl == "" {
		discord.EditMessage(progressMessage, commandMessage+": Invalid image")
		return errors.New("no image found")
	}

	imgExtension := imgwork.GetExtensionFromURL(imgUrl)
	if imgExtension == "" {
		discord.EditMessage(progressMessage, commandMessage+": Invalid image")
		return errors.New("invalid extension")
	}

	// 2. Check the image is a GIF
	if imgExtension != ".gif" {
		discord.EditMessage(progressMessage, commandMessage+": Image isn't a GIF")
		discord.SendUserMessageReply(message, false, "Can only reverse .gifs")
		return errors.New("image provided is not a gif")
	}

	// 3. Get the image as an io.Reader object
	discord.EditMessage(progressMessage, commandMessage+": Downloading GIF...")
	imageReader, err := imgwork.DownloadImageToReader(message.GuildID, imgUrl, true)
	if err != nil {
		discord.EditMessage(progressMessage, commandMessage+": Error downloading Image")
		return err
	}

	// 4. Change the GIF Speed
	var newImageBuffer bytes.Buffer
	discord.EditMessage(progressMessage, commandMessage+": Changing Speed...")
	err = changeSpeedGif(message.GuildID, imageReader, &newImageBuffer, s.SpeedUp)
	if err != nil {
		discord.EditMessage(progressMessage, commandMessage+": Error changing GIF Speed")
		return err
	}

	// 5. Return the reversed Image
	outputImageName := uuid.New().String() + ".gif"
	discord.DeleteMessageObject(progressMessage)
	discord.DeleteMessage(message)
	return discord.ReplyToMessageWithImageBuffer(message, true, outputImageName, &newImageBuffer)
}

func changeSpeedGif(guildId string, imageReader io.Reader, buffer *bytes.Buffer, speedUp bool) error {
	// 1. Decode the file into a GIF object
	timeStarted := time.Now()
	gifImage, err := gif.DecodeAll(imageReader)
	if err != nil {
		logger.Error(guildId, err)
		return err
	}

	// 2. Lower the Delay of the frames
	outGif := &gif.GIF{}
	var snippedDelay []int
	alreadySlowestCount := 0
	for i := range gifImage.Image {
		frameDelay := gifImage.Delay[i]
		if frameDelay <= 2 {
			alreadySlowestCount++
		}
		var newValue int
		if speedUp {
			newValue = int(math.Round(float64(frameDelay) / 2.0))
		} else {
			newValue = int(math.Round(float64(frameDelay) * 2.0))
		}
		if newValue < 2 {
			newValue = 2
		}
		snippedDelay = append(snippedDelay, newValue)
	}

	// 3. If it's already as slow as possible due to the delay, cut frames out
	var newPaletted []*image.Paletted
	var newDelay []int
	var newDisposal []byte
	if alreadySlowestCount == len(gifImage.Delay) {
		keptLast := false
		for i, frame := range gifImage.Image {
			if !keptLast {
				newPaletted = append(newPaletted, frame)
				newDelay = append(newDelay, snippedDelay[i])
				newDisposal = append(newDisposal, gifImage.Disposal[i])
				keptLast = true
			} else {
				keptLast = false
			}
		}
	} else {
		newPaletted = gifImage.Image
		newDelay = snippedDelay
		newDisposal = gifImage.Disposal
	}

	outGif.BackgroundIndex = gifImage.BackgroundIndex
	outGif.Config = gifImage.Config
	outGif.Disposal = newDisposal
	outGif.LoopCount = gifImage.LoopCount
	outGif.Image = newPaletted
	outGif.Delay = newDelay

	// 3. Write the new GIF to buffer arg
	err = gif.EncodeAll(buffer, outGif)
	if err != nil {
		logger.Error(guildId, err)
		return err
	}

	direction := "down"
	if speedUp {
		direction = "up"
	}

	logger.Info(guildId, "ChangeSpeed [%v] completed after [%v]", direction, time.Since(timeStarted))
	return nil
}
