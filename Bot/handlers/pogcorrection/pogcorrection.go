package pogcorrection

import (
	"strings"
	"time"

	"github.com/ZestHusky/femboy-control/Bot/audit"
	"github.com/ZestHusky/femboy-control/Bot/config"
	"github.com/ZestHusky/femboy-control/Bot/constants"
	"github.com/ZestHusky/femboy-control/Bot/friday"
	"github.com/ZestHusky/femboy-control/Bot/helpers"
	logger "github.com/ZestHusky/femboy-control/Bot/logging"
	"github.com/bwmarrin/discordgo"
)

var lastMsg1, lastMsg2, lastMsg3, lastMsg4, lastMsg5, lastMsg6 string = "", "", "", "", "", ""

var checkMinutes time.Duration = 5
var baitBottom time.Time = time.Now().Add(-checkMinutes * time.Minute)
var baitTop time.Time = time.Now().Add(-checkMinutes * time.Minute)

var pogCorrectText = []string{
	"No Pog, you are a bottom",
	"Stop denying it Pog, you're a bottom",
	"Incorrect, you're a bottom",
	"Nice try, but you're a bottom though",
	"Correction: Bottom",
	"Nice argument, unfortunately you're a bottom",
	"That's something a bottom would say",
}

var adjNegWords = []string{
	"not", "arent", "aren't", "is not", "isn't", "opposite", "isnt", "unlike",
}
var adjPosWords = []string{
	"am", "im", "i'm", "opposite", "exist", "i is", "me", "pog is", "i am", "myself", "i do", "of them", "undeniable", "to them", "its", "it's", "like",
}

func Detect(message string, fromUser string, repliedTo *discordgo.Message) string {

	// CHECK 1
	// Is it Fwiday?
	// ===========================================================
	if friday.IsItFwiday() {
		return ""
	}

	// CHECK 2
	// Was this even Pog?
	// ===========================================================
	var s string = ""
	msgContainsBottom := false
	msgContainsTop := false
	message = strings.ToLower(message)

	if fromUser != constants.USER_ID_POG {
		msgContainsBottom = helpers.DoesTextMentionWord(message, "bottom")
		if !msgContainsBottom {
			message = strings.Replace(message, "b word", "b-word", -1)
			msgContainsBottom = helpers.DoesTextMentionWord(message, "b-word")
		}

		if !msgContainsBottom {
			msgContainsTop = helpers.DoesTextMentionWord(message, "top")
			if msgContainsTop {
				audit.Log("Someone mentioned 'top'")
				baitTop = time.Now()
			}
		} else {
			audit.Log("Someone mentioned 'bottom'")
			baitBottom = time.Now()
		}

		return ""
	} else {
		s = strings.ToLower(message)
	}

	// CHECK 3
	// Shall we flat out just send a Toothjack?
	// ===========================================================
	adjSingleWordToothjacks := []string{
		"what", "how", "why", "where", "when", "whqt",
	}

	for _, tj := range adjSingleWordToothjacks {
		if s == tj {
			return "tooth"
		}
	}

	// CHECK 4
	// Do we have an "IM NOT A BOTTOM" type message?
	// ===========================================================
	s = GetPogString(s)

	// Do we have a Noun?
	// => From Pog?
	accuseBottom := msgContainsBottom
	accuseTop := msgContainsTop
	isFromPog := true

	// => From a Reply?
	if !accuseBottom && !accuseTop && repliedTo != nil {
		accuseBottom = helpers.DoesTextMentionWord(repliedTo.Content, "bottom")
		if accuseBottom {
			isFromPog = false
		}
		if !accuseBottom {
			accuseTop = helpers.DoesTextMentionWord(repliedTo.Content, "top")
			if accuseTop {
				isFromPog = false
			}
		}
	}

	// => From a previous Message?
	if !accuseBottom && !accuseTop {
		accuseBottom = WithinTime(baitBottom, time.Now())
		accuseTop = WithinTime(baitTop, time.Now())
	}

	// Do we have an Abjective?
	if accuseBottom || accuseTop {
		if accuseBottom {
			for _, adj := range adjNegWords {
				if helpers.DoesTextMentionWord(s, adj) {
					logStr := "Detected: [Bottom][" + adj + "]"
					if isFromPog {
						logStr += "[Pog instigated]"
					} else {
						logStr += "[Non-Pog instigated]"
					}
					audit.Log(logStr)
					ResetCache()
					return "correct"
				}
			}
		}

		if accuseTop {
			for _, adj := range adjPosWords {
				if helpers.DoesTextMentionWord(s, adj) {
					logStr := "Detected: [Top][" + adj + "]"
					if isFromPog {
						logStr += "[Pog instigated]"
					} else {
						logStr += "[Non-Pog instigated]"
					}
					audit.Log(logStr)
					ResetCache()
					return "correct"
				}
			}
		}

	}

	return ""

}

func SendCorrection(message *discordgo.MessageCreate) {
	showText := helpers.GetRandomText(pogCorrectText)
	messageObj := logger.MessageRefObj(message.Message)
	// Send Message
	_, err := config.Session.ChannelMessageSendReply(message.ChannelID, showText, &messageObj)
	if err == nil {
		audit.Log("Sent correction")
	} else {
		audit.Error(err)
		logger.SendError(message)
		return
	}
}

func CorrectionFromEdit(message *discordgo.MessageUpdate) {
	showText := helpers.GetRandomText(pogCorrectText)
	messageObj := logger.MessageRefObj(message.Message)
	// Send Message
	_, err := config.Session.ChannelMessageSendReply(message.ChannelID, showText, &messageObj)
	if err == nil {
		audit.Log("Corrected from Edit")
	} else {
		audit.Error(err)
		logger.SendErrorFromEdit(message)
		return
	}
}

func ResetCache() {
	lastMsg1 = ""
	lastMsg2 = ""
	lastMsg3 = ""
	lastMsg4 = ""
	lastMsg5 = ""
	lastMsg6 = ""
}

func GetPogString(s string) string {
	s = strings.ToLower(s)
	s = strings.Replace(s, "\n", " ", -1)
	lastMsg1 = lastMsg2
	lastMsg2 = lastMsg3
	lastMsg3 = lastMsg4
	lastMsg4 = lastMsg5
	lastMsg5 = lastMsg6
	lastMsg6 = s

	retString := lastMsg1 + " " + lastMsg2 + " " + lastMsg3 + " " + lastMsg4 + " " + lastMsg5 + " " + lastMsg6
	return strings.TrimSpace(retString)

}

func WithinTime(t1, t2 time.Time) bool {
	// Calculate the absolute difference between t1 and t2
	diff := t1.Sub(t2)
	if diff < 0 {
		diff = -diff
	}

	// Check if the difference is within the check minutes variable
	return diff <= checkMinutes*time.Minute
}
