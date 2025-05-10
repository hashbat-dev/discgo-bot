package wow

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"

	helpers "github.com/hashbat-dev/discgo-bot/Helpers"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

type NamedPokemonAPIResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type PokemonAbility struct {
	Ability  NamedPokemonAPIResource `json:"ability"`
	IsHidden bool                    `json:"is_hidden"`
	Slot     int                     `json:"slot"`
}

type PokemonStat struct {
	BaseStat int                     `json:"base_stat"`
	Effort   int                     `json:"effort"`
	Stat     NamedPokemonAPIResource `json:"stat"`
}

type PokemonTypeSlot struct {
	Slot int                     `json:"slot"`
	Type NamedPokemonAPIResource `json:"type"`
}

type PokemonData struct {
	Abilities      []PokemonAbility          `json:"abilities"`
	BaseExperience int                       `json:"base_experience"`
	Forms          []NamedPokemonAPIResource `json:"forms"`
	Height         int                       `json:"height"`
	ID             int                       `json:"id"`
	IsDefault      bool                      `json:"is_default"`
	Name           string                    `json:"name"`
	Order          int                       `json:"order"`
	Species        NamedPokemonAPIResource   `json:"species"`
	Stats          []PokemonStat             `json:"stats"`
	Types          []PokemonTypeSlot         `json:"types"`
	Weight         int                       `json:"weight"`
}

