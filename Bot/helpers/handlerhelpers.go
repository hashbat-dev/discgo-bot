package helpers

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

func SkipProcessing(session *discordgo.Session, newMsg *discordgo.MessageCreate, updMsg *discordgo.MessageUpdate) (bool, error) {
	// Get the Message object out of the appropriate Message object
	var inboundMessage *discordgo.Message
	if newMsg != nil {
		inboundMessage = newMsg.Message
	} else if updMsg != nil {
		inboundMessage = updMsg.Message
	} else {
		err := errors.New("no discord message object provided")
		return false, err
	}

	// Is the Message Author nil?
	if inboundMessage.Author == nil {
		return true, nil
	}

	// Is this a Bot posted Message?
	if inboundMessage.Author.ID == session.State.User.ID {
		return true, nil
	}

	return false, nil
}

func GetOptionMap(interaction *discordgo.InteractionCreate) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	options := interaction.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	return optionMap
}

func GetOptionUserIDValue(optionMap map[string]*discordgo.ApplicationCommandInteractionDataOption, property string) string {
	if opt, ok := optionMap[property]; ok {
		if opt.Value != "" {
			return opt.UserValue(nil).ID
		}
	}
	return ""
}

func GetOptionStringValue(optionMap map[string]*discordgo.ApplicationCommandInteractionDataOption, property string) string {
	if opt, ok := optionMap[property]; ok {
		if opt.Value != "" {
			return opt.StringValue()
		}
	}
	return ""
}

func GetOptionBoolValue(optionMap map[string]*discordgo.ApplicationCommandInteractionDataOption, property string) bool {
	if opt, ok := optionMap[property]; ok {
		return opt.BoolValue()
	}
	return false
}

func GetOptionIntValue(optionMap map[string]*discordgo.ApplicationCommandInteractionDataOption, property string) int {
	if opt, ok := optionMap[property]; ok {
		if opt.Value != "" {
			return int(opt.IntValue())
		}
	}
	return 0
}
