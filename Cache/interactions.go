package cache

import (
	"reflect"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	config "github.com/dabi-ngin/discgo-bot/Config"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
	"github.com/google/uuid"
)

var ActiveInteractions map[string]InteractionCache = make(map[string]InteractionCache)

type InteractionCache struct {
	StartInteraction discordgo.InteractionCreate
	Interactions     map[string]discordgo.InteractionCreate
	Values           InteractionValues
	Started          time.Time
}

type InteractionValues struct {
	String  map[string]string
	Integer map[string]int64
	Bool    map[string]bool
	User    map[string]*discordgo.User
	Channel map[string]*discordgo.Channel
	Role    map[string]*discordgo.Role
	Number  map[string]float64
}

func AddInteraction(i *discordgo.InteractionCreate, commandName string) string {
	// Create the Interaction object
	correlationId := uuid.New().String()

	logger.Info(i.GuildID, "Interaction ID: [%v] Started Processing, Command: [%v]", correlationId, commandName)
	ActiveInteractions[correlationId] = InteractionCache{
		StartInteraction: *i,
		Interactions: map[string]discordgo.InteractionCreate{
			i.ID: *i,
		},
		Values: InteractionValues{
			String:  make(map[string]string),
			Integer: make(map[string]int64),
			Bool:    make(map[string]bool),
			User:    make(map[string]*discordgo.User),
			Channel: make(map[string]*discordgo.Channel),
			Role:    make(map[string]*discordgo.Role),
			Number:  make(map[string]float64),
		},
		Started: time.Now(),
	}

	// Add any Options in the Interaction to the Map Cache
	if len(i.ApplicationCommandData().Options) > 0 {
		for _, option := range i.ApplicationCommandData().Options {
			switch option.Type {
			case discordgo.ApplicationCommandOptionString:
				ActiveInteractions[correlationId].Values.String[option.Name] = option.StringValue()
				logger.Info(i.GuildID, "Interaction ID: [%v] Obtained String Value for [%v]: %v", correlationId, option.Name, ActiveInteractions[correlationId].Values.String[option.Name])
			case discordgo.ApplicationCommandOptionInteger:
				ActiveInteractions[correlationId].Values.Integer[option.Name] = option.IntValue()
				logger.Info(i.GuildID, "Interaction ID: [%v] Obtained Integer Value for [%v]: %v", correlationId, option.Name, ActiveInteractions[correlationId].Values.Integer[option.Name])
			case discordgo.ApplicationCommandOptionBoolean:
				ActiveInteractions[correlationId].Values.Bool[option.Name] = option.BoolValue()
				logger.Info(i.GuildID, "Interaction ID: [%v] Obtained Bool Value for [%v]: %v", correlationId, option.Name, ActiveInteractions[correlationId].Values.Bool[option.Name])
			case discordgo.ApplicationCommandOptionUser:
				ActiveInteractions[correlationId].Values.User[option.Name] = option.UserValue(config.Session)
				logger.Info(i.GuildID, "Interaction ID: [%v] Obtained User Value for [%v]: %v", correlationId, option.Name, ActiveInteractions[correlationId].Values.User[option.Name])
			case discordgo.ApplicationCommandOptionChannel:
				ActiveInteractions[correlationId].Values.Channel[option.Name] = option.ChannelValue(config.Session)
				logger.Info(i.GuildID, "Interaction ID: [%v] Obtained Channel Value for [%v]: %v", correlationId, option.Name, ActiveInteractions[correlationId].Values.Channel[option.Name])
			case discordgo.ApplicationCommandOptionRole:
				ActiveInteractions[correlationId].Values.Role[option.Name] = option.RoleValue(config.Session, i.GuildID)
				logger.Info(i.GuildID, "Interaction ID: [%v] Obtained Role Value for [%v]: %v", correlationId, option.Name, ActiveInteractions[correlationId].Values.Role[option.Name])
			case discordgo.ApplicationCommandOptionNumber:
				ActiveInteractions[correlationId].Values.Number[option.Name] = option.FloatValue()
				logger.Info(i.GuildID, "Interaction ID: [%v] Obtained Number Value for [%v]: %v", correlationId, option.Name, ActiveInteractions[correlationId].Values.Number[option.Name])
			default:
				logger.ErrorText(i.GuildID, "AddInteraction encountered an unknown data type [%v]", option.Type.String())
			}
		}
	}

	return correlationId
}

