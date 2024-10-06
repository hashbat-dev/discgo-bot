package editmodule

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/png"
	"sync"
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

type FlipImage struct{}

func (s FlipImage) SelectName() string {
	return "Flip Image"
}

func (s FlipImage) Emoji() *discordgo.ComponentEmoji {
	return &discordgo.ComponentEmoji{Name: "🪞"}
}

func (s FlipImage) PermissionRequirement() int {
	return config.CommandLevelUser
}

func (s FlipImage) Complexity() int {
	return config.TRIVIAL_TASK
}

func (s FlipImage) Execute(i *discordgo.InteractionCreate, correlationId string) {
	cachedInteraction := cache.ActiveInteractions[correlationId]

	// 1. Get the Message object associated with the Interaction request
	message, err := discord.Message_GetObject(i.GuildID, i.ChannelID, cachedInteraction.Values.String["imgMessageId"])
	if err != nil {
		discord.Interactions_EditIntoError(i, "")
		cache.InteractionComplete(correlationId)
		return
	}

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

	// => Store these in the Interactions cache for later

	// 3. Create the Interaction Objects
	selectMenu := discord.CreateSelectMenu(discordgo.SelectMenu{
		CustomID: "flip-image_select-direction",
		Options: []discordgo.SelectMenuOption{
			{
				Label: "Left",
				Value: "left",
				Emoji: &discordgo.ComponentEmoji{Name: "⬅️"},
			}, {
				Label: "Right",
				Value: "right",
				Emoji: &discordgo.ComponentEmoji{Name: "➡️"},
			},
			{
				Label: "Left & Right",
				Value: "both",
				Emoji: &discordgo.ComponentEmoji{Name: "↔️"},
			}, {
				Label: "Up",
				Value: "up",
				Emoji: &discordgo.ComponentEmoji{Name: "⬆️"},
			}, {
				Label: "Down",
				Value: "down",
				Emoji: &discordgo.ComponentEmoji{Name: "⬇️"},
			}, {
				Label: "All",
				Value: "all",
				Emoji: &discordgo.ComponentEmoji{Name: "🔄"},
			},
		},
		Placeholder: "Choose Flip Direction...",
	}, correlationId, config.CPU_BOUND_TASK, FlipImageProcess)

	actionRow := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			selectMenu,
		},
	}

	// 4. Send the Select menu response
	err = config.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Flags:      discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{actionRow},
		},
	})
	if err != nil {
		logger.Error(i.GuildID, err)
	}

}

