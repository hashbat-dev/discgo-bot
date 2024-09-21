package database

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	logger "github.com/dabi-ngin/discgo-bot/Logger"
)

// Returns
func Guild_GetGuildDBInfo(GuildID string) (int, bool, string, string) {
	var ID sql.NullInt32
	var IsDev sql.NullBool
	var starUp, starDown sql.NullString
	err := Db.QueryRow("SELECT ID, IsDevServer, StarUpChannel, StarDownChannel FROM Guilds WHERE GuildID = ?", GuildID).Scan(&ID, &IsDev, &starUp, &starDown)
	if err != nil {
		if !strings.Contains(err.Error(), "no rows") {
			logger.Error(GuildID, err)
		}
		return 0, false, "", ""
	}

	if !ID.Valid {
		return 0, false, "", ""
	} else {
		isDev := false
		if IsDev.Valid {
			isDev = IsDev.Bool
		}
		starUpId := ""
		starDownId := ""
		if starUp.Valid {
			starUpId = starUp.String
		}
		if starDown.Valid {
			starDownId = starDown.String
		}
		return int(ID.Int32), isDev, starUpId, starDownId
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

func Guild_UpdateStarboardChannel(GuildID string, ChannelID string, IsUp bool) error {
	channel := "StarUpChannel"
	if !IsUp {
		channel = "StarDownChannel"
	}
	query := "UPDATE Guilds SET " + channel + " = ? WHERE GuildID = ?"
	_, err := Db.ExecContext(context.Background(), query, ChannelID, GuildID)
	if err != nil {
		logger.Error(GuildID, err)
		return err
	} else {
		return nil
	}

}
