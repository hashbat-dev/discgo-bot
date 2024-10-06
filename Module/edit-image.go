package module

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	cache "github.com/hashbat-dev/discgo-bot/Cache"
	config "github.com/hashbat-dev/discgo-bot/Config"
	discord "github.com/hashbat-dev/discgo-bot/Discord"
	editmodule "github.com/hashbat-dev/discgo-bot/EditModule"
	helpers "github.com/hashbat-dev/discgo-bot/Helpers"
	imgwork "github.com/hashbat-dev/discgo-bot/ImgWork"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

type Support struct{}

func (s Support) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name: "Edit/Meme Image",
		Type: discordgo.MessageApplicationCommand,
	}
}

func (s Support) PermissionRequirement() int {
	return config.CommandLevelUser
}

func (s Support) Complexity() int {
	return config.TRIVIAL_TASK
}

func (s Support) Execute(i *discordgo.InteractionCreate, correlationId string) {
	// 1. Get the Message object associated with the Interaction request
	messageID, message := discord.GetAssociatedMessageFromInteraction(i)

	// 2. Check there's an associated Image
	imgUrl := helpers.GetImageFromMessage(message, "")
	if imgUrl == "" {
		discord.Interactions_SendError(i, "No image found in message")
		return
	}

	imgExtension := imgwork.GetExtensionFromURL(imgUrl)
	if imgExtension == "" {
		discord.Interactions_SendError(i, fmt.Sprintf("Invalid image extension (%s)", imgExtension))
		return
	}

	// => Store these in the Interactions cache for later
	cache.ActiveInteractions[correlationId].Values.String["imgMessageId"] = messageID
	cache.ActiveInteractions[correlationId].Values.String["imgUrl"] = imgUrl
	cache.ActiveInteractions[correlationId].Values.String["imgExtension"] = imgExtension

	// 3. Get a list of Edit options
	var selectOptions []discordgo.SelectMenuOption
	for _, edit := range editmodule.EditList {
		selectOptions = append(selectOptions, discordgo.SelectMenuOption{
			Label: edit.SelectName(),
			Value: helpers.LettersNumbersAndDashesOnly(edit.SelectName()),
			Emoji: edit.Emoji(),
		})
	}

	// 4. Create the Interaction Objects
	selectMenu := discord.CreateSelectMenu(discordgo.SelectMenu{
		CustomID:    "edit-image_select",
		Options:     selectOptions,
		Placeholder: "Select Edit type...",
	}, correlationId, config.CPU_BOUND_TASK, SelectEditModule)

	actionRow := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			selectMenu,
		},
	}

	// 5. Send the Select menu response
	err := config.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Components: []discordgo.MessageComponent{actionRow},
			Flags:      discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		logger.Error(i.GuildID, err)
	}

}

func SelectEditModule(i *discordgo.InteractionCreate, correlationId string) {

}
