package wow

import (
	"fmt"

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
		s := "To create a Wow we keep rolling a 1 to 10 dice until you roll 5 or under. The total number rolled is how long your Wow is!"
		s += " There are various modifiers and effects which you may be lucky enough to trigger!"

		if len(wow.Generation.StaticEffects) > 0 {
			s += "\n"
			for _, effect := range wow.Generation.StaticEffects {
				if effect.SkipStatsOutput {
					continue
				}
				emoji := effect.Emoji
				if emoji == "" {
					emoji = DefaultEmoji
				}
				newLine := fmt.Sprintf("%s **%s**: %s", emoji, effect.Name, effect.Description)
				if (len(s) + len(newLine) + 2) > maxMsgLength {
					messages = append(messages, s)
					s = newLine
				} else {
					s += "\n" + newLine
				}

			}
		}

		rollHeader := "**Your Rolls...**"
		headerSent := false
		for i, roll := range wow.Generation.DiceRolls {
			emoji := "🎲"
			if i == len(wow.Generation.DiceRolls)-1 {
				emoji = "💀"
			}
			addText := ""
			if roll.AdditionalText != "" {
				addText = "\u00A0\u00A0" + roll.AdditionalText
			}
			newLine := fmt.Sprintf("%s **%d**%s", emoji, roll.Roll, addText)
			if len(roll.Effects) > 0 {
				for _, effect := range roll.Effects {
					if effect.SkipStatsOutput {
						continue
					}
					emoji := effect.Emoji
					if emoji == "" {
						emoji = DefaultEmoji
					}
					newLine += fmt.Sprintf("\n%s%s **%s**: %s", IndentPadding, emoji, effect.Name, effect.Description)
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
