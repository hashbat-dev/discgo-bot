package wow

import (
	"fmt"
	"strconv"
	"time"
)

type Effect struct {
	Name        string
	Description string
	Emoji       string
}

type EffectList func(*Generation) *Effect

var staticEffectList = []EffectList{
	staticSixetyNine,
	staticBlazeIt,
	staticWeekend,
	staticSpecificNumber,
	staticDayOfTheYear,
	staticMessageDubs,
	staticMessageIDSum,
	staticPityBonus,
	staticWeatherTemperature,
	staticWeatherClouds,
	staticWeatherWindy,
	staticWeatherNight,
	staticWeatherRaining,
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
	}
}

func staticWeekend(wow *Generation) *Effect {
	today := time.Now().Weekday()
	if today == time.Saturday || today == time.Sunday {
		wow.BonusRolls++
		return &Effect{
			Name:        "Weekend Sweetener",
			Description: "It's the Weekend! Have a roll on the house.",
		}
	}
	return nil
}

func staticSpecificNumber(wow *Generation) *Effect {
	if len(wow.Message.ID) <= 2 || wow.Message.ID[len(wow.Message.ID)-3:] != "365" {
		return nil
	}

	wow.BonusRolls += 3
	return &Effect{
		Name:        "Bad Vibes",
		Description: "Something about having your Message ID end in 365 lowered your multiplier by 0.8x.",
	}
}

func staticDayOfTheYear(wow *Generation) *Effect {
	yearDay := strconv.Itoa(time.Now().YearDay())

	if len(wow.Message.ID) <= 2 || wow.Message.ID[len(wow.Message.ID)-len(yearDay):] != yearDay {
		return nil
	}

	wow.BonusRolls += 5
	return &Effect{
		Name:        "Calendar Maxxing",
		Description: fmt.Sprintf("Today is the %s day of the year, the same number your Message ID ends in! Have 5 bonus rolls!", yearDay),
		Emoji:       "📆",
	}
}

func staticMessageDubs(wow *Generation) *Effect {
	matchingTrailingNumbers := countMatchingLastDigits(wow.Message.ID)
	var multi float64
	switch matchingTrailingNumbers {
	case 1:
		return nil
	case 2:
		multi = 1.2
	case 3:
		multi = 1.6
	case 4:
		multi = 2
	default:
		multi = 2.5
	}

	wow.Multiplier *= multi
	return &Effect{
		Name:        "Check 'em",
		Description: fmt.Sprintf("The last %d digits of your Message ID match! Get a %fx multiplier", matchingTrailingNumbers, multi),
		Emoji:       "👉",
	}
}

func staticMessageIDSum(wow *Generation) *Effect {
	idSum, err := sumDigits(wow.Message.ID)
	if err != nil {
		return nil
	}

	if !endsInZero(idSum) {
		return nil
	}
	wow.MinContinue--
	return &Effect{
		Name:        "O.C.D.",
		Description: fmt.Sprintf("Adding up all the digits of your Message ID makes a nice, tidy rounded number! Your min continue roll was reduced to %d", wow.MinContinue),
	}
}

func staticPityBonus(wow *Generation) *Effect {
	if lowerRankUserId, exists := dataLowestWowRank[wow.Message.GuildID]; exists {
		if lowerRankUserId != wow.Message.Author.ID {
			return nil
		}
	} else {
		return nil
	}

	wow.Multiplier *= 1.2
	wow.BonusRolls += 2
	return &Effect{
		Name:        "Pity Bonus",
		Description: "You're the lowest ranking Wower in the server. Damn bro, have a 1.2x multiplier and 2 free rolls",
		Emoji:       "🥺",
	}
}

