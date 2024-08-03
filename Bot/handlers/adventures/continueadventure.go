package adventures

import (
	"bufio"
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/dabi-ngin/discgo-bot/Bot/audit"
	"github.com/dabi-ngin/discgo-bot/Bot/config"
)

func ContinueAdventure(message *discordgo.MessageCreate) {
	refcode := GetRefCode(message)

	matchingAdventure := GetMatchingAdventureFromRefCode(refcode)

	selection := message.Content

	//just get the int
	selection = selection[len(selection)-1:]

	selectionInt, err1 := strconv.Atoi(selection)

	if err1 != nil {
		audit.Error(err1)
	}

	if len(matchingAdventure.Conclusions) < selectionInt-1 {
		_, err := config.Session.ChannelMessageSendReply(message.ChannelID, "Pick a number that's in the options, moron.", message.Reference())
		if err != nil {
			audit.Error(err)
		}
		return
	}

	replyMsg := matchingAdventure.Conclusions[selectionInt-1]

	_, err := config.Session.ChannelMessageSendReply(message.ChannelID, replyMsg, message.Reference())

	if err != nil {
		audit.Error(err)
	}

}

func GetMatchingAdventureFromRefCode(refcode string) Adventure {
	rarity := GetRarity(refcode)

	files, err := os.ReadDir("Bot/handlers/adventures/adventurefiles/" + rarity)

	if err != nil {
		audit.Error(err)
	}

	var newAdventure Adventure
	for _, file := range files {
		filePath := "Bot/handlers/adventures/adventurefiles/" + rarity + "/" + file.Name()
		fileText, err := os.ReadFile(filePath)

		if err != nil {
			audit.Error(err)
		}

		jsonerr := json.Unmarshal([]byte(fileText), &newAdventure)
		if err != nil {
			audit.Error(jsonerr)
		}

		if newAdventure.Identifier == refcode {
			return newAdventure
		}
	}
	return newAdventure
}

func GetRefCode(message *discordgo.MessageCreate) string {
	originalMsg := message.ReferencedMessage.Content

	var lines []string

	strScanner := bufio.NewScanner(strings.NewReader(originalMsg))
	for strScanner.Scan() {
		lines = append(lines, strScanner.Text())
	}

	return lines[len(lines)-1]
}

func GetRarity(refcode string) string {

	if refcode[0:1] == "c" {
		return "common"
	} else if refcode[0:1] == "u" {
		return "uncommon"
	} else if refcode[0:1] == "r" {
		return "rare"
	} else if refcode[0:1] == "l" {
		return "legendary"
	} else if refcode[0:1] == "g" {
		return "god-tier"
	}
	return "common"
}
