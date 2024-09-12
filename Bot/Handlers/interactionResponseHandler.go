package handlers

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	config "github.com/dabi-ngin/discgo-bot/Config"
	discord "github.com/dabi-ngin/discgo-bot/Discord"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
)

// Handles responses to Interactions

func HandleInteractionResponse(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionMessageComponent {
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
			DispatchTask(&WorkerItem{
				CommandType: config.CommandTypeSlashResponse,
				Complexity:  responseObject.Complexity,
				SlashCommandResponse: SlashCommandResponseWorker{
					Interaction:   i,
					ObjectID:      objectId,
					CorrelationID: correlationId,
				},
			})
		} else {
			logger.ErrorText(i.GuildID, "Unknown Interaction Response ObjectID [%v]", objectId)

			// Generic Error back to the user
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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
}
