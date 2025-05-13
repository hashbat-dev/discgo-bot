package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	helpers "github.com/hashbat-dev/discgo-bot/Helpers"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

func GetLastCheck(field string) (time.Time, error) {
	var time sql.NullTime
	query := fmt.Sprintf("SELECT %s FROM Config LIMIT 1", field)
	err := Db.QueryRow(query).Scan(&time)
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

func UpdateLastCheck(field string) {
	query := fmt.Sprintf("UPDATE Config SET %s = NOW()", field)
	_, err := Db.Exec(query, field)
	if err != nil {
		logger.Error("CONFIG", err)
	}
}