func UpdateInteraction(correlationId string, i *discordgo.InteractionCreate) {
	// Check we have the associated Interaction in the Cache
	if _, exists := ActiveInteractions[correlationId]; !exists {
		logger.ErrorText(i.GuildID, "Interaction Update could not find the associated CorrelationId [%v]", correlationId)
		return
	}

	// Add the Interaction
	ActiveInteractions[correlationId].Interactions[i.ID] = *i

	switch i.Type {
	// Slash Commands (directly) or Autocomplete
	case discordgo.InteractionApplicationCommand, discordgo.InteractionApplicationCommandAutocomplete:
		if len(i.ApplicationCommandData().Options) > 0 {
			for _, option := range i.ApplicationCommandData().Options {
				switch option.Type {
				case discordgo.ApplicationCommandOptionString:
					ActiveInteractions[correlationId].Values.String[option.Name] = option.StringValue()
					logger.Info(i.GuildID, "Interaction ID: [%v] Obtained String Value for [%v]: %v", correlationId, option.Name, ActiveInteractions[correlationId].Values.String[option.Name])
				case discordgo.ApplicationCommandOptionInteger:
					ActiveInteractions[correlationId].Values.Integer[option.Name] = option.IntValue()
					logger.Info(i.GuildID, "Interaction ID: [%v] Obtained Integer Value for [%v]: %v", correlationId, option.Name, ActiveInteractions[correlationId].Values.Integer[option.Name])
				case discordgo.ApplicationCommandOptionBoolean:
					ActiveInteractions[correlationId].Values.Bool[option.Name] = option.BoolValue()
					logger.Info(i.GuildID, "Interaction ID: [%v] Obtained Bool Value for [%v]: %v", correlationId, option.Name, ActiveInteractions[correlationId].Values.Bool[option.Name])
				case discordgo.ApplicationCommandOptionUser:
					ActiveInteractions[correlationId].Values.User[option.Name] = option.UserValue(config.Session)
					logger.Info(i.GuildID, "Interaction ID: [%v] Obtained User Value for [%v]: %v", correlationId, option.Name, ActiveInteractions[correlationId].Values.User[option.Name])
				case discordgo.ApplicationCommandOptionChannel:
					ActiveInteractions[correlationId].Values.Channel[option.Name] = option.ChannelValue(config.Session)
					logger.Info(i.GuildID, "Interaction ID: [%v] Obtained Channel Value for [%v]: %v", correlationId, option.Name, ActiveInteractions[correlationId].Values.Channel[option.Name])
				case discordgo.ApplicationCommandOptionRole:
					ActiveInteractions[correlationId].Values.Role[option.Name] = option.RoleValue(config.Session, i.GuildID)
					logger.Info(i.GuildID, "Interaction ID: [%v] Obtained Role Value for [%v]: %v", correlationId, option.Name, ActiveInteractions[correlationId].Values.Role[option.Name])
				case discordgo.ApplicationCommandOptionNumber:
					ActiveInteractions[correlationId].Values.Number[option.Name] = option.FloatValue()
					logger.Info(i.GuildID, "Interaction ID: [%v] Obtained Number Value for [%v]: %v", correlationId, option.Name, ActiveInteractions[correlationId].Values.Number[option.Name])
				default:
					logger.ErrorText(i.GuildID, "UpdateInteraction encountered an unknown CommandData data type: [%v]", option.Type.String())
				}
			}
		}

	// Message Component (selects/buttons etc.)
	case discordgo.InteractionMessageComponent:

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
				ActiveInteractions[correlationId].Values.String[objectID] = data.Values[0]
				logger.Info(i.GuildID, "Interaction ID: [%v] Obtained String Value for [%v]: %v", correlationId, objectID, ActiveInteractions[correlationId].Values.String[objectID])
			} else {
				// Handle button interactions
				ActiveInteractions[correlationId].Values.Bool[objectID] = true
				logger.Info(i.GuildID, "Interaction ID: [%v] Obtained Button Selected (Bool=True) Value for [%v]: %v", correlationId, objectID, ActiveInteractions[correlationId].Values.Bool[objectID])
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
							ActiveInteractions[correlationId].Values.String[objectID] = v
							logger.Info(i.GuildID, "Interaction ID: [%v] Obtained String Value for [%v]: %v", correlationId, objectID, ActiveInteractions[correlationId].Values.String[objectID])
						case int:
							ActiveInteractions[correlationId].Values.Integer[objectID] = int64(v)
							logger.Info(i.GuildID, "Interaction ID: [%v] Obtained Integer Value for [%v]: %v", correlationId, objectID, ActiveInteractions[correlationId].Values.Integer[objectID])
						case float64:
							ActiveInteractions[correlationId].Values.Number[objectID] = v
							logger.Info(i.GuildID, "Interaction ID: [%v] Obtained Number Value for [%v]: %v", correlationId, objectID, ActiveInteractions[correlationId].Values.Number[objectID])
						case bool:
							ActiveInteractions[correlationId].Values.Bool[objectID] = v
							logger.Info(i.GuildID, "Interaction ID: [%v] Obtained Bool Value for [%v]: %v", correlationId, objectID, ActiveInteractions[correlationId].Values.Bool[objectID])
						case *discordgo.User:
							ActiveInteractions[correlationId].Values.User[objectID] = v
							logger.Info(i.GuildID, "Interaction ID: [%v] Obtained User Value for [%v]: %v", correlationId, objectID, ActiveInteractions[correlationId].Values.User[objectID])
						case *discordgo.Channel:
							ActiveInteractions[correlationId].Values.Channel[objectID] = v
							logger.Info(i.GuildID, "Interaction ID: [%v] Obtained Channel Value for [%v]: %v", correlationId, objectID, ActiveInteractions[correlationId].Values.Channel[objectID])
						case *discordgo.Role:
							ActiveInteractions[correlationId].Values.Role[objectID] = v
							logger.Info(i.GuildID, "Interaction ID: [%v] Obtained Role Value for [%v]: %v", correlationId, objectID, ActiveInteractions[correlationId].Values.Role[objectID])
						default:
							logger.ErrorText(i.GuildID, "Interaction ID: [%v] Unable to obtain Type Assertion for value", i.GuildID)
						}
					}
				}
			}
		}

	default:
		logger.ErrorText(i.GuildID, "UpdateInteraction encountered an unknown data type: [%v]", i.Type.String())
	}
}

func InteractionComplete(correlationId string) {
	delete(ActiveInteractions, correlationId)
}
