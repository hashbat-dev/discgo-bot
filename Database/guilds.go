package database

import (
	"context"
	"database/sql"
	"errors"

	logger "github.com/dabi-ngin/discgo-bot/Logger"
)

func Guild_DoesGuildExist(GuildID string) (bool, error) {

	var gif sql.NullInt32
	err := db.QueryRow("SELECT COUNT(*) FROM Guilds WHERE GuildID = ?", GuildID).Scan(&gif)
	if err != nil {
		return false, err
	}

	if !gif.Valid {
		return false, nil
	} else {
		return gif.Int32 > 0, nil
	}

}

func Guild_UpdateMemberCount(GuildID string, MemberCount int) error {

	query := "UPDATE Guilds SET GuildMemberCount = ?, GuildMemberCountLastCheck = NOW() WHERE GuildID = ?"
	_, err := db.ExecContext(context.Background(), query, MemberCount, GuildID)
	if err != nil {
		return err
	} else {
		return nil
	}

}

func Guild_InsertNewEntry(GuildID string, GuildName string, MemberCount int, OwnerID string) error {
	query := "INSERT INTO Guilds (GuildID, GuildName, GuildMemberCount, GuildOwnerID) VALUES (?, ?, ?, ?)"
	insertResult, err := db.ExecContext(context.Background(), query, GuildID, GuildName, MemberCount, OwnerID)
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

	logger.Event(GuildID, "New Guild entry created: %v (%v members)", GuildName, MemberCount)
	return nil
}
