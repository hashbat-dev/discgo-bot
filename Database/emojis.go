package database

import (
	"database/sql"
	"errors"
	"strings"

	data "github.com/hashbat-dev/discgo-bot/Data"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

func GetAllGuildEmojis(guildId string) ([]data.GuildEmoji, error) {
	query := `SELECT
				L.ID, L.EmojiID, L.CategoryID, L.GuildID, L.AddedByUserID, 
    			L.AddedDateTime, C.EmojiCategory, S.Emoji
			  FROM 
			  	EmojiGuildLink L
			  INNER JOIN 
			  	EmojiCategories C ON C.ID = L.CategoryID
			  INNER JOIN 
			  	EmojiStorage S ON S.ID = L.EmojiID
			  WHERE
			  	L.GuildID = ?`
	rows, err := Db.Query(query, guildId)
	if err != nil {
		logger.Error(guildId, err)
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			logger.Error("DATABASE", err)
		}
	}()

	// Iterate over the rows
	var emojiList []data.GuildEmoji
	for rows.Next() {
		var ID, EmojiID, CategoryID sql.NullInt32
		var GuildID, AddedByUserID, EmojiCategory, Emoji sql.NullString
		var AddedDateTime sql.NullTime
		err := rows.Scan(&ID, &EmojiID, &CategoryID, &GuildID, &AddedByUserID, &AddedDateTime, &EmojiCategory, &Emoji)
		if err != nil {
			return nil, err
		}

		var r data.GuildEmoji
		if !ID.Valid {
			err = errors.New("invalid db id")
			logger.Error(guildId, err)
			continue
		} else {
			r.ID = int(ID.Int32)
		}
		if EmojiID.Valid {
			r.EmojiID = int(EmojiID.Int32)
		}
		if CategoryID.Valid {
			r.CategoryID = int(CategoryID.Int32)
		}
		if GuildID.Valid {
			r.GuildID = GuildID.String
		}
		if AddedByUserID.Valid {
			r.AddedByUserID = AddedByUserID.String
		}
		if AddedDateTime.Valid {
			r.AddedDateTime = AddedDateTime.Time
		}
		if EmojiCategory.Valid {
			r.EmojiCategory = EmojiCategory.String
		}
		if Emoji.Valid {
			r.Emoji = Emoji.String
		}
		emojiList = append(emojiList, r)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return emojiList, nil
}

func GetEmojiStorageID(guildId string, emoji string) (int, error) {
	//	=> Is it in the Database?
	emoji = strings.TrimSpace(emoji)
	var ID sql.NullInt32
	err := Db.QueryRow("SELECT ID FROM EmojiStorage WHERE Emoji = ?", emoji).Scan(&ID)
	noRows := false
	if err != nil {
		if err != sql.ErrNoRows {
			logger.Error(guildId, err)
			return 0, err
		} else {
			noRows = true
		}
	}

	if !noRows {
		if !ID.Valid {
			err = errors.New("db id was not valid")
			logger.Error(guildId, err)
			return 0, err
		}

		return int(ID.Int32), nil
	}

	// => Not in Database, create
	stmt, err := Db.Prepare("INSERT INTO EmojiStorage (Emoji) VALUES (?)")
	if err != nil {
		logger.Error(guildId, err)
		return 0, err
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
			logger.Error("DATABASE", err)
		}
	}()

	res, err := stmt.Exec(emoji)
	if err != nil {
		logger.Error(guildId, err)
		return 0, err
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		logger.Error(guildId, err)
		return 0, err
	}

	return int(lastInsertID), nil
}

func InsertEmojiGuildLink(emojiId int, categoryId int, guildId string, addedByUserId string) error {
	query := `INSERT INTO EmojiGuildLink (EmojiID, CategoryID, GuildID, AddedByUserID)
          SELECT * FROM (SELECT ? AS EmojiID, ? AS CategoryID, ? AS GuildID, ? AS AddedByUserID) AS tmp
          WHERE NOT EXISTS (
              SELECT ID FROM EmojiGuildLink WHERE EmojiID = ? AND GuildID = ?
          ) LIMIT 1`
	_, err := Db.Exec(query, emojiId, categoryId, guildId, addedByUserId, emojiId, guildId)
	if err != nil {
		logger.Error(guildId, err)
		return err
	}

	return nil
}

func DeleteAllEmojiLinks(guildId string) error {
	query := `DELETE FROM
				EmojiGuildLink
			  WHERE 
			  	GuildID = ?`
	_, err := Db.Exec(query, guildId)
	if err != nil {
		logger.Error(guildId, err)
		return err
	}
	return TidyEmojiStorage(guildId)
}

func TidyEmojiStorage(guildId string) error {
	queryStorage := `SELECT * FROM EmojiStorage
	WHERE ID NOT IN (
		SELECT EmojiID FROM EmojiGuildLink
	)`
	_, err := Db.Exec(queryStorage)
	if err != nil {
		logger.Error(guildId, err)
		return err
	}

	return err
}
