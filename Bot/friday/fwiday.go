package friday

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ZestHusky/femboy-control/Bot/audit"
	"github.com/ZestHusky/femboy-control/Bot/config"
	"github.com/ZestHusky/femboy-control/Bot/constants"
	"github.com/ZestHusky/femboy-control/Bot/handlers/userstats"
	"github.com/ZestHusky/femboy-control/Bot/helpers"
	logger "github.com/ZestHusky/femboy-control/Bot/logging"
	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
)

func Fwiday(discord *discordgo.Session, interaction *discordgo.InteractionCreate, useChar string) {
	fmt.Println("[FWIDAY] Received command")

	var embedContent []*discordgo.MessageEmbed
	embedContent = append(embedContent, GetFwiday(discord, interaction, useChar))
	if embedContent[0].Title == "Error" {
		return
	}

	discord.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: embedContent,
		},
	})
}

func IsItFwiday() bool {

	if config.FwidayCancelled {
		return false
	}

	locEST, _ := time.LoadLocation("America/New_York")
	locUK, _ := time.LoadLocation("Europe/London")
	nowEST := time.Now().In(locEST)
	nowUK := time.Now().In(locUK)

	// Is it Friday in EST or UK?
	if nowEST.Weekday() == 5 || nowUK.Weekday() == 5 {
		return true
	} else {
		return false
	}
}

func GetFwiday(discord *discordgo.Session, interaction *discordgo.InteractionCreate, useChar string) *discordgo.MessageEmbed {
	isFriday := false
	locEST, _ := time.LoadLocation("America/New_York")
	locUK, _ := time.LoadLocation("Europe/London")
	nowEST := time.Now().In(locEST)
	nowUK := time.Now().In(locUK)

	// What Fwiday was used?
	inputDay := "F" + useChar + "iday"

	// What's Pog's opinion of the bot?
	pogOpinion, err := userstats.GetBotRating(constants.USER_ID_POG)
	if err != nil {
		audit.Error(err)
	}

	// Is it Friday in EST or UK?
	pogBlock := false

	actuallyIsFwiday := false
	if nowUK.Weekday() == 5 {
		isFriday = true
	} else {
		if nowEST.Weekday() == 5 {
			if pogOpinion >= 0 {
				isFriday = true
			} else {
				pogBlock = true
			}
		}
	}

	if isFriday && config.FwidayCancelled {
		actuallyIsFwiday = true
		isFriday = false
	}

	if isFriday {
		return ItsFwidayEmbed(inputDay)
	} else {
		if actuallyIsFwiday {
			e := embed.NewEmbed()
			e.SetTitle(inputDay + " is Currently Cancelled")
			e.SetDescription(helpers.GetText("cancelledfwiday", true, "", ""))
			return e.MessageEmbed
		} else if pogBlock {
			// It WOULD be Fwiday, if Pog liked the bot >w<
			e := embed.NewEmbed()
			e.SetTitle("Hmm, Nope! It's not " + inputDay)
			e.SetDescription("It might be " + inputDay + " in the US, but... Pog has a " + fmt.Sprint(pogOpinion) + " rating of me! " + helpers.GetEmote("nope", true) + "\n\nSo nope, it's not " + inputDay + ".")
			return e.MessageEmbed
		} else {
			// It's not Friday! :ae_cry:
			// Work out how long until Friday (in the UK, as this will always be first)
			addDays := 0
			switch nowUK.Weekday() {
			case 0:
				addDays = 5
			case 1:
				addDays = 4
			case 2:
				addDays = 3
			case 3:
				addDays = 2
			case 4:
				addDays = 1
			case 6:
				addDays = 6
			}

			nextFri, err := time.Parse("2006-01-02", nowUK.AddDate(0, 0, addDays).Format("2006-01-02"))
			if err != nil {
				fmt.Println("[FWIDAY] ERROR GETTING nextFri DATE")
				fmt.Println("[FWIDAY] " + err.Error())
				logger.SendErrorInteraction(interaction)
				return embed.NewEmbed().SetTitle("Error").MessageEmbed
			}

			t := GetTimeDifference(nowUK, nextFri)

			e := embed.NewEmbed()
			e.SetTitle("It's not " + inputDay)
			setDesc := inputDay + " is in " + t
			if pogOpinion <= 0 {
				setDesc += "\n\nWoooow.. Pog's opinion of me is " + fmt.Sprint(pogOpinion) + "?! Looks like it'll be UK " + inputDay + " only this week " + helpers.GetEmote("yep", true)
			}
			e.SetDescription(setDesc)
			return e.MessageEmbed
		}
	}

}

func ItsFwidayEmbed(inputDay string) *discordgo.MessageEmbed {
	e := embed.NewEmbed()
	e.SetTitle("IT'S " + strings.ToUpper(inputDay))
	e.SetDescription(helpers.GetText("itsfwiday", true, "", ""))
	return e.MessageEmbed
}

func GetTimeDifference(a, b time.Time) string {
	days, hours, minutes, seconds := getDifference(a, b)
	t := ""

	if days > 0 {
		t += strconv.Itoa(days) + " day"
		if days != 1 {
			t += "s"
		}
		t += ", "
	}
	if hours > 0 {
		t += strconv.Itoa(hours) + " hour"
		if hours != 1 {
			t += "s"
		}
		t += ", "
	}
	if minutes > 0 {
		t += strconv.Itoa(minutes) + " minute"
		if minutes != 1 {
			t += "s"
		}
		t += ", "
	}
	if seconds > 0 {
		t += strconv.Itoa(seconds) + " second"
		if minutes != 1 {
			t += "s"
		}
		t += ", "
	}

	t = strings.TrimSuffix(t, ", ")

	return t
}

func getDifference(a, b time.Time) (days, hours, minutes, seconds int) {
	monthDays := [12]int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	y1, m1, d1 := a.Date()
	y2, m2, d2 := b.Date()
	h1, min1, s1 := a.Clock()
	h2, min2, s2 := b.Clock()
	totalDays1 := y1*365 + d1
	for i := 0; i < (int)(m1)-1; i++ {
		totalDays1 += monthDays[i]
	}
	totalDays1 += leapYears(a)
	totalDays2 := y2*365 + d2
	for i := 0; i < (int)(m2)-1; i++ {
		totalDays2 += monthDays[i]
	}
	totalDays2 += leapYears(b)
	days = totalDays2 - totalDays1
	hours = h2 - h1
	minutes = min2 - min1
	seconds = s2 - s1
	if seconds < 0 {
		seconds += 60
		minutes--
	}
	if minutes < 0 {
		minutes += 60
		hours--
	}
	if hours < 0 {
		hours += 24
		days--
	}
	return days, hours, minutes, seconds
}

func leapYears(date time.Time) (leaps int) {
	y, m, _ := date.Date()
	if m <= 2 {
		y--
	}
	leaps = y/4 + y/400 - y/100
	return leaps
}
