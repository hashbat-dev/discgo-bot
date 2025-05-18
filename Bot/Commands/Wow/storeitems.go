package wow

import (
	"fmt"
	"time"

	helpers "github.com/hashbat-dev/discgo-bot/Helpers"
)

const (
	shopItemTypeStatic = iota
	shopItemTypeRoll
)

type WowShopItem struct {
	ID          int //--> DO NOT CHANGE ONCE LIVE <-- These are mapped to the Database
	TypeID      int // shopItemType from the iota above
	Name        string
	Description string
	Emoji       string
	Cost        int
	Duration    time.Duration
	OneTimeUse  bool
	MaxAtOnce   int
	Apply       func(*Generation) []*Effect
}

var (
	ShopItems = []WowShopItem{
		{
			ID:          1,
			TypeID:      shopItemTypeStatic,
			Name:        "Fart Jar",
			Description: "Open a fart jar before your next Wow, get a starting boost of 5-500 O's depending on the length of the fart.",
			Emoji:       "🧪",
			Cost:        5000,
			OneTimeUse:  true,
			Apply: func(g *Generation) []*Effect {
				rand := helpers.GetRandomNumber(5, 500)
				g.OCount += rand
				return []*Effect{{
					Name:        "Fart Jar",
					Description: "Open a fart jar before your next Wow, get a starting boost of 5-500 O's depending on the length of the fart.",
					Emoji:       "🧪",
					FromShop:    true,
				}}
			},
		}, {
			ID:          2,
			TypeID:      shopItemTypeStatic,
			Name:        "Bojangles",
			Description: "Makes you gassy, for the next hour rip a fart at the start of your Wow to get a starting boost of 5-50 O's.",
			Emoji:       "🍗",
			Cost:        20000,
			Duration:    time.Hour * 1,
			Apply: func(g *Generation) []*Effect {
				rand := helpers.GetRandomNumber(5, 50)
				g.OCount += rand
				return []*Effect{{
					Name:        "Bojangles",
					Description: fmt.Sprintf("Makes you gassy, for the next hour rip a fart at the start of your Wow to get a starting boost of 5-50 O's. (You got %d)", rand),
					Emoji:       "🍗",
					FromShop:    true,
				}}
			},
		}, {
			ID:          3,
			TypeID:      shopItemTypeStatic,
			Name:        "Chippy Tea",
			Description: "Nothing beats it, for the next hour every Dice Roll has its maximum possible roll increased by +1 (+2 if it's Friday).",
			Emoji:       "🐟",
			Cost:        40000,
			Duration:    time.Hour * 1,
			MaxAtOnce:   3,
			Apply: func(g *Generation) []*Effect {
				var countText string
				if time.Now().Weekday() == time.Friday {
					countText = "+2"
					g.MaxRollValue += 2
					g.MinContinue += 2
				} else {
					countText = "+1"
					g.MaxRollValue++
					g.MinContinue++
				}
				return []*Effect{{
					Name:        "Chippy Tea",
					Description: fmt.Sprintf("Nothing beats it, for the next hour every Dice Roll's maximum roll is %s", countText),
					Emoji:       "🐟",
					FromShop:    true,
				}}
			},
		},
		{
			ID:          4,
			TypeID:      shopItemTypeStatic,
			Name:        "Tobacco Mystery Pack",
			Description: "For an hour you're thrown a random Tobacco pouch at the start of your next Wow, Drum: +1 to your roll min continue, Gold Leaf: -1",
			Emoji:       "🚬",
			Cost:        15000,
			Duration:    time.Hour * 1,
			MaxAtOnce:   3,
			Apply: func(g *Generation) []*Effect {
				isDrum := helpers.GetRandomNumber(0, 1) == 1
				var desc string
				if isDrum {
					desc = "Bad luck, you got given Drum. +1 to your roll min continue."
					g.MinContinue++
				} else {
					desc = "You got Gold Leaf! -1 to your roll min continue."
					g.MinContinue--
				}
				return []*Effect{{
					Name:        "Tobacco Mystery Pack",
					Description: desc,
					Emoji:       "🚬",
					FromShop:    true,
				}}
			},
		},
		{
			ID:          5,
			TypeID:      shopItemTypeStatic,
			Name:        "Cigarette",
			Description: "You roll a cigarette, get +1 roll of the bonus variety on your next Wow.",
			Emoji:       "🚬",
			Cost:        1000,
			OneTimeUse:  true,
			Apply: func(g *Generation) []*Effect {
				g.BonusRolls++
				return []*Effect{{
					Name:        "Cigarette",
					Description: "You roll a cigarette, get +1 roll of the bonus variety on your next Wow.",
					Emoji:       "🚬",
					FromShop:    true,
				}}
			},
		},
		{
			ID:          6,
			TypeID:      shopItemTypeRoll,
			Name:        "Funny Number Bolt-on",
			Description: "For the next hour if you roll a 6, add a 9 onto the roll.",
			Emoji:       "💞",
			Cost:        3000,
			Duration:    time.Hour * 1,
			MaxAtOnce:   3,
			Apply: func(g *Generation) []*Effect {
				if g.CurrentRoll == 6 {
					g.CurrentRoll += 9
					return []*Effect{{
						Name:        "Funny Number",
						Description: "For the next hour if you roll a 6, add a 9 onto the roll.",
						Emoji:       "💞",
						FromShop:    true,
					}}
				}
				return nil
			},
		},
		{
			ID:          7,
			TypeID:      shopItemTypeRoll,
			Name:        "Close Shave Scissors",
			Description: "For the next hour if you roll your min roll continue value your roll value is doubled.",
			Emoji:       "✂️",
			Cost:        10000,
			Duration:    time.Hour * 1,
			MaxAtOnce:   3,
			Apply: func(g *Generation) []*Effect {
				if g.CurrentRoll == g.MinContinue {
					g.CurrentRoll *= 2
					return []*Effect{{
						Name:        "Close Shave",
						Description: "For the next hour if you roll your min roll continue value your roll value is doubled.",
						Emoji:       "✂️",
						FromShop:    true,
					}}
				}
				return nil
			},
		},
		{
			ID:          8,
			TypeID:      shopItemTypeStatic,
			Name:        "Lucky 115 Briefcase",
			Description: "For the next 15 hours, roll a dice on every Wow roll between 1-115. If you get 115 you gain a x15 multiplier. If you roll a 1 this item is deleted.",
			Emoji:       "💙",
			Cost:        50000,
			Duration:    time.Hour * 15,
			MaxAtOnce:   1,
			Apply: func(g *Generation) []*Effect {
				rand := helpers.GetRandomNumber(1, 115)
				switch rand {
				case 115:
					g.Multiplier *= 15
					return []*Effect{{
						Name:        "Lucky 115",
						Description: "You rolled 115! Get a x15 multiplier.",
						Emoji:       "💙",
						FromShop:    true,
					}}
				case 1:
					return []*Effect{{
						Name:         "Lucky 115",
						Description:  "You rolled 1, this item is now deleted.",
						Emoji:        "💙",
						FromShop:     true,
						SelfDestruct: true,
					}}
				default:
					return nil
				}
			},
		},
		{
			ID:          9,
			TypeID:      shopItemTypeRoll,
			Name:        "Double or Nothing Shot",
			Description: "The next time you hit a 10, gain a 5x multiplier. However if you hit a 1 before then gain a 0.5x multiplier. Item is deleted on use.",
			Emoji:       "🥃",
			Cost:        20000,
			Duration:    time.Hour * 200,
			MaxAtOnce:   1,
			Apply: func(g *Generation) []*Effect {
				switch g.CurrentRoll {
				case 10:
					g.Multiplier *= 5
					return []*Effect{{
						Name:         "Double or Nothing Shot",
						Description:  "You rolled a 10! Gain a 5x multiplier, this item is now deleted.",
						Emoji:        "🥃",
						FromShop:     true,
						SelfDestruct: true,
					}}
				case 1:
					g.Multiplier *= 0.5
					return []*Effect{{
						Name:         "Double or Nothing Shot",
						Description:  "You rolled 1, applying a 0.5x multiplier, this item is now deleted.",
						Emoji:        "🥃",
						FromShop:     true,
						SelfDestruct: true,
					}}
				default:
					return nil
				}
			},
		},
		{
			ID:          10,
			TypeID:      shopItemTypeRoll,
			Name:        "Loaded Dice",
			Description: "For the next hour, rolling a 7, 8 or 9 will bump it to a 10 and double it.",
			Emoji:       "🪛",
			Cost:        20000,
			Duration:    time.Hour * 1,
			MaxAtOnce:   1,
			Apply: func(g *Generation) []*Effect {
				switch g.CurrentRoll {
				case 7:
					g.CurrentRoll = 20
					return []*Effect{{
						Name:        "Loaded Dice",
						Description: "For the next hour, rolling a 7, 8 or 9 will bump it to a 10 and double it.",
						Emoji:       "🪛",
						FromShop:    true,
					}}
				case 8:
					g.CurrentRoll = 20
					return []*Effect{{
						Name:        "Loaded Dice",
						Description: "For the next hour, rolling a 7, 8 or 9 will bump it to a 10 and double it.",
						Emoji:       "🪛",
						FromShop:    true,
					}}
				case 9:
					g.CurrentRoll = 20
					return []*Effect{{
						Name:        "Loaded Dice",
						Description: "For the next hour, rolling a 7, 8 or 9 will bump it to a 10 and double it.",
						Emoji:       "🪛",
						FromShop:    true,
					}}
				default:
					return nil
				}
			},
		},
	}
)
