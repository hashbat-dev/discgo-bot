package module

import (
	"github.com/bwmarrin/discordgo"
	cache "github.com/hashbat-dev/discgo-bot/Cache"
	config "github.com/hashbat-dev/discgo-bot/Config"
	database "github.com/hashbat-dev/discgo-bot/Database"
	discord "github.com/hashbat-dev/discgo-bot/Discord"
	helpers "github.com/hashbat-dev/discgo-bot/Helpers"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

type DeleteImage struct {
	ImageCategory string
}

func NewDeleteImage(imageCategory string) *DeleteImage {
	return &DeleteImage{
		ImageCategory: imageCategory,
	}
}

func (s DeleteImage) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name: "Delete from !" + s.ImageCategory + " Bank",
		Type: discordgo.MessageApplicationCommand,
	}
}

func (s DeleteImage) PermissionRequirement() int {
	return config.CommandLevelBotAdmin
}

func (s DeleteImage) Complexity() int {
	return config.TRIVIAL_TASK
}

func (s DeleteImage) Execute(i *discordgo.InteractionCreate, correlationId string) {
	_, message := discord.GetAssociatedMessageFromInteraction(i)

	imgUrl := helpers.GetImageFromMessage(message, "")
	if imgUrl == "" {
		logger.ErrorText(i.GuildID, "no image found")
		cache.InteractionComplete(correlationId)
		return
	}

	imgCat, err := database.GetImgCategory(message.GuildID, s.ImageCategory)
	if err != nil {
		logger.ErrorText(i.GuildID, "unable to get image category")
		cache.InteractionComplete(correlationId)
		return
	}

	imgStorage, err := database.GetImgStorage(message.GuildID, imgUrl)
	if err != nil {
		logger.ErrorText(i.GuildID, "unable to get image storage")
		cache.InteractionComplete(correlationId)
		return
	}

	imgGuildLink, err := database.GetImgGuildLink(message.GuildID, imgCat, imgStorage)
	if err != nil {
		logger.ErrorText(i.GuildID, "unable to get image guild link")
		cache.InteractionComplete(correlationId)
		return
	}

	err = database.DeleteGuildLink(imgGuildLink)
	if err != nil {
		logger.ErrorText(i.GuildID, "unable to delete image")
		cache.InteractionComplete(correlationId)
		return
	}

	discord.Interactions_SendMessage(i, "Delete Image", "Image deleted from the !"+s.ImageCategory+" bank")
	cache.InteractionComplete(correlationId)
}
