package sentience

import (
	"math/rand"
	"strings"
	"time"

	"github.com/ZestHusky/femboy-control/Bot/audit"
	"github.com/ZestHusky/femboy-control/Bot/config"
	dbhelper "github.com/ZestHusky/femboy-control/Bot/dbhelpers"
	"github.com/ZestHusky/femboy-control/Bot/helpers"
	"github.com/ZestHusky/femboy-control/Bot/logging"
	"github.com/bwmarrin/discordgo"
)

func ResponseToMessage(message *discordgo.MessageCreate, isPraise bool) {

	// Get Gif
	gifCat := "badbot"
	if isPraise {
		gifCat = "goodbot"
	}

	gif, err := dbhelper.GetRandGifURL(gifCat)
	if err != nil {
		logging.SendError(message)
		return
	}

	// Send Message
	_, err = config.Session.ChannelMessageSend(message.ChannelID, gif)
	if err != nil {
		audit.Error(err)
		logging.SendError(message)
		return
	}

	// Count this?
	count := true
	if gifCat == "goodbot" {
		exists := false
		for _, praise := range praiseCache {
			if praise.UserId == message.Author.ID {
				exists = true
				if time.Since(praise.LastPraise).Minutes() < float64(praiseCooldownMins) {
					count = false
				}
				break
			}
		}

		if !exists {
			praiseCache = append(praiseCache, PraiseCache{
				UserId:     message.Author.ID,
				LastPraise: time.Now(),
			})
		}
	}

	if count {
		dbhelper.CountCommand(gifCat, message.Author.ID)
	}

}

func DetectBotTalk(message *discordgo.MessageCreate, lowerMsg string) {

	botSaid := false
	if strings.Contains(" "+lowerMsg+" ", " bot ") {
		botSaid = true
	}

	if !botSaid {
		botSaid = helpers.DoesTextMentionWord(lowerMsg, "bot")
	}

	if !botSaid {
		return
	}

	for _, m := range abuse {
		if helpers.DoesTextMentionWord(lowerMsg, m) {
			ResponseToMessage(message, false)
			return
		}
	}

	for _, m := range praise {
		if helpers.DoesTextMentionWord(lowerMsg, m) {
			ResponseToMessage(message, true)
			return
		}
	}

}

func RandomReplies(message *discordgo.MessageCreate) bool {

	chance := 600
	random := rand.Intn(chance + 1)
	chanceYes := 200
	chanceNo := 400

	replyText := ""
	if random == chanceYes {
		replyText = helpers.GetRandomText(yesQuotes)
		audit.Log("[!!!] Random YES sent to User: " + helpers.GetNicknameFromID(message.GuildID, message.Author.ID))
		dbhelper.CountCommand("random-yes", message.Message.Author.ID)
	} else if random == chanceNo {
		replyText = helpers.GetRandomText(noQuotes)
		audit.Log("[!!!] Random NO sent to User: " + helpers.GetNicknameFromID(message.GuildID, message.Author.ID))
		dbhelper.CountCommand("random-no", message.Message.Author.ID)
	}

	if replyText != "" {
		_, err := config.Session.ChannelMessageSendReply(message.ChannelID, replyText, message.Reference())
		if err != nil {
			audit.Error(err)
		}
		return true
	}

	return false
}
