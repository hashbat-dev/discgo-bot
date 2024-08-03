package dbhelper

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ZestHusky/femboy-control/Bot/audit"
)

func GetGIFCategoryCount(category string) (int64, error) {
	var count int64
	err := db.QueryRow("SELECT COUNT(*) AS Count FROM GifBank WHERE Category = ?", category).Scan(&count)
	if err != nil && strings.Contains(err.Error(), "no rows") {
		count = 0
		err = nil
	}

	return count, err
}

func InsertGIF(category string, gifUrl string, userId string) (int64, error) {
	query := "INSERT INTO GifBank (Category, GifURL, DateTimeAdded, AddedByID) VALUES (?, ?, NOW(), ?)"
	insertResult, err := db.ExecContext(context.Background(), query, category, gifUrl, userId)
	if err != nil {
		return 0, err
	}

	id, err := insertResult.LastInsertId()
	if err != nil {
		return 0, err
	} else if id == 0 {
		err = errors.New("returned id insert was 0")
		return 0, err
	}

	audit.Log("Inserted ID: " + fmt.Sprint(id))
	return id, nil

}

func DeleteGIF(category string, gifUrl string) error {
	query := "DELETE FROM GifBank WHERE Category = ? AND GifURL = ?"
	deleteResult, err := db.ExecContext(context.Background(), query, category, gifUrl)
	if err != nil {
		return err
	}

	var check int64
	check, err = deleteResult.RowsAffected()

	if err != nil {
		return err
	} else if check == 0 {
		return errors.New("GifBank_Delete - 0 rows affected on update")
	} else {
		audit.Log("Deleted from Category: " + category + ", GifURL: " + gifUrl)
		return nil
	}
}

func GetRandGifURL(category string) (string, error) {
	var gif string
	err := db.QueryRow("SELECT GifURL FROM GifBank WHERE Category = ? ORDER BY RAND() LIMIT 1; ", category).Scan(&gif)
	if err != nil {
		audit.Error(err)
		gif = ""
	}

	return gif, err
}

type GifBankEntry struct {
	ID            int
	Category      string
	GifURL        string
	DateTimeAdded time.Time
	AddedByID     string
}

func GetAllGifs(category string) ([]GifBankEntry, error) {

	if category == "" {
		return []GifBankEntry{}, errors.New("no category provided")
	}

	db := GetDB()
	if db == nil {
		return []GifBankEntry{}, errors.New("db object was nil")
	}

	var retArray []GifBankEntry
	query := "SELECT Category, GifURL, DateTimeAdded, AddedByID FROM GifBank WHERE Category = '" + category + "' ORDER BY DateTimeAdded DESC;"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows
	for rows.Next() {

		var gifBankEntry GifBankEntry
		var category, gifUrl, addedById sql.NullString
		var dateTimeAdded sql.NullTime

		err := rows.Scan(&category, &gifUrl, &dateTimeAdded, &addedById)
		if err != nil {
			return nil, err
		}

		if category.Valid {
			gifBankEntry.Category = category.String
		} else {
			gifBankEntry.Category = ""
		}

		if gifUrl.Valid {
			gifBankEntry.GifURL = gifUrl.String
		} else {
			gifBankEntry.GifURL = ""
		}

		if dateTimeAdded.Valid {
			gifBankEntry.DateTimeAdded = dateTimeAdded.Time
		} else {
			gifBankEntry.DateTimeAdded = time.Now().AddDate(-20, 0, 0)
		}

		if addedById.Valid {
			gifBankEntry.AddedByID = addedById.String
		} else {
			gifBankEntry.AddedByID = ""
		}

		retArray = append(retArray, gifBankEntry)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return retArray, nil

}
