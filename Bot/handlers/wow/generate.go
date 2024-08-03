package wow

import (
	"fmt"
	"strings"
	"time"

	"github.com/ZestHusky/femboy-control/Bot/friday"
	"github.com/ZestHusky/femboy-control/Bot/helpers"
	"github.com/bwmarrin/discordgo"
)

type GeneratedWow struct {
	WowText     string
	MiddleCount int
	Rolls       []DiceRoll
	LogText     string
}

func GenerateWow(message *discordgo.MessageCreate) GeneratedWow {

	maxRollWow := 10
	maxRollCont := 10
	maxRolls := 50
	maxSize := 1000
	maxContinue := 5
	middleCount := 1

	wowMiddle := "o"
	wowEnd := ""
	makeUppercase := false
	minZero := false
	bonusRolls := 0
	tempBonusRolls := 0
	deductRolls := 0

	var diceRolls []DiceRoll

	// Get any UserID Effects =============================================================================
	currentEffects := GetActiveEffectsForUser(message.Author.ID)
	for _, e := range currentEffects {

		diceRolls = append(diceRolls, DiceRoll{
			RollContinue:  0,
			RollLength:    0,
			Special:       "effect",
			SpecialString: e.EffectDescription,
		})

		maxContinue += e.ContinueModifier
		maxRollCont += e.ContinueRollModifier
		maxRollWow += e.WowRollModifier

		if e.FreeRolls != 0 {
			if e.FreeRolls < 0 {
				deductRolls += e.FreeRolls * -1
			} else {
				bonusRolls += e.FreeRolls
			}

			if e.FreeRolls == 0 {
				if e.ContinueModifier == 0 && e.ContinueRollModifier == 0 && e.WowRollModifier == 0 {
					DeleteActiveEffect(e)
				}
			}
		}

		if e.TempFreeRolls != 0 {
			if e.TempFreeRolls < 0 {
				deductRolls += e.TempFreeRolls * -1
			} else {
				tempBonusRolls += e.TempFreeRolls
			}

			if e.ContinueModifier == 0 && e.ContinueRollModifier == 0 && e.WowRollModifier == 0 {
				DeleteActiveEffect(e)
			} else {
				e.TempFreeRolls = 0
				UpdateActiveEffect(e)
			}
		}
	}

	// Start ROLLING ROLLING ROLLING ======================================================================
	for i := 0; i < maxRolls; i++ {

		// Start the Dice Roll ============================================================================
		var diceRoll DiceRoll
		currTime := time.Now()

		// Calculations ONLY on the First Roll ---------------------------------------
		if i == 0 {

			// Doubles or Triples?
			last1 := message.ID[len(message.ID)-1]
			last2 := message.ID[len(message.ID)-2]
			last3 := message.ID[len(message.ID)-3]
			if last1 == last2 && last1 == last3 {
				diceRoll.Special = "trips"
				trip1 := helpers.GetRandomNumber(0, maxRollCont)
				middleCount += trip1
				wowMiddle += AddCharacter("o", trip1)

				trip2 := helpers.GetRandomNumber(0, maxRollCont)
				middleCount += trip2
				wowMiddle += AddCharacter("o", trip2)

				trip3 := helpers.GetRandomNumber(0, maxRollCont)
				middleCount += trip3
				wowMiddle += AddCharacter("o", trip3)

				diceRoll.RollSpecial = []int{trip1, trip2, trip3}
			} else if last1 == last2 {
				diceRoll.Special = "dubs"
				dub1 := helpers.GetRandomNumber(0, maxRollCont)
				middleCount += dub1
				wowMiddle += AddCharacter("o", dub1)

				dub2 := helpers.GetRandomNumber(0, maxRollCont)
				middleCount += dub2
				wowMiddle += AddCharacter("o", dub2)

				diceRoll.RollSpecial = []int{dub1, dub2}
			}
		}

		// Generate a Random Special event? -----------------------------------
		if diceRoll.Special == "" {
			addedSpecial := false

			// Higher Max?
			if maxContinue <= 8 {
				lowerMax := helpers.GetRandomNumber(1, 100)
				if lowerMax == 50 {
					diceRoll.Special = "highermax"
					maxContinue += 1
					addedSpecial = true
				}
			}

			// Pro Gambler?
			if !addedSpecial {
				proGamble := helpers.GetRandomNumber(1, 100)
				if proGamble == 100 {
					diceRoll.Special = "progamble"
					maxRollWow += 10
				}
			}

			// Cheeky Chap?
			if !addedSpecial {
				proGamble := helpers.GetRandomNumber(1, 50)
				if proGamble == 25 {
					diceRoll.Special = "cheeky"
					bonusRolls++
				}
			}

			// Fibonacci?
			if !addedSpecial {

				doFib := helpers.GetRandomNumber(1, 34)
				if doFib == 21 {

					fibVal := FibonacciValue()
					middleCount += fibVal
					wowMiddle += AddCharacter("o", fibVal)

					diceRoll.Special = "fibonacci"
					diceRoll.RollSpecial = append(diceRoll.RollSpecial, fibVal)
				}
			}

			//oneArmBandit
			if !addedSpecial {
				doOneArmBandit := helpers.GetRandomNumber(1, 20)
				if doOneArmBandit == 1 {
					oneArmBanditVal, oneArmBanditResult := NewOneArmBandit()
					middleCount += oneArmBanditVal
					wowMiddle += AddCharacter("o", oneArmBanditVal)

					diceRoll.Special = "oneArmBandit"
					diceRoll.RollSpecial = append(diceRoll.RollSpecial, oneArmBanditVal)
					diceRoll.SpecialString = "\n**Your Spin**: " + oneArmBanditResult
					if oneArmBanditVal == 0 {
						diceRoll.SpecialString += " (**You Lost**!)"
					} else {
						diceRoll.SpecialString += " (**You Won** " + fmt.Sprint(oneArmBanditVal) + "!)"
					}
				}
			}
		}

		// [CONTINUE] Roll -----------------------------------------------------
		roll := helpers.GetRandomNumber(0, maxRollCont)

		// Crit Success?
		if roll == 0 {
			freeRollAdd := helpers.GetRandomNumber(0, maxRollWow)
			middleCount += freeRollAdd
			wowMiddle += AddCharacter("o", freeRollAdd)
			diceRoll.RollSpecial = append(diceRoll.RollSpecial, freeRollAdd)
		}

		// [WOOOW] Roll --------------------------------------------------------
		addWows := helpers.GetRandomNumber(0, maxRollWow)
		middleCount += addWows
		wowMiddle += AddCharacter("o", addWows)

		// Do we have any Bonus Rolls to also add?
		if bonusRolls > 0 {
			for r := 0; r < bonusRolls; r++ {
				bonusAdd := helpers.GetRandomNumber(0, maxRollWow)
				middleCount += bonusAdd
				wowMiddle += AddCharacter("o", bonusAdd)
				diceRoll.RollSpecial = append(diceRoll.RollSpecial, bonusAdd)
			}
		}

		// Any Temp rolls?
		if tempBonusRolls > 0 {
			for r := 0; r < tempBonusRolls; r++ {
				bonusAdd := helpers.GetRandomNumber(0, maxRollWow)
				middleCount += bonusAdd
				wowMiddle += AddCharacter("o", bonusAdd)
				diceRoll.RollSpecial = append(diceRoll.RollSpecial, bonusAdd)
			}
			tempBonusRolls = 0
		}

		diceRoll.RollContinue = roll
		diceRoll.RollLength = addWows

		// Update the End character of the Wooow string? ---------------------
		if i == 1 {
			wowEnd = "."
		} else if i == 2 {
			wowEnd = "!"
		} else if i == 3 {
			makeUppercase = true
		} else if i == 4 {
			wowEnd = "!!!"
		} else {
			wowEnd += " " + helpers.GetEmote("ae_cry", false)
		}

		// Is it time to End this Tomfoolery? --------------------------------
		if roll > maxContinue {

			// Do they get the Minute Zero Continue buff?
			if !minZero {
				minZero = true
				if currTime.Second() == 0 {
					diceRoll.Special = "minzero"
					diceRolls = append(diceRolls, diceRoll)
					continue
				}
			}

			// Is it Fwiday?
			if friday.IsItFwiday() {
				addFwiday := helpers.GetRandomNumber(0, maxRollWow)
				wowMiddle += AddCharacter("o", addFwiday)
				middleCount += addWows
				diceRoll.RollLength += addFwiday
				diceRoll.Special = "fwiday"
				diceRoll.RollSpecial = append(diceRoll.RollSpecial, addFwiday)
			}

			// Any Deductions?
			if deductRolls > 0 {
				for d := 0; d < deductRolls; d++ {
					deductInt := helpers.GetRandomNumber(0, maxRollWow)
					wowMiddle = DelCharacters(wowMiddle, deductInt)
					deductInt = deductInt * -1
					middleCount += deductInt
					diceRoll.RollSpecial = append(diceRoll.RollSpecial, deductInt)
				}
			}

			diceRolls = append(diceRolls, diceRoll)
			break // YEET - Exit Rolling (Anti-Levi Moment) ==================
		}

		if len(wowMiddle) > maxSize {
			diceRolls = append(diceRolls, diceRoll)
			break
		}

		diceRolls = append(diceRolls, diceRoll)

	}

	// Log the Rolls ----------------------------------------------------------
	rollText := ""
	for _, r := range diceRolls {
		if r.Special == "effect" && r.RollContinue == 0 && r.RollLength == 0 {
			if strings.Contains(r.SpecialString, "**") {
				splitEffect := strings.Split(r.SpecialString, "**")
				if splitEffect[1] != "" {
					rollText += "[" + splitEffect[1] + "]"
				} else {
					rollText += "[" + r.SpecialString + "]"
				}
			} else {
				rollText += "[" + r.SpecialString + "]"
			}
		} else {
			rollText += "["
			rollText += ">" + fmt.Sprint(r.RollContinue)
			rollText += "/+" + fmt.Sprint(r.RollLength)
			if r.Special != "" {
				rollText += "(" + r.Special
				if len(r.RollSpecial) > 0 {
					for _, x := range r.RollSpecial {
						rollText += "[" + fmt.Sprint(x) + "]"
					}
				}

				rollText += ")"
			}
			rollText += "]"
		}

	}

	retText := "W" + wowMiddle + "w" + wowEnd
	if makeUppercase {
		retText = strings.ToUpper(retText)
	}

	return GeneratedWow{
		WowText:     retText,
		MiddleCount: middleCount,
		Rolls:       diceRolls,
		LogText:     rollText,
	}
}
