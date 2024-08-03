package adventures

import (
	"encoding/json"
	"os"

	"github.com/dabi-ngin/discgo-bot/Bot/audit"
)

type Adventure struct {
	Name        string
	Rarity      string
	StartPrompt string
	Options     []string
	Conclusions []string
	Identifier  string
}

func GetNewAdventure(rarity string, index int) Adventure {

	files, err := os.ReadDir("Bot/handlers/adventures/adventurefiles/" + rarity)

	if err != nil {
		audit.Error(err)
	}

	var fileList []string

	for _, file := range files {
		fileList = append(fileList, "Bot/handlers/adventures/adventurefiles/"+rarity+"/"+file.Name())
	}

	fileText, err := os.ReadFile(fileList[index])

	if err != nil {
		audit.Error(err)
	}

	var newAdventure Adventure

	jsonerr := json.Unmarshal([]byte(fileText), &newAdventure)
	if jsonerr != nil {
		audit.Error(jsonerr)
	}

	newAdventure = FormatAdventureText(newAdventure)

	return newAdventure

}

func FormatAdventureText(newAdventure Adventure) Adventure {

	newAdventure.Name = "# :placard:**" + newAdventure.Name + "**"

	switch newAdventure.Rarity {
	case "common":
		newAdventure.Rarity = ":poop: Common :poop:"
	case "uncommon":
		newAdventure.Rarity = ":coin: Uncommon :coin:"
	case "rare":
		newAdventure.Rarity = ":medal: Rare :medal:"
	case "legendary":
		newAdventure.Rarity = ":gem: Legendary :gem:"
	case "god-tier":
		newAdventure.Rarity = ":trophy: HOLY :fire: FUCKING :fire: SHIT :fire: IT'S :fire: GOD-TIER! :fire: :trophy:"
	}

	newAdventure.Rarity = "### Rarity: " + newAdventure.Rarity

	for i, option := range newAdventure.Options {
		newAdventure.Options[i] = "*" + option + "*"
	}

	return newAdventure
}
