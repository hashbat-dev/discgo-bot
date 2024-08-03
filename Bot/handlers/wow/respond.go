package wow

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/ZestHusky/femboy-control/Bot/audit"
	"github.com/ZestHusky/femboy-control/Bot/constants"
	"github.com/ZestHusky/femboy-control/Bot/helpers"
	"github.com/bwmarrin/discordgo"
)

func WowDetection(message *discordgo.MessageCreate) {
	go DoDetection(message)
}

func DoDetection(message *discordgo.MessageCreate) {

	tooSoon := false
	found := false
	timeNow := time.Now()
	cacheItem := WowSpamCache{}
	for _, cache := range wowSpamCache {
		if cache.UserID == message.Author.ID {
			found = true

			// Expired Session?
			if time.Since(cache.LastUsed) >= 60*time.Minute {
				DeleteAllSessionEvents(cache.UserID)
				cache.SessionCount++
				cache.SessionStart = timeNow
			}

			if message.ChannelID != constants.CHANNEL_BOT_SPAM {
				// Not #bot-spam Spam check
				if time.Since(cache.LastUsed) <= 30*time.Second {
					tooSoon = true
				} else {
					cache.LastUsed = timeNow
					cache.SessionCount++
				}
			} else {
				cache.LastUsed = timeNow
				cache.SessionCount++
			}

			cacheItem = cache
			UpdateSpamCache(cache)
		}
	}

	if !found {
		// Not in the Cache, add them!
		newCache := WowSpamCache{
			UserID:       message.Author.ID,
			LastUsed:     timeNow,
			SessionCount: 0,
			SessionStart: timeNow,
		}
		wowSpamCache = append(wowSpamCache, newCache)
		cacheItem = newCache
	}

	if message.ChannelID == constants.CHANNEL_BOT_SPAM {
		// Apply any Cache based effects? ============================================

		if timeNow == cacheItem.SessionStart {

			// => New Session!
			InsertIntoEffectCache(message, WowEffect{
				UserID:            message.Author.ID,
				EffectName:        "sessionstart",
				EffectDescription: "**Hello Bonus** 😁: (1 turn) Haven't seen you for a while! Have 2 bonus rolls.",
				ActiveUntil:       timeNow.Add(1 * time.Minute),
				TempFreeRolls:     2,
			})
		} else {

			// => Number Based
			if cacheItem.SessionCount == 18 {
				InsertIntoEffectCache(message, WowEffect{
					UserID:            message.Author.ID,
					EffectName:        "cachenum18",
					EffectDescription: "**Now Legal!** ❤️: (1 turn) Have a bonus roll on the house.",
					TempFreeRolls:     1,
					SessionBased:      true,
				})
			} else if cacheItem.SessionCount == 50 {
				InsertIntoEffectCache(message, WowEffect{
					UserID:            message.Author.ID,
					EffectName:        "cachenum50",
					EffectDescription: "**Absolute Spammer** 😎: +1 bonus roll per turn.",
					FreeRolls:         1,
					SessionBased:      true,
				})
			} else if cacheItem.SessionCount == 100 {
				InsertIntoEffectCache(message, WowEffect{
					UserID:            message.Author.ID,
					EffectName:        "cachenum100",
					EffectDescription: "**That's Cap frfr** 💯: +1 to all Wooow rolls, ",
					WowRollModifier:   1,
					TempFreeRolls:     2,
					SessionBased:      true,
				})
			} else if cacheItem.SessionCount == 200 {
				InsertIntoEffectCache(message, WowEffect{
					UserID:            message.Author.ID,
					EffectName:        "cachenum200",
					EffectDescription: "**Bruh Chill** 😟: Hopefully +3 bonuses this turn helps you calm down",
					TempFreeRolls:     3,
					SessionBased:      true,
				})
			} else if cacheItem.SessionCount == 400 {
				InsertIntoEffectCache(message, WowEffect{
					UserID:            message.Author.ID,
					EffectName:        "cachenum400",
					EffectDescription: "**400 Jeremy? 400?** 🤓: -1 to all your Continue rolls, lucky you!",
					ContinueModifier:  -1,
					SessionBased:      true,
				})
			} else if cacheItem.SessionCount == 500 {
				InsertIntoEffectCache(message, WowEffect{
					UserID:            message.Author.ID,
					EffectName:        "cachenum500",
					EffectDescription: "**What are you doing?** 🧐: Fuck it, have 5 free rolls.. and +1 per turn!",
					FreeRolls:         1,
					TempFreeRolls:     5,
					SessionBased:      true,
				})
			}
		}
	}

	if tooSoon && message.ChannelID != constants.CHANNEL_BOT_SPAM {
		InsertIntoEffectCache(message, WowEffect{
			UserID:               message.Author.ID,
			EffectName:           "wowspam",
			EffectDescription:    "**Absolute Spamming Samuel** 📫: (60s) Your Continue roll numbers are reduced by 1.",
			ActiveUntil:          timeNow.Add(1 * time.Minute),
			ContinueRollModifier: -1,
		})
	}

	// 1. Get the Message
	msg := strings.ToLower(message.Content)
	if msg == "" {
		return
	}

	// 2. Did they "WOOOOW"?
	foundWow := false
	regCheck := regexp.MustCompile(`(?i)^w+o{1,}w+$`)
	if regCheck.MatchString(message.Content) {
		foundWow = true
	}

	// 3. If so, reply
	if foundWow {

		timeGenStart := time.Now()
		newWow := GenerateWow(message)
		timeGenerate := time.Since(timeGenStart)

		timeSpecialStart := time.Now()
		specialCount := 0
		for _, r := range newWow.Rolls {
			if r.Special != "" {
				specialCount++
			}
		}

		if specialCount > 0 {
			newWow.WowText += " " + helpers.GetSuperscriptNumber(specialCount, true)
		}
		timeSpecial := time.Since(timeSpecialStart)

		timeCacheStart := time.Now()
		SendCache = append(SendCache, WowSendCache{
			ChannelID:   message.ChannelID,
			ReplyText:   newWow.WowText,
			MessageRef:  *message.Reference(),
			WowRolls:    newWow.Rolls,
			WowCount:    newWow.MiddleCount,
			UserID:      message.Author.ID,
			MessageSent: false,
		})
		timeCache := time.Since(timeCacheStart)

		logText := fmt.Sprintf("Added to Send Cache: [Total: %v][Generate: %v][Special: %v][Cache: %v], Generate Log: %v", time.Since(timeGenStart), timeGenerate, timeSpecial, timeCache, newWow.LogText)
		audit.Log(logText)

	} else if message.ChannelID == constants.CHANNEL_BOT_SPAM {

		// 4. Did they fail pasting?
		charsOnly := strings.ToLower(strings.Replace(strings.TrimSpace(message.Content), " ", "", -1))
		failPaste := true
		for _, ch := range charsOnly {
			if ch != 'v' && ch != 'V' {
				failPaste = false
				break
			}
		}

		if failPaste {
			InsertIntoEffectCache(message, WowEffect{
				UserID:            message.Author.ID,
				EffectName:        "failpaste",
				EffectDescription: "**Absolute Pasting Pete** 🐍: (60s) Every response has 1 negative Dice Roll deducted from its final score",
				ActiveUntil:       time.Now().Add(1 * time.Minute),
				FreeRolls:         -1,
			})
		}
	}
}

