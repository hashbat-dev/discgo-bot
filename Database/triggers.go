package database

import (
	"database/sql"
	"errors"

	triggers "github.com/hashbat-dev/discgo-bot/Bot/Commands/Triggers"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

func GetAllGuildPhrases(guildId string) ([]triggers.PhraseLink, error) {
	var phraseLinks []triggers.PhraseLink

	query := `	SELECT 
				L.ID, L.PhraseID, P.Phrase, L.GuildID, L.AddedByUserID, L.AddedDateTime, L.NotifyOnDetection, L.WordOnlyMatch
				FROM TriggerGuildLink L		
				INNER JOIN DiscGo.TriggerPhrases P ON P.ID = L.PhraseID
			  	WHERE L.GuildID = ?
			 `

	rows, err := Db.Query(query, guildId)
	if err != nil {
		logger.Error(guildId, err)
		return nil, err
	}
	defer rows.Close()

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
