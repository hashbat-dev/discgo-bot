package bang

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	config "github.com/dabi-ngin/discgo-bot/Config"
	database "github.com/dabi-ngin/discgo-bot/Database"
	discord "github.com/dabi-ngin/discgo-bot/Discord"
	helpers "github.com/dabi-ngin/discgo-bot/Helpers"
)

type AddImage struct {
	ImageCategory string
}

func NewAddImage(imageCategory string) *AddImage {
	return &AddImage{
		ImageCategory: imageCategory,
	}
}

func (s AddImage) Name() string {
	return "addimage"
}

func (s AddImage) PermissionRequirement() int {
	return config.CommandLevelUser
}

func (s AddImage) Complexity() int {
	return config.TRIVIAL_TASK
}

func (s AddImage) Execute(message *discordgo.MessageCreate, command string) error {
	imgUrl := helpers.GetImageFromMessage(message.Message, "")
	if imgUrl == "" {
		return errors.New("no image found")
	}

	err := database.AddImg(message, helpers.RemoveStartingXCharacters(command, 3), imgUrl)
	if err != nil {
		discord.SendUserMessageReply(message, true, "Error adding Image")
		return err
	}

	discord.SendUserMessageReply(message, true, "Image successfully added")
	return err
}
