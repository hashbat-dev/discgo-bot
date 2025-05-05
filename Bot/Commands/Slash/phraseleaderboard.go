package slash

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	cache "github.com/hashbat-dev/discgo-bot/Cache"
	config "github.com/hashbat-dev/discgo-bot/Config"
	database "github.com/hashbat-dev/discgo-bot/Database"
	discord "github.com/hashbat-dev/discgo-bot/Discord"
	helpers "github.com/hashbat-dev/discgo-bot/Helpers"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

var (
	trophyImgUrl = "https://i.imgur.com/XK0Fqvr.png"
)

func PhraseLeaderboard(i *discordgo.InteractionCreate, correlationId string) {
	cachedInteraction := cache.ActiveInteractions[correlationId]
	phrase := cachedInteraction.Values.String["phrase"]
	phrase = strings.ToLower(strings.TrimSpace(phrase))

	// 1. Validate
	if phrase == "" {
		discord.SendEmbedFromInteraction(i, "Error", "No Phrase entered!")
		cache.InteractionComplete(correlationId)
		return
	}
	if len(phrase) > 50 {
		discord.SendEmbedFromInteraction(i, "Error", "Phrase too long! Phrases can be a maximum of 50 characters")
		cache.InteractionComplete(correlationId)
		return
	}

	triggerPhrase, err := database.GetTriggerPhrase(i.GuildID, phrase)
	if err != nil {
		discord.SendGenericErrorFromInteraction(i)
		cache.InteractionComplete(correlationId)
		return
	}

	linkExists, err := database.DoesPhraseLinkExist(i.GuildID, triggerPhrase.ID)
	if err != nil {
		discord.SendGenericErrorFromInteraction(i)
		cache.InteractionComplete(correlationId)
		return
	}

	if !linkExists {
		discord.SendEmbedFromInteraction(i, "Error", fmt.Sprintf("The phrase '%s' isn't being tracked! You can ask a bot admin to create it with /add-trigger", phrase))
		cache.InteractionComplete(correlationId)
		return
	}

	// 2. Get stats
	ranks, err := database.GetPhraseLeaderboard(i.GuildID, phrase)
	if err != nil {
		discord.SendGenericErrorFromInteraction(i)
		cache.InteractionComplete(correlationId)
		return
	}

	// 3. Display Board
	e := embed.NewEmbed()
	e.SetTitle(helpers.CapitaliseWords(phrase) + " Leaderboard")
	e.SetDescription(getLeaderboardText(ranks))
	e.SetThumbnail(trophyImgUrl)
	e.SetColor(discord.EmbedColourGold)
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

func getLeaderboardText(ranks []database.PhraseLeaderboardUser) string {
	if len(ranks) == 0 {
		return "Nobody has said this yet! Maybe you'll be the first?"
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

		s += fmt.Sprintf(" **%d** <@%s>", rank.Count, rank.UserID)
	}

	return s
}
