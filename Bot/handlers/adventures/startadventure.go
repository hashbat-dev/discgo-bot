package adventures

import (
	"os"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/dabi-ngin/discgo-bot/Bot/audit"
	"github.com/dabi-ngin/discgo-bot/Bot/config"
	"github.com/dabi-ngin/discgo-bot/Bot/helpers"
)

func StartAdventure(message *discordgo.MessageCreate) {
	rarity := DetermineRarity()
	index := DetermineIndex(rarity)

	newAdventure := GetNewAdventure(rarity, index)

	var replyMsg string

	replyMsg += newAdventure.Name
	replyMsg = Newline(replyMsg)

	replyMsg += newAdventure.Rarity
	replyMsg = NewParagraph(replyMsg)
	replyMsg = Newline(replyMsg)

	replyMsg += newAdventure.StartPrompt
	replyMsg = NewParagraph(replyMsg)

	i := 1

	for _, option := range newAdventure.Options {
		replyMsg += strconv.Itoa(i) + ": "
		replyMsg += option
		replyMsg = Newline(replyMsg)
		i += 1
	}

	replyMsg = Newline(replyMsg)
	replyMsg += "Select an option with !adv x"
	replyMsg = Newline(replyMsg)

	replyMsg = Newline(replyMsg)
	replyMsg += newAdventure.Identifier

	_, err := config.Session.ChannelMessageSendReply(message.ChannelID, replyMsg, message.Reference())

	if err != nil {
		audit.Error(err)
	}
}

func Newline(text string) string {
	return text + "\n"
}

func NewParagraph(text string) string {
	text = Newline(text)
	text = Newline(text)
	return text
}

func DetermineIndex(rarity string) int {
	files, err := os.ReadDir("Bot/handlers/adventures/adventurefiles/" + rarity)

	if err != nil {
		audit.Error(err)
	}

	return helpers.GetRandomNumber(0, len(files))
}

func DetermineRarity() string {
	roll := helpers.GetRandomNumber(1, 100)

	if helpers.IsIntBetweenXandY(roll, 0, 50) {
		return "common"
	} else if helpers.IsIntBetweenXandY(roll, 51, 75) {
		return "uncommon"
	} else if helpers.IsIntBetweenXandY(roll, 76, 90) {
		return "rare"
	} else if helpers.IsIntBetweenXandY(roll, 91, 99) {
		return "legendary"
	} else if roll == 100 {
		return "god-tier"
	}
	return "common"
}
