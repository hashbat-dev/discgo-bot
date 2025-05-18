package slash

import (
	"github.com/bwmarrin/discordgo"
	wow "github.com/hashbat-dev/discgo-bot/Bot/Commands/Wow"
	cache "github.com/hashbat-dev/discgo-bot/Cache"
	discord "github.com/hashbat-dev/discgo-bot/Discord"
)

func WowInventory(i *discordgo.InteractionCreate, correlationId string) {
	cachedInteraction := cache.ActiveInteractions[correlationId]
	user := cachedInteraction.Values.User["user"]

	s := wow.GetWowInventoryText(i.GuildID, user.ID)
	discord.SendEmbedFromInteraction(i, "Wow Inventory", s, 0)

	cache.InteractionComplete(correlationId)
}