func FlipImageProcess(i *discordgo.InteractionCreate, correlationId string) {
	cachedInteraction := cache.ActiveInteractions[correlationId]
	discord.Interactions_EditText(&cachedInteraction.StartInteraction, "Flip Image", "Getting image...")

	// 1. Check we have a valid Image and Extension
	imgUrl := cache.ActiveInteractions[correlationId].Values.String["imgUrl"]
	imgExtension := cache.ActiveInteractions[correlationId].Values.String["imgExtension"]
	isAnimated := imgExtension == ".gif"

	// 2. Get the image as an io.Reader object
	discord.Interactions_EditText(&cachedInteraction.StartInteraction, "Flip Image", "Downloading image...")
	imageReader, err := imgwork.DownloadImageToReader(i.GuildID, imgUrl, isAnimated)
	if err != nil {
		logger.ErrorText(i.GuildID, "Couldn't download image")
		discord.Interactions_EditIntoError(&cachedInteraction.StartInteraction, "")
		cache.InteractionComplete(correlationId)
		return
	}

	// 3. Convert it to a buffer, this is needed for multiple operations
	var imageBuffer bytes.Buffer
	_, err = imageBuffer.ReadFrom(imageReader)
	if err != nil {
		logger.Error(i.GuildID, err)
		discord.Interactions_EditIntoError(&cachedInteraction.StartInteraction, "")
		cache.InteractionComplete(correlationId)
		return
	}

	// 4. Work out which Directions we want flipping
	FlipDirection := cache.ActiveInteractions[correlationId].Values.String["flip-image_select-direction"]
	if FlipDirection == "" {
		logger.ErrorText(i.GuildID, "Flip Direction was blank")
		discord.Interactions_EditIntoError(&cachedInteraction.StartInteraction, "")
		cache.InteractionComplete(correlationId)
		return
	}

	var flipDirections []string
	switch FlipDirection {
	case "left":
		flipDirections = append(flipDirections, "left")
	case "right":
		flipDirections = append(flipDirections, "right")
	case "up":
		flipDirections = append(flipDirections, "up")
	case "down":
		flipDirections = append(flipDirections, "down")
	case "both":
		flipDirections = append(flipDirections, "left")
		flipDirections = append(flipDirections, "right")
	case "all":
		flipDirections = append(flipDirections, "left")
		flipDirections = append(flipDirections, "right")
		flipDirections = append(flipDirections, "up")
		flipDirections = append(flipDirections, "down")
	default:
		logger.ErrorText(i.GuildID, "Unknown Flip Direction: %s", FlipDirection)
		discord.Interactions_EditIntoError(&cachedInteraction.StartInteraction, "")
		cache.InteractionComplete(correlationId)
		return
	}

	// 5. Perform the Flips and write back each Image individually
	message, err := config.Session.ChannelMessage(i.ChannelID, cache.ActiveInteractions[correlationId].Values.String["imgMessageId"])
	if err != nil {
		logger.Error(i.GuildID, err)
		discord.Interactions_EditIntoError(&cachedInteraction.StartInteraction, "")
		cache.InteractionComplete(correlationId)
		return
	}

	for _, flip := range flipDirections {
		discord.Interactions_EditText(&cachedInteraction.StartInteraction, "Flip Image", "Flipping "+flip+"...")
		outputImageName := uuid.New().String()
		if imgExtension == ".gif" {
			outputImageName += ".gif"
		} else {
			outputImageName += ".png"
		}

		var newImageBuffer bytes.Buffer
		if isAnimated {
			newImageBuffer, err = flipImageGif(i.GuildID, imageBuffer, flip)
		} else {
			newImageBuffer, err = flipImageStatic(i.GuildID, imageBuffer, flip)
		}

		if err != nil {
			continue
		}

		err = discord.Message_ReplyWithImage(message, false, outputImageName, &newImageBuffer)
		if err != nil {
			logger.Error(message.GuildID, err)
		}
	}

	// 6. Delete the calling Message
	discord.Interactions_EditText(&cachedInteraction.StartInteraction, "Flip Image Completed", "")
	cache.InteractionComplete(correlationId)
}

func flipImageGif(guildId string, imageReader bytes.Buffer, flipDirection string) (bytes.Buffer, error) {
	timeStarted := time.Now()
	imageReaderGifObject, err := gif.DecodeAll(&imageReader)
	if err != nil {
		logger.Error(guildId, err)
		return bytes.Buffer{}, err
	}

	mirroredGif := &gif.GIF{
		Image:           make([]*image.Paletted, len(imageReaderGifObject.Image)),
		Delay:           imageReaderGifObject.Delay,
		LoopCount:       imageReaderGifObject.LoopCount,
		Disposal:        imageReaderGifObject.Disposal,
		Config:          imageReaderGifObject.Config,
		BackgroundIndex: imageReaderGifObject.BackgroundIndex,
	}

	var wg sync.WaitGroup
	bounds := imageReaderGifObject.Image[0].Bounds()
	for i, frame := range imageReaderGifObject.Image {
		wg.Add(1)
		go func(flip string, x int, f *image.Paletted, b image.Rectangle) {
			defer wg.Done()
			mirroredGif.Image[x] = flipGifFrame(f, flip, b)
		}(flipDirection, i, frame, bounds)
	}
	wg.Wait()

	for i := range mirroredGif.Disposal {
		mirroredGif.Disposal[i] = gif.DisposalNone
	}

	var buf bytes.Buffer
	err = gif.EncodeAll(&buf, mirroredGif)
	if err != nil {
		logger.Error(guildId, err)
	}

	logger.Info(guildId, "Flip Gif completed after [%v]", time.Since(timeStarted))
	return buf, err
}

