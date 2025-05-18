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
			Name:        "Funny Number",
			Description: "For the next hour if you roll a 6, add a 9 onto the roll.",
			Emoji:       "💞",
			Cost:        3000,
			Duration:    time.Hour * 1,
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
			Name:        "Close Shave",
			Description: "For the next hour if you roll your min roll continue value your roll value is doubled.",
			Emoji:       "✂️",
			Cost:        10000,
			Duration:    time.Hour * 1,
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
	}
)
