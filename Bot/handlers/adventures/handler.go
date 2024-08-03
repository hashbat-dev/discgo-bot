package adventures

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	startAdventure    = iota
	continueAdvanture = iota
)

func HandleMessage(message *discordgo.MessageCreate) {
	var messageType int
	msgContent := strings.ToLower(message.Content)

	//get message content
	if msgContent == "!a" || msgContent == "!adventure" || msgContent == "!adv" || msgContent == "!quest" {
		messageType = startAdventure
	} else {
		messageType = continueAdvanture
	}

	//if it's just !adventure or !a and not a reply then start a new adventure
	switch messageType {
	case startAdventure:
		StartAdventure(message)
	case continueAdvanture:
		ContinueAdventure(message)
	}

	//if it's !adventure / !a followed by a number and IS a reply then activate that option
}
