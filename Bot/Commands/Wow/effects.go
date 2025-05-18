package wow

import (
	"fmt"
	"strconv"
	"time"

	helpers "github.com/hashbat-dev/discgo-bot/Helpers"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

type Effect struct {
	Name            string
	Description     string
	SkipStatsOutput bool
	Emoji           string
	FromShop        bool
}

type EffectList func(*Generation) []*Effect

var staticEffectList = []EffectList{
	staticMessageID,
	staticTimeAndDate,
	staticWeather,
	staticPityBonus,
	staticFactWithStats,
	staticPokemon,
}

var rollEffectList = []EffectList{
	rollRollNumber,
	rollRandomExtras,
	rollStreakCheck,
}

// Static Effects ========================================================================
func staticMessageID(wow *Generation) []*Effect {
	var ret []*Effect
	if len(wow.Message.ID) <= 5 {
		return nil
	}

	// MessageID endings
	if wow.Message.ID[len(wow.Message.ID)-2:] == "69" {
		wow.BonusRolls += 3
		ret = append(ret, &Effect{
			Name:        "Niceeee",
			Description: "Message ID ended in 69, get 3 free rolls.",
		})
	}

	if wow.Message.ID[len(wow.Message.ID)-3:] == "420" {
		i := helpers.GetRandomNumber(6, 9)
		wow.BonusRolls += i
		ret = append(ret, &Effect{
			Name:        "Blaze it",
			Description: "Message ID ended in 420! Get a random number of free rolls between 6 and 9.",
		})
	}

	if wow.Message.ID[len(wow.Message.ID)-3:] == "365" {
		wow.Multiplier *= 0.8
		ret = append(ret, &Effect{
			Name:        "Bad Vibes",
			Description: "Something about having your Message ID end in 365 lowered your multiplier by 0.8x.",
		})
	}

	// MessageID dubs/trips
	matchingTrailingNumbers := countMatchingLastDigits(wow.Message.ID)
	var multi float64
	multiTriggered := false
	switch matchingTrailingNumbers {
	case 1:
		multiTriggered = false
	case 2:
		multiTriggered = true
		multi = 1.2
	case 3:
		multiTriggered = true
		multi = 1.6
	case 4:
		multiTriggered = true
		multi = 2
	default:
		multiTriggered = true
		multi = 2.5
	}

	if multiTriggered {
		wow.Multiplier *= multi
		ret = append(ret, &Effect{
			Name:        "Check 'em",
			Description: fmt.Sprintf("The last %d digits of your Message ID match! Get a %.1fx multiplier", matchingTrailingNumbers, multi),
			Emoji:       "👉",
		})
	}

	// MessageID sum
	idSum, err := sumDigits(wow.Message.ID)
	if err != nil && endsInZero(idSum) {
		wow.MinContinue--
		ret = append(ret, &Effect{
			Name:        "O.C.D.",
			Description: fmt.Sprintf("Adding up all the digits of your Message ID makes a nice, tidy rounded number! Your min continue roll was reduced to %d", wow.MinContinue),
		})
	}

	return ret
}

func staticTimeAndDate(wow *Generation) []*Effect {
	var ret []*Effect
	timeNow := time.Now()

	if timeNow.Weekday() == time.Saturday || timeNow.Weekday() == time.Sunday {
		wow.BonusRolls++
		ret = append(ret, &Effect{
			Name:        "Weekend Sweetener",
			Description: "It's the Weekend! Have a roll on the house.",
		})
	}

	if len(wow.Message.ID) <= 5 {
		return ret
	}

	yearDay := strconv.Itoa(time.Now().YearDay())
	if wow.Message.ID[len(wow.Message.ID)-len(yearDay):] == yearDay {
		wow.BonusRolls += 5
		ret = append(ret, &Effect{
			Name:        "Calendar Maxxing",
			Description: fmt.Sprintf("Today is the %s day of the year, the same number your Message ID ends in! Have 5 bonus rolls!", yearDay),
			Emoji:       "📆",
		})
	}

	return ret
}

func staticPityBonus(wow *Generation) []*Effect {
	dataLowRankLock.RLock()
	defer dataLowRankLock.RUnlock()
	if lowerRankUserId, exists := dataLowestWowRank[wow.Message.GuildID]; exists {
		if lowerRankUserId != wow.Message.Author.ID {
			return nil
		}
	} else {
		return nil
	}

	wow.Multiplier *= 1.2
	wow.BonusRolls += 2
	return []*Effect{
		{
			Name:        "Pity Bonus",
			Description: "You're the lowest ranking Wower in the server. Damn bro, have a 1.2x multiplier and 2 free rolls",
			Emoji:       "🥺",
		},
	}
}

func staticWeather(wow *Generation) []*Effect {
	var ret []*Effect

	// Temperature
	if dataCurrentWeather.Current.Temperature2m > 18.0 {
		wow.Multiplier *= 0.8
		ret = append(ret, &Effect{
			Name: "Go Outside",
			Description: fmt.Sprintf("It's %.1f%s in Manchester right now, you should be touching grass. 0.8x multiplier.",
				dataCurrentWeather.Current.Temperature2m, dataCurrentWeather.CurrentUnits.Temperature2m),
			Emoji: "🌞",
		})
	} else if dataCurrentWeather.Current.Temperature2m < 5.0 {
		wow.Multiplier *= 1.2
		ret = append(ret, &Effect{
			Name: "Bit Nippy",
			Description: fmt.Sprintf("It's %.1f%s in Manchester right now, stay inside and stay warm. 1.2x multiplier.",
				dataCurrentWeather.Current.Temperature2m, dataCurrentWeather.CurrentUnits.Temperature2m),
			Emoji: "❄️",
		})
	}

	// Clouds
	if dataCurrentWeather.Current.CloudCover == 0 {
		wow.BonusRolls++
		ret = append(ret, &Effect{
			Name:        "Cloudless",
			Description: "There's not a cloud in the sky in Manchester right now, how neat! Have a bonus roll.",
			Emoji:       "☀️",
		})
	}

	// Wind
	if dataCurrentWeather.Current.WindSpeed10m > 10.0 {
		wow.BonusRolls += 3
		ret = append(ret, &Effect{
			Name:        "Wimdy",
			Description: fmt.Sprintf("The wind is %.1f%s in Manchester right now, it blew 3 bonus rolls your way!", dataCurrentWeather.Current.WindSpeed10m, dataCurrentWeather.CurrentUnits.WindSpeed10m),
			Emoji:       "💨",
		})
	}

	// Night
	if dataCurrentWeather.Current.IsDay == 0 {
		wow.BonusRolls++
		ret = append(ret, &Effect{
			Name:        "Night night",
			Description: "It's nighttime in Manchester, will a bonus roll help get you into bed?",
			Emoji:       "🌜",
		})
	}

	// Rain
	if dataCurrentWeather.Current.Rain > 0.0 {
		wow.MinContinue--
		ret = append(ret, &Effect{
			Name:        "Pretty Moist",
			Description: fmt.Sprintf("There's %.1f%s in Manchester right now (of course), lets lower your min continue roll to %d", dataCurrentWeather.Current.Rain, dataCurrentWeather.CurrentUnits.Rain, wow.MinContinue),
			Emoji:       "🌧️",
		})
	}

	return ret
}

func staticFactWithStats(wow *Generation) []*Effect {
	fact, hasStats := checkIfRandomFactHasStats()
	if !hasStats {
		return nil
	}
	return []*Effect{
		{
			Name:        "Number Nerd",
			Description: fmt.Sprintf("Looked up a random fact and it contained a number, have a 1.2x multiplier! The fact was: _%s_", fact),
			Emoji:       "🤓",
		},
	}
}

func staticPokemon(wow *Generation) []*Effect {
	if !pokeInit {
		return nil
	}

	randomNumber := helpers.GetRandomNumber(1, len(dataPokemon)) - 1
	var pokemon *PokemonData

	if poke, exists := dataPokemon[randomNumber]; exists {
		pokemon = &poke
	} else {
		logger.ErrorText("WOW", "Pokémon ID %d was not in the data cache", randomNumber)
		return nil
	}

	ret := getPokemonEffects(wow, pokemon)
	pokemonName := getPokemonName(pokemon.Name)
	title := fmt.Sprintf("A wild %s appeared!", pokemonName)
	description := ""
	if len(ret) > 0 {
		description += " Gain the following effects..."
		for _, effect := range ret {
			ret = append(ret, effect)
			description += fmt.Sprintf("\n%s%s %s", IndentPadding, effect.Emoji, effect.Description)
		}
	}
	if description == "" {
		description = "You didn't get any effects from it, but hey.. still pretty neat right?"
	}

	ret = append(ret, &Effect{
		Name:        title,
		Description: description,
		Emoji:       "🔴",
	})

	return ret
}

// Roll Based Effects ====================================================================
func rollRollNumber(wow *Generation) []*Effect {
	var ret []*Effect

	if wow.CurrentRoll >= 10 {
		wow.BonusRolls++
		ret = append(ret, &Effect{
			Name:        "Crit Roll",
			Description: fmt.Sprintf("You rolled a %d! Have another roll on the house", wow.CurrentRoll),
			Emoji:       "🎯",
		})
	}

	if wow.CurrentRoll == 1 {
		wow.Multiplier *= 0.9
		ret = append(ret, &Effect{
			Name:        "Bruh",
			Description: "This Wow is ass. Session terminated. 0.9x multiplier",
			Emoji:       "⛔",
		})
	}

	return ret
}

func rollRandomExtras(wow *Generation) []*Effect {
	var ret []*Effect

	// Current Roll randomiser
	i := helpers.GetRandomNumber(1, 50)
	description := ""
	switch i {
	case 50:
		description = "Your current roll got Quadrupled! (1/50 chance)"
		wow.CurrentRoll = wow.CurrentRoll * 4
	case 1, 25:
		description = "Your current total got Doubled! (1/25 chance)"
		wow.CurrentRoll = wow.CurrentRoll * 2
	}
	if description != "" {
		ret = append(ret, &Effect{
			Name:        "Rare Roll",
			Description: description,
		})
	}

	// Current Total randomiser
	i = helpers.GetRandomNumber(1, 100)
	description = ""
	switch i {
	case 100:
		description = "Your current total got Quadrupled! (1/100 chance)"
		wow.OCount = ((wow.OCount + wow.CurrentRoll) * 4) - wow.CurrentRoll
	case 1, 25, 50, 75:
		description = "Your current total got Doubled! (1/25 chance)"
		wow.OCount = ((wow.OCount + wow.CurrentRoll) * 2) - wow.CurrentRoll
	}

	if description != "" {
		ret = append(ret, &Effect{
			Name:        "Rare Roll",
			Description: description,
		})
	}

	// Death Dodger
	i = helpers.GetRandomNumber(1, 300)
	description = ""
	switch i {
	case 300:
		oldMin := wow.MinContinue
		newMin := wow.MinContinue - 2
		if newMin < 1 {
			newMin = 1
		}
		wow.MinContinue = newMin
		description = fmt.Sprintf("Your min continue roll decreased from %d to %d! (1/300 chance)", oldMin, newMin)
	case 100, 200:
		oldMin := wow.MinContinue
		newMin := wow.MinContinue - 1
		if newMin < 1 {
			newMin = 1
		}
		wow.MinContinue = newMin
		description = fmt.Sprintf("Your min continue roll decreased from %d to %d! (1/150 chance)", oldMin, newMin)
	}

	if description != "" {
		ret = append(ret, &Effect{
			Name:        "Death Dodger",
			Description: description,
		})
	}
	return ret
}

func rollStreakCheck(wow *Generation) []*Effect {
	var ret []*Effect

	if len(wow.DiceRolls) >= 2 &&
		wow.DiceRolls[len(wow.DiceRolls)-1].Roll == wow.DiceRolls[len(wow.DiceRolls)-2].Roll &&
		wow.DiceRolls[len(wow.DiceRolls)-1].Roll == wow.CurrentRoll {
		// Check for Triple
		wow.Multiplier *= 1.4
		ret = append(ret, &Effect{
			Name:        "Oh Baby a Triple",
			Description: "You rolled the same number 3 times in a row, get a 1.4x multiplier!",
			Emoji:       "🎳",
		})
	} else if len(wow.DiceRolls) >= 1 && wow.DiceRolls[len(wow.DiceRolls)-1].Roll == wow.CurrentRoll {
		// Check for Double
		wow.Multiplier *= 1.2
		ret = append(ret, &Effect{
			Name:        "Dirty Double",
			Description: "You rolled the same number twice, get a 1.2x multiplier!",
			Emoji:       "🎳",
		})
	}

	return ret
}
