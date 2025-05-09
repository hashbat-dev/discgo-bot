package bang

import (
	"bytes"
	"errors"
	"image"
	"image/gif"
	"io"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	config "github.com/hashbat-dev/discgo-bot/Config"
	discord "github.com/hashbat-dev/discgo-bot/Discord"
	helpers "github.com/hashbat-dev/discgo-bot/Helpers"
	imgwork "github.com/hashbat-dev/discgo-bot/ImgWork"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

type Shuffle struct{}

func (s Shuffle) Name() string {
	return "shuffle"
}

func (s Shuffle) PermissionRequirement() int {
	return config.CommandLevelUser
}

func (s Shuffle) Complexity() int {
	return config.CPU_BOUND_TASK
}

func (s Shuffle) Execute(message *discordgo.MessageCreate, command string) error {
	progressMessage := discord.SendUserMessageReply(message, false, "Shuffle: Finding GIF...")

	// 1. Check we have a valid Image and Extension
	imgUrl := helpers.GetImageFromMessage(message.Message, "")
	if imgUrl == "" {
		discord.EditMessage(progressMessage, "Shuffle: Invalid image")
		return errors.New("no image found")
	}

	imgExtension := imgwork.GetExtensionFromURL(imgUrl)
	if imgExtension == "" {
		discord.EditMessage(progressMessage, "Shuffle: Invalid image")
		return errors.New("invalid extension")
	}

	// 2. Check the image is a GIF
	if imgExtension != ".gif" {
		discord.EditMessage(progressMessage, "Shuffle: Image was not a GIF")
		return errors.New("image provided is not a gif")
	}

	// 3. Get the image as an io.Reader object
	discord.EditMessage(progressMessage, "Shuffle: Downloading GIF...")
	imageReader, err := imgwork.DownloadImageToReader(message.GuildID, imgUrl, true)
	if err != nil {
		return err
	}

	// 4. Shuffle the GIF
	var buf bytes.Buffer
	discord.EditMessage(progressMessage, "Shuffle: Shuffling Frames...")
	err = shuffleGif(message.GuildID, imageReader, &buf)
	if err != nil {
		discord.SendUserMessageReply(message, false, "Error Shuffling GIF")
		return err
	}

	// 5. Return the reversed Image
	outputImageName := uuid.New().String() + ".gif"
	discord.DeleteMessageObject(progressMessage)
	discord.DeleteMessage(message)
	return discord.ReplyToMessageWithImageBuffer(message, true, outputImageName, &buf)
}

func shuffleGif(guildId string, imageReader io.Reader, buffer *bytes.Buffer) error {
	timeStarted := time.Now()
	gifImage, err := gif.DecodeAll(imageReader)
	if err != nil {
		logger.Error(guildId, err)
		return err
	}

	outGif := &gif.GIF{
		Image:           make([]*image.Paletted, 0, len(gifImage.Image)),
		Delay:           make([]int, 0, len(gifImage.Delay)),
		Disposal:        make([]byte, 0, len(gifImage.Disposal)),
		BackgroundIndex: gifImage.BackgroundIndex,
		LoopCount:       gifImage.LoopCount,
		Config:          gifImage.Config,
	}

	for i, frame := range gifImage.Image {
		if rand.Int() % 2 == 0 {
			tempFrames := []*image.Paletted{frame}
			tempDelay := []int{gifImage.Delay[i]}
			tempDisposal := []byte{gifImage.Disposal[i]}
			tempFrames = append(tempFrames, outGif.Image...)
			tempDelay = append(tempDelay, outGif.Delay...)
			tempDisposal = append(tempDisposal, outGif.Disposal...)
			outGif.Image = tempFrames
			outGif.Delay = tempDelay
			outGif.Disposal = tempDisposal
		} else {
			outGif.Image = append(outGif.Image, gifImage.Image[i])
			outGif.Delay = append(outGif.Delay, gifImage.Delay[i])
			if len(gifImage.Disposal) > 0 {
				outGif.Disposal = append(outGif.Disposal, gifImage.Disposal[i])
			} else {
				outGif.Disposal = append(outGif.Disposal, gif.DisposalNone) // Default to no disposal
			}
		}
	}

	// 5. Encode the shuffled GIF to the buffer
	err = gif.EncodeAll(buffer, outGif)
	if err != nil {
		logger.Error(guildId, err)
		return err
	}

	logger.Info(guildId, "Shuffle GIF completed after [%v]", time.Since(timeStarted))
	return nil
}
