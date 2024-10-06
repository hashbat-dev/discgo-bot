package editmodule

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
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	cache "github.com/hashbat-dev/discgo-bot/Cache"
	config "github.com/hashbat-dev/discgo-bot/Config"
	discord "github.com/hashbat-dev/discgo-bot/Discord"
	helpers "github.com/hashbat-dev/discgo-bot/Helpers"
	imgwork "github.com/hashbat-dev/discgo-bot/ImgWork"
)

type DeepFry struct{}

func (s DeepFry) SelectName() string {
	return "Deep Fry"
}

func (s DeepFry) Emoji() *discordgo.ComponentEmoji {
	return &discordgo.ComponentEmoji{Name: "🧑‍🍳"}
}

func (s DeepFry) PermissionRequirement() int {
	return config.CommandLevelUser
}

func (s DeepFry) Complexity() int {
	return config.CPU_BOUND_TASK
}

func (s DeepFry) Execute(i *discordgo.InteractionCreate, correlationId string) {
	msgTitle := "Deep Fry"
	discord.Interactions_SendMessage(i, msgTitle, "Decoding media...")

	// 1. Check we have a valid Image and Extension
	_, message := discord.GetAssociatedMessageFromInteraction(i)
	imgUrl := helpers.GetImageFromMessage(message, "")
	if imgUrl == "" {
		discord.Interactions_EditIntoError(i, "No image found in Message")
		cache.InteractionComplete(correlationId)
		return
	}

	imgExtension := imgwork.GetExtensionFromURL(imgUrl)
	if imgExtension == "" {
		discord.Interactions_EditIntoError(i, "Invalid image extension")
		cache.InteractionComplete(correlationId)
		return
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
	discord.Interactions_EditText(i, msgTitle, "Downloading image...")
	imageReader, downloadErr := imgwork.DownloadImageToReader(message.GuildID, imgUrl, isAnimated)
	if downloadErr != nil {
		discord.Interactions_EditIntoError(i, "")
		cache.InteractionComplete(correlationId)
		return
	}

	// 4. Write the new Image to a Bytes Buffer
	var newImageBuffer bytes.Buffer
	discord.Interactions_EditText(i, msgTitle, "Frying :cook:...")
	deepfryImageErr := deepfryImage(imageReader, &newImageBuffer, isAnimated)
	if deepfryImageErr != nil {
		discord.Interactions_EditIntoError(i, "")
		cache.InteractionComplete(correlationId)
		return
	}

	// 5. Send the new Image back to the User
	replyErr := discord.Message_ReplyWithImage(message, true, outputImageName, &newImageBuffer)
	if replyErr != nil {
		discord.Interactions_EditIntoError(i, "")
		cache.InteractionComplete(correlationId)
		return
	}

	discord.Interactions_EditText(i, msgTitle+" Completed", "")
	cache.InteractionComplete(correlationId)
}

func deepfryImage(
	imageReader io.Reader,
	newImgBuffer *bytes.Buffer,
	isAnimated bool,
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
	// 2. Write the resulting Bytes to the OutputFilePath as a file
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
