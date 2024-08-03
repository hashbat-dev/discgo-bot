package imagework

import (
	"bytes"
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/bwmarrin/discordgo"
	"github.com/chai2010/webp"
	"github.com/dabi-ngin/discgo-bot/Bot/audit"
	"github.com/dabi-ngin/discgo-bot/Bot/config"
	"github.com/dabi-ngin/discgo-bot/Bot/helpers"
	"github.com/dabi-ngin/discgo-bot/Bot/logging"
	"github.com/google/uuid"
)

var ResizedWidth = 300

func HandleMessage(message *discordgo.MessageCreate, requestType string) {
	// 1. Get the Source Image
	imageUrl, imageExtension, err := GetImageURLAndExtension(message)
	if err != nil {
		logging.SendErrorMsg(message, err.Error())
		audit.Error(err)
		return
	}
	if imageUrl == "" || imageExtension == "" {
		err = errors.New("could not find image url/extension")
		logging.SendErrorMsg(message, err.Error())
		audit.Error(err)
		return
	}

	// 2. Generate the output image name - shows as attachment in discord
	newImgName := uuid.New().String()
	isGif := false
	if imageExtension == ".gif" {
		newImgName += ".gif"
		isGif = true
	} else {
		newImgName += ".png"
	}

	// 3. Confirm we have the correct File Extension for .gif only functions
	if !isGif {
		canContinue := true
		if requestType == "reverse" || requestType == "speedup" {
			canContinue = false
		}
		if !canContinue {
			err = errors.New("how am I supposed to reverse something that isn't a gif? stupid")
		}
	}

	if err != nil {
		logging.SendErrorMsg(message, err.Error())
		audit.Error(err)
		return
	}

	// 4. Download the source image
	var downloadedStaticImg *image.Image
	var downloadedGifImg *gif.GIF

	if isGif {
		downloadedGifImg, err = helpers.DownloadGif(imageUrl)
	} else {
		downloadedStaticImg, err = DownloadImage(imageUrl)
	}

	if err != nil {
		logging.SendErrorMsg(message, err.Error())
		audit.Error(err)
		return
	}

	// 5. Resize the downloaded image if the function requires it
	resizeImage := false
	if requestType == "makespeech" {
		resizeImage = true
	}
	var imageReader io.Reader

	if resizeImage {
		// Resizing is needed
		if isGif {
			imageReader, err = ResizeImageGif(downloadedGifImg, uint(ResizedWidth))
		} else {
			imageReader, err = ResizeImageStatic(*downloadedStaticImg, uint(ResizedWidth))
		}
	} else {
		// Write the downloaded file as is to the io.Reader
		var buf bytes.Buffer
		if isGif {
			err = gif.EncodeAll(&buf, downloadedGifImg)
			if err == nil {
				imageReader = &buf
			}
		} else {
			switch imageExtension {
			case ".png":
				err = png.Encode(&buf, *downloadedStaticImg)
			case ".jpg":
				err = jpeg.Encode(&buf, *downloadedStaticImg, nil)
			case ".webp":
				webpBytes, errWebp := webp.EncodeRGBA(*downloadedStaticImg, 100.0)
				if errWebp != nil {
					err = errWebp
				} else {
					buf = *bytes.NewBuffer(webpBytes)
				}
			default:
				err = errors.New("unsupported file extension of " + imageExtension)
			}
			if err == nil {
				imageReader = &buf
			}
		}
	}

	if err != nil {
		logging.SendErrorMsg(message, err.Error())
		audit.Error(err)
		return
	}

	// 6. Perform our Function's operation
	responseBuffer := new(bytes.Buffer)
	responseChan := make(chan error)

	req := FlipRequest{
		Message:        message,
		ImageReader:    imageReader,
		IsGif:          isGif,
		RequestType:    requestType,
		ResponseBuffer: responseBuffer,
		ResponseChan:   responseChan,
	}

	EnqueueFlipRequest(req)

	go func() {
		err := <-responseChan
		if err != nil {
			logging.SendErrorMsg(message, err.Error())
			audit.Error(err)
			return
		}

		// 7. Send the generated image as a Message
		fileObj := &discordgo.File{
			Name:   newImgName,
			Reader: responseBuffer,
		}

		_, err = config.Session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
			Files:     []*discordgo.File{fileObj},
			Reference: message.ReferencedMessage.MessageReference,
		})

		if err != nil {
			audit.Error(err)
		} else {
			audit.Log("Successfully created and sent image work, RequestType: " + requestType)
		}
	}()
}
