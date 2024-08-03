package userstats

import (
	dbhelper "github.com/dabi-ngin/discgo-bot/Bot/dbhelpers"
)

func GetBotRating(inUser string) (int, error) {
	goodBot, err := dbhelper.CommandLogGetCountForUser("goodbot", inUser)
	if err != nil {
		return 0, err
	}

	badBot, err := dbhelper.CommandLogGetCountForUser("badbot", inUser)
	if err != nil {
		return 0, err
	}
	return goodBot - badBot, nil
}