func getPokemonData(url string) *PokemonData {
	resp, err := http.Get(url)
	if err != nil {
		logger.Error("WOW", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("WOW", err)
		return nil
	}

	var data PokemonData
	if err := json.Unmarshal(body, &data); err != nil {
		logger.Error("WOW", err)
		return nil
	}

	return &data
}

func getPokemonEffects(wow *Generation, pokemon *PokemonData) []*Effect {
	var ret []*Effect

	// Type based ===========================================================
	// All types: https://pokeapi.co/api/v2/type/

	// Work out what effects we have in play
	hot := dataCurrentWeather.Current.Temperature2m > 18.0
	icy := dataCurrentWeather.Current.Temperature2m < 2.0
	rain := dataCurrentWeather.Current.Rain > 0
	night := dataCurrentWeather.Current.IsDay != 0
	cloudy := dataCurrentWeather.Current.CloudCover > 50
	noClouds := dataCurrentWeather.Current.CloudCover == 0
	windy := dataCurrentWeather.Current.WindSpeed10m > 15

	// Apply based on Type
	for _, pokeType := range pokemon.Types {
		switch pokeType.Type.Name {
		case "normal":
		case "fighting":
		case "flying":
			if cloudy {
				wow.Multiplier *= 0.8
				ret = append(ret, &Effect{
					Name: "Pokémon Flying Debuff",
					Description: fmt.Sprintf("Flying Type Debuff! There's currently %d%s cloud cover in Manchester. 0.8x multiplier.",
						dataCurrentWeather.Current.CloudCover, dataCurrentWeather.CurrentUnits.CloudCover),
					Emoji:           "☁️",
					SkipStatsOutput: true,
				})
			}
			if noClouds {
				wow.Multiplier *= 1.2
				ret = append(ret, &Effect{
					Name:            "Pokémon Flying Boost",
					Description:     "Flying Type Boost! There's currently no clouds in Manchester. 1.2x multiplier.",
					Emoji:           "☁️",
					SkipStatsOutput: true,
				})
			}
			if windy {
				wow.Multiplier *= 0.8
				ret = append(ret, &Effect{
					Name: "Pokémon Flying Debuff",
					Description: fmt.Sprintf("Flying Type Debuff! There's currently %.1f%s wind in Manchester. 0.8x multiplier.",
						dataCurrentWeather.Current.WindSpeed10m, dataCurrentWeather.CurrentUnits.WindSpeed10m),
					Emoji:           "💨",
					SkipStatsOutput: true,
				})
			}
		case "bug":
			if hot {
				wow.Multiplier *= 1.2
				ret = append(ret, &Effect{
					Name: "Pokémon Bug Boost",
					Description: fmt.Sprintf("Bug Type Boost! It's currently %.1f%s in Manchester. 1.2x multiplier.",
						dataCurrentWeather.Current.Temperature2m, dataCurrentWeather.CurrentUnits.Temperature2m),
					Emoji:           "☁️",
					SkipStatsOutput: true,
				})
			}
		case "ghost":
			if night {
				wow.Multiplier *= 1.2
				ret = append(ret, &Effect{
					Name:            "Pokémon Ghost Boost",
					Description:     "Ghost Type Boost! It's currently night time in Manchester. 1.2x multiplier.",
					Emoji:           "🌑",
					SkipStatsOutput: true,
				})
			} else {
				wow.Multiplier *= 0.8
				ret = append(ret, &Effect{
					Name:            "Pokémon Ghost Debuff",
					Description:     "Ghost Type Debuff! It's currently day time in Manchester. 0.8x multiplier.",
					Emoji:           "🌞",
					SkipStatsOutput: true,
				})
			}
		case "steel":
		case "fire":
			if hot {
				wow.Multiplier *= 1.2
				ret = append(ret, &Effect{
					Name: "Pokémon Fire Boost",
					Description: fmt.Sprintf("Fire Type Boost! It's currently %.1f%s in Manchester. 1.2x multiplier.",
						dataCurrentWeather.Current.Temperature2m, dataCurrentWeather.CurrentUnits.Temperature2m),
					Emoji:           "🔥",
					SkipStatsOutput: true,
				})
			}
			if rain {
				wow.Multiplier *= 0.8
				ret = append(ret, &Effect{
					Name:            "Pokémon Fire Debuff",
					Description:     "Fire Type Debuff! It's currently raining in Manchester. 0.8x multiplier.",
					Emoji:           "🌧️",
					SkipStatsOutput: true,
				})
			}
			if icy {
				wow.Multiplier *= 0.8
				ret = append(ret, &Effect{
					Name:            "Pokémon Fire Debuff",
					Description:     "Fire Type Debuff! It's currently Icy in Manchester. 0.8x multiplier.",
					Emoji:           "❄️",
					SkipStatsOutput: true,
				})
			}
		case "water":
			if hot {
				wow.Multiplier *= 0.8
				ret = append(ret, &Effect{
					Name: "Pokémon Water Debuff",
					Description: fmt.Sprintf("Water Type Debuff! It's currently  %.1f%s in Manchester. 1.2x multiplier.",
						dataCurrentWeather.Current.Temperature2m, dataCurrentWeather.CurrentUnits.Temperature2m),
					Emoji:           "🔥",
					SkipStatsOutput: true,
				})
			}
		case "grass":
			if hot && rain {
				wow.Multiplier *= 1.6
				ret = append(ret, &Effect{
					Name:            "Pokémon Grass SUPER Boost",
					Description:     "Grass Type SUPER Debuff! It's currently hot AND raining in Manchester! 1.6x multiplier.",
					Emoji:           "🌦️",
					SkipStatsOutput: true,
				})
			} else if rain {
				wow.Multiplier *= 1.2
				ret = append(ret, &Effect{
					Name:            "Pokémon Grass Boost",
					Description:     "Grass Type Boost! It's currently raining in Manchester. 1.2x multiplier.",
					Emoji:           "🌧️",
					SkipStatsOutput: true,
				})
			}
		case "electric":
			if rain {
				wow.Multiplier *= 1.2
				ret = append(ret, &Effect{
					Name:            "Pokémon Electric Boost",
					Description:     "Electric Type Boost! It's currently raining in Manchester. 1.2x multiplier.",
					Emoji:           "🌧️",
					SkipStatsOutput: true,
				})
			}
		case "psychic":
		case "ice":
			if hot {
				wow.Multiplier *= 0.8
				ret = append(ret, &Effect{
					Name: "Pokémon Ice Debuff",
					Description: fmt.Sprintf("Ice Type Debuff! It's currently  %.1f%s in Manchester. 0.8x multiplier.",
						dataCurrentWeather.Current.Temperature2m, dataCurrentWeather.CurrentUnits.Temperature2m),
					Emoji:           "🔥",
					SkipStatsOutput: true,
				})
			}
		case "dragon":
		case "dark":
			if night {
				wow.Multiplier *= 1.2
				ret = append(ret, &Effect{
					Name:            "Pokémon Dark Boost",
					Description:     "Dark Type Boost! It's currently night time in Manchester. 1.2x multiplier.",
					Emoji:           "🌑",
					SkipStatsOutput: true,
				})
			}
		case "fairy":
			if !night {
				wow.Multiplier *= 1.2
				ret = append(ret, &Effect{
					Name:            "Pokémon Fairy Boost",
					Description:     "Fairy Type Boost! It's currently day time in Manchester. 1.2x multiplier.",
					Emoji:           "🌞",
					SkipStatsOutput: true,
				})
			}
		case "stellar":
		}
	}

	// Furry Check ==========================================================
	furryBonus := []int{
		6,    // Charizard
		31,   // Nidoqueen
		38,   // Ninetales
		108,  // Lickitung
		131,  // Lapras
		136,  // Flareon
		149,  // Dragonite
		157,  // Typhlosion
		241,  // Miltank
		249,  // Lugia
		282,  // Gardevoir
		335,  // Zangoose
		359,  // Absol
		405,  // Luxray
		407,  // Roserade
		418,  // Buizel
		419,  // Floatzel
		428,  // Lopbunny
		431,  // Glameow
		445,  // Garchomp
		448,  // Lucario
		471,  // Glaceon
		509,  // Purrlion
		510,  // Liepard
		531,  // Audino
		571,  // Zoroark
		573,  // Cinccino
		620,  // Mienshao
		621,  // Druddigon
		654,  // Braixen
		655,  // Delphox
		700,  // Sylveon
		706,  // Goodra
		727,  // Incineroar
		745,  // Lycanroc
		758,  // Salazzle
		763,  // Tsareena
		815,  // Cinderace
		818,  // Inteleon
		827,  // Nickit
		836,  // Boltund
		858,  // Hatterene
		862,  // Obstagoon
		876,  // Indeedee
		888,  // Zacian
		889,  // Zamazenta
		908,  // Meowscarada
		1007, // Koraidon
	}

	if slices.Contains(furryBonus, pokemon.ID) {
		wow.BonusRolls += 3
		ret = append(ret, &Effect{
			Name:            "Pokémon Furry Alert",
			Description:     "OwO? Uh oh, let's not think too much into that Pokémon. Here's 3 bonus rolls to distract you.",
			Emoji:           "👀",
			SkipStatsOutput: true,
		})
	}

	// Pokemon Specific =====================================================
	switch pokemon.ID {
	case 70, 71, 779: // Weepinbell, Victreebel, Bruxish
		wow.MinContinue++
		ret = append(ret, &Effect{
			Name:            "Pokémon Mouth Warning",
			Description:     "GET THAT OUT OF IT'S MOUTH! +1 to your min continue roll",
			Emoji:           "💋",
			SkipStatsOutput: true,
		})
	case 119: // Seaking
		wow.Multiplier *= 2.0
		ret = append(ret, &Effect{
			Name:            "Pokémon F.Yeah Seaking",
			Description:     "FUCK YEAH SEAKING! 2.0x MULTIPLIER",
			Emoji:           "🐟",
			SkipStatsOutput: true,
		})
	case 124: // Jynx
		wow.MinContinue++
		ret = append(ret, &Effect{
			Name:            "Pokémon Jynx Vibes",
			Description:     "Yeah no, Jynx ain't giving those good vibes. +1 to your min continue roll to get you out of here ASAP.",
			Emoji:           "😩",
			SkipStatsOutput: true,
		})
	case 134: // Vaporeon
		wow.Multiplier *= 2.0
		ret = append(ret, &Effect{
			Name:            "Pokémon Vapoeron Copypasta",
			Description:     "Did you know in terms of Pokémon Wow rolls Vapoeron is the most compatiable Pokémon with a 2.0x multiplier?",
			Emoji:           "👀",
			SkipStatsOutput: true,
		})
	case 143, 321: // Snorlax
		wow.MinContinue += 3
		ret = append(ret, &Effect{
			Name:            "Pokémon Block",
			Description:     "Damn, it's blocking your rolls, +3 to your min continue roll",
			Emoji:           "⛔",
			SkipStatsOutput: true,
		})
	case 258: // Mudkip
		wow.BonusRolls++
		ret = append(ret, &Effect{
			Name:            "Pokémon Mudkip Meme",
			Description:     "Mud.. Kip.. Old meme, get +1 bonus roll",
			Emoji:           "💧",
			SkipStatsOutput: true,
		})
	case 339: // Bidoof
		wow.MinContinue++
		ret = append(ret, &Effect{
			Name:            "Pokémon Bidoof oof",
			Description:     "BidOOF. +1 to your min continue roll",
			Emoji:           "🤓",
			SkipStatsOutput: true,
		})
	case 495: // Snivy
		wow.BonusRolls += 2
		ret = append(ret, &Effect{
			Name:            "Pokémon Snivy Smug",
			Description:     "Don't act TOO Smug about your +2 bonus rolls",
			Emoji:           "🤓",
			SkipStatsOutput: true,
		})
	case 568, 569: // Trubbish, Garbodor
		wow.Multiplier *= 0.5
		ret = append(ret, &Effect{
			Name:            "Pokémon Garbage",
			Description:     "These Pokémon are TRASH. So is your Wow, 0.5x multiplier",
			Emoji:           "🤓",
			SkipStatsOutput: true,
		})
	case 591: // Amoonguss
		wow.BonusRolls += 3
		ret = append(ret, &Effect{
			Name:            "Pokémon Sus",
			Description:     "There's something Sus about this Pokémon. Have 3 bonus rolls.",
			Emoji:           "📪",
			SkipStatsOutput: true,
		})
	case 677: // Espurr
		wow.MinContinue -= 2
		ret = append(ret, &Effect{
			Name:            "Pokémon Focus",
			Description:     "Those staring eyes are focusing on the best Wows, -2 to your min continue roll.",
			Emoji:           "🥺",
			SkipStatsOutput: true,
		})
	case 734, 735: // Yungoos,Gumshoos
		wow.BonusRolls += 4
		ret = append(ret, &Effect{
			Name:            "Pokémon Orange",
			Description:     "This is a great Pokémon, truly, one of the best. Let me tell you, I know Pokémon, I know bonus rolls, have 4.",
			Emoji:           "🍊",
			SkipStatsOutput: true,
		})
	case 750: // Mudsdale
		wow.MinContinue--
		ret = append(ret, &Effect{
			Name:            "Pokémon White Woman",
			Description:     "Keep it away from the White Woman! Quick, -1 min continue roll to get you out of here",
			Emoji:           "🐴",
			SkipStatsOutput: true,
		})
	case 810, 811, 812: // Grookey, Thwackey, Rillaboom
		wow.BonusRolls += 3
		ret = append(ret, &Effect{
			Name:            "Pokémon Mmm Monke",
			Description:     "Mmm.. Monke. +3 bonus rolls.",
			Emoji:           "🐒",
			SkipStatsOutput: true,
		})
	case 872: // Snom
		wow.BonusRolls++
		ret = append(ret, &Effect{
			Name:            "Pokémon :3",
			Description:     "Look at it's face :3 +1 bonus roll",
			Emoji:           "🐱",
			SkipStatsOutput: true,
		})
	case 906, 907: // Sprigatito, Floragato
		wow.MinContinue -= 2
		ret = append(ret, &Effect{
			Name:            "Pokémon Jamacian",
			Description:     "Damn, this Pokémon seems to have some kind of Jamacian vibe to it. -2 to your min continue roll",
			Emoji:           "☘️",
			SkipStatsOutput: true,
		})
	case 959: // Tinkaton
		wow.MinContinue--
		ret = append(ret, &Effect{
			Name:            "Pokémon Smash",
			Description:     "Watch it SMASH your min continue roll down by 1!",
			Emoji:           "🔨",
			SkipStatsOutput: true,
		})
	}

	return ret
}

func getPokemonName(name string) string {
	name = helpers.CapitaliseWords(name)
	if !strings.Contains(name, "-") {
		return name
	}
	parts := strings.SplitN(name, "-", 2)
	return strings.TrimSpace(parts[0])
}
