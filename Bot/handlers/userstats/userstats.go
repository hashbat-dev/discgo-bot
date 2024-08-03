package userstats

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	"github.com/dabi-ngin/discgo-bot/Bot/audit"
	"github.com/dabi-ngin/discgo-bot/Bot/config"
	"github.com/dabi-ngin/discgo-bot/Bot/constants"
	dbhelper "github.com/dabi-ngin/discgo-bot/Bot/dbhelpers"
	"github.com/dabi-ngin/discgo-bot/Bot/helpers"
	"github.com/dabi-ngin/discgo-bot/Bot/logging"
)

func GetStats(interaction *discordgo.InteractionCreate) {

	// Put a Loading Interaction in
	embedStart := embed.NewEmbed()
	embedStart.SetDescription("Obtaining user information...")

	err := config.Session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embedStart.MessageEmbed},
		},
	})

	if err != nil {
		audit.Error(err)
		logging.SendErrorInteraction(interaction)
		return
	}

	errorEmbed := embed.NewEmbed()
	errorEmbed.SetTitle("Oh noes!")
	errorEmbed.SetDescription("An error occurred getting user stats. Pwease don't be mad >w<")
	errorEmbed.SetImage(constants.GIF_AE_CRY)
	var webHook discordgo.WebhookEdit
	webHook.Embeds = &[]*discordgo.MessageEmbed{errorEmbed.MessageEmbed}

	// Get the user we're looking for
	optionMap := helpers.GetOptionMap(interaction)
	inUser := helpers.GetOptionUserIDValue(optionMap, "user")
	if inUser == "" {
		inUser = interaction.Member.User.ID
	}
	inAvatar := helpers.GetAvatarURLFromID(interaction.GuildID, inUser)
	userName := helpers.GetNicknameFromID(interaction.GuildID, inUser)
	embedText := ""

	// Get Bot rating
	botRating, err := GetBotRating(inUser)
	if err != nil {
		config.Session.InteractionResponseEdit(interaction.Interaction, &webHook)
		audit.Error(err)
		return
	}

	if botRating > 0 {
		embedText += constants.EMOTE_THUMB_UP + " They're a Level " + fmt.Sprint(botRating) + " Bot Lover! ❤️"
	} else {
		embedText += constants.EMOTE_THUMB_DOWN + " They're a Level " + fmt.Sprint(botRating*-1) + " Bot Hater 😢"
	}

	// Have they been given any Cigarettes?
	cigCount, _, _, err := dbhelper.GetRankFromUser(inUser, "cigarette")
	if err != nil {
		audit.Error(err)
	} else if cigCount > 0 {
		embedText += "\n"
		embedText += "🚬 They've been given " + fmt.Sprint(cigCount*500) + " cigarettes."
	}

	// Have they been Random Toothjacked?
	randomToothCount, err := dbhelper.CommandLogGetCountForUser("randomtooth", inUser)
	if err != nil {
		audit.Error(err)
	} else if randomToothCount > 0 {
		embedText += "\n"
		embedText += ":toothbrush: They've been Random Toothjacked " + fmt.Sprint(randomToothCount) + " time"
		if randomToothCount != 1 {
			embedText += "s"
		}
		embedText += "!"
	}

	// Have they been Random Abused?
	randomRacism, err := dbhelper.CommandLogGetCountForUser("randomslur", inUser)
	if err != nil {
		audit.Error(err)
	} else if randomRacism > 0 {
		embedText += "\n"
		embedText += ":person_bald::skin-tone-5: They've been Randomly Racially Abused " + fmt.Sprint(randomRacism) + " time"
		if randomRacism != 1 {
			embedText += "s"
		}
		embedText += "!"
	}

	// Have they been Woooow'd?
	maxWow, totalWow, err := dbhelper.CountWow_GetCounts(inUser)
	if err != nil {
		audit.Error(err)
	} else if totalWow > 0 {
		embedText += "\n"
		embedText += "🤓 They've been Wooow'd at " + fmt.Sprint(totalWow) + " time"
		if totalWow != 1 {
			embedText += "s"
		}
		embedText += ", their longest was " + fmt.Sprint(maxWow) + "!"
	}

	// Get Ranks for all tracked words
	embedText += "\n\n"

	rankCats, err := dbhelper.GetDistinctRanks()
	if err != nil {
		config.Session.InteractionResponseEdit(interaction.Interaction, &webHook)
		audit.Error(err)
		return
	}

	if len(rankCats) == 0 {
		embedText += userName + " hasn't interacted much else " + helpers.GetEmote("ae_cry", true)
	} else {
		embedText += userName + " is a..."
		for _, rankCat := range rankCats {
			rankCount, rankText, _, err := dbhelper.GetRankFromUser(inUser, rankCat)
			if err != nil {
				config.Session.InteractionResponseEdit(interaction.Interaction, &webHook)
				audit.ErrorWithText("Category: "+rankCat, err)
				return
			} else if rankCount == 0 {
				continue
			} else {
				embedText += "\n" + helpers.GetEmote("medal", true) + rankText + " with a " + rankCat + " count of " + fmt.Sprint(rankCount)
			}
		}
	}

	e := embed.NewEmbed()
	e.SetTitle(userName + "'s Server Stats")
	e.SetDescription(embedText)
	e.SetImage(inAvatar)

	var successHook discordgo.WebhookEdit
	successHook.Embeds = &[]*discordgo.MessageEmbed{e.MessageEmbed}
	config.Session.InteractionResponseEdit(interaction.Interaction, &successHook)
}

func WowBoard(interaction *discordgo.InteractionCreate) {

	wowRanks, err := dbhelper.CountWow_Ranking()
	if err != nil {
		audit.Error(err)
		logging.SendErrorInteraction(interaction)
		return
	}

	text := "Who's had the longest WOOOOW's from the Bot?\n"

	leaderId := ""
	for i, rank := range wowRanks {

		inString := "\n"

		if i == 0 {
			leaderId = rank.UserID
			inString += "🥇 "
		} else if i == 1 {
			inString += "🥈 "
		} else if i == 2 {
			inString += "🥉 "
		} else {
			inString += "🏅 "
		}

		avg := float64(rank.TotalCount) / float64(rank.MaxWow)
		opw := fmt.Sprintf("%.2f", avg)
		inString += fmt.Sprint(rank.MaxWow) + " (" + opw + "opw) - " + helpers.GetNicknameFromID(interaction.GuildID, rank.UserID)
		text += inString
	}

	title := "Server WOOOW Leaderboard"
	logging.SendMessageInteraction(interaction, title, text, helpers.GetAvatarURLFromID(interaction.GuildID, leaderId), "OPW: Average number of O's per WOW", false)
}
