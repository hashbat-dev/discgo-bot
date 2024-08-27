package bang

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	config "github.com/dabi-ngin/discgo-bot/Config"
	database "github.com/dabi-ngin/discgo-bot/Database"
	discord "github.com/dabi-ngin/discgo-bot/Discord"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
)

type GetImage struct{}

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
	imgCat, err := database.GetImgCategory(message.GuildID, command)
	if err != nil {
		discord.SendUserError(message, "Invalid Category")
		return errors.New("unable to get gif category")
	}

	imgUrl, err := database.GetRandomImage(message.GuildID, imgCat.ID)
	if err != nil {
		discord.SendUserError(message, "Couldn't find an Image")
		return err
	}

	_, err = config.Session.ChannelMessageSend(message.ChannelID, imgUrl)
	if err != nil {
		logger.Error(message.GuildID, err)
		discord.SendUserError(message, "Couldn't send Image")
		return err
	}

	return nil
}
