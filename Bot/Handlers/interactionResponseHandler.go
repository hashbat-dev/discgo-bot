package handlers

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	config "github.com/hashbat-dev/discgo-bot/Config"
	discord "github.com/hashbat-dev/discgo-bot/Discord"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

func HandleInteractionResponse(s *discordgo.Session, i *discordgo.InteractionCreate) {
	handleInteractionResponseQueue <- i
}

// Handles responses to Interactions
func ProcessInteractionResponse(i *discordgo.InteractionCreate) {
	logger.Debug(i.GuildID, "HandleInteractionResponse")
	switch i.Type {
	case discordgo.InteractionMessageComponent:
		handleInteractionMessageComponent(i)
	case discordgo.InteractionModalSubmit:
		handleInteractionModalSubmit(i)
	default:
		// Not an issue, this will be handled elsewhere
	}
}

func handleInteractionMessageComponent(i *discordgo.InteractionCreate) {
	inboundObjectId := i.MessageComponentData().CustomID

	// Make sure we have the format of <ObjectID>|<CorrelationID>
	if !strings.Contains(inboundObjectId, "|") {
		logger.ErrorText(i.GuildID, "Interaction response in wrong format, no CorrelationId provided")
		return
	}

	// Get the ObjectID / CorrelationID from the assigned ID
	splitObjectId := strings.Split(inboundObjectId, "|")
	objectId := splitObjectId[0]
	correlationId := splitObjectId[1]
	if responseObject, exists := discord.InteractionResponseHandlers[objectId]; exists {
		DispatchTask(&Task{
			CommandType: config.CommandTypeSlashResponse,
			Complexity:  responseObject.Complexity,
			SlashResponseDetails: &SlashResponseDetails{
				Interaction:   i,
				ObjectID:      objectId,
				CorrelationID: correlationId,
			},
		})
	} else {
		logger.ErrorText(i.GuildID, "Unknown Interaction Response ObjectID [%v]", objectId)

		// Generic Error back to the user
		err := config.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error finding Command",
			},
		})
		if err != nil {
			logger.Error(i.GuildID, err)
		}
	}
}

func handleInteractionModalSubmit(i *discordgo.InteractionCreate) {
	inboundObjectId := i.ModalSubmitData().CustomID

	// Make sure we have the format of <ObjectID>|<CorrelationID>
	if !strings.Contains(inboundObjectId, "|") {
		logger.ErrorText(i.GuildID, "Interaction response in wrong format, no CorrelationId provided")
		return
	}

	// Get the ObjectID / CorrelationID from the assigned ID
	splitObjectId := strings.Split(inboundObjectId, "|")
	objectId := splitObjectId[0]
	correlationId := splitObjectId[1]
	if responseObject, exists := discord.InteractionResponseHandlers[objectId]; exists {
		DispatchTask(&Task{
			CommandType: config.CommandTypeSlashResponse,
			Complexity:  responseObject.Complexity,
			SlashResponseDetails: &SlashResponseDetails{
				Interaction:   i,
				ObjectID:      objectId,
				CorrelationID: correlationId,
			},
		})
	} else {
		logger.ErrorText(i.GuildID, "Unknown Interaction Response ObjectID [%v]", objectId)

		// Generic Error back to the user
		err := config.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error finding Command",
			},
		})
		if err != nil {
			logger.Error(i.GuildID, err)
		}
	}
}
