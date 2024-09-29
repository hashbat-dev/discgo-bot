package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

type FakeYouModel struct {
	ID             int
	Title          string
	ModelToken     string
	UpdateDateTime time.Time
	AddedDateTime  time.Time
}

func AddOrUpdateFakeYouModel(Title string, ModelToken string) error {
	query := `
		INSERT INTO FakeYouModels (Title, ModelToken, UpdatedDateTime, AddedDateTime)
		VALUES (?, ?, NOW(), NOW())
		ON DUPLICATE KEY UPDATE 
			ModelToken = IF(ModelToken <> VALUES(ModelToken), VALUES(ModelToken), ModelToken),
			UpdatedDateTime = IF(ModelToken <> VALUES(ModelToken), NOW(), UpdatedDateTime);
		`
	_, err := Db.Exec(query, Title, ModelToken)
	if err != nil {
		logger.Error("FAKEYOU", err)
	}

	return err
}

func GetFakeYouModels(searchTerm string) (map[string]FakeYouModel, error) {
	var retArray = make(map[string]FakeYouModel)
	var rows *sql.Rows
	var err error

	// 1. Get rows from the Database
	if searchTerm == "" {
		query := "SELECT ID, Title, ModelToken, UpdatedDateTime, AddedDateTime FROM FakeYouModels ORDER BY ID"
		rows, err = Db.Query(query)
		if err != nil {
			return nil, err
		}
	} else {
		searchTerm = "%" + searchTerm + "%"
		query := "SELECT ID, Title, ModelToken, UpdatedDateTime, AddedDateTime FROM FakeYouModels WHERE Title LIKE ? ORDER BY ID"
		rows, err = Db.Query(query, searchTerm)
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()

	// 2. Iterate over the returned rows
	for rows.Next() {

		var newEntry FakeYouModel
		var ID sql.NullInt32
		var Title, ModelToken sql.NullString
		var UpdatedDateTime, AddedDateTime sql.NullTime

		err := rows.Scan(&ID, &Title, &ModelToken, &UpdatedDateTime, &AddedDateTime)
		if err != nil {
			return nil, err
		}

		if !ID.Valid {
			logger.ErrorText("FAKEYOU", "DB Item has invalid ID")
		}

		if !Title.Valid {
			logger.ErrorText("FAKEYOU", "DB Item has invalid Title")
		}

		newEntry.ID = int(ID.Int32)
		newEntry.Title = Title.String
		if ModelToken.Valid {
			newEntry.ModelToken = ModelToken.String
		}
		if UpdatedDateTime.Valid {
			newEntry.UpdateDateTime = UpdatedDateTime.Time
		}
		if AddedDateTime.Valid {
			newEntry.AddedDateTime = AddedDateTime.Time
		}

		retArray[newEntry.Title] = newEntry
	}

	// 3. Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return retArray, nil

}

func GetModelNameFromToken(guildId string, modelToken string) (string, error) {
	var title sql.NullString
	err := Db.QueryRow("SELECT Title FROM FakeYouModels WHERE ModelToken = ? LIMIT 1; ", modelToken).Scan(&title)
	if err != nil {
		logger.Error(guildId, err)
		return "", err
	}

	if !title.Valid {
		err = fmt.Errorf("model token [%v] returned invalid title", modelToken)
		logger.Error(guildId, err)
		return "", err
	}

	return title.String, nil
}

func DeleteFakeYouModel(model FakeYouModel) error {
	query := "DELETE FROM FakeYouModels WHERE ID = ?"
	deleteResult, err := Db.ExecContext(context.Background(), query, model.ID)
	if err != nil {
		return err
	}

	var check int64
	check, err = deleteResult.RowsAffected()

	if err != nil {
		logger.Error("FAKEYOU", err)
		return err
	} else if check == 0 {
		err = fmt.Errorf("tried to delete FakeYouModel Id %v but no rows affected", model.ID)
		logger.Error("FAKEYOU", err)
		return err
	} else {
		logger.Info("FAKEYOU", "Deleted orphaned Model, ID:[%v], Title:[%v] Model:[%v]", model.ID, model.Title, model.ModelToken)
		return nil
	}
}

func InsertFakeYouLog(guildId string, userId string, correlationId string, messageId string, modelName string, modelId string, requestText string, attempts int) {
	query := `INSERT INTO FakeYouRequestLog
				(GuildID, UserID, CorrelationID, ResultMessageID, ModelName, ModelID, RequestText, Attempts)
				VALUES
				(?, ?, ?, ?, ?, ?, ?, ?)`
	insertResult, err := Db.ExecContext(context.Background(), query, guildId, userId, correlationId, messageId, modelName, modelId, requestText, attempts)
	if err != nil {
		logger.Error(guildId, err)
		return
	}

	id, err := insertResult.LastInsertId()
	if err != nil {
		logger.Error(guildId, err)
		return
	} else if id == 0 {
		err = errors.New("insertFakeYouLog returned 0 as inserted id")
		logger.Error(guildId, err)
		return
	}

	logger.Debug(guildId, "Inserted FakeYouRequestLog ID: %v", id)
}

type FakeYouLog struct {
	ID                int
	GuildID           string
	UserID            string
	CorrelationID     string
	ResultMessageID   string
	ModelName         string
	ModelID           string
	RequestText       string
	Attempts          int
	CompletedDateTime time.Time
}

func GetFakeYouLog(guildId string, messageId string) (FakeYouLog, error) {

	var ID, Attempts sql.NullInt32
	var GuildID, UserID, CorrelationID, ResultMessageID, ModelName, ModelID, RequestText sql.NullString
	var CompletedDateTime sql.NullTime

	query := `SELECT ID, GuildID, UserID, CorrelationID, ResultMessageID, ModelName, ModelID, RequestText, Attempts, CompletedDateTime
				FROM FakeYouRequestLog WHERE GuildID = ? AND ResultMessageID = ?`
	err := Db.QueryRow(query, guildId, messageId).Scan(&ID, &GuildID, &UserID, &CorrelationID, &ResultMessageID, &ModelName, &ModelID, &RequestText, &Attempts, &CompletedDateTime)
	if err != nil {
		logger.Error(guildId, err)
		return FakeYouLog{}, err
	}

	retLog := FakeYouLog{}
	if ID.Valid {
		retLog.ID = int(ID.Int32)
	}
	if GuildID.Valid {
		retLog.GuildID = GuildID.String
	}
	if UserID.Valid {
		retLog.UserID = UserID.String
	}
	if CorrelationID.Valid {
		retLog.CorrelationID = CorrelationID.String
	}
	if ResultMessageID.Valid {
		retLog.ResultMessageID = ResultMessageID.String
	}
	if ModelName.Valid {
		retLog.ModelName = ModelName.String
	}
	if ModelID.Valid {
		retLog.ModelID = ModelID.String
	}
	if RequestText.Valid {
		retLog.RequestText = RequestText.String
	}
	if Attempts.Valid {
		retLog.Attempts = int(Attempts.Int32)
	}

	return retLog, nil

}
