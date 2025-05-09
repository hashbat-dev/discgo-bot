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

	description := fmt.Sprintf("Max Wow: **%d**", stats.MaxWow)
	description += fmt.Sprintf("\nObtained: **%s**", helpers.NiceDateFormat(stats.MaxWowUpdated))

	if len(stats.Effects) > 0 {
		description += "\n\n**Effect History**"
		lastType := ""
		for _, effect := range stats.Effects {
			if lastType != effect.Type {
				description += fmt.Sprintf("\n**%s**", helpers.CapitaliseWords(effect.Type))
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
			description += fmt.Sprintf("\n%s%s %s - **%d** %s, last used **%s**", wow.IndentPadding, emoji, effect.Name, effect.Count, time, helpers.NiceDateFormat(effect.LastTriggered))
		}
	}

	errEmbed := embed.NewEmbed()
	errEmbed.SetTitle(title)
	errEmbed.SetDescription(description)
	err = config.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{errEmbed.MessageEmbed},
		},
	})

	if err != nil {
		logger.Error(i.GuildID, err)
	}
	cache.InteractionComplete(correlationId)
}
