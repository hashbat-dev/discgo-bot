package wow

import (
	"fmt"
	"slices"

	"github.com/bwmarrin/discordgo"
	cache "github.com/hashbat-dev/discgo-bot/Cache"
	config "github.com/hashbat-dev/discgo-bot/Config"
	database "github.com/hashbat-dev/discgo-bot/Database"
	discord "github.com/hashbat-dev/discgo-bot/Discord"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

var (
	CurrencyIcon = "₩"
)

func BuyItem(i *discordgo.InteractionCreate, correlationId string, ID int) {
	discord.InteractionLoadingStart(i)

	// Get shop item
	var item WowShopItem
	found := false
	for _, x := range ShopItems {
		if x.ID == ID {
			item = x
			found = true
			break
		}
	}
	if !found {
		logger.ErrorText(i.GuildID, "Couldn't find shop item")
		return
	}

	userId := cache.GetUserIDFromCorrelationID(correlationId)

	// Have they hit the max item count?
	if item.MaxAtOnce > 0 && wowInventoryItemCount(i.GuildID, userId, item.ID) >= item.MaxAtOnce {
		discord.UpdateInteractionResponse(i, "Max Item Count", fmt.Sprintf("You cannot have more than %d %s items at once", item.MaxAtOnce, item.Name), config.EmbedColourRed)
		return
	}

	// Can they afford it?
	userBalance, err := database.GetUserWowBalance(i.GuildID, userId)
	if err != nil {
		discord.SendGenericErrorFromInteraction(i)
		return
	}

	if userBalance < item.Cost {
		text := fmt.Sprintf("Your balance is %d%s", userBalance, CurrencyIcon)
		if userBalance < 10000 {
			text += ", you're broke."
		}
		discord.UpdateInteractionResponse(i, "Transaction Declined", text, config.EmbedColourRed)
		return
	}

	// Deduct the Cost
	err = database.UpdateWowBalance(i.GuildID, userId, item.Cost, false)
	if err != nil {
		discord.SendGenericErrorFromInteraction(i)
		return
	}

	// Count the Purchase
	database.CountWowPurchase(i.GuildID, userId, item.ID)

	// Add item
	err = addToWowInventory(i.GuildID, userId, item)
	if err != nil {
		discord.SendGenericErrorFromInteraction(i)
		return
	}

	text := fmt.Sprintf("You bought %s for `%d%s`\nYour new balance is `%d%s`", item.Name, item.Cost, CurrencyIcon, userBalance-item.Cost, CurrencyIcon)
	discord.UpdateInteractionResponse(i, "Transaction Successful", text, config.EmbedColourGreen)
}

func GetWowInventoryText(guildId string, userId string) string {
	var userInv []InventoryItem
	found := false
	dataInventoryLock.RLock()
	if inv, exists := dataUserInventories[fmt.Sprintf("%s|%s", guildId, userId)]; exists {
		userInv = inv
		found = true
	}
	dataInventoryLock.RUnlock()

	if !found {
		return "This user doesn't have any Wow items at the moment!"
	}

	itemCounts := make(map[string]int)
	for _, inv := range userInv {
		itemCounts[inv.ShopItem.Name]++
	}
	var nameCheck []string

	s := fmt.Sprintf("<@%s> currently has the following Wow Items in their inventory...", userId)
	for _, inv := range userInv {
		if slices.Contains(nameCheck, inv.ShopItem.Name) {
			continue
		}
		emoji := inv.ShopItem.Emoji
		if emoji == "" {
			emoji = DefaultEmoji
		}
		countText := ""
		if itemCounts[inv.ShopItem.Name] > 1 {
			countText = fmt.Sprintf(" x%d", itemCounts[inv.ShopItem.Name])
		}
		s += fmt.Sprintf("\n\n%s **%s%s**: %s", emoji, inv.ShopItem.Name, countText, inv.ShopItem.Description)

		nameCheck = append(nameCheck, inv.ShopItem.Name)
	}

	return s

}
