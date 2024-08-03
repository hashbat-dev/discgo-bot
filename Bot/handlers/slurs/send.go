package slurs

import (
	"errors"
	"strings"

	"github.com/ZestHusky/femboy-control/Bot/audit"
	"github.com/ZestHusky/femboy-control/Bot/config"
	dbhelper "github.com/ZestHusky/femboy-control/Bot/dbhelpers"
	"github.com/ZestHusky/femboy-control/Bot/helpers"
	"github.com/ZestHusky/femboy-control/Bot/logging"
	"github.com/bwmarrin/discordgo"
)

func SendASlur(message *discordgo.MessageCreate, command string) {

	var nBomb string = ""
	if command == "job" || command == "jobslur" {
		getBomb, err := GetRandomJobSlur()
		if err != nil {
			audit.Error(err)
			logging.SendErrorMsg(message, "Sorry, I tried sending a slur but I just couldn't find it in me to abuse them >w<")
			return
		} else {
			nBomb = getBomb
		}
	} else {
		getBomb, err := GetRandomSlur()
		if err != nil {
			audit.Error(err)
			logging.SendErrorMsg(message, "Sorry, I tried sending a slur but I just couldn't find it in me to abuse them >w<")
			return
		} else {
			nBomb = getBomb
		}
	}

	if nBomb == "" {
		err := errors.New("no slur found")
		audit.Error(err)
		logging.SendErrorMsg(message, "Sorry, I tried sending a slur but I just couldn't find it in me to abuse them >w<")
		return
	}

	nBomb = strings.ReplaceAll(nBomb, "  ", " ")

	// Did they reply to a message?
	if message.ReferencedMessage != nil {
		_, err := config.Session.ChannelMessageSendReply(message.ChannelID, nBomb, message.ReferencedMessage.Reference())
		if err != nil {
			audit.Error(err)
			logging.SendErrorMsg(message, "Sorry, I tried sending a slur but I just couldn't find it in me to abuse them >w<")
			return
		}
	} else {
		_, err := config.Session.ChannelMessageSend(message.ChannelID, nBomb)
		if err != nil {
			audit.Error(err)
			logging.SendErrorMsg(message, "Sorry, I tried sending a slur but I just couldn't find it in me to abuse them >w<")
			return
		}
	}

}

var prefixes []string = []string{
	"You absolute ",
	"You're a fucking",
	"You fucking",
	"You turbo",
	"You degenerated",
	"You dirty",
	"You decrepid",
	"You negrodian",
	"You retarded",
	"You stupid",
	"You moronic",
	"You insane",
	"You submissive little",
	"You useless",
	"Alright calm down you absolute",
	"Chill out you",
}

func GetRandomSlur() (string, error) {

	// Let's use the N-Word
	slur, err := dbhelper.GetARandomSlur()
	if err != nil {
		return "", err
	}

	// Build the slur, now we're cooking
	nBomb := helpers.GetRandomText(prefixes) + " " + slur.Slur
	return nBomb, nil

}

func GetRandomJobSlur() (string, error) {

	// Let's use the N-Word
	nation, err := dbhelper.GetARandomNationality()
	if err != nil {
		return "", err
	}

	job, err := dbhelper.GetARandomJobTitle()
	if err != nil {
		return "", err
	}

	// Build the slur, now we're cooking
	nBomb := helpers.GetRandomText(prefixes) + " " + nation.Nationality + " " + job.JobTitle
	return nBomb, nil

}
