package fart

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ZestHusky/femboy-control/Bot/audit"
	"github.com/ZestHusky/femboy-control/Bot/config"
	dbhelpers "github.com/ZestHusky/femboy-control/Bot/dbhelpers"
	"github.com/ZestHusky/femboy-control/Bot/helpers"
	logger "github.com/ZestHusky/femboy-control/Bot/logging"
	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
)

type Farter struct {
	userId     string
	count      int
	lastUpdate time.Time
}

var fartCache []Farter
var fartBuffer float64 = 15

func LogFart(message *discordgo.MessageCreate) error {
	// Okay, so someone has farted. Don't panic!
	// You open the window while this function does its work

	// 1. Is this for the sender, or did they reply to/mention someone?
	inUserId := message.Author.ID

	// 2. Do they exist in the Cache? This will save a few DB checks
	var cacheIndex int = -1
	if len(fartCache) > 0 {
		for i, cachedFarter := range fartCache {
			if cachedFarter.userId == inUserId {
				cacheIndex = i
			}
		}
	}

	// 2A. If not, do they exist in the Database?
	freshFarter := false
	dbPhrase := "fart"
	if cacheIndex < 0 {
		var userId string = ""
		var count int32
		var lastUpdate time.Time
		err := dbhelpers.GetDB().QueryRow("SELECT UserID, Count, LastTracked FROM TrackWords WHERE UserID = ? AND Phrase = ?", inUserId, dbPhrase).Scan(&userId, &count, &lastUpdate)
		if err != nil && strings.Contains(err.Error(), "no rows") {
			userId = ""
			err = nil
		} else if err != nil {
			return err
		}

		if userId != "" {

			// Yup! They're a pre-existing farter
			newEntry := Farter{
				userId:     userId,
				count:      int(count),
				lastUpdate: lastUpdate,
			}

			fartCache = append(fartCache, newEntry)

		} else {

			// Nope, create a new database row for them
			initCount := 0
			query := "INSERT INTO TrackWords (Phrase, UserID, Count, LastTracked) VALUES (?, ?, ?, NOW())"
			insertResult, err := dbhelpers.GetDB().ExecContext(context.Background(), query, dbPhrase, inUserId, initCount)
			if err != nil {
				return err
			}

			InsertedId, err := insertResult.LastInsertId()
			if err != nil {
				return err
			} else if InsertedId == 0 {
				err = errors.New("returned id insert was 0")
				return err
			}

			// We have a Fresh Farter!
			newEntry := Farter{
				userId:     inUserId,
				count:      0,
				lastUpdate: time.Now().UTC(),
			}

			fartCache = append(fartCache, newEntry)
			freshFarter = true
		}

		for i, cachedFarter := range fartCache {
			if cachedFarter.userId == inUserId {
				cacheIndex = i
			}
		}
	}

	// 2B. If we don't have a cache item for them by now, we're fucked
	if cacheIndex < 0 {
		err := errors.New("couldn't create cache item")
		return err
	}

	// 2A. How recently have they used this? Check they're not spamming
	secondsSince := time.Since(fartCache[cacheIndex].lastUpdate).Seconds()
	if !freshFarter && secondsSince <= fartBuffer {
		// Either shenanigans, or someone has probably shit their pants
		audit.Log("Blocked !fart by " + inUserId + ", only " + fmt.Sprint(int(secondsSince)) + " seconds since last try")
		logger.SendErrorMsg(message, "You farted "+fmt.Sprint(int(secondsSince))+" seconds ago?! Either you've shit your pants or that's a continuation of your previous fart.")
		return nil
	}

	// 3. LET 'ER RIP!
	fartCache[cacheIndex].count++
	fartCache[cacheIndex].lastUpdate = time.Now().UTC()

	// 4. Now update the Database
	query := "UPDATE TrackWords SET Count=?, LastTracked = ? WHERE UserID = ? AND Phrase = ?"
	_, err := dbhelpers.GetDB().ExecContext(context.Background(), query, fartCache[cacheIndex].count, fartCache[cacheIndex].lastUpdate, fartCache[cacheIndex].userId, dbPhrase)
	if err != nil {
		return err
	}

	// 5. Did they hit a new rank?
	rankName, rankImg, newRank, err := dbhelpers.GetRankFromCount(fartCache[cacheIndex].count, "fart")
	if err != nil {
		return err
	}

	if newRank {
		// New rank!
		e := embed.NewEmbed()
		e.SetTitle(rankName)
		e.SetDescription("<@" + inUserId + "> hit the rank of " + rankName + "!")
		if rankImg != "" {
			e.SetImage(rankImg)
		}
		e.SetFooter("They've farted " + fmt.Sprint(fartCache[cacheIndex].count) + " times")
		config.Session.ChannelMessageSendEmbed(message.ChannelID, e.MessageEmbed)

	} else {
		// Regular ol' fart
		fartText := helpers.GetText("fart", true, "<@"+inUserId+">", "")

		fartGIFURL, err := dbhelpers.GetRandGifURL("fart")
		if err != nil {
			return err
		}

		e := embed.NewEmbed()
		e.SetTitle("Someone Farted")
		e.SetDescription(fartText)
		e.SetImage(fartGIFURL)
		e.SetFooter("They've farted " + fmt.Sprint(fartCache[cacheIndex].count) + " times")
		config.Session.ChannelMessageSendEmbed(message.ChannelID, e.MessageEmbed)
	}

	return nil
}
