package bang

import (
	"fmt"
	"strings"

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
	length, msg := wow.GetStatsText(message.ReferencedMessage.ID)
	if length == 0 || msg == "" {
		return discord.ReplyToMessage(message, "Couldn't find stats for this Wow! We keep detailed statistics for 60 minutes, though every Wow still counts towards your personal statistics!")
	}

	// Split msg by lines and group into chunks
	lines := strings.Split(msg, "\n")
	var chunks []string
	currentChunk := ""

	for _, line := range lines {
		// +1 for the newline being added back in
		if len(currentChunk)+len(line)+1 > config.MAX_EMBED_DESC_LENGTH {
			chunks = append(chunks, currentChunk)
			currentChunk = ""
		}
		if currentChunk != "" {
			currentChunk += "\n"
		}
		currentChunk += line
	}
	if currentChunk != "" {
		chunks = append(chunks, currentChunk)
	}

	// Build embeds from chunks
	title := fmt.Sprintf("Level %d Wow", length)
	var embeds []*discordgo.MessageEmbed

	for i, chunk := range chunks {
		e := embed.NewEmbed()
		if i == 0 {
			e.SetTitle(title)
		} else {
			e.SetTitle(fmt.Sprintf("Level %d Wow (continued)", length))
		}
		e.SetDescription(chunk)
		embeds = append(embeds, e.MessageEmbed)
	}

	// Send first embed as reply
	_, err := config.Session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
		Reference: message.Reference(),
		Embed:     embeds[0],
	})
	if err != nil {
		return err
	}

	// Send remaining embeds as follow-ups
	for _, e := range embeds[1:] {
		_, err = config.Session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
			Reference: message.Reference(),
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
