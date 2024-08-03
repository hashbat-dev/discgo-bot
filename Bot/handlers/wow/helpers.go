package wow

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ZestHusky/femboy-control/Bot/audit"
	"github.com/ZestHusky/femboy-control/Bot/config"
	dbhelper "github.com/ZestHusky/femboy-control/Bot/dbhelpers"
	"github.com/ZestHusky/femboy-control/Bot/helpers"
	"github.com/bwmarrin/discordgo"
)

func AddCharacter(char string, times int) string {
	finalString := ""
	for i := 0; i < times; i++ {
		finalString += char
	}
	return finalString
}

func DelCharacters(currString string, times int) string {

	retValue := ""
	if !(times >= len(currString)) {
		startIndex := len(currString) - times
		retValue = currString[:startIndex]
	}

	if retValue == "" {
		retValue = "o"
	}

	return retValue
}

func GetIntsFromRoll(roll string) (int, int) {
	// Split the string by the '/' character
	parts := strings.Split(roll, "/")

	// Convert the split parts to integers
	num1, err := strconv.Atoi(parts[0])
	if err != nil {
		audit.Error(err)
		return 0, 0
	}

	num2, err := strconv.Atoi(parts[1])
	if err != nil {
		audit.Error(err)
		return 0, 0
	}

	// Return the two integers
	return num1, num2
}

func InsertIntoCache(msg *discordgo.Message, rolls []DiceRoll, middleCount int, userId string) {

	// Delete Cache items
	var newCache []WowStatItem
	for _, w := range WowStatCache {
		if time.Since(w.Added) <= 30*time.Minute {
			newCache = append(newCache, w)
		}
	}

	// Add new Item
	var newStat WowStatItem = WowStatItem{
		MessageID:   msg.ID,
		Rolls:       rolls,
		MiddleCount: middleCount,
		UserID:      userId,
		Added:       time.Now(),
	}

	newCache = append(newCache, newStat)
	WowStatCache = newCache

}

func InsertIntoEffectCache(message *discordgo.MessageCreate, effect WowEffect) {

	var newCache []WowEffect

	timeNow := time.Now()
	needToAdd := true
	for _, e := range wowEffects {

		if e.ActiveUntil.Before(timeNow) {
			continue
		}

		if e.UserID == effect.UserID && e.EffectName == effect.EffectName {
			e.ActiveUntil = effect.ActiveUntil
			needToAdd = false
			newCache = append(newCache, e)
		}
	}

	if needToAdd {
		newCache = append(newCache, effect)
	}

	wowEffects = newCache

	if effect.EffectName == "failpaste" && needToAdd {
		msgText := "You dirty cheater! " + helpers.GetEmote("ae_cry", false) + " For the next 60 seconds you're getting a debuff."
		msgText += "\n(*You get 1 deducting Dice Roll added to each Wooow*)"
		_, err := config.Session.ChannelMessageSendReply(message.ChannelID, msgText, message.Reference())
		if err != nil {
			audit.Error(err)
		} else {
			audit.Log("Gave FailPaste debuff to UserID: " + message.Author.ID)
		}
	}

}

func GetActiveEffectsForUser(userId string) []WowEffect {

	var active []WowEffect

	timeNow := time.Now()
	for _, e := range wowEffects {
		if e.UserID == userId && !e.ActiveUntil.Before(timeNow) {
			active = append(active, e)
		}
	}

	return active
}

func DeleteActiveEffect(effect WowEffect) {

	var newCache []WowEffect
	for _, e := range wowEffects {

		if e.UserID == effect.UserID && e.ActiveUntil == effect.ActiveUntil && e.EffectName == effect.EffectName {
			continue
		}

		newCache = append(newCache, e)
	}
	wowEffects = newCache

}

func UpdateActiveEffect(effect WowEffect) {

	var newCache []WowEffect
	for _, e := range wowEffects {

		if e.UserID == effect.UserID && e.EffectName == effect.EffectName {
			newCache = append(newCache, effect)
		} else {
			newCache = append(newCache, e)
		}

	}
	wowEffects = newCache

}

func UpdateSpamCache(spam WowSpamCache) {

	var newCache []WowSpamCache
	for _, e := range wowSpamCache {

		if e.UserID == spam.UserID {
			newCache = append(newCache, spam)
		} else {
			newCache = append(newCache, e)
		}

	}
	wowSpamCache = newCache

}

func DeleteAllSessionEvents(userId string) {
	var newCache []WowEffect
	for _, e := range wowEffects {

		if e.UserID == userId {
			if !e.SessionBased {
				newCache = append(newCache, e)
			}
		} else {
			newCache = append(newCache, e)
		}

	}
	wowEffects = newCache
}

func SendNextCacheItem() {
	for i, c := range SendCache {

		// Skip sent items
		if c.MessageSent {
			continue
		}

		// Send Message
		go func() {
			timeSendStart := time.Now()
			msg, err := config.Session.ChannelMessageSendReply(c.ChannelID, c.ReplyText, &c.MessageRef)
			if err != nil {
				audit.Error(err)
			} else {
				InsertIntoCache(msg, c.WowRolls, c.WowCount, c.UserID)
				dbhelper.CountWow(c.UserID, c.WowCount)
				timeSending := time.Since(timeSendStart)

				audit.Log("Sent from WowCache, " + fmt.Sprint(len(SendCache)) + " left, took " + fmt.Sprint(timeSending))
			}
		}()
		SendCache = DeleteCacheAtIndex(SendCache, i)
		return
	}
}

func DeleteCacheAtIndex(slice []WowSendCache, index int) []WowSendCache {
	if index < 0 || index >= len(slice) {
		// Index out of range, return original slice
		return slice
	}

	return append(slice[:index], slice[index+1:]...)
}
