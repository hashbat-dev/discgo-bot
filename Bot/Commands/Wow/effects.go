package wow

import (
	"fmt"
	"time"
)

type Effect struct {
	Name        string
	Description string
	Emoji       string
	RollNumber  int
	BonusRolls  int
	Multiplier  float32
}

type EffectList func(*Generation) *Effect

var staticEffectList = []EffectList{
	staticTest,
	staticSixetyNine,
	staticBlazeIt,
	staticWeekend,
}

var rollEffectList = []EffectList{
	rollRolled10,
	rollCountMultiplier,
	rollTotalIncreaseChance,
	rollDeathDodger,
	rollStreakCheck,
}

// Static Effects ========================================================================
func staticSixetyNine(wow *Generation) *Effect {
	if len(wow.Message.ID) <= 2 || wow.Message.ID[len(wow.Message.ID)-2:] != "69" {
		return nil
	}

	wow.BonusRolls += 3
	return &Effect{
		Name:        "Niceeee",
		Description: "Message ID ended in 69, get 3 free rolls.",
		BonusRolls:  3,
	}
}

func staticBlazeIt(wow *Generation) *Effect {
	if len(wow.Message.ID) <= 1 || wow.Message.ID[len(wow.Message.ID)-3:] != "420" {
		return nil
	}

	i := getRandomNumber(6, 9)
	wow.BonusRolls += i
	return &Effect{
		Name:        "Blaze it",
		Description: "Message ID ended in 420! Get a random number of free rolls between 6 and 9.",
		BonusRolls:  i,
	}
}

func staticWeekend(wow *Generation) *Effect {
	today := time.Now().Weekday()
	if today == time.Saturday || today == time.Sunday {
		wow.BonusRolls++
		return &Effect{
			Name:        "Weekend Sweetener",
			Description: "It's the Weekend! Have a roll on the house.",
			BonusRolls:  1,
		}
	}
	return nil
}

func staticTest(wow *Generation) *Effect {
	return &Effect{
		Name:        "Testing",
		Description: "AHHHHHHHHHH",
		BonusRolls:  1,
	}
}

// Roll Based Effects ====================================================================
func rollRolled10(wow *Generation) *Effect {
	if wow.CurrentRoll < 10 {
		return nil
	}

	return &Effect{
		Name:        "Crit Roll",
		Description: "You rolled a 10! Have another roll on the house",
		BonusRolls:  1,
		Emoji:       "🎯",
	}
}

func rollCountMultiplier(wow *Generation) *Effect {
	i := getRandomNumber(1, 50)
	description := ""
	switch i {
	case 50:
		description = "Your current roll got Quadrupled! (1/50 chance)"
		wow.CurrentRoll = wow.CurrentRoll * 4
	case 1, 25:
		description = "Your current total got Doubled! (1/25 chance)"
		wow.CurrentRoll = wow.CurrentRoll * 2
	}
	if description == "" {
		return nil
	} else {
		return &Effect{
			Name:        "Rare Roll",
			Description: description,
		}
	}
}

func rollTotalIncreaseChance(wow *Generation) *Effect {
	i := getRandomNumber(1, 100)
	description := ""

	if i == 100 {
		description = "Your current total got Quadrupled! (1/100 chance)"
		wow.OCount = ((wow.OCount + wow.CurrentRoll) * 4) - wow.CurrentRoll
	} else if i == 1 || i == 25 || i == 50 || i == 75 {
		description = "Your current total got Doubled! (1/25 chance)"
		wow.OCount = ((wow.OCount + wow.CurrentRoll) * 2) - wow.CurrentRoll
	}
	if description == "" {
		return nil
	} else {
		return &Effect{
			Name:        "Rare Roll",
			Description: description,
		}
	}
}

func rollDeathDodger(wow *Generation) *Effect {
	i := getRandomNumber(1, 300)
	description := ""

	if i == 300 {
		oldMin := wow.MinContinue
		newMin := wow.MinContinue + 2
		if newMin >= wow.MaxRollValue {
			newMin = wow.MaxRollValue - 1
		}
		wow.MinContinue = newMin
		description = fmt.Sprintf("Your min continue roll increased from %d to %d! (1/300 chance)", oldMin, newMin)
	} else if i == 100 || i == 200 {
		oldMin := wow.MinContinue
		newMin := wow.MinContinue + 1
		if newMin >= wow.MaxRollValue {
			newMin = wow.MaxRollValue - 1
		}
		wow.MinContinue = newMin
		description = fmt.Sprintf("Your min continue roll increased from %d to %d! (1/150 chance)", oldMin, newMin)
	}

	if description == "" {
		return nil
	} else {
		return &Effect{
			Name:        "Death Dodger",
			Description: description,
		}
	}
}

func rollStreakCheck(wow *Generation) *Effect {
	if len(wow.DiceRolls) >= 2 &&
		wow.DiceRolls[len(wow.DiceRolls)-1].Roll == wow.DiceRolls[len(wow.DiceRolls)-2].Roll &&
		wow.DiceRolls[len(wow.DiceRolls)-1].Roll == wow.CurrentRoll {
		// Check for Triple
		wow.Multiplier = wow.Multiplier * 1.4
		return &Effect{
			Name:        "Oh Baby a Triple",
			Description: "You rolled the same number 3 times in a row, get a 1.4x multiplier!",
			Emoji:       "🎳",
		}
	} else if len(wow.DiceRolls) >= 1 && wow.DiceRolls[len(wow.DiceRolls)-1].Roll == wow.CurrentRoll {
		// Check for Double
		wow.Multiplier = wow.Multiplier * 1.2
		return &Effect{
			Name:        "Dirty Double",
			Description: "You rolled the same number twice, get a 1.2x multiplier!",
			Emoji:       "🎳",
		}
	}

	return nil
}
