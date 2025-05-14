package bang

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	wow "github.com/hashbat-dev/discgo-bot/Bot/Commands/Wow"
	config "github.com/hashbat-dev/discgo-bot/Config"
	discord "github.com/hashbat-dev/discgo-bot/Discord"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

type WowStats struct{}

func (s WowStats) Name() string {
	return "wowstats"
}

func (s WowStats) PermissionRequirement() int {
	return config.CommandLevelUser
}

func (s WowStats) Complexity() int {
	return config.TRIVIAL_TASK
}

func (s WowStats) Execute(message *discordgo.MessageCreate, command string) error {
	// Replied to a message?
	if message.ReferencedMessage == nil {
		return discord.ReplyToMessage(message, "You haven't replied to anything, dummy!")
	}

	// Is message in Wow cache?
	length, messages := wow.GetStatsText(message.ReferencedMessage.ID)
	if length == 0 || len(messages) == 0 {
		return discord.ReplyToMessage(message, "Couldn't find stats for this Wow! We keep detailed statistics for 60 minutes, though every Wow still counts towards your personal statistics!")
	}

	// Build embeds from chunks
	var embeds []*discordgo.MessageEmbed
	for i, chunk := range messages {
		e := embed.NewEmbed()
		if i == 0 {
			e.SetTitle(fmt.Sprintf("Level %d Wow", length))
		} else {
			e.SetTitle(fmt.Sprintf("Level %d Wow (continued)", length))
		}
		e.SetDescription(chunk)
		embeds = append(embeds, e.MessageEmbed)
	}

	// Send first embed as reply
	_, err := config.Session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
		Reference: message.ReferencedMessage.Reference(),
		Embed:     embeds[0],
	})
	if err != nil {
		return err
	}

	// Send remaining embeds as follow-ups
	for _, e := range embeds[1:] {
		_, err = config.Session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
			Reference: message.ReferencedMessage.Reference(),
			Embed:     e,
		})
		if err != nil {
			logger.Error(message.GuildID, err)
			break
		}
	}

	discord.DeleteMessage(message)
	return nil
}
