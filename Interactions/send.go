package interactions

import (
	"github.com/bwmarrin/discordgo"
	config "github.com/hashbat-dev/discgo-bot/Config"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

func SendUserText(i *discordgo.InteractionCreate, updateText string) {
	_, err := config.Session.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: updateText,
	})
	if err != nil {
		logger.Error(i.GuildID, err)
	}
}
