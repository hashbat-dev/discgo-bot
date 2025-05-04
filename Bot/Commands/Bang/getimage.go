package bang

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	config "github.com/hashbat-dev/discgo-bot/Config"
	database "github.com/hashbat-dev/discgo-bot/Database"
	discord "github.com/hashbat-dev/discgo-bot/Discord"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

type GetImage struct {
	ImageCategory string
}

func NewGetImage(imageCategory string) *GetImage {
	return &GetImage{
		ImageCategory: imageCategory,
	}
}

func (s GetImage) Name() string {
	return "getimage"
}

func (s GetImage) PermissionRequirement() int {
	return config.CommandLevelUser
}

func (s GetImage) Complexity() int {
	return config.TRIVIAL_TASK
}

func (s GetImage) Execute(message *discordgo.MessageCreate, command string) error {
	imgCat, err := database.GetImgCategory(message.GuildID, s.ImageCategory)
	if err != nil {
		discord.SendUserError(message, "Invalid Category")
		return errors.New("unable to get gif category")
	}

	imgUrl, err := database.GetRandomImage(message.GuildID, imgCat.ID)
	if err != nil {
		discord.SendUserError(message, "Couldn't find an Image")
		return err
	}

	msg := &discordgo.MessageSend{
		Content: imgUrl,
	}

	if message.MessageReference != nil {
		// Preserve the original reply reference
		msg.Reference = message.MessageReference
	}

	_, err = config.Session.ChannelMessageSendComplex(message.ChannelID, msg)
	if err != nil {
		logger.Error(message.GuildID, err)
		discord.SendUserError(message, "Couldn't send image")
		return err
	}

	discord.DeleteMessage(message)
	return nil
}
