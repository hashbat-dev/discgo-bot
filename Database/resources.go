package database

import (
	"database/sql"
	"errors"

	logger "github.com/dabi-ngin/discgo-bot/Logger"
)

func GetRandomResource(guildId string, resourceTypeId int) (string, error) {

	query := `SELECT Resource FROM ResourceStorage WHERE ResourceTypeID = ?	ORDER BY RAND() LIMIT 1;`
	var dbRes sql.NullString

	err := Db.QueryRow(query, resourceTypeId).Scan(&dbRes)
	if err != nil {
		logger.Error(guildId, err)
		return "", err
	}

	if !dbRes.Valid {
		err = errors.New("invalid resource returned from db")
		logger.Error(guildId, err)
		return "", err
	}

	return dbRes.String, nil

}
