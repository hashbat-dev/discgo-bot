package database

import (
	"database/sql"
	"errors"
	"fmt"

	triggers "github.com/hashbat-dev/discgo-bot/Bot/Commands/Triggers"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

type PhraseLeaderboardUser struct {
	UserID string
	Count  int
}

func GetGuildPhrases(guildId string, phraseId int) ([]triggers.PhraseLink, error) {
	var phraseLinks []triggers.PhraseLink

	var rows *sql.Rows
	var err error
	if phraseId == 0 {
		// Return ALL
		query := `	SELECT 
					L.ID, L.PhraseID, P.Phrase, L.GuildID, L.AddedByUserID, L.AddedDateTime, L.NotifyOnDetection, L.WordOnlyMatch
					FROM TriggerGuildLink L		
					INNER JOIN TriggerPhrases P ON P.ID = L.PhraseID
					WHERE L.GuildID = ?
				`
		rows, err = Db.Query(query, guildId)
	} else {
		// Return single instance
		query := `	SELECT 
					L.ID, L.PhraseID, P.Phrase, L.GuildID, L.AddedByUserID, L.AddedDateTime, L.NotifyOnDetection, L.WordOnlyMatch
					FROM TriggerGuildLink L		
					INNER JOIN TriggerPhrases P ON P.ID = L.PhraseID
					WHERE L.GuildID = ? AND L.PhraseID = ?
				`
		rows, err = Db.Query(query, guildId, phraseId)
	}

	if err != nil {
		logger.Error(guildId, err)
		return nil, err
	}
	defer func(g string) {
		err := rows.Close()
		if err != nil {
			logger.Error(g, err)
		}
	}(guildId)

	// Iterate over the rows
	for rows.Next() {

		var ID, PhraseID sql.NullInt32
		var NotifyOnDetection, WordOnlyMatch sql.NullBool
		var GuildID, AddedByUserID, Phrase sql.NullString
		var AddedDateTime sql.NullTime

		err := rows.Scan(&ID, &PhraseID, &Phrase, &GuildID, &AddedByUserID, &AddedDateTime, &NotifyOnDetection, &WordOnlyMatch)
		if err != nil {
			return nil, err
		}

		if !ID.Valid {
			logger.Error(guildId, errors.New("invalid db id"))
			continue
		}

		if !PhraseID.Valid {
			logger.Error(guildId, errors.New("invalid db phraseId"))
			continue
		}

		if !Phrase.Valid {
			logger.Error(guildId, errors.New("invalid db phrase"))
			continue
		}

		if !GuildID.Valid {
			logger.Error(guildId, errors.New("invalid db guildId"))
			continue
		}

		if !AddedByUserID.Valid {
			logger.Error(guildId, errors.New("invalid db addedByUserId"))
			continue
		}

		if !AddedDateTime.Valid {
			logger.Error(guildId, errors.New("invalid db addedDateTime"))
			continue
		}

		if !NotifyOnDetection.Valid {
			logger.Error(guildId, errors.New("invalid db notifyOnDetection"))
			continue
		}

		if !WordOnlyMatch.Valid {
			logger.Error(guildId, errors.New("invalid db wordOnlyMatch"))
			continue
		}

		phraseLinks = append(phraseLinks, triggers.PhraseLink{
			ID: int(ID.Int32),
			Phrase: triggers.Phrase{
				ID:                int(PhraseID.Int32),
				Phrase:            Phrase.String,
				NotifyOnDetection: NotifyOnDetection.Bool,
				WordOnlyMatch:     WordOnlyMatch.Bool,
			},
			GuildID:       GuildID.String,
			AddedByUserID: AddedByUserID.String,
			AddedDateTime: AddedDateTime.Time,
		})
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return phraseLinks, nil

}

func GetTriggerPhrase(guildId string, phrase string) (triggers.DbTriggerPhrase, error) {
	var triggerPhrase = triggers.DbTriggerPhrase{}

	// 1. Is it in the Database?
	var ID sql.NullInt32
	var Phrase sql.NullString

	err := Db.QueryRow("SELECT ID, Phrase FROM TriggerPhrases WHERE Phrase = ?", phrase).Scan(&ID, &Phrase)
	noRows := false
	if err != nil {
		if err != sql.ErrNoRows {
			logger.Error(guildId, err)
			return triggerPhrase, err
		} else {
			noRows = true
		}
	}

	if !noRows {
		if !ID.Valid {
			err = errors.New("db id was not valid")
			logger.Error(guildId, err)
			return triggerPhrase, err
		}

		if !Phrase.Valid {
			err = errors.New("db phrase was not valid")
			logger.Error(guildId, err)
			return triggerPhrase, err
		}

		triggerPhrase.ID = int(ID.Int32)
		triggerPhrase.Phrase = Phrase.String

		return triggerPhrase, nil
	}

	// 2. Not in Database, create
	stmt, err := Db.Prepare("INSERT INTO TriggerPhrases (Phrase) VALUES (?)")
	if err != nil {
		logger.Error(guildId, err)
		return triggerPhrase, err
	}
	defer func(g string) {
		err := stmt.Close()
		if err != nil {
			logger.Error(g, err)
		}
	}(guildId)

	res, err := stmt.Exec(phrase)
	if err != nil {
		logger.Error(guildId, err)
		return triggerPhrase, err
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		logger.Error(guildId, err)
		return triggerPhrase, err
	} else if lastInsertID == 0 {
		logger.Error(guildId, errors.New("last insert id was 0"))
		return triggerPhrase, err
	}

	triggerPhrase.ID = int(lastInsertID)
	triggerPhrase.Phrase = phrase
	return triggerPhrase, nil
}

func DoesPhraseLinkExist(guildId string, phraseId int) (bool, error) {
	query := `SELECT COUNT(*) AS Count FROM DiscGo.TriggerGuildLink WHERE GuildID = ? AND PhraseID = ?`
	var count sql.NullInt32

	err := Db.QueryRow(query, guildId, phraseId).Scan(&count)
	if err != nil {
		logger.Error(guildId, err)
		return false, err
	}

	if !count.Valid {
		err = errors.New("invalid count returned from db")
		logger.Error(guildId, err)
		return false, err
	}

	return count.Int32 > 0, nil
}

func InsertPhraseGuildLink(phraseId int, guildId string, userId string, notifyOnDetection bool, wordOnlyMatch bool) error {
	query := `INSERT INTO TriggerGuildLink (PhraseID, GuildID, AddedByUserID, NotifyOnDetection, WordOnlyMatch)
          SELECT * FROM (SELECT ? AS PhraseID, ? AS GuildID, ? AS AddedByUserID, ? AS NotifyOnDetection, ? AS WordOnlyMatch) AS tmp
          WHERE NOT EXISTS (
              SELECT ID FROM TriggerGuildLink WHERE PhraseID = ? AND GuildID = ?
          ) LIMIT 1`

	_, err := Db.Exec(query, phraseId, guildId, userId, notifyOnDetection, wordOnlyMatch, phraseId, guildId)
	if err != nil {
		logger.Error(guildId, err)
		return err
	}

	return nil
}

func UpdatePhraseGuildLink(phraseId int, guildId string, newNotify int, newWordOnly int, newDelete int) error {
	// Delete
	if newDelete == 1 {
		query := `DELETE FROM TriggerGuildLink WHERE PhraseID = ? AND GuildID = ?`
		_, err := Db.Exec(query, phraseId, guildId)
		if err != nil {
			logger.Error(guildId, err)
			return err
		}

		return nil
	}

	// Update
	query := "UPDATE TriggerGuildLink SET "
	updates := false
	if newNotify > -1 {
		query += fmt.Sprintf("NotifyOnDetection = %d", newNotify)
		updates = true
	}
	if newWordOnly > -1 {
		if updates {
			query += ", "
		}
		query += fmt.Sprintf("WordOnlyMatch = %d", newWordOnly)
		updates = true
	}
	if !updates {
		return errors.New("no updates to perform")
	}

	query += " WHERE PhraseID = ? AND GuildID = ?"
	_, err := Db.Exec(query, phraseId, guildId)
	if err != nil {
		logger.Error(guildId, err)
		return err
	}

	return nil
}

func GetPhraseLeaderboard(guildId string, phrase string) ([]PhraseLeaderboardUser, error) {
	var ranks []PhraseLeaderboardUser
	query := `	SELECT 
					C.UserID, C.Count
				FROM 
					DiscGo.TriggerGuildLink L 
				INNER JOIN
					DiscGo.TriggerPhrases T ON T.ID = L.PhraseID
				INNER JOIN
					DiscGo.CommandLog C ON C.GuildID = L.GuildID
					AND C.CommandTypeID = 2
					AND C.Command = T.Phrase
				WHERE 
					L.GuildID = ?
					AND T.Phrase = ? 
				ORDER BY
					C.Count DESC
				LIMIT 6
				`
	rows, err := Db.Query(query, guildId, phrase)
	if err != nil {
		logger.Error(guildId, err)
		return nil, err
	}
	defer func(g string) {
		err := rows.Close()
		if err != nil {
			logger.Error(g, err)
		}
	}(guildId)

	// Iterate over the rows
	for rows.Next() {

		var UserID sql.NullString
		var Count sql.NullInt32

		err := rows.Scan(&UserID, &Count)
		if err != nil {
			return nil, err
		}

		if !UserID.Valid {
			logger.Error(guildId, errors.New("invalid userid"))
			continue
		}

		if !Count.Valid {
			logger.Error(guildId, errors.New("invalid db wordOnlyMatch"))
			continue
		}

		ranks = append(ranks, PhraseLeaderboardUser{
			UserID: UserID.String,
			Count:  int(Count.Int32),
		})
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ranks, nil

}
