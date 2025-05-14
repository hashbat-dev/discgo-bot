package wow

import (
	"fmt"
	"math"
	"strings"

	"github.com/bwmarrin/discordgo"
	discord "github.com/hashbat-dev/discgo-bot/Discord"
	helpers "github.com/hashbat-dev/discgo-bot/Helpers"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

type Generation struct {
	Message       *discordgo.MessageCreate
	OCount        int
	Rolls         int
	CurrentRoll   int
	MinContinue   int
	MaxRollValue  int
	BonusRolls    int
	Multiplier    float64
	DiceRolls     []DiceRoll
	StaticEffects []Effect
	Effects       []Effect
	EffectCount   int
	Output        string
	WowMessageIDs []string
}

type DiceRoll struct {
	Number         int
	Roll           int
	Effects        []Effect
	AdditionalText string
}

func generate(message *discordgo.MessageCreate) {
	if !dataInit {
		err := discord.ReplyToMessage(message, "Bot is starting, try again shortly!")
		if err != nil {
			logger.Error(message.GuildID, err)
		}
		return
	}

	wow := Generation{
		Message:      message,
		MinContinue:  6,
		MaxRollValue: 10,
		Multiplier:   1.0,
		Rolls:        0,
	}

	// Apply Static effects, some of these include free rolls
	for _, fn := range staticEffectList {
		effects := fn(&wow)
		if len(effects) > 0 {
			for _, effect := range effects {
				if !effect.SkipStatsOutput {
					wow.EffectCount++
				}
				wow.StaticEffects = append(wow.StaticEffects, *effect)
				wow.Effects = append(wow.Effects, *effect)
			}
		}
	}

	// Apply Dice Rolls
	fullDebuffMsgGiven := false
	for {
		// 1. Roll
		wow.Rolls++
		wow.CurrentRoll = helpers.GetRandomNumber(1, wow.MaxRollValue)

		// 2. Loop through any Roll based effects
		var rollEffects []Effect
		for _, fn := range rollEffectList {
			effects := fn(&wow)
			if len(effects) > 0 {
				for _, effect := range effects {
					if !effect.SkipStatsOutput {
						wow.EffectCount++
					}
					rollEffects = append(rollEffects, *effect)
					wow.Effects = append(wow.Effects, *effect)
				}
			}
		}

		// 3. See whether to break off from future Rolls
		finished := false
		var addText []string

		minContinue := wow.MinContinue
		reduction := wow.BonusRolls / 3
		if reduction > 0 {
			if !fullDebuffMsgGiven {
				addText = append(addText, fmt.Sprintf("continue roll debuffed due to %d bonus rolls", wow.BonusRolls))
				fullDebuffMsgGiven = true
			} else {
				addText = append(addText, fmt.Sprintf("debuff - %d bonus rolls", wow.BonusRolls))
			}
		}
		minContinue -= reduction
		if minContinue < 1 {
			minContinue = 1
		}

		if wow.CurrentRoll < minContinue {
			if wow.BonusRolls > 0 {
				wow.BonusRolls--
				addText = append(addText, fmt.Sprintf("bonus roll avoided death - %d left", wow.BonusRolls))
			} else {
				finished = true
			}
		}

		additionalText := ""
		if len(addText) > 0 {
			additionalText = strings.Join(addText, ", ")
		}

		// 4. Add to the Roll cache
		wow.OCount += wow.CurrentRoll
		wow.DiceRolls = append(wow.DiceRolls, DiceRoll{
			Number:         wow.Rolls,
			Roll:           wow.CurrentRoll,
			Effects:        rollEffects,
			AdditionalText: additionalText,
		})

		if finished {
			break
		}

	}

	if wow.OCount < 1 {
		wow.OCount = 1
	}

	if wow.Multiplier > 1.0 {
		wow.OCount = int(math.Ceil(float64(wow.OCount) * wow.Multiplier))
	}

	subText := " " + intAsSubscript(wow.OCount) + "." + intAsSubscript(wow.EffectCount)
	wowText := fmt.Sprintf("W%sw%s", getOs(wow.OCount), subText)
	if wow.OCount > 75 {
		wowText = strings.ToUpper(wowText)
	}
	wow.Output = wowText
	pushWow(wow)
}

func pushWow(wow Generation) {
	logger.Info(wow.Message.GuildID, "Wow response queued, MessageID: %s", wow.Message.ID)
	queueRespond <- &wow
	queueDatabase <- &wow
}
