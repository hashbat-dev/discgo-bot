package wow

import (
	"fmt"
	"math"
	"strings"

	"github.com/bwmarrin/discordgo"
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
	Output        string
	WowMessageID  string
}

type DiceRoll struct {
	Number         int
	Roll           int
	Effects        []Effect
	AdditionalText string
}

func generate(message *discordgo.MessageCreate) {
	if !dataInit {
		GetEffectData()
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
				wow.StaticEffects = append(wow.StaticEffects, *effect)
				wow.Effects = append(wow.Effects, *effect)
			}
		}
	}

	// Apply Dice Rolls
	for {
		// 1. Roll
		wow.Rolls++
		wow.CurrentRoll = getRandomNumber(1, wow.MaxRollValue)

		// 2. Loop through any Roll based effects
		var rollEffects []Effect
		for _, fn := range rollEffectList {
			effects := fn(&wow)
			if len(effects) > 0 {
				for _, effect := range effects {
					rollEffects = append(rollEffects, *effect)
					wow.Effects = append(wow.Effects, *effect)
				}
			}
		}

		// 3. See whether to break off from future Rolls
		finished := false
		addText := ""

		if wow.CurrentRoll < wow.MinContinue {
			if wow.BonusRolls > 0 {
				wow.BonusRolls--
				addText = fmt.Sprintf("bonus roll used to avoid death, %d left", wow.BonusRolls)
			} else {
				finished = true
			}
		}

		// 3. Add to the Roll cache
		wow.OCount += wow.CurrentRoll
		wow.DiceRolls = append(wow.DiceRolls, DiceRoll{
			Number:         wow.Rolls,
			Roll:           wow.CurrentRoll,
			Effects:        rollEffects,
			AdditionalText: addText,
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

	effectCountText := ""
	if len(wow.Effects) > 0 {
		effectCountText = " " + intAsSubscript(len(wow.Effects))
	}

	wowText := fmt.Sprintf("W%sw%s", getOs(wow.OCount), effectCountText)
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
