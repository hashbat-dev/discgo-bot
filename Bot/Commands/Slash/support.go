package slash

import (
	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	discord "github.com/dabi-ngin/discgo-bot/Discord"
)

func SupportInfo(i *discordgo.InteractionCreate, correlationId string) {
	embed := embed.NewEmbed()
	embed.SetTitle("Help & Support")
	embed.SetDescription("For help and support on Discgo Bot please contact us.")
	discord.ReplyToInteractionWithEmbed(i, embed, true)
}
