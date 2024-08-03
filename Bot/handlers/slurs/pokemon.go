package slurs

import (
	"fmt"
	"strings"

	"github.com/ZestHusky/femboy-control/Bot/audit"
	dbhelper "github.com/ZestHusky/femboy-control/Bot/dbhelpers"
	"github.com/ZestHusky/femboy-control/Bot/handlers/fakeyou"
	"github.com/ZestHusky/femboy-control/Bot/helpers"
	"github.com/bwmarrin/discordgo"
)

var pokemonTTSName string = "pokemon"

func PokeSlur(message *discordgo.MessageCreate) bool {

	// 1. Get a Random Slur entry from the database
	slur, err := dbhelper.GetARandomSlur()
	if err != nil {
		helpers.ReplyToMessageWithText(message, "Hey come on, let's be nice...")
		return false
	}

	// 2. Remove the S off the end of the SlurTarget
	textToConvert := slur.Slur + "-. The " + RemoveTrailingS(slur.SlurTarget) + " Pokémon. " + slur.SlurDescription

	// 3. Sent the TTS!
	delOrig, err := fakeyou.RequestTTS(message, pokemonTTSName, textToConvert, "PokeSlur Entry "+fmt.Sprintf("%04d", slur.ID)+" "+slur.Slur, "PokéSlur Entry")
	if err != nil {
		audit.Error(err)
	}
	return delOrig
}

func RemoveTrailingS(s string) string {
	// Check if the string is empty or has a length of 1
	if len(s) == 0 || len(s) == 1 {
		return s
	}

	// Get the last character of the string
	lastChar := s[len(s)-1:]

	// Check if the last character is 's' or 'S'
	if strings.ToLower(lastChar) == "s" {
		// Return the substring excluding the last character
		return s[:len(s)-1]
	}

	// Return the original string if the last character is not 's'
	return s
}
