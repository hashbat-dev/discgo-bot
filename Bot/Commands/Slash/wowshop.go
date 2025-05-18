package slash

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	wow "github.com/hashbat-dev/discgo-bot/Bot/Commands/Wow"
	cache "github.com/hashbat-dev/discgo-bot/Cache"
	config "github.com/hashbat-dev/discgo-bot/Config"
	discord "github.com/hashbat-dev/discgo-bot/Discord"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

var (
	maxPerPage = 5
	pageBtnId  = "wowpage_"
	wowBuyId   = "wowbuy_"
)

func WowShop(i *discordgo.InteractionCreate, correlationId string) {
	ServeWowShop(i, correlationId, 0, true)
}

func PageButton(i *discordgo.InteractionCreate, correlationId string) {
	inboundId := i.MessageComponentData().CustomID
	splitStr := strings.Split(inboundId, "|")
	minIndexStr := strings.ReplaceAll(splitStr[0], pageBtnId, "")
	minIndex, err := strconv.Atoi(minIndexStr)
	if err != nil {
		logger.Error(i.GuildID, err)
		discord.SendGenericErrorFromInteraction(i)
		return
	}
	ServeWowShop(i, correlationId, minIndex, false)
}

func BuyItem(i *discordgo.InteractionCreate, correlationId string) {
	inboundId := i.MessageComponentData().CustomID
	splitStr := strings.Split(inboundId, "|")
	minIndexStr := strings.ReplaceAll(splitStr[0], wowBuyId, "")
	itemId, err := strconv.Atoi(minIndexStr)
	if err != nil {
		logger.Error(i.GuildID, err)
		discord.SendGenericErrorFromInteraction(i)
		return
	}
	wow.BuyItem(i, correlationId, itemId)
}

func ServeWowShop(i *discordgo.InteractionCreate, correlationId string, minIndex int, first bool) {
	if !first {
		discord.InteractionLoadingStart(i)
	}
	var buyButtons []discordgo.MessageComponent

	// Get items to show
	if minIndex < 0 {
		minIndex = 0
	} else if minIndex >= len(wow.ShopItems) {
		minIndex = len(wow.ShopItems)
	}

	endIndex := minIndex + maxPerPage
	if endIndex > len(wow.ShopItems) {
		endIndex = len(wow.ShopItems)
	}

	itemsToShow := wow.ShopItems
	if len(itemsToShow) > maxPerPage {
		itemsToShow = itemsToShow[minIndex:endIndex]
	}

	// Get Page buttons
	var navButtons []discordgo.MessageComponent
	buttonPrevIndex := minIndex - maxPerPage
	buttonNextIndex := minIndex + maxPerPage

	if buttonPrevIndex < 0 {
		buttonPrevIndex = 0
	}
	if buttonNextIndex > len(wow.ShopItems) {
		buttonNextIndex = minIndex
	}

	prevButton := discord.CreateButton(discordgo.Button{
		CustomID: fmt.Sprintf("%s%d", pageBtnId, buttonPrevIndex),
		Label:    "◀️ Prev Page",
		Style:    discordgo.SecondaryButton,
	}, correlationId, config.IO_BOUND_TASK, PageButton)
	nextButton := discord.CreateButton(discordgo.Button{
		CustomID: fmt.Sprintf("%s%d", pageBtnId, buttonNextIndex),
		Label:    "Next Page ▶️",
		Style:    discordgo.SecondaryButton,
	}, correlationId, config.IO_BOUND_TASK, PageButton)

	if buttonPrevIndex < minIndex {
		navButtons = append(navButtons, prevButton)
	}
	if buttonNextIndex > minIndex {
		navButtons = append(navButtons, nextButton)
	}

	// Loop through items, create the text and the associated buy button
	shopText := ""
	for _, item := range itemsToShow {
		emoji := "⭐"
		if item.Emoji != "" {
			emoji = item.Emoji
		}

		shopText += fmt.Sprintf("%s **%s** `%d%s`", emoji, item.Name, item.Cost, wow.CurrencyIcon)
		shopText += "\n" + item.Description
		shopText += "\n\n"

		buyButtons = append(buyButtons, discord.CreateButton(discordgo.Button{
			CustomID: fmt.Sprintf("%s%d", wowBuyId, item.ID),
			Label:    fmt.Sprintf("%s Buy %s", emoji, item.Name),
			Style:    discordgo.PrimaryButton,
		}, correlationId, config.IO_BOUND_TASK, BuyItem))
	}

	// Create an Embed for the items
	shopEmbed := embed.NewEmbed()
	shopEmbed.SetDescription(shopText)

	// Send response
	var err error
	if first {
		err = config.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:  discordgo.MessageFlagsEphemeral, // Makes it ephemeral
				Embeds: []*discordgo.MessageEmbed{shopEmbed.MessageEmbed},
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: buyButtons,
					},
					discordgo.ActionsRow{
						Components: navButtons,
					},
				},
			},
		})
	} else {
		_, err = config.Session.InteractionResponseEdit(cache.ActiveInteractions[correlationId].StartInteraction.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{shopEmbed.MessageEmbed},
			Components: &[]discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: buyButtons,
				},
				discordgo.ActionsRow{
					Components: navButtons,
				},
			},
		})
	}

	if err != nil {
		logger.Error(i.GuildID, err)
	}

	if !first {
		discord.InteractionLoadingFinish(i)
	}
}
