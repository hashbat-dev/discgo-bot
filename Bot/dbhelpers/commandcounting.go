package dbhelper

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/dabi-ngin/discgo-bot/Bot/audit"
	"github.com/dabi-ngin/discgo-bot/Bot/config"
)

func CountCommand(command string, user string) {

	if config.IsDev {
		return
	}

	if strings.TrimSpace(command) == "" {
		return
	}

	// Get the ID
	id, err := CommandLogGetID(command, user)
	if err != nil {
		audit.Error(err)
	}

	if id > 0 {
		// Update
		err = CommandLog_Update(id)
	} else {
		// Create
		err = CommandLogInsert(command, user)
	}

	if err != nil {
		audit.ErrorWithText("Command: "+command+", UserID: "+user, err)
	}
}

func CommandLogInsert(command string, user string) error {
	query := "INSERT INTO CommandLog (Command, UserID) VALUES (?, ?)"
	result, err := db.ExecContext(context.Background(), query, command, user)
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
		return errors.New("commandLog_Insert - Blank ID returned on insert")
	} else {
		return nil
	}
}

func CommandLog_Update(id int64) error {
	query := "UPDATE CommandLog SET Count=Count+1, LastUsed=NOW() WHERE ID = ?"
	result, err := db.ExecContext(context.Background(), query, id)

	if err != nil {
		return err
	}

	var check int64
	check, err = result.RowsAffected()

	if err != nil {
		return err
	} else if check == 0 {
		return errors.New("CommandLog_Update - 0 rows affected on update")
	} else {
		return nil
	}
}

func CommandLogGetID(command string, user string) (int64, error) {
	var id int64
	err := db.QueryRow("SELECT ID FROM CommandLog WHERE Command = ? AND UserID = ?", command, user).Scan(&id)
	if err != nil && strings.Contains(err.Error(), "no rows") {
		id = 0
		err = nil
	}

	return id, err
}

func CommandLogGetCountForUser(command string, user string) (int, error) {

	retVal := 0
	var count sql.NullInt32
	err := db.QueryRow("SELECT SUM(Count) AS Cnt FROM BotDB.CommandLog WHERE UserID = ? AND Command = ?;", user, command).Scan(&count)
	if err != nil && strings.Contains(err.Error(), "no rows") {
		retVal = 0
		err = nil
	} else if err != nil {
		audit.Error(err)
		retVal = 0
	} else if count.Valid {
		retVal = int(count.Int32)
	}

	return retVal, err
}
