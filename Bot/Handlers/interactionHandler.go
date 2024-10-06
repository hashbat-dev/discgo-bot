package handlers

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	cache "github.com/hashbat-dev/discgo-bot/Cache"
	config "github.com/hashbat-dev/discgo-bot/Config"
	discord "github.com/hashbat-dev/discgo-bot/Discord"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
	module "github.com/hashbat-dev/discgo-bot/Module"
)

func HandleInteractionResponse(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if config.ServiceSettings.ISDEV != cache.ActiveGuilds[i.GuildID].IsDev {
		return
	}

	switch i.Type {
	case discordgo.InteractionMessageComponent:
		handleInteractionMessageComponent(i)
	case discordgo.InteractionModalSubmit:
		handleInteractionModalSubmit(i)
	default:
		handleGenericInteraction(i)
	}
}

func handleGenericInteraction(i *discordgo.InteractionCreate) {
	logger.Debug(i.GuildID, "Interaction Recieved: handleGenericInteraction")
	cmdName := i.ApplicationCommandData().Name

	foundCmd := false
	for _, mod := range module.ModuleList {
		if mod.Command().Name == cmdName {
			foundCmd = true

			if !ConfirmPermissions(i, mod.PermissionRequirement()) {
				logger.Event(i.GuildID, "User [ID: %s, UserName: %s] was blocked from using command [%s]", i.Member.User.ID, i.Member.User.Username, cmdName)
				discord.Interactions_SendMessage(i, "Permission Denied", "You do not have permission to use this command.")
			} else {
				DispatchTask(&Task{
					CommandType: config.CommandTypeModule,
					Complexity:  mod.Complexity(),
					ModuleDetails: &ModuleDetails{
						Interaction: i,
						Module:      mod,
					},
				})
			}

			break
		}
	}

	if !foundCmd {
		logger.ErrorText(i.GuildID, "No Handler found for Slash Command: %v", cmdName)
	}
}

func handleInteractionMessageComponent(i *discordgo.InteractionCreate) {
	logger.Debug(i.GuildID, "Interaction Recieved: handleInteractionMessageComponent")
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
			CommandType: config.CommandTypeModuleResponse,
			Complexity:  responseObject.Complexity,
			ModuleResponseDetails: &ModuleResponseDetails{
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
	logger.Debug(i.GuildID, "Interaction Recieved: handleInteractionModalSubmit")
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
			CommandType: config.CommandTypeModuleResponse,
			Complexity:  responseObject.Complexity,
			ModuleResponseDetails: &ModuleResponseDetails{
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
