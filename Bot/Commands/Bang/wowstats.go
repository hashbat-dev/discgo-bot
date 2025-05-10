package bang

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	wow "github.com/hashbat-dev/discgo-bot/Bot/Commands/Wow"
	config "github.com/hashbat-dev/discgo-bot/Config"
	discord "github.com/hashbat-dev/discgo-bot/Discord"
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
	length, msg := wow.GetStatsText(message.ReferencedMessage.ID)
	if length == 0 || msg == "" {
		return discord.ReplyToMessage(message, "Couldn't find stats for this Wow! We keep detailed statistics for 60 minutes, though every Wow sitll counts towards your personal statistics!")
	}

	// Return stats
	title := fmt.Sprintf("Level %d Wow", length)
	e := embed.NewEmbed()
	e.SetTitle(title)
	e.SetDescription(msg)
	discord.ReplyToMessageWithEmbed(message.ReferencedMessage, *e)
	discord.DeleteMessage(message)
	return nil
}
