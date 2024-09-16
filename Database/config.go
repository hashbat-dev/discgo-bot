package database

import (
	"database/sql"
	"errors"
	"time"

	helpers "github.com/dabi-ngin/discgo-bot/Helpers"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
)

func GetLastFakeYouCheck() (time.Time, error) {
	var time sql.NullTime
	err := Db.QueryRow("SELECT LastFakeYouCheck FROM Config LIMIT 1").Scan(&time)
	if err != nil {
		logger.Error("CONFIG", err)
		return helpers.GetNullDateTime(), err
	}

	if !time.Valid {
		return helpers.GetNullDateTime(), errors.New("invalid date time")
	} else {
		return time.Time, nil
	}
}

func UpdateLastFakeYouCheck() {
	query := "UPDATE Config SET LastFakeYouCheck = NOW()"
	_, err := Db.Exec(query)
	if err != nil {
		logger.Error("CONFIG", err)
	}
}
