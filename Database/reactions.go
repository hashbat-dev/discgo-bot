package database

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

type StarboardMessage struct {
	ID                 int
	GuildID            string
	UserID             string
	OriginalMessageID  string
	StarboardMessageID string
	IsUpChannel        bool
	Score              int
	EmojiString        string
	CreatedDateTime    time.Time
}

func Starboard_Get(guildId string, originalMessageId string) StarboardMessage {
	var ID, Score sql.NullInt32
	var GuildID, UserID, OriginalMessageID, StarboardMessageID, EmojiString sql.NullString
	var IsUpChannel sql.NullBool
	var CreatedDateTime sql.NullTime
	var r StarboardMessage
	query := `SELECT
				ID, GuildID, UserID, OriginalMessageID, StarboardMessageID, IsUpChannel,
				Score, EmojiString, CreatedDateTime
				FROM StarboardMessages
				WHERE GuildID = ? AND OriginalMessageID = ?
			`
	err := Db.QueryRow(query, guildId, originalMessageId).Scan(&ID, &GuildID, &UserID, &OriginalMessageID,
		&StarboardMessageID, &IsUpChannel, &Score, &EmojiString, &CreatedDateTime)
	if err != nil {
		if !strings.Contains(err.Error(), "no rows") {
			logger.Error(guildId, err)
		}
		return r
	}

	if !ID.Valid {
		return r
	} else {
		r.ID = int(ID.Int32)
	}

	if GuildID.Valid {
		r.GuildID = GuildID.String
	}
	if UserID.Valid {
		r.UserID = UserID.String
	}
	if OriginalMessageID.Valid {
		r.OriginalMessageID = OriginalMessageID.String
	}
	if StarboardMessageID.Valid {
		r.StarboardMessageID = StarboardMessageID.String
	}
	if IsUpChannel.Valid {
		r.IsUpChannel = IsUpChannel.Bool
	}
	if Score.Valid {
		r.Score = int(Score.Int32)
	}
	if EmojiString.Valid {
		r.EmojiString = EmojiString.String
	}
	if CreatedDateTime.Valid {
		r.CreatedDateTime = CreatedDateTime.Time
	}
	return r
}

func Starboard_Delete(guildId string, dbId int) {
	query := `DELETE FROM StarboardMessages WHERE ID = ?`
	_, err := Db.Exec(query, dbId)
	if err != nil {
		logger.Error(guildId, err)
	}
}

func Starboard_InsertUpdate(newObj StarboardMessage) {
	if newObj.ID == 0 {
		// Insert
		query := `INSERT INTO
					StarboardMessages
					(GuildID, UserID, OriginalMessageID, StarboardMessageID, IsUpChannel, Score, EmojiString)
					VALUES
					(?, ?, ?, ?, ?, ?, ?)
		`
		insertResult, err := Db.ExecContext(context.Background(), query, newObj.GuildID, newObj.UserID, newObj.OriginalMessageID, newObj.StarboardMessageID,
			newObj.IsUpChannel, newObj.Score, newObj.EmojiString)
		if err != nil {
			logger.Error(newObj.GuildID, err)
			return
		}
		id, err := insertResult.LastInsertId()
		if err != nil {
			logger.Error(newObj.GuildID, err)
			return
		} else if id == 0 {
			err = errors.New("starboard insert returned id = 0")
			logger.Error(newObj.GuildID, err)
			return
		}
	} else {
		// Update
		updateQuery := `UPDATE 
						StarboardMessages
						SET
						GuildID = ?,
						UserID = ?,
						OriginalMessageID = ?,
						StarboardMessageID = ?,
						IsUpChannel = ?,
						Score = ?,
						EmojiString = ?
						WHERE ID = ?
		`
		result, err := Db.Exec(updateQuery, newObj.GuildID, newObj.UserID, newObj.OriginalMessageID, newObj.StarboardMessageID,
			newObj.IsUpChannel, newObj.Score, newObj.EmojiString, newObj.ID)
		if err != nil {
			logger.Error(newObj.GuildID, err)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			logger.Error(newObj.GuildID, err)
			return
		}

		if rowsAffected == 0 {
			err = errors.New("starboard update returned 0 rows affected")
			logger.Error(newObj.GuildID, err)
			return
		}

	}

}
