package wow

import (
	"fmt"
	"math"
	"strings"
	"time"

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
	HardLimitHit  bool
}

type DiceRoll struct {
	Number         int
	Roll           int
	Effects        []Effect
	AdditionalText string
}

var (
	hardLimitMultiplier = 1000.0
	hardLimitBonusRolls = 100
	hardLimitRolls      = 10000
	hardLimitMaxRoll    = 50
)

func checkHardLimitHit(g *Generation) bool {
	if g.Multiplier > hardLimitMultiplier {
		g.Multiplier = hardLimitMultiplier
	}
	if g.BonusRolls > hardLimitBonusRolls {
		g.BonusRolls = hardLimitBonusRolls
	}
	if g.MaxRollValue > hardLimitMaxRoll {
		g.MaxRollValue = hardLimitMaxRoll
	}
	if g.MinContinue < g.MaxRollValue/4 {
		g.MinContinue = g.MaxRollValue / 4
	}
	if g.Rolls >= hardLimitRolls {
		g.HardLimitHit = true
		return true
	} else {
		return false
	}
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

	// Get user's inventory
	var userInv []InventoryItem
	var userInvStatic []InventoryItem
	var userInvRoll []InventoryItem
	wowTime := time.Now()

	logger.Dev("GENERATE", "Locking Inventory cache")
	dataInventoryLock.RLock()
	if inv, exists := dataUserInventories[fmt.Sprintf("%s|%s", message.GuildID, message.Author.ID)]; exists {
		userInv = inv
	}
	dataInventoryLock.RUnlock()
	logger.Dev("GENERATE", "Unlocking Inventory cache")

	for _, inv := range userInv {
		switch inv.ShopItem.TypeID {
		case shopItemTypeStatic:
			userInvStatic = append(userInvStatic, inv)
		case shopItemTypeRoll:
			userInvRoll = append(userInvRoll, inv)
		}
	}
	logger.Dev("GENERATE", "Sorted inventory into Static/Roll slices")

	// Static Effects ----------------------------------------------------------
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
		if checkHardLimitHit(&wow) {
			break
		}
	}
	logger.Dev("GENERATE", "Applied static effects")

	// Apply static shop effects
	for _, inv := range userInvStatic {
		// Has it expired?
		if !inv.ShopItem.OneTimeUse && wowTime.After(inv.Expires) {
			deleteFromWowInventory(message.GuildID, message.Author.ID, inv.DatabaseID)
			continue
		}
		// Apply effect
		effects := inv.ShopItem.Apply(&wow)
		if len(effects) > 0 {
			for _, effect := range effects {
				if !effect.SkipStatsOutput {
					wow.EffectCount++
				}
				wow.StaticEffects = append(wow.StaticEffects, *effect)
				wow.Effects = append(wow.Effects, *effect)
				if effect.SelfDestruct {
					deleteFromWowInventory(message.GuildID, message.Author.ID, inv.DatabaseID)
				}
			}
		}
		// Was it one time use?
		if inv.ShopItem.OneTimeUse {
			deleteFromWowInventory(message.GuildID, message.Author.ID, inv.DatabaseID)
		}
		if checkHardLimitHit(&wow) {
			break
		}
	}
	logger.Dev("GENERATE", "Applied static store effects")

	// Dice Rolls ---------------------------------------------------------------
	// Apply Dice Rolls
	fullDebuffMsgGiven := false
	for {
		// 1. Roll
		wow.Rolls++
		wow.CurrentRoll = helpers.GetRandomNumber(1, wow.MaxRollValue)
		logger.Dev("GENERATE", "Beginning roll %d - Rolled %d of max %d", wow.Rolls, wow.CurrentRoll, wow.MaxRollValue)

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
		logger.Dev("GENERATE", "Roll %d - Applied roll effects", wow.Rolls)

		// 3. Loop through any Roll based shop items
		for _, inv := range userInvRoll {
			// Has it expired?
			if !inv.ShopItem.OneTimeUse && wowTime.After(inv.Expires) {
				deleteFromWowInventory(message.GuildID, message.Author.ID, inv.DatabaseID)
				continue
			}
			// Apply effect
			effects := inv.ShopItem.Apply(&wow)
			if len(effects) > 0 {
				for _, effect := range effects {
					if !effect.SkipStatsOutput {
						wow.EffectCount++
					}
					rollEffects = append(rollEffects, *effect)
					wow.Effects = append(wow.Effects, *effect)
					if effect.SelfDestruct {
						deleteFromWowInventory(message.GuildID, message.Author.ID, inv.DatabaseID)
					}
				}
			}
			// Was it one time use?
			if inv.ShopItem.OneTimeUse {
				deleteFromWowInventory(message.GuildID, message.Author.ID, inv.DatabaseID)
			}
			if checkHardLimitHit(&wow) {
				break
			}
		}
		logger.Dev("GENERATE", "Roll %d - Applied store roll effects", wow.Rolls)

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
		logger.Dev("GENERATE", "Roll %d - Calculated deductions: %d", wow.Rolls, reduction)
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
		logger.Dev("GENERATE", "Roll %d - Additional text: %s", wow.Rolls, additionalText)

		// 4. Add to the Roll cache
		wow.OCount += wow.CurrentRoll
		wow.DiceRolls = append(wow.DiceRolls, DiceRoll{
			Number:         wow.Rolls,
			Roll:           wow.CurrentRoll,
			Effects:        rollEffects,
			AdditionalText: additionalText,
		})
		logger.Dev("GENERATE", "Roll %d - Added to Roll Cache", wow.Rolls)

		if checkHardLimitHit(&wow) {
			break
		}
		if finished {
			break
		}
	}

	logger.Dev("GENERATE", "Final Wow Count: %d", wow.OCount)
	if wow.OCount < 1 {
		wow.OCount = 1
	}

	logger.Dev("GENERATE", "Final Wow Multiplier: %f", wow.Multiplier)
	if wow.Multiplier > 1.0 {
		wow.OCount = int(math.Ceil(float64(wow.OCount) * wow.Multiplier))
	}

	subText := " " + intAsSubscript(wow.OCount) + "." + intAsSubscript(wow.EffectCount)
	wowText := fmt.Sprintf("W%sw%s", getOs(wow.OCount), subText)
	if wow.OCount > 75 {
		wowText = strings.ToUpper(wowText)
	}
	wow.Output = wowText
	logger.Dev("GENERATE", "Pushing Wow")
	pushWow(wow)
}

func pushWow(wow Generation) {
	logger.Info(wow.Message.GuildID, "Wow response queued, MessageID: %s", wow.Message.ID)
	queueRespond <- &wow
	queueDatabase <- &wow
}
