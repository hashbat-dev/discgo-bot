package wow

import (
	"fmt"
	"slices"

	config "github.com/hashbat-dev/discgo-bot/Config"
)

var (
	IndentPadding = "\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0"
	DefaultEmoji  = "⭐"
)

func GetStatsText(messageId string) (int, []string) {
	var messages []string
	maxMsgLength := config.MAX_EMBED_DESC_LENGTH
	if wow, ok := cacheBank[messageId]; ok {
		s := "We roll a 1 to 10 dice, if you roll 6 or under we roll again. The total number rolled is how long your Wow is!"
		s += " Stats are changed by real world effects, shop items + more!"

		if len(wow.Generation.StaticEffects) > 0 {
			s += "\n"
			statShop := false

			staticEffectCounts := make(map[string]int)
			for _, effect := range wow.Generation.StaticEffects {
				staticEffectCounts[effect.Name]++
			}

			var staticNameCheck []string
			for _, effect := range wow.Generation.StaticEffects {
				if effect.SkipStatsOutput {
					continue
				}
				if slices.Contains(staticNameCheck, effect.Name) {
					continue
				}
				// COUNT INSTANCES IN THE wow.Generation.StaticEffects WHICH HAVE THE SAME NAME
				if effect.FromShop && !statShop {
					s += "\n\n**From your Purchased Items...**"
					statShop = true
				}
				emoji := effect.Emoji
				if emoji == "" {
					emoji = DefaultEmoji
				}
				countText := ""
				if staticEffectCounts[effect.Name] > 1 {
					countText = fmt.Sprintf(" x%d", staticEffectCounts[effect.Name])
				}
				newLine := fmt.Sprintf("%s **%s%s**: %s", emoji, effect.Name, countText, effect.Description)
				if (len(s) + len(newLine) + 2) > maxMsgLength {
					messages = append(messages, s)
					s = newLine
				} else {
					s += "\n" + newLine
				}

				staticNameCheck = append(staticNameCheck, effect.Name)
			}
		}

		rollHeader := "**Your Rolls...**"
		headerSent := false

		var rollNameCheck []string
		for i, roll := range wow.Generation.DiceRolls {
			emoji := "🎲"
			if i == len(wow.Generation.DiceRolls)-1 {
				emoji = "💀"
			}
			addText := ""
			if roll.AdditionalText != "" {
				addText = "\u00A0\u00A0" + roll.AdditionalText
			}
			rollEffectCounts := make(map[string]int)
			for _, effect := range roll.Effects {
				rollEffectCounts[effect.Name]++
			}
			newLine := fmt.Sprintf("%s **%d**%s", emoji, roll.Roll, addText)
			if len(roll.Effects) > 0 {
				for _, effect := range roll.Effects {
					if effect.SkipStatsOutput {
						continue
					}
					if slices.Contains(rollNameCheck, effect.Name) {
						continue
					}
					emoji := effect.Emoji
					if emoji == "" {
						emoji = DefaultEmoji
					}
					storeText := ""
					if effect.FromShop {
						storeText = " (Item)"
					}
					countText := ""
					if rollEffectCounts[effect.Name] > 1 {
						countText = fmt.Sprintf(" x%d", rollEffectCounts[effect.Name])
					}
					newLine += fmt.Sprintf("\n%s%s **%s%s**%s: %s", IndentPadding, emoji, effect.Name, countText, storeText, effect.Description)

					rollNameCheck = append(rollNameCheck, effect.Name)
				}
			}

			if headerSent {
				if (len(s) + len(newLine) + 2) > maxMsgLength {
					messages = append(messages, s)
					s = newLine
				} else {
					s += "\n" + newLine
				}
			} else {
				if (len(s) + len(newLine) + 6 + len(rollHeader)) > maxMsgLength {
					messages = append(messages, s)
					s = rollHeader + "\n" + newLine
				} else {
					s += "\n\n" + rollHeader + "\n" + newLine
				}
				headerSent = true
			}
		}

		if len(s) > 0 {
			messages = append(messages, s)
		}

		return wow.Generation.OCount, messages
	} else {
		return 0, nil
	}
}
