package dbhelper

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ZestHusky/femboy-control/Bot/audit"
	"github.com/ZestHusky/femboy-control/Bot/constants"
)

type ToDoEntry struct {
	Category         string
	AssignedID       int
	ToDoText         string
	CreatedBy        string
	CreatedDateTime  time.Time
	StartedBy        string
	StartedDateTime  time.Time
	FinishedBy       string
	FinishedDateTime time.Time
	Version          string
}

func ToDoDBAdd(category string, item string, userId string) (string, error) {

	// Get the ID to use
	itemId, err := ToDoDBGetNextID(category)
	if err != nil {
		return "", err
	} else if itemId == 0 {
		return "", errors.New("failed to get next ID for item")
	}

	query := "INSERT INTO ToDoList (Category, AssignedID, ToDoText, CreatedBy, CreatedDateTime) VALUES (?, ?, ?, ?, NOW())"
	insertResult, err := db.ExecContext(context.Background(), query, category, itemId, item, userId)
	if err != nil {
		return "", err
	}

	id, err := insertResult.LastInsertId()
	if err != nil {
		return "", err
	} else if id == 0 {
		err = errors.New("returned id insert was 0")
		return "", err
	}

	audit.Log("Inserted ID: " + fmt.Sprint(id))
	return category + "-" + fmt.Sprint(itemId), nil

}

func ToDoDBGetNextID(category string) (int, error) {
	var lastId int
	err := GetDB().QueryRow("SELECT AssignedID FROM ToDoList WHERE Category = ? ORDER BY AssignedID DESC LIMIT 1;", category).Scan(&lastId)
	if err != nil && strings.Contains(err.Error(), "no rows") {
		lastId = constants.TODO_FIRST_ID
		err = nil
	} else if err != nil {
		audit.Error(err)
		return 0, err
	}

	return lastId + 1, err
}

func ToDoDBIsIDValid(category string, assignedId int) (bool, error) {

	var count int
	err := GetDB().QueryRow("SELECT COUNT(*) AS Count FROM ToDoList WHERE Category = ? AND AssignedID = ?;", category, assignedId).Scan(&count)
	if err != nil && strings.Contains(err.Error(), "no rows") {
		count = 0
		err = nil
	} else if err != nil {
		audit.Error(err)
	}

	return count > 0, err
}

func ToDoDBUpdate(category string, assignedId int, started string, finished string, newText string, newCategory string, version string) error {

	var params []interface{}
	addedCount := 0
	query := "UPDATE ToDoList SET "

	if started != "" {
		query += "StartedBy = ?, StartedDateTime = NOW()"
		params = append(params, started)
		addedCount++
	}

	if finished != "" {
		if addedCount > 0 {
			query += ", "
		}
		query += "FinishedBy = ?, FinishedDateTime = NOW()"
		params = append(params, finished)
		addedCount++
	}

	if newText != "" {
		if addedCount > 0 {
			query += ", "
		}
		query += "ToDoText = ?"
		params = append(params, newText)
		addedCount++
	}

	if newCategory != "" {
		if addedCount > 0 {
			query += ", "
		}
		query += "Category = ?"
		params = append(params, newCategory)
		addedCount++
	}

	if version != "" {
		if addedCount > 0 {
			query += ", "
		}
		query += "Version = ?"
		params = append(params, version)
		addedCount++
	}

	query += " WHERE Category = ? AND AssignedID = ?"
	params = append(params, category)
	params = append(params, assignedId)

	_, err := GetDB().ExecContext(context.Background(), query, params...)
	if err != nil {
		return err
	} else {
		audit.Log("Updated Item: " + query)
	}

	return nil
}

func ToDoDelete(category string, assignedId int) error {

	query := "UPDATE ToDoList SET DeletedDateTime = NOW() WHERE Category = ? AND AssignedID = ?"

	result, err := GetDB().ExecContext(context.Background(), query, category, assignedId)
	if err != nil {
		return err
	}

	rowCount, err := result.RowsAffected()
	if err != nil {
		return err
	} else if rowCount == 0 {
		return errors.New("no rows were affected")
	}

	audit.Log(fmt.Sprintf("Item deleted %v-%v", category, assignedId))
	return nil
}

func ToDoGetList(category string) ([]ToDoEntry, error) {

	var retArray []ToDoEntry
	query := "SELECT Category, AssignedID, ToDoText, CreatedBy, CreatedDateTime, "
	query += "StartedBy, StartedDateTime, FinishedBy, FinishedDateTime, Version FROM ToDoList WHERE DeletedDateTime IS NULL AND (FinishedDateTime IS NULL OR Version IS NULL)"
	if category != "" {
		query += " AND Category = '" + category + "'"
	}
	query += " ORDER BY AssignedID"
	rows, err := GetDB().Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows
	for rows.Next() {

		var newEntry ToDoEntry
		var category, todotext, createdby, startedby, finishedby, version sql.NullString
		var createddatetime, starteddatetime, finisheddatetime sql.NullTime
		var assignedId int

		err := rows.Scan(&category, &assignedId, &todotext, &createdby, &createddatetime, &startedby, &starteddatetime, &finishedby, &finisheddatetime, &version)
		if err != nil {
			return nil, err
		}

		newEntry.Category = category.String
		newEntry.AssignedID = assignedId
		newEntry.ToDoText = todotext.String
		newEntry.CreatedBy = createdby.String
		newEntry.CreatedDateTime = createddatetime.Time
		if startedby.Valid {
			newEntry.StartedBy = startedby.String
		} else {
			newEntry.StartedBy = ""
		}
		if starteddatetime.Valid {
			newEntry.StartedDateTime = starteddatetime.Time
		}
		if finishedby.Valid {
			newEntry.FinishedBy = finishedby.String
		} else {
			newEntry.FinishedBy = ""
		}
		if finisheddatetime.Valid {
			newEntry.FinishedDateTime = finisheddatetime.Time
		}
		if version.Valid {
			newEntry.Version = version.String
		} else {
			newEntry.Version = ""
		}

		retArray = append(retArray, newEntry)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return retArray, nil

}