func staticWeatherTemperature(wow *Generation) *Effect {
	if dataCurrentWeather.Current.Temperature2m > 18.0 {
		wow.Multiplier *= 0.8
		return &Effect{
			Name: "Go Outside",
			Description: fmt.Sprintf("It's %f%s in Manchester right now, you should be touching grass. 0.8x multiplier.",
				dataCurrentWeather.Current.Temperature2m, dataCurrentWeather.CurrentUnits.Temperature2m),
			Emoji: "🌞",
		}
	} else if dataCurrentWeather.Current.Temperature2m < 5.0 {
		wow.Multiplier *= 1.2
		return &Effect{
			Name: "Bit Nippy",
			Description: fmt.Sprintf("It's %f%s in Manchester right now, stay inside and stay warm. 1.2x multiplier.",
				dataCurrentWeather.Current.Temperature2m, dataCurrentWeather.CurrentUnits.Temperature2m),
			Emoji: "❄️",
		}
	}
	return nil
}

func staticWeatherClouds(wow *Generation) *Effect {
	if dataCurrentWeather.Current.CloudCover == 0 {
		wow.BonusRolls++
		return &Effect{
			Name:        "Cloudless",
			Description: "There's not a cloud in the sky in Manchester right now, how neat! Have a bonus roll.",
			Emoji:       "☀️",
		}
	}
	return nil
}

func staticWeatherWindy(wow *Generation) *Effect {
	if dataCurrentWeather.Current.WindSpeed10m > 10.0 {
		wow.BonusRolls += 3
		return &Effect{
			Name:        "Wimdy",
			Description: fmt.Sprintf("The wind is %f%s in Manchester right now, it blew 3 bonus rolls your way!", dataCurrentWeather.Current.WindSpeed10m, dataCurrentWeather.CurrentUnits.WindSpeed10m),
			Emoji:       "💨",
		}
	}
	return nil
}

func staticWeatherNight(wow *Generation) *Effect {
	if dataCurrentWeather.Current.IsDay == 0 {
		wow.BonusRolls++
		return &Effect{
			Name:        "Night night",
			Description: "It's nighttime in Manchester, will a bonus roll help get you into bed?",
			Emoji:       "🌜",
		}
	}
	return nil
}

func staticWeatherRaining(wow *Generation) *Effect {
	if dataCurrentWeather.Current.Rain > 0.0 {
		wow.MinContinue--
		return &Effect{
			Name:        "Pretty Moist",
			Description: fmt.Sprintf("There's %f%s in Manchester right now (of course), lets lower your min continue roll to %d", dataCurrentWeather.Current.Rain, dataCurrentWeather.CurrentUnits.Rain, wow.MinContinue),
			Emoji:       "🌧️",
		}
	}
	return nil
}

// Roll Based Effects ====================================================================
func rollRolled10(wow *Generation) *Effect {
	if wow.CurrentRoll < 10 {
		return nil
	}

	wow.BonusRolls++
	return &Effect{
		Name:        "Crit Roll",
		Description: "You rolled a 10! Have another roll on the house",
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

	switch i {
	case 100:
		description = "Your current total got Quadrupled! (1/100 chance)"
		wow.OCount = ((wow.OCount + wow.CurrentRoll) * 4) - wow.CurrentRoll
	case 1, 25, 50, 75:
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

	switch i {
	case 300:
		oldMin := wow.MinContinue
		newMin := wow.MinContinue + 2
		if newMin >= wow.MaxRollValue {
			newMin = wow.MaxRollValue - 1
		}
		wow.MinContinue = newMin
		description = fmt.Sprintf("Your min continue roll increased from %d to %d! (1/300 chance)", oldMin, newMin)
	case 100, 200:
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
		wow.Multiplier *= 1.4
		return &Effect{
			Name:        "Oh Baby a Triple",
			Description: "You rolled the same number 3 times in a row, get a 1.4x multiplier!",
			Emoji:       "🎳",
		}
	} else if len(wow.DiceRolls) >= 1 && wow.DiceRolls[len(wow.DiceRolls)-1].Roll == wow.CurrentRoll {
		// Check for Double
		wow.Multiplier *= 1.2
		return &Effect{
			Name:        "Dirty Double",
			Description: "You rolled the same number twice, get a 1.2x multiplier!",
			Emoji:       "🎳",
		}
	}

	return nil
}
