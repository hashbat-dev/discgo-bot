package handlers

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	config "github.com/hashbat-dev/discgo-bot/Config"
	interactions "github.com/hashbat-dev/discgo-bot/Interactions"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

func WorkerNewInteraction(i *discordgo.InteractionCreate) {
	interactions.SendUserText(i, "Starting Request...")
	objectId, correlationId := getNameAndCorrelationId(i)
	if objectId == "" || correlationId == "" {
		logger.ErrorText(i.GuildID, "[Interactions] No ObjectID or CorrelationID, exiting HandleInteraction [ObjectID: %s] [CorrelationID: %s]", objectId, correlationId)
		interactions.SendUserText(i, "Error Processing Request [Code: 1001]")
		return
	}

	interactions.SendUserText(i, "Getting Details...")
	interactions.UpdateCache(correlationId, i)
	var handler interactions.HandleInteraction
	if foundHandle, exists := interactions.Handlers[objectId]; exists {
		handler = foundHandle
	} else {
		logger.ErrorText(i.GuildID, "[Interactions] No Handler found for the ObjectID: %s", objectId)
		interactions.SendUserText(i, "Error Processing Request [Code: 1002]")
		return
	}

	interactions.SendUserText(i, "Begin Processing...")
	DispatchTask(&Task{
		CommandType: config.CommandTypeInteractionCall,
		Complexity:  handler.Complexity,
		InteractionCall: InteractionCall{
			ObjectID:      objectId,
			CorrelationID: correlationId,
			Execute:       handler.Execute,
		},
	})
}

func getNameAndCorrelationId(i *discordgo.InteractionCreate) (string, string) {
	var objectId, correlationId string
	switch i.Type {
	case discordgo.InteractionMessageComponent:
		// Message response object (Selects, Buttons, etc.)
		objectId = i.MessageComponentData().CustomID
		if !strings.Contains(objectId, "|") {
			correlationId = uuid.New().String()
			logger.Info(i.GuildID, "[Interactions] [%v] New CorrelationID Generated: %s", i.Type, correlationId)
		} else {
			splitObjectId := strings.Split(objectId, "|")
			objectId = splitObjectId[0]
			correlationId = splitObjectId[1]
		}
	case discordgo.InteractionModalSubmit:
		// Modals (pop-up form with multiple input options)
		objectId = i.ModalSubmitData().CustomID
		if !strings.Contains(objectId, "|") {
			correlationId = uuid.New().String()
			logger.Info(i.GuildID, "[Interactions] [%v] New CorrelationID Generated: %s", i.Type, correlationId)
		} else {
			splitObjectId := strings.Split(objectId, "|")
			objectId = splitObjectId[0]
			correlationId = splitObjectId[1]
		}
	case discordgo.InteractionApplicationCommand:
		// New Interaction Request
		objectId = i.ApplicationCommandData().Name
		if !strings.Contains(objectId, "|") {
			correlationId = uuid.New().String() // Create a new CorrelationID
			logger.Info(i.GuildID, "[Interactions] [%v] New CorrelationID Generated: %s", i.Type, correlationId)
		} else {
			splitObjectId := strings.Split(objectId, "|")
			objectId = splitObjectId[0]
			correlationId = splitObjectId[1]
		}
	default:
		logger.ErrorText(i.GuildID, "[Interactions] Unhandled Interaction Type: %s", i.Type.String())
	}
	return objectId, correlationId
}
