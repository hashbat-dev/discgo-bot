package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

type ImgurLog struct {
	ID              int
	GuildID         string
	UserID          string
	ImgurID         string
	ImgurType       string
	ImgurTitle      string
	ImgurLink       string
	ImgurDeleteHash string
	CreatedDateTime time.Time
}

func InsertImgurLog(GuildID string, UserID string, ImgurID string, ImgurType string, ImgurTitle string, ImgurLink string, ImgurDeleteHash string) error {
	query := "INSERT INTO ImgurLog (GuildID, UserID, ImgurID, ImgurType, ImgurTitle, ImgurLink, ImgurDeleteHash) VALUES (?, ?, ?, ?, ?, ?, ?)"
	insertResult, err := Db.ExecContext(context.Background(), query, GuildID, UserID, ImgurID, ImgurType, ImgurTitle, ImgurLink, ImgurDeleteHash)
	if err != nil {
		return err
	}

	id, err := insertResult.LastInsertId()
	if err != nil {
		return err
	} else if id == 0 {
		err = errors.New("returned id insert was 0")
		return err
	}

	logger.Debug(GuildID, "New ImgurLog inserted: %v [%s]", id, ImgurLink)
	return nil
}

func DeleteImgurLog(guildId string, deleteHash string) error {
	query := `DELETE FROM ImgurLog WHERE ImgurDeleteHash = ?`
	_, err := Db.Exec(query, deleteHash)
	if err != nil {
		logger.Error(guildId, err)
		return err
	}

	return nil
}

func GetAllImgurLogs(guildId string) ([]ImgurLog, error) {
	query := "SELECT ID, GuildID, UserID, ImgurID, ImgurType, ImgurTitle, ImgurLink, ImgurDeleteHash, CreatedDateTime FROM ImgurLog ORDER BY ID"
	rows, err := Db.Query(query)
	if err != nil {
		logger.Error(guildId, err)
		return nil, err
	}
	defer rows.Close()

	var retSlice []ImgurLog
	for rows.Next() {

		var newEntry ImgurLog
		var ID sql.NullInt32
		var GuildID, UserID, ImgurID, ImgurType, ImgurTitle, ImgurLink, ImgurDeleteHash sql.NullString
		var CreatedDateTime sql.NullTime

		err := rows.Scan(&ID, &GuildID, &UserID, &ImgurID, &ImgurType, &ImgurTitle, &ImgurLink, &ImgurDeleteHash, &CreatedDateTime)
		if err != nil {
			logger.Error(guildId, err)
			return nil, err
		}

		if ID.Valid {
			newEntry.ID = int(ID.Int32)
		}
		if GuildID.Valid {
			newEntry.GuildID = GuildID.String
		}
		if UserID.Valid {
			newEntry.UserID = UserID.String
		}
		if ImgurID.Valid {
			newEntry.ImgurID = ImgurID.String
		}
		if ImgurType.Valid {
			newEntry.ImgurType = ImgurType.String
		}
		if ImgurTitle.Valid {
			newEntry.ImgurTitle = ImgurTitle.String
		}
		if ImgurLink.Valid {
			newEntry.ImgurLink = ImgurLink.String
		}
		if ImgurDeleteHash.Valid {
			newEntry.ImgurDeleteHash = ImgurDeleteHash.String
		}
		if CreatedDateTime.Valid {
			newEntry.CreatedDateTime = CreatedDateTime.Time
		}

		retSlice = append(retSlice, newEntry)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return retSlice, nil
}
