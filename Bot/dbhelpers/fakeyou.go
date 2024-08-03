package dbhelper

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dabi-ngin/discgo-bot/Bot/audit"
)

type TTSEntry struct {
	Command     string
	Description string
	Model       string
}

func GetTTSModel(command string) (string, error) {

	var model string
	err := GetDB().QueryRow("SELECT Model FROM TTSModels WHERE Command = ?;", command).Scan(&model)
	if err != nil && strings.Contains(err.Error(), "no rows") {
		model = ""
		err = nil
	} else if err != nil {
		audit.Error(err)
	}

	return model, err

}

func DoesTTSModelExist(command string, model string) string {

	var count int
	err := GetDB().QueryRow("SELECT COUNT(*) FROM TTSModels WHERE Command = ?;", command).Scan(&count)
	if err != nil && strings.Contains(err.Error(), "no rows") {
		count = 0
		err = nil
	} else if err != nil {
		audit.Error(err)
	}

	if err != nil {
		return err.Error()
	} else if count > 0 {
		return "Command already exists"
	}

	if model != "" {
		err = GetDB().QueryRow("SELECT COUNT(*) FROM TTSModels WHERE Model = ?;", model).Scan(&count)
		if err != nil && strings.Contains(err.Error(), "no rows") {
			count = 0
			err = nil
		} else if err != nil {
			audit.Error(err)
		}

		if err != nil {
			return err.Error()
		} else if count > 0 {
			return "Model already exists"
		}
	}

	return ""
}

func InsertTTSModel(command string, model string, description string, userId string) error {

	query := "INSERT INTO TTSModels (Command, Model, Description, UpdatedBy) VALUES (?, ?, ?, ?)"
	insertResult, err := db.ExecContext(context.Background(), query, command, model, description, userId)
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

	audit.Log(fmt.Sprintf("Insert ID: %v, Command: %v", id, command))
	return nil

}

//goland:noinspection Annotator,Annotator
func UpdateTTSModel(command string, newCommand string, newModel string, newDescription string, userId string) error {
	var params []interface{}
	addedCount := 0
	//goland:noinspection Annotator
	query := "UPDATE TTSModels SET "

	if newCommand != "" {
		query += "Command = ?"
		params = append(params, newCommand)
		addedCount++
	}

	if newModel != "" {
		if addedCount > 0 {
			query += ", "
		}
		query += "Model = ?"
		params = append(params, newModel)
		addedCount++
	}

	if newDescription != "" {
		if addedCount > 0 {
			query += ", "
		}
		query += "Description = ?"
		params = append(params, newDescription)
	}

	//goland:noinspection Annotator,Annotator
	query += ", UpdatedBy = ?, UpdatedDateTime = NOW() WHERE Command = ?"
	params = append(params, userId)
	params = append(params, command)

	_, err := GetDB().ExecContext(context.Background(), query, params...)
	if err != nil {
		return err
	} else {
		audit.Log("Updated TTSModel Item: " + query)
	}

	return nil
}

func GetTTSList() ([]TTSEntry, error) {

	var retArray []TTSEntry
	query := "SELECT Command, Description, Model FROM BotDB.TTSModels ORDER BY Command;"
	rows, err := GetDB().Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows
	for rows.Next() {

		var newEntry TTSEntry
		var command, description, model string

		err := rows.Scan(&command, &description, &model)
		if err != nil {
			return nil, err
		}

		newEntry.Command = command
		newEntry.Description = description
		newEntry.Model = model

		retArray = append(retArray, newEntry)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return retArray, nil

}
