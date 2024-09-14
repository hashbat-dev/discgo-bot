package bang

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	config "github.com/dabi-ngin/discgo-bot/Config"
	database "github.com/dabi-ngin/discgo-bot/Database"
	discord "github.com/dabi-ngin/discgo-bot/Discord"
	helpers "github.com/dabi-ngin/discgo-bot/Helpers"
)

type DelImage struct {
	ImageCategory string
}

func NewDelImage(imageCategory string) *DelImage {
	return &DelImage{
		ImageCategory: imageCategory,
	}
}

func (s DelImage) Name() string {
	return "delimage"
}

func (s DelImage) PermissionRequirement() int {
	return config.CommandLevelUser
}

func (s DelImage) Complexity() int {
	return config.TRIVIAL_TASK
}

func (s DelImage) Execute(message *discordgo.MessageCreate, command string) error {
	imgUrl := helpers.GetImageFromMessage(message.Message, "")
	if imgUrl == "" {
		return errors.New("no image found")
	}

	imgCat, err := database.GetImgCategory(message.GuildID, command)
	if err != nil {
		return errors.New("unable to get gif category")
	}

	imgStorage, err := database.GetImgStorage(message.GuildID, imgUrl)
	if err != nil {
		return errors.New("unable to get gif category")
	}

	imgGuildLink, err := database.GetImgGuildLink(message.GuildID, imgCat, imgStorage)
	if err != nil {
		return errors.New("unable to get guild link")
	}

	err = database.DeleteGuildLink(imgGuildLink)
	if err != nil {
		discord.SendUserMessageReply(message, true, "Unable to Delete Image")
		return err
	}

	discord.SendUserMessageReply(message, true, "Image successfully Deleted")
	discord.DeleteMessage(message)
	return nil
}
