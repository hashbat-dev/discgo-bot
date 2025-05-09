package wow

import "fmt"

var (
	IndentPadding = "\u00A0\u00A0\u00A0\u00A0\u00A0\u00A0"
	DefaultEmoji  = "⭐"
)

func GetStatsText(messageId string) (int, string) {
	if wow, ok := cacheBank[messageId]; ok {
		s := "To createa Wow we keep rolling a 1 to 10 dice until you roll 5 or under. The total number rolled is how long your Wow is!"
		s += " There are various modifiers and effects which you may be lucky enough to trigger!"

		if len(wow.Generation.StaticEffects) > 0 {
			s += "\n"
			for _, effect := range wow.Generation.StaticEffects {
				emoji := effect.Emoji
				if emoji == "" {
					emoji = DefaultEmoji
				}
				s += fmt.Sprintf("\n%s **%s**: %s", emoji, effect.Name, effect.Description)
			}
		}

		s += "\n\n**Your Rolls...**"
		for i, roll := range wow.Generation.DiceRolls {
			emoji := "🎲"
			if i == len(wow.Generation.DiceRolls)-1 {
				emoji = "💀"
			}
			s += fmt.Sprintf("\n%s **%d**", emoji, roll.Roll)
			if len(roll.Effects) > 0 {
				for _, effect := range roll.Effects {
					emoji := effect.Emoji
					if emoji == "" {
						emoji = DefaultEmoji
					}
					s += fmt.Sprintf("\n%s%s **%s**: %s", IndentPadding, emoji, effect.Name, effect.Description)
				}
			}
		}
		return wow.Generation.OCount, s
	} else {
		return 0, ""
	}
}
