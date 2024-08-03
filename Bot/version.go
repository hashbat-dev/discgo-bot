package bot

import (
	"context"
	"errors"
	"strings"

	"github.com/ZestHusky/femboy-control/Bot/audit"
	"github.com/ZestHusky/femboy-control/Bot/config"
	"github.com/ZestHusky/femboy-control/Bot/constants"
	dbhelpers "github.com/ZestHusky/femboy-control/Bot/dbhelpers"
	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
)

// When releasing a new version...
// 1. Increment the version number
// 2. Enter new features as string entries in the array
//
// The Bot will automatically pickup, log and notify

var botVersionNumber string = "1.9.15"
var botVersionFeatures []string = []string{
	"Fixed adventure implementation!",
}

func VersionCheck(discord *discordgo.Session) {

	// Are we in Dev?
	if config.IsDev {
		return
	}

	// Has this Version been recorded yet?
	verCount := 0
	err := dbhelpers.GetDB().QueryRow("SELECT COUNT(*) AS Count FROM BotDB.BotFeatures WHERE Version = ?;", botVersionNumber).Scan(&verCount)
	if err != nil && strings.Contains(err.Error(), "no rows") {
		verCount = 0
		err = nil
	} else if err != nil {
		audit.Error(err)
	}

	// Already done, exit
	if verCount > 0 {
		return
	}

	// New Version! ============================================

	// 1. Loop and Insert into Database / Build Text for Discord message
	verText := ""
	for _, feature := range botVersionFeatures {
		verText += "* " + feature + "\n"
		insertResult, err := dbhelpers.GetDB().ExecContext(context.Background(), "INSERT INTO BotFeatures (Version, Feature) VALUES (?, ?)", botVersionNumber, feature)
		if err != nil {
			audit.Error(err)
			continue
		}

		InsertedId, err := insertResult.LastInsertId()
		if err != nil {
			audit.Error(err)
			continue
		} else if InsertedId == 0 {
			audit.Error(errors.New("[VERSION] Returned id insert was 0"))
			continue
		}
	}

	// 2. Write to Discord
	e := embed.NewEmbed()
	e.SetTitle("Bottom Bot - v" + botVersionNumber)
	e.SetDescription(verText)
	e.SetImage(constants.GIF_UPDATE)
	discord.ChannelMessageSendEmbed(constants.CHANNEL_BOT_FEEDBACK, e.MessageEmbed)
}
