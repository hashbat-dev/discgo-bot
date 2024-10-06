package editmodule

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"io"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	cache "github.com/hashbat-dev/discgo-bot/Cache"
	config "github.com/hashbat-dev/discgo-bot/Config"
	discord "github.com/hashbat-dev/discgo-bot/Discord"
	helpers "github.com/hashbat-dev/discgo-bot/Helpers"
	imgwork "github.com/hashbat-dev/discgo-bot/ImgWork"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

type Reverse struct{}

func (s Reverse) SelectName() string {
	return "Reverse GIF"
}

func (s Reverse) Emoji() *discordgo.ComponentEmoji {
	return &discordgo.ComponentEmoji{Name: "⏪"}
}

func (s Reverse) PermissionRequirement() int {
	return config.CommandLevelUser
}

func (s Reverse) Complexity() int {
	return config.CPU_BOUND_TASK
}

func (s Reverse) Execute(i *discordgo.InteractionCreate, correlationId string) {
	msgTitle := "Reverse"
	discord.Interactions_SendMessage(i, msgTitle, "Decoding media...")

	// 1. Check we have a valid Image and Extension
	_, message := discord.GetAssociatedMessageFromInteraction(i)
	imgUrl := helpers.GetImageFromMessage(message, "")
	if imgUrl == "" {
		discord.Interactions_EditIntoError(i, "No image found in Message")
		cache.InteractionComplete(correlationId)
		return
	}

	// 2. Check the image is a GIF
	imgExtension := imgwork.GetExtensionFromURL(imgUrl)
	if imgExtension != ".gif" {
		discord.Interactions_EditIntoError(i, fmt.Sprintf("Can't reverse %s's!", imgExtension))
		cache.InteractionComplete(correlationId)
		return
	}

	// 3. Get the image as an io.Reader object
	discord.Interactions_EditText(i, msgTitle, "Downloading GIF...")
	imageReader, err := imgwork.DownloadImageToReader(message.GuildID, imgUrl, true)
	if err != nil {
		discord.Interactions_EditIntoError(i, "")
		cache.InteractionComplete(correlationId)
		return
	}

	// 4. Reverse the GIF
	var newImageBuffer bytes.Buffer
	discord.Interactions_EditText(i, msgTitle, "Reversing GIF...")
	err = reverseGif(message.GuildID, imageReader, &newImageBuffer)
	if err != nil {
		discord.Interactions_EditIntoError(i, "")
		cache.InteractionComplete(correlationId)
		return
	}

	// 5. Return the reversed Image
	outputImageName := uuid.New().String() + ".gif"
	discord.Interactions_EditText(i, msgTitle+" Completed", "")
	if discord.Message_ReplyWithImage(message, true, outputImageName, &newImageBuffer) != nil {
		logger.Info(i.GuildID, "Error sending ReplyWithImage")
	}
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
