package handlers

import (
	"github.com/bwmarrin/discordgo"
	config "github.com/hashbat-dev/discgo-bot/Config"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

func HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := config.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		logger.Error(i.GuildID, err)
		return
	}
	DispatchTask(&Task{
		CommandType: config.CommandTypeNewInteraction,
		Complexity:  config.TRIVIAL_TASK,
		Interaction: i,
	})
}
