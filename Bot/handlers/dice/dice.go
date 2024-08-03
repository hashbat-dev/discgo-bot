package dice

import (
	"fmt"

	"github.com/ZestHusky/femboy-control/Bot/audit"
	"github.com/ZestHusky/femboy-control/Bot/helpers"
	"github.com/bwmarrin/discordgo"
)

type DiceRolls struct {
	Die  string
	Roll int
	Max  int
}

func RollDice(interaction *discordgo.InteractionCreate) {

	optionMap := helpers.GetOptionMap(interaction)
	inD2 := helpers.GetOptionIntValue(optionMap, "d2")
	inD6 := helpers.GetOptionIntValue(optionMap, "d6")
	inD10 := helpers.GetOptionIntValue(optionMap, "d10")
	inD20 := helpers.GetOptionIntValue(optionMap, "d20")
	inD50 := helpers.GetOptionIntValue(optionMap, "d50")
	inD100 := helpers.GetOptionIntValue(optionMap, "d100")

	audit.Log(fmt.Sprintf("Dice Roll: (D2 %v), (D6 %v), (D10 %v), (D20 %v), (D50 %v), (D100 %v)", inD2, inD6, inD10, inD20, inD50, inD100))

	// Do we have at least one dice?
	responseTitle := "🎲 Roll Dice 🎲"
	totalDice := inD2 + inD6 + inD10 + inD20 + inD50 + inD100
	if totalDice == 0 {
		audit.SendInteractionResponse(interaction, responseTitle, "", "You didn't give me any dice! >w< What am I going to roll now? I'm not Levi..", true, true, "")
		return
	} else if totalDice > 50 {
		audit.SendInteractionResponse(interaction, responseTitle, "", "Bruh, I'm not rolling more than 50 Dice. My hands will get tired >w<", true, true, "")
		return
	}

	// Roll the Dice!
	var diceRolls []DiceRolls

	for i := 0; i < inD2; i++ {
		diceRolls = append(diceRolls, DiceRolls{
			Die:  "D2",
			Roll: helpers.GetRandomNumber(1, 2),
			Max:  2,
		})
	}

	for i := 0; i < inD6; i++ {
		diceRolls = append(diceRolls, DiceRolls{
			Die:  "D6",
			Roll: helpers.GetRandomNumber(1, 6),
			Max:  6,
		})
	}

	for i := 0; i < inD10; i++ {
		diceRolls = append(diceRolls, DiceRolls{
			Die:  "D10",
			Roll: helpers.GetRandomNumber(1, 10),
			Max:  10,
		})
	}

	for i := 0; i < inD20; i++ {
		diceRolls = append(diceRolls, DiceRolls{
			Die:  "D20",
			Roll: helpers.GetRandomNumber(1, 20),
			Max:  20,
		})
	}

	for i := 0; i < inD50; i++ {
		diceRolls = append(diceRolls, DiceRolls{
			Die:  "D50",
			Roll: helpers.GetRandomNumber(1, 50),
			Max:  50,
		})
	}

	for i := 0; i < inD100; i++ {
		diceRolls = append(diceRolls, DiceRolls{
			Die:  "D100",
			Roll: helpers.GetRandomNumber(1, 100),
			Max:  100,
		})
	}

	// Get the Output
	outputString := "You rolled " + fmt.Sprint(totalDice) + " dice..."

	totalRoll := 0
	totalMax := 0
	for _, roll := range diceRolls {
		totalRoll += roll.Roll
		totalMax += roll.Max
		outputString += "\n🎲 **" + roll.Die + "**: Rolled " + fmt.Sprint(roll.Roll)
	}

	footerText := "🎲 Total roll: " + fmt.Sprint(totalRoll) + ", Max Possible: " + fmt.Sprint(totalMax)

	audit.SendInteractionResponse(interaction, responseTitle, outputString, footerText, false, false, "")

}
