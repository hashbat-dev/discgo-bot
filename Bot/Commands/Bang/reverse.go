package bang

import (
	"bytes"
	"errors"
	"image"
	"image/gif"
	"io"
	"time"

	"github.com/bwmarrin/discordgo"
	config "github.com/dabi-ngin/discgo-bot/Config"
	discord "github.com/dabi-ngin/discgo-bot/Discord"
	helpers "github.com/dabi-ngin/discgo-bot/Helpers"
	imgwork "github.com/dabi-ngin/discgo-bot/ImgWork"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
	"github.com/google/uuid"
)

type Reverse struct{}

func (s Reverse) Name() string {
	return "reverse"
}

func (s Reverse) PermissionRequirement() int {
	return config.CommandLevelUser
}

func (s Reverse) Complexity() int {
	return config.CPU_BOUND_TASK
}

func (s Reverse) Execute(message *discordgo.MessageCreate, command string) error {
	progressMessage := discord.SendUserMessageReply(message, false, "Reverse: Finding GIF...")

	// 1. Check we have a valid Image and Extension
	imgUrl := helpers.GetImageFromMessage(message.Message, "")
	if imgUrl == "" {
		discord.EditMessage(progressMessage, "Reverse: Invalid image")
		return errors.New("no image found")
	}

	imgExtension := imgwork.GetExtensionFromURL(imgUrl)
	if imgExtension == "" {
		discord.EditMessage(progressMessage, "Reverse: Invalid image")
		return errors.New("invalid extension")
	}

	// 2. Check the image is a GIF
	if imgExtension != ".gif" {
		discord.EditMessage(progressMessage, "Reverse: Image was not a GIF")
		return errors.New("image provided is not a gif")
	}

	// 3. Get the image as an io.Reader object
	discord.EditMessage(progressMessage, "Reverse: Downloading GIF...")
	imageReader, _, err := imgwork.DownloadImageToReader(message.GuildID, imgUrl, true, 0)
	if err != nil {
		return err
	}

	// 4. Reverse the GIF
	var newImageBuffer bytes.Buffer
	discord.EditMessage(progressMessage, "Reverse: Reversing GIF...")
	err = reverseGif(message.GuildID, imageReader, &newImageBuffer)
	if err != nil {
		discord.SendUserMessageReply(message, false, "Error reversing GIF")
		return err
	}

	// 5. Return the reversed Image
	outputImageName := uuid.New().String() + ".gif"
	discord.DeleteMessageObject(progressMessage)
	discord.DeleteMessage(message)
	return discord.ReplyToMessageWithImageBuffer(message, true, outputImageName, &newImageBuffer)
}

func reverseGif(guildId string, imageReader io.Reader, buffer *bytes.Buffer) error {
	// 1. Decode the file into a GIF object
	timeStarted := time.Now()
	gifImage, err := gif.DecodeAll(imageReader)
	if err != nil {
		logger.Error(guildId, err)
		return err
	}

	// 2. Create a new GIF object for the reversed frames
	outGif := &gif.GIF{
		Image:           make([]*image.Paletted, 0, len(gifImage.Image)),
		Delay:           make([]int, 0, len(gifImage.Delay)),
		Disposal:        make([]byte, 0, len(gifImage.Disposal)),
		BackgroundIndex: gifImage.BackgroundIndex,
		LoopCount:       gifImage.LoopCount,
		Config:          gifImage.Config,
	}

	frameCount := len(gifImage.Image)

	// 3. Reverse the frames and handle disposal methods and transparency
	for i := frameCount - 1; i >= 0; i-- {
		outGif.Image = append(outGif.Image, gifImage.Image[i])
		outGif.Delay = append(outGif.Delay, gifImage.Delay[i])

		// Add disposal method if it exists (important for handling transparency)
		if len(gifImage.Disposal) > 0 {
			outGif.Disposal = append(outGif.Disposal, gifImage.Disposal[i])
		} else {
			outGif.Disposal = append(outGif.Disposal, gif.DisposalNone) // Default to no disposal
		}
	}

	// 4. Handle transparency - We need to check and reset the transparency in frames with DisposalBackground (2)
	for i := frameCount - 1; i >= 0; i-- {
		if len(gifImage.Disposal) > 0 && gifImage.Disposal[i] == gif.DisposalBackground {
			clearBackground(outGif.Image[frameCount-1-i], gifImage.BackgroundIndex)
		}
	}

	// 5. Encode the reversed GIF to the buffer
	err = gif.EncodeAll(buffer, outGif)
	if err != nil {
		logger.Error(guildId, err)
		return err
	}

	logger.Info(guildId, "Reverse GIF completed after [%v]", time.Since(timeStarted))
	return nil
}

// clearBackground ensures the transparent pixels are handled correctly when the DisposalBackground method is used.
func clearBackground(img *image.Paletted, backgroundIndex byte) {
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			if img.ColorIndexAt(x, y) == backgroundIndex {
				img.SetColorIndex(x, y, 0) // Reset to transparent pixel
			}
		}
	}
}
