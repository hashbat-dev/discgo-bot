package database

import (
	"context"
	"database/sql"
	"errors"

	logger "github.com/dabi-ngin/discgo-bot/Logger"
)

func Guild_DoesGuildExist(GuildID string) (int, bool, error) {

	var ID sql.NullInt32
	var IsDev sql.NullBool
	err := Db.QueryRow("SELECT ID, IsDevServer FROM Guilds WHERE GuildID = ?", GuildID).Scan(&ID, &IsDev)
	if err != nil {
		return 0, false, err
	}

	if !ID.Valid {
		return 0, false, nil
	} else {
		isDev := false
		if IsDev.Valid {
			isDev = IsDev.Bool
		}
		return int(ID.Int32), isDev, nil
	}

}

func Guild_UpdateMemberCount(GuildID string, MemberCount int) error {

	query := "UPDATE Guilds SET GuildMemberCount = ?, GuildMemberCountLastCheck = NOW() WHERE GuildID = ?"
	_, err := Db.ExecContext(context.Background(), query, MemberCount, GuildID)
	if err != nil {
		return err
	} else {
		return nil
	}

}

func Guild_InsertNewEntry(GuildID string, GuildName string, MemberCount int, OwnerID string) (int, error) {
	query := "INSERT INTO Guilds (GuildID, GuildName, GuildMemberCount, GuildOwnerID) VALUES (?, ?, ?, ?)"
	insertResult, err := Db.ExecContext(context.Background(), query, GuildID, GuildName, MemberCount, OwnerID)
	if err != nil {
		return 0, err
	}

	id, err := insertResult.LastInsertId()
	if err != nil {
		return 0, err
	} else if id == 0 {
		err = errors.New("returned id insert was 0")
		return 0, err
	}

	logger.Event(GuildID, "New Guild entry created: %v (%v members)", GuildName, MemberCount)
	return int(id), nil
}
