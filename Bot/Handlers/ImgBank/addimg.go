package imgbank

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	database "github.com/dabi-ngin/discgo-bot/Database"
	helpers "github.com/dabi-ngin/discgo-bot/Helpers"
	structs "github.com/dabi-ngin/discgo-bot/Structs"
)

func AddImg(message *discordgo.MessageCreate, self structs.BangCommand) error {

	bangCommand := helpers.CheckForBangCommand(message.Content)
	if bangCommand == "" {
		helpers.SendUserError(message, "Invalid Command")
		return errors.New("couldn't obtain bang command")
	}

	command := self.ImgCategory

	imageUrl := helpers.GetImageFromMessage(message.ReferencedMessage, "")
	if imageUrl == "" {
		helpers.SendUserError(message, "Couldn't find an Image")
		return errors.New("did not get image from message")
	}

	err := database.AddImg(message, command, imageUrl)
	if err != nil {
		helpers.SendUserError(message, "Couldn't add new Image")
		return err
	}

	return nil

}
