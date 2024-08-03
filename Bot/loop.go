package bot

import (
	"os"
	"path/filepath"
	"time"

	"github.com/ZestHusky/femboy-control/Bot/audit"
	"github.com/ZestHusky/femboy-control/Bot/config"
	"github.com/ZestHusky/femboy-control/Bot/constants"
	"github.com/ZestHusky/femboy-control/Bot/friday"
	"github.com/ZestHusky/femboy-control/Bot/handlers/wow"
	"github.com/ZestHusky/femboy-control/Bot/helpers"
	embed "github.com/clinet/discordgo-embed"
)

var LastCheckedDate int = 0
var SecondCheck2 int = 0
var SecondCheck10 int = 0
var MinuteCheck5 int = 0
var firstLoop bool = true
var devSkippedFridayCheck = false

func Loop(stop chan bool) {
	for {
		select {
		case <-stop:
			// Received stop signal, exit the loop
			audit.Log("Loop Stopped")
			return
		case <-time.After(time.Duration(500) * time.Millisecond):
			SecondCheck2++
			SecondCheck10++
			MinuteCheck5++

			// Perform every Second
			wow.SendNextCacheItem()

			// Only perform every 2 seconds
			if firstLoop || SecondCheck2 == 4 {
				Second2()
				SecondCheck2 = 0
			}

			// Only perform every 10 Seconds
			if firstLoop || SecondCheck10 == 20 {
				Second10()
				SecondCheck10 = 0
			}

			// Only perform every 5 minutes
			if firstLoop || MinuteCheck5 == 600 {
				Minute5()
				MinuteCheck5 = 0
			}

			firstLoop = false
		}
	}
}

func Second2() {
	audit.SendNextLogBatch()
}

func Second10() {
	// 1. Is it Fwiday? ----------------------------
	// 1A. Check if the Date is different to the last one we checked
	locUK, _ := time.LoadLocation("Europe/London")
	_, _, d := time.Now().In(locUK).Date()
	if d != LastCheckedDate {

		if config.IsDev {
			return
		}

		// 1B. It is, do the check
		if time.Now().In(locUK).Weekday() == time.Friday {
			if config.FwidayCancelled {
				audit.Log("Fwiday Cancelled")
				e := embed.NewEmbed()
				e.SetTitle("Fwiday is Cancelled")
				e.SetDescription(helpers.GetText("cancelledfwiday", true, "", ""))
				config.Session.ChannelMessageSendEmbed(constants.CHANNEL_GENERAL, e.MessageEmbed)
			} else {
				audit.Log("It's Fwiday!")
				embed := friday.ItsFwidayEmbed("FWIDAY")
				config.Session.ChannelMessageSendEmbed(constants.CHANNEL_GENERAL, embed)
			}
		}
		LastCheckedDate = d
	} else {
		if !devSkippedFridayCheck {
			audit.LogDevOnly("Skipping Fwiday check!")
			devSkippedFridayCheck = true
		}
	}
}

func Minute5() {
	go ClearTempFolder()
}

func ClearTempFolder() {

	currentTime := time.Now()
	thresholdTime := currentTime.Add(-5 * time.Minute)

	// Walk through the directory and its subdirectories
	err := filepath.Walk(constants.TEMP_DIRECTORY, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file is a regular file and older than 5 minutes
		if info.Mode().IsRegular() && info.ModTime().Before(thresholdTime) {
			// Delete the file
			err := os.Remove(path)
			if err != nil {
				audit.Error(err)
			} else {
				audit.Log("Deleted Temp file: " + path)
			}
		}

		return err
	})

	if err != nil {
		audit.Error(err)
	}
}