var fibCache []int = []int{}

func FibonacciValue() int {
	timeStart := time.Now()
	genCache := false
	fibs := []int{1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233}
	if len(fibCache) == 0 {
		fibsRev := []int{1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233}
		helpers.ReverseIntArray(fibsRev)
		for i, f := range fibsRev {
			for l := 0; l < f; l++ {
				fibCache = append(fibCache, fibs[i])
			}
		}
		genCache = true
	}
	arrayTimes := time.Since(timeStart)
	randStart := time.Now()
	randomNumber := helpers.GetRandomInt(fibCache)
	randTimes := time.Since(randStart)

	logText := fmt.Sprintf("Time taken: [Total: %v][Generated Cache: %v][Array Gen: %v]", time.Since(timeStart), genCache, arrayTimes)
	logText += fmt.Sprintf("[Random Number: %v] Rolled: %v", randTimes, randomNumber)
	audit.Log(logText)
	return randomNumber

}

func NewOneArmBandit() (int, string) {
	timeStart := time.Now()

	// Get our Roll ints
	randomRoll := helpers.GetRandomNumber(0, 999)
	individualInts := helpers.GetIndividualInts(randomRoll)
	getRolls := time.Since(timeStart)

	initStart := time.Now()
	draw1 := individualInts[0]
	draw2 := 0
	draw3 := 0
	if len(individualInts) > 1 {
		draw2 = individualInts[1]
	}
	if len(individualInts) > 2 {
		draw3 = individualInts[2]
	}

	var result string = fmt.Sprintf("[ %d | %d | %d ]", draw1, draw2, draw3)

	var matches int = 0
	var matchedValue int = 0
	initTime := time.Since(initStart)

	normDrawStart := time.Now()
	if draw1 == draw2 && draw1 == draw3 {
		matches = 3
		matchedValue = draw1
	} else if draw1 == draw2 || draw1 == draw3 || draw2 == draw3 {
		matches = 2
		if draw1 == draw2 {
			matchedValue = draw1
		} else if draw1 == draw3 {
			matchedValue = draw1
		} else if draw2 == draw3 {
			matchedValue = draw2
		}
	}
	normDrawTime := time.Since(normDrawStart)

	comboStart := time.Now()
	if matches == 0 {
		sortArr := []int{draw1, draw2, draw3}
		sort.Ints(sortArr)
		if sortArr[2] == sortArr[1]+1 && sortArr[1] == sortArr[0]+1 {
			matches = 5
			matchedValue = sortArr[2]
		}
	}
	comboTime := time.Since(comboStart)

	logText := fmt.Sprintf("Time taken: [Total: %v][Get Rolls: %v][Init Time: %v][Norm Draw: %v]", time.Since(timeStart), getRolls, initTime, normDrawTime)
	logText += fmt.Sprintf("[Combo Time: %v] Rolled: %v", comboTime, result)
	audit.Log(logText)

	if matches >= 2 {
		return matchedValue * matches, result
	}
	return 0, result
}
