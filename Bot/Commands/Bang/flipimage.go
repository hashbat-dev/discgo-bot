package bang

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/png"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	config "github.com/dabi-ngin/discgo-bot/Config"
	discord "github.com/dabi-ngin/discgo-bot/Discord"
	helpers "github.com/dabi-ngin/discgo-bot/Helpers"
	imgwork "github.com/dabi-ngin/discgo-bot/ImgWork"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
	"github.com/google/uuid"
)

type FlipImage struct {
	FlipDirection string
}

func NewFlipImage(flipDirection string) *FlipImage {
	return &FlipImage{
		FlipDirection: flipDirection,
	}
}

func (s FlipImage) Name() string {
	return "flipimage"
}

func (s FlipImage) PermissionRequirement() int {
	return config.CommandLevelUser
}

func (s FlipImage) Complexity() int {
	return config.CPU_BOUND_TASK
}

func (s FlipImage) Execute(message *discordgo.MessageCreate, command string) error {
	// 1. Check we have a valid Image and Extension
	imgUrl := helpers.GetImageFromMessage(message.Message, "")
	if imgUrl == "" {
		return errors.New("no image found")
	}

	imgExtension := imgwork.GetExtensionFromURL(imgUrl)
	if imgExtension == "" {
		return errors.New("invalid extension")
	}

	isAnimated := imgExtension == ".gif"

	// 2. Get the image as an io.Reader object
	imageReader, _, err := imgwork.DownloadImageToReader(message.GuildID, imgUrl, isAnimated, 0)
	if err != nil {
		return err
	}

	// 3. Convert it to a buffer, this is needed for multiple operations
	var imageBuffer bytes.Buffer
	_, err = imageBuffer.ReadFrom(imageReader)
	if err != nil {
		logger.Error(message.GuildID, err)
		return err
	}

	// 4. Work out which Directions we want flipping
	var flipDirections []string
	switch s.FlipDirection {
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
		err = fmt.Errorf("Unknown Flip Direction [%v]", s.FlipDirection)
		logger.Error(message.GuildID, err)
		return err
	}

	// 5. Perform the Flips and write back each Image individually
	for _, flip := range flipDirections {
		outputImageName := uuid.New().String()
		if imgExtension == ".gif" {
			outputImageName += ".gif"
		} else {
			outputImageName += ".png"
		}

		var newImageBuffer bytes.Buffer
		if isAnimated {
			newImageBuffer, err = flipImageGif(message.GuildID, imageBuffer, flip)
		} else {
			newImageBuffer, err = flipImageStatic(message.GuildID, imageBuffer, flip)
		}

		if err != nil {
			continue
		}

		err = discord.ReplyToMessageWithImageBuffer(message, true, outputImageName, &newImageBuffer)
		if err != nil {
			logger.Error(message.GuildID, err)
		}
	}

	// 6. Delete the calling Message
	discord.DeleteMessage(message)
	return nil
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
