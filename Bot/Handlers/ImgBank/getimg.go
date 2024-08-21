package imgbank

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	config "github.com/dabi-ngin/discgo-bot/Config"
	database "github.com/dabi-ngin/discgo-bot/Database"
	helpers "github.com/dabi-ngin/discgo-bot/Helpers"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
	structs "github.com/dabi-ngin/discgo-bot/Structs"
)

func GetImg(message *discordgo.MessageCreate, self structs.BangCommand) error {

	bangCommand := helpers.CheckForBangCommand(message.Content)
	if bangCommand == "" {
		helpers.SendUserError(message, "Invalid Command")
		return errors.New("couldn't obtain bang command")
	}

	imgCat, err := database.GetImgCategory(message.GuildID, bangCommand)
	if err != nil {
		helpers.SendUserError(message, "Invalid Category")
		return errors.New("unable to get gif category")
	}

	imgUrl, err := database.GetRandomImage(message.GuildID, imgCat.ID)
	if err != nil {
		helpers.SendUserError(message, "Couldn't find an Image")
		return err
	}

	_, err = config.Session.ChannelMessageSend(message.ChannelID, imgUrl)
	if err != nil {
		logger.Error(message.GuildID, err)
		helpers.SendUserError(message, "Couldn't send Image")
		return err
	}

	return nil

}
