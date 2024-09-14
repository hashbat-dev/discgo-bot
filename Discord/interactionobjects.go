package discord

import (
	"github.com/bwmarrin/discordgo"
)

type SlashResponse struct {
	Complexity int
	Execute    func(i *discordgo.InteractionCreate, correlationId string)
}

var InteractionResponseHandlers map[string]SlashResponse = make(map[string]SlashResponse)

func CreateSelectMenu(selectMenu discordgo.SelectMenu, correlationId string, complexity int, executeFunction func(i *discordgo.InteractionCreate, correlationId string)) discordgo.SelectMenu {
	ObjectID := selectMenu.CustomID
	selectMenu.CustomID = selectMenu.CustomID + "|" + correlationId
	if _, exists := InteractionResponseHandlers[ObjectID]; !exists {
		InteractionResponseHandlers[ObjectID] = SlashResponse{
			Complexity: complexity,
			Execute:    executeFunction,
		}
	}
	return selectMenu
}
