package slash

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	wow "github.com/hashbat-dev/discgo-bot/Bot/Commands/Wow"
	cache "github.com/hashbat-dev/discgo-bot/Cache"
	config "github.com/hashbat-dev/discgo-bot/Config"
	database "github.com/hashbat-dev/discgo-bot/Database"
	discord "github.com/hashbat-dev/discgo-bot/Discord"
	helpers "github.com/hashbat-dev/discgo-bot/Helpers"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

func WowStats(i *discordgo.InteractionCreate, correlationId string) {
	cachedInteraction := cache.ActiveInteractions[correlationId]
	user := cachedInteraction.Values.User["user"]

	if user == nil {
		discord.SendEmbedFromInteraction(i, "Error", "No User provided!")
		cache.InteractionComplete(correlationId)
		return
	}

	if user.ID == config.Session.State.User.ID {
		discord.SendEmbedFromInteraction(i, "Error", "I don't Wow! Well, I did then, but I don't count.")
		cache.InteractionComplete(correlationId)
		return
	}

	stats, err := database.GetUserWowStats(i.GuildID, i.Member.User.ID)
	if err != nil {
		discord.SendGenericErrorFromInteraction(i)
		cache.InteractionComplete(correlationId)
		return
	}

	title := fmt.Sprintf("%s's Wow Stats", helpers.CapitaliseWords(i.Member.User.Username))
	baseDescription := fmt.Sprintf("Max Wow: **%d**", stats.MaxWow)
	baseDescription += fmt.Sprintf("\nObtained: **%s**", helpers.NiceDateFormat(stats.MaxWowUpdated))

	embeds := []*discordgo.MessageEmbed{}
	firstEmbed := embed.NewEmbed()
	firstEmbed.SetTitle(title)
	firstEmbed.SetDescription(baseDescription)
	embeds = append(embeds, firstEmbed.MessageEmbed)

	// Build effect history text
	if len(stats.Effects) > 0 {
		var effectChunks []string
		lastType := ""
		currentChunk := "**Effect History**"

		for _, effect := range stats.Effects {
			if lastType != effect.Type {
				currentChunk += fmt.Sprintf("\n\n**%s**", helpers.CapitaliseWords(effect.Type))
				lastType = effect.Type
			}
			time := "time"
			if effect.Count != 1 {
				time += "s"
			}
			emoji := effect.Emoji
			if emoji == "" {
				emoji = wow.DefaultEmoji
			}
			line := fmt.Sprintf("\n%s%s %s - **%d** %s, last used **%s**", wow.IndentPadding, emoji, effect.Name, effect.Count, time, helpers.NiceDateFormat(effect.LastTriggered))

			// If adding this line would overflow currentChunk, start a new chunk
			if len(currentChunk)+len(line) > config.MAX_EMBED_DESC_LENGTH {
				effectChunks = append(effectChunks, currentChunk)
				currentChunk = ""
			}
			currentChunk += line
		}
		if currentChunk != "" {
			effectChunks = append(effectChunks, currentChunk)
		}

		// Create an embed for each effect chunk
		for _, chunk := range effectChunks {
			e := embed.NewEmbed()
			e.SetDescription(chunk)
			embeds = append(embeds, e.MessageEmbed)
		}
	}

	// Limit to Discord's 10-embed maximum
	if len(embeds) > 10 {
		embeds = embeds[:10]
	}

	// Respond with the first embed
	err = config.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embeds[0]},
		},
	})

	if err != nil {
		logger.Error(i.GuildID, err)
		cache.InteractionComplete(correlationId)
		return
	}

	// Follow up with additional embeds
	for _, embed := range embeds[1:] {
		_, err := config.Session.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{embed},
		})
		if err != nil {
			logger.Error(i.GuildID, err)
			break
		}
	}

	cache.InteractionComplete(correlationId)
}
