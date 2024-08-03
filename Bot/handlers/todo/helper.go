package todo

import (
	"strings"

	"github.com/ZestHusky/femboy-control/Bot/helpers"
	"github.com/bwmarrin/discordgo"
)

func ToDoTypeSelector() []*discordgo.ApplicationCommandOptionChoice {

	var returnValue []*discordgo.ApplicationCommandOptionChoice

	for _, toDoType := range ToDoTypes {
		returnValue = append(returnValue, &discordgo.ApplicationCommandOptionChoice{
			Name:  toDoType.Description + " (" + toDoType.Name + ", Abbreviation: " + toDoType.Abbreviator + ")",
			Value: toDoType.Name,
		})
	}

	return returnValue

}

func GetCategoryAndAssignedIDFromID(id string) (string, int) {

	foundCat := ""
	foundID := 0

	inIDAlpha := strings.ToUpper(helpers.GetLettersOnlyCharactersFromString(id))
	inIDNumbers := helpers.GetNumbersOnlyCharactersFromString(id)
	for _, cat := range ToDoTypes {
		if inIDAlpha == cat.Name || inIDAlpha == cat.Abbreviator {
			foundCat = cat.Name
			foundID = inIDNumbers
			break
		}
	}

	if inIDAlpha == "" || foundID == 0 {
		return "", 0
	}

	return foundCat, foundID
}
