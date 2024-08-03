package dbhelper

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dabi-ngin/discgo-bot/Bot/audit"
	"github.com/dabi-ngin/discgo-bot/Bot/config"
)

func CountWow(user string, count int) {

	if config.IsDev {
		return
	}

	// Get current Max (if any)
	currMax, _, err := CountWow_GetCounts(user)
	if err != nil {
		return
	}

	if currMax > 0 {
		// Update
		if count <= currMax {
			count = 0
		}
		err = CountWow_Update(user, count)
	} else {
		// Create
		err = CountWow_Insert(user, count)
	}

	if err != nil {
		audit.ErrorWithText(fmt.Sprintf("Length: %v, UserID: %v", count, user), err)
	}

}

func CountWow_Insert(user string, maxWow int) error {
	query := "INSERT INTO WowCount (UserID, TotalCount, MaxWow, LastUsed) VALUES (?, 1, ?, NOW())"
	result, err := db.ExecContext(context.Background(), query, user, maxWow)
	if err != nil {
		return err
	}

	var id int64
	id, err = result.LastInsertId()
	if err != nil {
		return err
	}

	audit.Log("Inserted ID: " + strconv.FormatInt(id, 10))
	if id == 0 {
		return errors.New("CountWow_Insert - Blank ID returned on insert")
	} else {
		return nil
	}
}

func CountWow_Update(user string, maxWow int) error {
	query := "UPDATE WowCount SET TotalCount=TotalCount+1"
	if maxWow > 0 {
		query += ", MaxWow = " + fmt.Sprint(maxWow)
	}
	query += ", LastUsed=NOW() WHERE UserID = ?"
	result, err := db.ExecContext(context.Background(), query, user)

	if err != nil {
		return err
	}

	var check int64
	check, err = result.RowsAffected()

	if err != nil {
		return err
	} else if check == 0 {
		return errors.New("CountWow_Update - 0 rows affected on update")
	} else {
		return nil
	}
}

func CountWow_GetCounts(user string) (int, int, error) {
	var wowCount int
	var totalCount int
	err := db.QueryRow("SELECT MaxWow, TotalCount FROM WowCount WHERE UserID = ?", user).Scan(&wowCount, &totalCount)
	if err != nil && strings.Contains(err.Error(), "no rows") {
		wowCount = 0
		totalCount = 0
		err = nil
	} else if err != nil {
		audit.Error(err)
	}

	return wowCount, totalCount, err
}

type WowCountEntry struct {
	ID         int
	UserID     string
	TotalCount int
	MaxWow     int
	LastUsed   time.Time
}

func CountWow_Ranking() ([]WowCountEntry, error) {
	var retArray []WowCountEntry
	query := "SELECT ID, UserID, TotalCount, MaxWow, LastUsed FROM BotDB.WowCount ORDER BY MaxWow DESC;"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows
	for rows.Next() {

		var id, totalCount, maxWow sql.NullInt32
		var userId sql.NullString
		var lastUsed sql.NullTime
		var entry WowCountEntry

		err := rows.Scan(&id, &userId, &totalCount, &maxWow, &lastUsed)
		if err != nil {
			return nil, err
		}

		if id.Valid {
			entry.ID = int(id.Int32)
		} else {
			entry.ID = 0
		}

		if userId.Valid {
			entry.UserID = userId.String
		} else {
			entry.UserID = ""
		}

		if totalCount.Valid {
			entry.TotalCount = int(totalCount.Int32)
		} else {
			entry.TotalCount = 0
		}

		if maxWow.Valid {
			entry.MaxWow = int(maxWow.Int32)
		} else {
			entry.MaxWow = 0
		}

		if lastUsed.Valid {
			entry.LastUsed = lastUsed.Time
		} else {
			entry.LastUsed = time.Now().AddDate(-20, 0, 0)
		}

		retArray = append(retArray, entry)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return retArray, nil

}
