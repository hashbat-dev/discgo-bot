package translate

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/ZestHusky/femboy-control/Bot/audit"
	"github.com/ZestHusky/femboy-control/Bot/config"
	"github.com/ZestHusky/femboy-control/Bot/friday"
	"github.com/ZestHusky/femboy-control/Bot/logging"
	"github.com/bwmarrin/discordgo"
)

func TranslateBottomSpeech(message *discordgo.MessageCreate) {

	if message.Message.ReferencedMessage == nil {
		logging.SendErrorMsgReply(message, "You've not replied to a message doofus!")
		return
	}

	if friday.IsItFwiday() || message.Message.ReferencedMessage.Author.ID != config.BullyTarget {
		logging.SendErrorMsgReply(message, "**Translation from Bottom to English:** Error. This isn't bottom text.")
		return
	}

	// 1. Have we already translated this today?
	var newCache []Translation
	for _, trans := range translationCache {
		if time.Since(trans.Added).Hours() <= 24 {
			newCache = append(newCache, trans)
		}
	}

	for _, trans := range newCache {
		if trans.Original == strings.TrimSpace(message.Message.ReferencedMessage.Content) {
			SendResponse(message, trans.Original, trans.New)
			return
		}
	}

	// 2. Pick a random number of sentences
	returnString := ""

	var includedIndexes []int = []int{}
	messageLength := len(message.Message.ReferencedMessage.Content)
	chopSize := 10
	loopLimit := int(messageLength / chopSize)
	audit.Log(fmt.Sprintf("LoopLimit: %v, from Message Length: %v, Chop Size: %v", loopLimit, messageLength, chopSize))
	if loopLimit > len(bottomSentences)-1 {
		loopLimit = loopLimit - 2
	} else if loopLimit == 0 {
		loopLimit = 1
	}

	for i := 0; i < loopLimit; i++ {
		if i > 0 {
			returnString += " "
		}

		continueLoop := true
		newIndex := 0
		loops := 1
		maxLoops := 25
		for continueLoop {
			randomNo := rand.Intn(len(bottomSentences))
			alreadyHave := false
			for _, i := range includedIndexes {
				if i == randomNo {
					alreadyHave = true
					break
				}
			}

			continueLoop = alreadyHave
			if !continueLoop {
				newIndex = randomNo
			}
			loops++
			if loops == maxLoops {
				break
			}
		}

		returnString += bottomSentences[newIndex]
	}

	// 3. Sprinkle it with an uwu or owo?
	uwuChance := rand.Intn(5)
	if uwuChance == 1 {
		returnString += " uwu"
	} else if uwuChance == 2 {
		returnString += " owo"
	}

	// 4. Send the translation
	newCache = append(newCache, Translation{
		Original: strings.TrimSpace(message.Message.ReferencedMessage.Content),
		New:      returnString,
		Added:    time.Now(),
	})

	translationCache = newCache
	SendResponse(message, strings.TrimSpace(message.Message.ReferencedMessage.Content), returnString)
}

func SendResponse(message *discordgo.MessageCreate, original string, new string) {
	sendText := "**Original:** " + original + "\n\n" + "**Translated:** " + new

	_, err := config.Session.ChannelMessageSendReply(message.ChannelID, sendText, message.ReferencedMessage.Reference())
	if err != nil {
		logging.SendErrorMsgReply(message, "**Translation from Bottom to English:** Error. Couldn't correctly process the associated bottom energy.")
		audit.Error(err)
	}
}

var bottomSentences []string = []string{
	"I just really love kissing boys.",
	"I wish someone would just shove my face into a pillow.",
	"Man, I really love Fallout: New Vegas.",
	"My legs feel so soft!",
	"I just really love boi smell.",
	"I've worn these stockings for 3 days straight now.",
	"If I don't get what I want I will sulk, throw a tantrum and demand kisses.",
	"By the way, are you done with those socks? Can I have them?",
	"I'm feeling submissive and breedable right now.",
	"Who's down to clown in my downtown?",
	"Can't I just drink Monster and fuck bois?",
	"Oh man League of Legends really makes me want a dick in my mouth.",
	"I wanna be spanked.",
	"Tristana really makes my panties wet.",
	"Pillow Princess moment.",
	"Don't blame me, blame the girl shots.",
	"I'm such a siwwy biwwy.",
	"I want to be daddy's uwurinal.",
	"I need to go make stinkies.",
	"I wish someone would use me like a beaten hoe.",
}

type Translation struct {
	Original string
	New      string
	Added    time.Time
}

var translationCache []Translation
