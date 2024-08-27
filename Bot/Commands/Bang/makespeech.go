package bang

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	config "github.com/dabi-ngin/discgo-bot/Config"
	helpers "github.com/dabi-ngin/discgo-bot/Helpers"
	img "github.com/dabi-ngin/discgo-bot/Img"
)

type MakeSpeech struct{}

func (s MakeSpeech) Name() string {
	return "makespeech"
}

func (s MakeSpeech) PermissionRequirement() int {
	return config.CommandLevelUser
}

func (s MakeSpeech) Complexity() int {
	return config.IO_BOUND_TASK
}

func (s MakeSpeech) Execute(message *discordgo.MessageCreate, command string) error {
	// Check we have an Image and that it's a valid extension
	imgUrl := helpers.GetImageFromMessage(message.Message, "")
	if imgUrl == "" {
		return errors.New("no image found")
	}

	imgExtension := img.GetExtensionFromURL(imgUrl)
	if imgExtension == "" {
		return errors.New("invalid extension")
	}

	// Is it Animated? If so bounce it back into the Animated Channel
	if imgExtension == ".gif" {
		//TODO I hate import cycles
	}

	return nil
}
