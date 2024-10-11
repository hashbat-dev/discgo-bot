package interactions

import (
	"reflect"
	"strings"

	"github.com/bwmarrin/discordgo"
	config "github.com/hashbat-dev/discgo-bot/Config"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

func ExtractValues(correlationId string, i *discordgo.InteractionCreate) {
	if _, exists := Active[correlationId]; !exists {
		logger.ErrorText(i.GuildID, "[Interactions] ExtractValues could not find CorrelationID in Active cache [%v]", correlationId)
		return
	}

	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		// Slash Commands - Initial creation
		if len(i.ApplicationCommandData().Options) > 0 {
			for _, option := range i.ApplicationCommandData().Options {
				switch option.Type {
				case discordgo.ApplicationCommandOptionString:
					Active[correlationId].Values.String[option.Name] = option.StringValue()
					logger.Info(i.GuildID, "[Interactions] %s :: String [%v]: %v", correlationId, option.Name, Active[correlationId].Values.String[option.Name])
				case discordgo.ApplicationCommandOptionInteger:
					Active[correlationId].Values.Integer[option.Name] = option.IntValue()
					logger.Info(i.GuildID, "[Interactions] %s :: Integer [%v]: %v", correlationId, option.Name, Active[correlationId].Values.Integer[option.Name])
				case discordgo.ApplicationCommandOptionBoolean:
					Active[correlationId].Values.Bool[option.Name] = option.BoolValue()
					logger.Info(i.GuildID, "[Interactions] %s :: Bool [%v]: %v", correlationId, option.Name, Active[correlationId].Values.Bool[option.Name])
				case discordgo.ApplicationCommandOptionUser:
					Active[correlationId].Values.User[option.Name] = option.UserValue(config.Session)
					logger.Info(i.GuildID, "[Interactions] %s :: User [%v]: %v", correlationId, option.Name, Active[correlationId].Values.User[option.Name])
				case discordgo.ApplicationCommandOptionChannel:
					Active[correlationId].Values.Channel[option.Name] = option.ChannelValue(config.Session)
					logger.Info(i.GuildID, "[Interactions] %s :: Channel [%v]: %v", correlationId, option.Name, Active[correlationId].Values.Channel[option.Name])
				case discordgo.ApplicationCommandOptionRole:
					Active[correlationId].Values.Role[option.Name] = option.RoleValue(config.Session, i.GuildID)
					logger.Info(i.GuildID, "[Interactions] %s :: Role [%v]: %v", correlationId, option.Name, Active[correlationId].Values.Role[option.Name])
				case discordgo.ApplicationCommandOptionNumber:
					Active[correlationId].Values.Number[option.Name] = option.FloatValue()
					logger.Info(i.GuildID, "[Interactions] %s :: Number [%v]: %v", correlationId, option.Name, Active[correlationId].Values.Number[option.Name])
				default:
					logger.ErrorText(i.GuildID, "[Interactions] %s :: ExtractValues encountered an unknown CommandData data type: [%v]", correlationId, option.Type.String())
				}
			}
		}

	case discordgo.InteractionModalSubmit:
		// Modal Submission
		data := i.ModalSubmitData()

		for _, actionRow := range data.Components {
			if row, ok := actionRow.(*discordgo.ActionsRow); ok {
				for _, comp := range row.Components {
					if input, ok := comp.(*discordgo.TextInput); ok {
						objectId := input.CustomID
						if strings.Contains(objectId, "|") {
							objectId = strings.Split(objectId, "|")[0]
						}
						Active[correlationId].Values.String[objectId] = input.Value
						logger.Info(i.GuildID, "[Interactions] %s :: String [%s]: %s", correlationId, objectId, Active[correlationId].Values.String[objectId])
					}
				}
			}
		}
	case discordgo.InteractionMessageComponent:
		// Message Component - Selects, buttons, etc.
		switch data := i.Interaction.Data.(type) {
		case *discordgo.MessageComponentInteractionData:

			var objectID string
			if strings.Contains(data.CustomID, "|") {
				objectID = strings.Split(data.CustomID, "|")[0]
			} else {
				objectID = data.CustomID
			}

			if len(data.Values) > 0 {
				// Handle select menu values
				Active[correlationId].Values.String[objectID] = data.Values[0]
				logger.Info(i.GuildID, "[Interactions] %s :: String [%v]: %v", correlationId, objectID, Active[correlationId].Values.String[objectID])
			} else {
				// Handle button interactions
				Active[correlationId].Values.Bool[objectID] = true
				logger.Info(i.GuildID, "[Interactions] %s :: Button [%v]: %v", correlationId, objectID, Active[correlationId].Values.Bool[objectID])
			}
		default:
			// Handle situations where a type could not be asserted yet we can access the returned value.
			// This will always be added as a string.
			val := reflect.ValueOf(i.Interaction.Data)
			if val.Kind() == reflect.Ptr && !val.IsNil() {
				val = val.Elem()
			}

			if val.Kind() == reflect.Struct {
				// Extract CustomID
				var providedObjectID string
				if field := val.FieldByName("CustomID"); field.IsValid() && field.Kind() == reflect.String {
					providedObjectID = field.String()
				}

				var objectID string
				if strings.Contains(providedObjectID, "|") {
					objectID = strings.Split(providedObjectID, "|")[0]
				} else {
					objectID = providedObjectID
				}

				// Extract Values
				if field := val.FieldByName("Values"); field.IsValid() && field.Kind() == reflect.Slice {
					for j := 0; j < field.Len(); j++ {
						value := field.Index(j).Interface()
						switch v := value.(type) {
						case string:
							Active[correlationId].Values.String[objectID] = v
							logger.Info(i.GuildID, "[Interactions] %s :: String [%v]: %v", correlationId, objectID, Active[correlationId].Values.String[objectID])
						case int:
							Active[correlationId].Values.Integer[objectID] = int64(v)
							logger.Info(i.GuildID, "[Interactions] %s :: Integer [%v]: %v", correlationId, objectID, Active[correlationId].Values.Integer[objectID])
						case float64:
							Active[correlationId].Values.Number[objectID] = v
							logger.Info(i.GuildID, "[Interactions] %s :: Number [%v]: %v", correlationId, objectID, Active[correlationId].Values.Number[objectID])
						case bool:
							Active[correlationId].Values.Bool[objectID] = v
							logger.Info(i.GuildID, "[Interactions] %s :: Bool [%v]: %v", correlationId, objectID, Active[correlationId].Values.Bool[objectID])
						case *discordgo.User:
							Active[correlationId].Values.User[objectID] = v
							logger.Info(i.GuildID, "[Interactions] %s :: User [%v]: %v", correlationId, objectID, Active[correlationId].Values.User[objectID])
						case *discordgo.Channel:
							Active[correlationId].Values.Channel[objectID] = v
							logger.Info(i.GuildID, "[Interactions] %s :: Channel [%v]: %v", correlationId, objectID, Active[correlationId].Values.Channel[objectID])
						case *discordgo.Role:
							Active[correlationId].Values.Role[objectID] = v
							logger.Info(i.GuildID, "[Interactions] %s :: Role [%v]: %v", correlationId, objectID, Active[correlationId].Values.Role[objectID])
						default:
							logger.ErrorText(i.GuildID, "[Interactions] %s :: Unable to obtain Type Assertion for value", i.GuildID)
						}
					}
				}
			}
		}

	default:
		logger.ErrorText(i.GuildID, "UpdateInteraction encountered an unknown data type: [%v]", i.Type.String())
	}
}
