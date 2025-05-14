package wow

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	config "github.com/hashbat-dev/discgo-bot/Config"
	helpers "github.com/hashbat-dev/discgo-bot/Helpers"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

func respond(wow *Generation) {
	messageTexts := helpers.SplitText(wow.Output, config.MAX_MESSAGE_LENGTH)
	first := true
	var firstWowMsg *discordgo.Message
	for _, msgText := range messageTexts {
		var msgRef *discordgo.MessageReference
		if first {
			msgRef = wow.Message.Reference()
		}

		msg := &discordgo.MessageSend{
			Content:   msgText,
			Reference: msgRef,
		}

		wowMsg, err := config.Session.ChannelMessageSendComplex(wow.Message.ChannelID, msg)
		if err != nil {
			logger.Error(wow.Message.GuildID, err)
		}

		if first {
			firstWowMsg = wowMsg
			first = false
		}
		wow.WowMessageIDs = append(wow.WowMessageIDs, wowMsg.ID)
	}

	addToCache(wow)

	// New Record?
	if wow.OCount > dataHighestWowInGuild[wow.Message.GuildID] {
		emb := embed.NewEmbed()
		emb.SetTitle("New Wow Record")
		emb.SetDescription(fmt.Sprintf("Ladies and Gentlemen, <@%s> has just broken the all time Wow record!", wow.Message.Author.ID))
		emb.SetThumbnail(config.TROPHY_IMG_URL)
		emb.SetFooter(fmt.Sprintf("%d level Wow", wow.OCount))
		_, err := config.Session.ChannelMessageSendComplex(wow.Message.ChannelID, &discordgo.MessageSend{
			Reference: firstWowMsg.Reference(),
			Embed:     emb.MessageEmbed,
		})

		if err != nil {
			logger.Error(wow.Message.GuildID, err)
		}

		dataHighestWowInGuild[wow.Message.GuildID] = wow.OCount
	}
}
