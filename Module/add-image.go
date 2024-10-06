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

type AddImage struct {
	ImageCategory string
}

func NewAddImage(imageCategory string) *AddImage {
	return &AddImage{
		ImageCategory: imageCategory,
	}
}

func (s AddImage) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name: "Add to !" + s.ImageCategory + " Bank",
		Type: discordgo.MessageApplicationCommand,
	}
}

func (s AddImage) PermissionRequirement() int {
	return config.CommandLevelUser
}

func (s AddImage) Complexity() int {
	return config.TRIVIAL_TASK
}

func (s AddImage) Execute(i *discordgo.InteractionCreate, correlationId string) {
	_, message := discord.GetAssociatedMessageFromInteraction(i)

	imgUrl := helpers.GetImageFromMessage(message, "")
	if imgUrl == "" {
		logger.ErrorText(i.GuildID, "no image found")
		cache.InteractionComplete(correlationId)
		return
	}

	err := database.AddImg(message, s.ImageCategory, imgUrl)
	if err != nil {
		discord.Interactions_SendError(i, "")
		cache.InteractionComplete(correlationId)
		return
	}

	discord.Interaction_SendPublicMessage(i, "Image Added to !"+s.ImageCategory+" Bank", "")
	cache.InteractionComplete(correlationId)
}
