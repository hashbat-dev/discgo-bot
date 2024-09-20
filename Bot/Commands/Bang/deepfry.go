package bang

import (
	"bytes"
	"errors"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"io"
	"sync"

	"github.com/bwmarrin/discordgo"
	config "github.com/dabi-ngin/discgo-bot/Config"
	discord "github.com/dabi-ngin/discgo-bot/Discord"
	helpers "github.com/dabi-ngin/discgo-bot/Helpers"
	imgwork "github.com/dabi-ngin/discgo-bot/ImgWork"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

type DeepFry struct{}

func (s DeepFry) Name() string {
	return "deepfry"
}

func (s DeepFry) PermissionRequirement() int {
	return config.CommandLevelUser
}

func (s DeepFry) Complexity() int {
	return config.CPU_BOUND_TASK
}

func (s DeepFry) Execute(message *discordgo.MessageCreate, command string) error {
	progressMessage := discord.SendUserMessageReply(message, false, "Deepfry: Decoding media...")

	// 1. Check we have a valid Image and Extension
	imgUrl := helpers.GetImageFromMessage(message.Message, "")
	if imgUrl == "" {
		discord.EditMessage(progressMessage, "Deepfry: Invalid media")
		return errors.New("no media found")
	}

	imgExtension := imgwork.GetExtensionFromURL(imgUrl)
	if imgExtension == "" {
		discord.EditMessage(progressMessage, "Deepfry: Invalid image")
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
	discord.EditMessage(progressMessage, "Deepfry: Downloading Media...")
	imageReader, downloadErr := imgwork.DownloadImageToReader(message.GuildID, imgUrl, isAnimated)
	if downloadErr != nil {
		discord.SendUserMessageReply(message, false, "Error creating Media")
		return downloadErr
	}

	// 4. Write the new Image to a Bytes Buffer
	var newImageBuffer bytes.Buffer
	discord.EditMessage(progressMessage, "Deepfry: Frying :cook:")
	deepfryImageErr := deepfryImage(message.GuildID, imageReader, &newImageBuffer, isAnimated, imgExtension)
	if deepfryImageErr != nil {
		discord.SendUserMessageReply(message, false, "Error creating Media")
		return deepfryImageErr
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

func deepfryImage(
	guildId string,
	imageReader io.Reader,
	newImgBuffer *bytes.Buffer,
	isAnimated bool,
	imgExtension string,
) error {

	var err error

	if isAnimated {
		err = DeepFryGIF(imageReader, newImgBuffer)
		if err != nil {
			return errors.New(err.Error())
		}
	} else {
		err = DeepFryIMG(imageReader, newImgBuffer)
		if err != nil {
			return errors.New(err.Error())
		}
	}
	if err != nil {
		return err
	}

	// 2. Write the resulting Bytes to the OutputFilePath as a file
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

func DeepFryIMG(imageReader io.Reader, newImgBuffer *bytes.Buffer) error {

	var contrastOffset float64 = 85
	var brightnessOffset float64 = -20
	var sigmaOffset float64 = 5

	img, _, err := image.Decode(imageReader)
	if err != nil {
		return err
	}

	img = imaging.AdjustContrast(img, contrastOffset)
	img = imaging.AdjustBrightness(img, brightnessOffset)
	img = imaging.Sharpen(img, sigmaOffset)

	err = png.Encode(newImgBuffer, img)
	return err
}

func DeepFryGIF(imageReader io.Reader, newImgBuffer *bytes.Buffer) error {
	var contrastOffset float64 = 55
	var brightnessOffset float64 = -5
	var sigmaOffset float64 = 3.5

	img, err := gif.DecodeAll(imageReader)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for i := range img.Image {
		frame := img.Image[i]

		wg.Add(1)
		go func(gi int, gframe image.Image, gcontrastOffset float64, gbrightnessOffset float64, gSigma float64) {

			defer wg.Done()

			gbounds := frame.Bounds()

			// Convert GIF frame to RGBA
			rgba := image.NewRGBA(gbounds)
			draw.Draw(rgba, gbounds, gframe, gbounds.Min, draw.Src)

			// Apply effects
			nrgba := imaging.AdjustContrast(rgba, gcontrastOffset)
			nrgba = imaging.AdjustBrightness(nrgba, gbrightnessOffset)
			nrgba = imaging.Sharpen(nrgba, gSigma)

			// Convert *image.NRGBA to *image.RGBA
			rgba = nrgbaToRGBA(nrgba)

			// Convert back to GIF frame
			palettedImage := image.NewPaletted(gbounds, palette.Plan9)
			draw.FloydSteinberg.Draw(palettedImage, gbounds, rgba, image.Point{})
			img.Image[gi] = palettedImage
		}(i, frame, contrastOffset, brightnessOffset, sigmaOffset)

	}

	wg.Wait()

	// Encode the modified GIF
	err = gif.EncodeAll(newImgBuffer, img)
	if err != nil {
		return err
	}

	return nil
}

func nrgbaToRGBA(nrgba *image.NRGBA) *image.RGBA {
	bounds := nrgba.Bounds()
	rgba := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba.Set(x, y, nrgba.At(x, y))
		}
	}
	return rgba
}
