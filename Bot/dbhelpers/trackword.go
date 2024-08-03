package dbhelper

import (
	"context"
	"strings"

	"github.com/ZestHusky/femboy-control/Bot/audit"
)

func CountWord(userId string, keyword string) error {

	rowId := 0

	err := GetDB().QueryRow("SELECT ID FROM TrackWords WHERE Phrase = ? AND UserID = ?", keyword, userId).Scan(&rowId)
	if err != nil && strings.Contains(err.Error(), "no rows") {
		rowId = 0
		err = nil
	} else if err != nil {
		audit.Error(err)
		return err
	}

	if rowId == 0 {
		query := "INSERT INTO TrackWords (Phrase, UserID, Count) VALUES (?, ?, ?)"
		_, err := GetDB().ExecContext(context.Background(), query, keyword, userId, 1)
		if err != nil {
			audit.Error(err)
			return err
		}
	} else {
		query := "UPDATE TrackWords SET Count=Count+1 WHERE Phrase = ? AND UserID = ?"
		_, err := GetDB().ExecContext(context.Background(), query, keyword, userId)
		if err != nil {
			audit.Error(err)
			return err
		}
	}

	return nil

}
