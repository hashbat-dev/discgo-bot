package slash

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	cache "github.com/hashbat-dev/discgo-bot/Cache"
	config "github.com/hashbat-dev/discgo-bot/Config"
	database "github.com/hashbat-dev/discgo-bot/Database"
	discord "github.com/hashbat-dev/discgo-bot/Discord"
	helpers "github.com/hashbat-dev/discgo-bot/Helpers"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

func WowLeaderboard(i *discordgo.InteractionCreate, correlationId string) {

	ranks, err := database.GetWowLeaderboard(i.GuildID)
	if err != nil {
		discord.SendGenericErrorFromInteraction(i)
		cache.InteractionComplete(correlationId)
		return
	}

	e := embed.NewEmbed()
	e.SetTitle("Wow Leaderboard")
	e.SetDescription(getWowLeaderboardText(ranks))
	e.SetThumbnail(config.TROPHY_IMG_URL)
	e.SetColor(config.EmbedColourGold)
	err = config.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{e.MessageEmbed},
		},
	})
	if err != nil {
		logger.Error(i.GuildID, err)
	}

	cache.InteractionComplete(correlationId)

}

func getWowLeaderboardText(ranks []database.WowLeaderboard) string {
	if len(ranks) == 0 {
		return "Nobody has Woooow'd here, what's wrong with you all?"
	}

	s := ""
	for i, rank := range ranks {
		if i > 0 {
			s += "\n"
		}

		switch i {
		case 0:
			s += "🥇"
		case 1:
			s += "🥈"
		case 2:
			s += "🥉"
		default:
			s += "🏅"
		}

		s += fmt.Sprintf("<@%s> **%d** - %s", rank.UserID, rank.MaxWow, helpers.NiceDateFormat(rank.MaxWowUpdated))
	}

	return s
}
