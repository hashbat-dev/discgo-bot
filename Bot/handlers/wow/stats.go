package wow

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	"github.com/dabi-ngin/discgo-bot/Bot/audit"
	"github.com/dabi-ngin/discgo-bot/Bot/config"
	"github.com/dabi-ngin/discgo-bot/Bot/helpers"
	"github.com/dabi-ngin/discgo-bot/Bot/logging"
)

func GetWowStat(message *discordgo.MessageCreate) bool {

	// Did they reply to a message?
	if message.ReferencedMessage == nil {
		logging.SendErrorMsgReply(message, "You didn't reply to one of my Woooows! "+helpers.GetEmote("ae_cry", false))
		return false
	}

	if message.ReferencedMessage.Author.ID != config.Session.State.User.ID {
		logging.SendErrorMsgReply(message, "You didn't reply to one of my Woooows! "+helpers.GetEmote("ae_cry", false))
		return false
	}

	// Is the Message in our Cache?
	wowMessageId := message.ReferencedMessage.ID
	foundItem := false
	var wowCacheItem WowStatItem
	for _, w := range WowStatCache {
		if w.MessageID == wowMessageId {
			wowCacheItem = w
			foundItem = true
			break
		}
	}

	if !foundItem {
		logging.SendErrorMsgReply(message, "That Wooow's not in my cache anymore! "+helpers.GetEmote("ae_cry", false)+" I keep Woooow stats stored for 30 minutes")
		return false
	}

	// Get our Cache Items
	var wowEffects []DiceRoll
	var diceRolls []DiceRoll
	for _, roll := range wowCacheItem.Rolls {
		if roll.Special == "effect" {
			wowEffects = append(wowEffects, roll)
		} else {
			diceRolls = append(diceRolls, roll)
		}
	}

	// Output the Stats
	conEmoji := "♻️"
	wowEmoji := "😩"

	statText := "I have 2 dice that I roll (0 to 10), they're named the **Continue** and **Wooow** Dice. What do they mean?\n\n"
	statText += conEmoji + "** Continue:** 1 to 5: You get another roll! 0: You get a **bonus** roll!\n"
	statText += wowEmoji + "** Wooow:** How many O's to add to the Woooow (You start with 1!)\n\n"

	// Any Effects?
	if len(wowEffects) > 0 {
		statText += "You had the following effects in play for this Roll..."
		for _, e := range wowEffects {
			statText += "\n" + e.SpecialString
		}
		statText += "\n\n"
	}

	// Get the Rolls
	iRoll := 0
	baseContinue := 5
	baseWowRoll := 10

	statText += "The Rolls for this Wooow were..."
	for _, roll := range diceRolls {

		// Specials to show BEFORE Roll
		if roll.Special == "highermax" {
			baseContinue++
			statText += "\n**Higher Continue!** ⏫: Rolls of " + fmt.Sprint(baseContinue) + " and lower now continue! (1% Chance)"
		} else if roll.Special == "progamble" {
			baseWowRoll += 10
			statText += "\n**Pro Gambler!** 🎲: You can now get up to " + fmt.Sprint(baseWowRoll) + " on a Wooow roll! (1% Chance)"
		} else if roll.Special == "cheeky" {
			statText += "\n**Cheeky Lad!** 😳: You now get an extra Wooow roll per turn! (These do not grant Critical Success) (2% Chance)"
		} else if roll.Special == "fibonacci" {
			statText += "\n**Fibonacci Fam!** 🍷: I'm adding a randomly compounded Fibonacci Sequence to your next roll for you! (2.9% Chance)"
		} else if roll.Special == "dubs" || roll.Special == "trips" {
			statInt := "2"
			statWord := "Dubs"
			if roll.Special == "trips" {
				statWord = "Trips"
				statInt = "3"
			}
			statText += "\n**Nice " + statWord + "!** ↗️: Your Discord Message ID had some hecking " + statWord + "! "
			statText += "I've added " + statInt + " free rolls to your first Wooow roll!"

		} else if roll.Special == "oneArmBandit" {
			statText += "\n**Spin the Wheel!** 🎰: Two or more matching symbols guarantees payout! (5% Chance)"
			statText += roll.SpecialString
		}

		iRoll++
		statText += "\n"
		statText += StandardLine(roll, iRoll, conEmoji, wowEmoji)

		// Specials to show AFTER Roll
		if roll.Special == "minzero" {
			statText += "\n**Zero Chance!** 🤯: Something about the time was funky, get a free continue!"
		} else if roll.Special == "fwiday" {
			statText += "\n**It's Fwidayyy!** 😊: You got one more free roll!"
		}

	}

	e := embed.NewEmbed()
	e.SetTitle("Level " + fmt.Sprint(wowCacheItem.MiddleCount) + " WOW!")
	e.SetDescription(statText)
	if wowCacheItem.UserID != "" {
		e.SetFooter("This was " + helpers.GetNicknameFromID(message.GuildID, wowCacheItem.UserID) + "'s Wooow")
	}

	_, err := config.Session.ChannelMessageSendEmbedReply(message.ChannelID, e.MessageEmbed, message.ReferencedMessage.Reference())
	if err != nil {
		audit.Error(err)
		logging.SendErrorMsgReply(message, "Something went wrong! >w<")
		return false
	}
	return true

}

func StandardLine(roll DiceRoll, iRoll int, conEmoji string, wowEmoji string) string {
	statText := "**Roll " + fmt.Sprint(iRoll) + "**: " + conEmoji + " " + fmt.Sprint(roll.RollContinue)
	if roll.RollContinue == 0 {
		statText += "!"
	}
	statText += ", " + wowEmoji + " " + fmt.Sprint(roll.RollLength)

	if len(roll.RollSpecial) > 0 {
		statText += " (**Bonus Rolls**: "
		for x, i := range roll.RollSpecial {
			if x > 0 {
				statText += ", "
			}
			if i >= 0 {
				statText += wowEmoji
			} else {
				statText += "🪓"
			}
			statText += " " + fmt.Sprint(i)
		}
		statText += ")"
	}
	return statText
}