func flipGifFrame(img *image.Paletted, flipDirection string, bounds image.Rectangle) *image.Paletted {
	width := bounds.Dx()
	height := bounds.Dy()

	newImg := image.NewPaletted(bounds, img.Palette)
	var wg sync.WaitGroup
	if flipDirection == "left" || flipDirection == "right" {
		for y := 0; y < height; y++ {
			wg.Add(1)
			if flipDirection == "left" {
				go func(y int) {
					defer wg.Done()
					for x := 0; x < width/2; x++ {
						leftColor := img.At(x, y)
						rightX := width - x - 1
						newImg.Set(x, y, leftColor)
						newImg.Set(rightX, y, leftColor)
					}
				}(y)
			} else {
				go func(y int) {
					defer wg.Done()
					for x := width / 2; x < width; x++ {
						rightColor := img.At(x, y)
						leftX := width - x - 1
						newImg.Set(leftX, y, rightColor)
						newImg.Set(x, y, rightColor)
					}
				}(y)
			}
		}
	} else {
		for x := 0; x < width; x++ {
			wg.Add(1)
			if flipDirection == "up" {
				go func(x int) {
					defer wg.Done()
					for y := 0; y < height/2; y++ {
						leftColor := img.At(x, y)
						downY := height - y - 1
						newImg.Set(x, y, leftColor)
						newImg.Set(x, downY, leftColor)
					}
				}(x)
			} else {
				go func(x int) {
					defer wg.Done()
					for y := height / 2; y < height; y++ {
						rightColor := img.At(x, y)
						upY := height - y - 1
						newImg.Set(x, upY, rightColor)
						newImg.Set(x, y, rightColor)
					}
				}(x)
			}
		}
	}
	wg.Wait()
	return newImg
}

func flipImageStatic(guildId string, imageReader bytes.Buffer, flipDirection string) (bytes.Buffer, error) {
	timeStarted := time.Now()
	img, _, err := image.Decode(&imageReader)
	if err != nil {
		logger.Error(guildId, err)
		return bytes.Buffer{}, err
	}

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	newImg := image.NewRGBA(bounds)

	if flipDirection == "left" || flipDirection == "right" {
		if flipDirection == "left" {
			// Mirror left side to the right
			for y := 0; y < height; y++ {
				for x := 0; x <= width/2; x++ {
					newImg.Set(x, y, img.At(x, y))
					if x != width/2 || width%2 == 0 {
						newImg.Set(width-x-1, y, img.At(x, y))
					}
				}
			}

			// If the width is odd, shift left-hand pixels right by one pixel
			if width%2 != 0 {
				for y := 0; y < height; y++ {
					for x := width/2 - 1; x >= 0; x-- {
						newImg.Set(x+1, y, newImg.At(x, y))
					}
					newImg.Set(0, y, image.Transparent)
				}
			}
		} else {
			// Mirror right side to the left
			for y := 0; y < height; y++ {
				for x := width / 2; x < width; x++ {
					newImg.Set(x, y, img.At(x, y))
					newImg.Set(width-x-1, y, img.At(x, y))
				}
			}

			// If the width is odd, shift right-hand pixels left by one pixel
			if width%2 != 0 {
				for y := 0; y < height; y++ {
					for x := width/2 + 1; x < width; x++ {
						newImg.Set(x-1, y, newImg.At(x, y))
					}
					newImg.Set(width-1, y, image.Transparent)
				}
			}
		}
	} else {
		if flipDirection == "up" {
			// Mirror left side to the right
			for x := 0; x < width; x++ {
				for y := 0; y <= height/2; y++ {
					newImg.Set(x, y, img.At(x, y))
					if x != height/2 || height%2 == 0 {
						newImg.Set(x, height-y-1, img.At(x, y))
					}
				}
			}

			// If the height is odd, shift left-hand pixels right by one pixel
			if height%2 != 0 {
				for x := 0; x < width; x++ {
					for y := height/2 - 1; y >= 0; y-- {
						newImg.Set(x+1, y, newImg.At(x, y))
					}
					newImg.Set(x, 0, image.Transparent)
				}
			}
		} else {
			// Mirror right side to the left
			for x := 0; x < width; x++ {
				for y := height / 2; y < height; y++ {
					newImg.Set(x, y, img.At(x, y))
					newImg.Set(x, height-y-1, img.At(x, y))
				}
			}

			// If the width is odd, shift right-hand pixels left by one pixel
			if height%2 != 0 {
				for x := 0; x < width; x++ {
					for y := height/2 + 1; y < height; y++ {
						newImg.Set(x-1, y, newImg.At(x, y))
					}
					newImg.Set(x, height-1, image.Transparent)
				}
			}
		}
	}

	var buf bytes.Buffer
	err = png.Encode(&buf, newImg)
	if err != nil {
		logger.Error(guildId, err)
	}
	logger.Info(guildId, "Flip Image completed after [%v]", time.Since(timeStarted))
	return buf, err
}
