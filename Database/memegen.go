package database

import (
	"context"
	"errors"

	logger "github.com/dabi-ngin/discgo-bot/Logger"
)

func InsertMemeGenLog(guildId string, userId string, correlationId string, messageId string, sourceUrl string, sourceExt string, memeGenUrl string) {
	query := `INSERT INTO MemeGenLog
				(GuildID, UserID, CorrelationID, ResultMessageID, SourceURL, SourceExt, MemeGenURL)
				VALUES
				(?, ?, ?, ?, ?, ?, ?)`
	insertResult, err := Db.ExecContext(context.Background(), query, guildId, userId, correlationId, messageId, sourceUrl, sourceExt, memeGenUrl)
	if err != nil {
		logger.Error(guildId, err)
		return
	}

	id, err := insertResult.LastInsertId()
	if err != nil {
		logger.Error(guildId, err)
		return
	} else if id == 0 {
		err = errors.New("InsertMemeGenLog returned 0 as inserted id")
		logger.Error(guildId, err)
		return
	}

	logger.Debug(guildId, "Inserted MemeGenLog ID: %v", id)
}
