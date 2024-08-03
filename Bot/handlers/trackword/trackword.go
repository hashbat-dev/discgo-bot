package trackword

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dabi-ngin/discgo-bot/Bot/audit"
	dbhelpers "github.com/dabi-ngin/discgo-bot/Bot/dbhelpers"
	logger "github.com/dabi-ngin/discgo-bot/Bot/logging"
)

type Mention struct {
	user    string
	keyword string
	time    time.Time
}

var lastMentions []Mention
var secondBuffer float64 = 30

var comboDetection []Mention
var comboMaxCache time.Duration = 65 * time.Minute

func TrackWord(discord *discordgo.Session, message *discordgo.MessageCreate, keyword string) {

	// 1. Get the Current Mention as an object.
	userId := message.Author.ID
	currMention := Mention{
		user:    userId,
		keyword: keyword,
		time:    time.Now(),
	}

	comboDetection = append(comboDetection, currMention)

	// 2C. Clean up Combos
	var newComboDetection []Mention
	for _, combo := range comboDetection {
		if time.Since(combo.time) <= comboMaxCache {
			newComboDetection = append(newComboDetection, combo)
		}
	}
	comboDetection = newComboDetection

	// 3. Check for someone being a Spamming Samantha
	// 3A. Are they in the Cache?
	tooQuick := false
	cacheIndex := -1

	for i, cache := range lastMentions {
		if cache.user == userId {
			timeSince := time.Since(cache.time).Seconds()
			if timeSince <= secondBuffer {
				tooQuick = true
			}
			cacheIndex = i
			break
		}
	}

	// 3B. Return?
	if tooQuick {
		audit.Log(fmt.Sprintf("Keyword: %v, Enacting buffer: %v", keyword, secondBuffer))
		return
	}

	// 4. Update our records
	// 4A. Cache
	if cacheIndex >= 0 {
		lastMentions[cacheIndex].time = time.Now()
	} else {
		lastMentions = append(lastMentions, currMention)
	}

	// 4B. Database
	err := dbhelpers.CountWord(userId, keyword)
	if err != nil {
		audit.Error(err)
		return
	}

	// 5. Get our Message Text to send
	// 5A. Basic bitch response
	mentionText := strings.ToUpper(keyword) + " MENTIONED"

	// 5B. Specific overrides
	if keyword == "mischief" {
		mentionText = "https://tenor.com/view/helldivers-helldivers-2-mischief-mischief-mode-mischief-mode-activated-gif-15827932577000824336"
	} else if keyword == "goon" {
		getGif, err := dbhelpers.GetRandGifURL("goon")
		if err != nil {
			audit.Error(err)
		} else {
			mentionText = getGif
		}
	} else if keyword == "cigarette" {
		mentionText = "https://tenor.com/view/500-cigarettes-the-orville-gif-7658479820042125764"
	}

	// 6. Send the Message
	messageObj := logger.MessageRefObj(message.Message)
	_, err = discord.ChannelMessageSendReply(message.ChannelID, mentionText, &messageObj)
	if err == nil {
		audit.Log("Logged Mention for: " + keyword)
	} else {
		audit.Error(err)
		logger.SendError(message)
		return
	}
}
