package database

import (
	"database/sql"
	"time"
)

type BangCommand struct {
	ID            int
	Command       string
	PackageName   string
	FunctionName  string
	Active        bool
	DateTimeAdded time.Time
}

func GetBangs() ([]BangCommand, error) {

	var bangCommands []BangCommand
	query := `	SELECT 
				ID, Command, PackageName, FunctionName, Active, DateTimeAdded
				FROM BangCommands			
			  	WHERE Active = 1
			 `

	rows, err := Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows
	for rows.Next() {

		currentRow := BangCommand{
			ID:           0,
			Command:      "",
			PackageName:  "",
			FunctionName: "",
		}
		var id sql.NullInt32
		var active sql.NullBool
		var command, packageName, functionName sql.NullString
		var dateTimeAdded sql.NullTime

		err := rows.Scan(&id, &command, &packageName, &functionName, &active, &dateTimeAdded)
		if err != nil {
			return nil, err
		}

		if id.Valid {
			currentRow.ID = int(id.Int32)
		}

		if command.Valid {
			currentRow.Command = command.String
		}

		if packageName.Valid {
			currentRow.PackageName = packageName.String
		}

		if functionName.Valid {
			currentRow.FunctionName = functionName.String
		}

		if active.Valid {
			currentRow.Active = active.Bool
		}

		if dateTimeAdded.Valid {
			currentRow.DateTimeAdded = dateTimeAdded.Time
		}

		bangCommands = append(bangCommands, currentRow)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return bangCommands, nil
}
