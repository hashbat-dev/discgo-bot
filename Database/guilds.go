package database

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	logger "github.com/dabi-ngin/discgo-bot/Logger"
)

type GuildInfo struct {
	ID               int
	GuildID          string
	GuildName        string
	GuildMemberCount int
	GuildOwnerID     string
	GuildAdminRole   string
	StarUpChannel    string
	StarDownChannel  string
	UpdatedDateTime  time.Time
	CreatedDateTime  time.Time
	IsDevServer      bool
}

func Guild_Get(guildId string) (GuildInfo, error) {
	var ID, GuildMemberCount sql.NullInt32
	var IsDevServer sql.NullBool
	var GuildID, GuildName, GuildOwnerID, GuildAdminRole, StarUpChannel, StarDownChannel sql.NullString
	var UpdatedDateTime, CreatedDateTime sql.NullTime
	var r GuildInfo
	query := `SELECT
					ID, GuildID, GuildName, GuildMemberCount, GuildOwnerID, GuildAdminRole,
					StarUpChannel, StarDownChannel, UpdatedDateTime, CreatedDateTime, IsDevServer
				FROM
					Guilds
				WHERE
					GuildID = ?	`

	err := Db.QueryRow(query, guildId).Scan(&ID, &GuildID, &GuildName, &GuildMemberCount, &GuildOwnerID, &GuildAdminRole,
		&StarUpChannel, &StarDownChannel, &UpdatedDateTime, &CreatedDateTime, &IsDevServer)
	if err != nil {
		if !strings.Contains(err.Error(), "no rows") {
			logger.Error(guildId, err)
			err = nil
		}
		return r, err
	}

	if !ID.Valid {
		return r, err
	} else {
		r.ID = int(ID.Int32)
	}

	if GuildID.Valid {
		r.GuildID = GuildID.String
	}
	if GuildName.Valid {
		r.GuildName = GuildName.String
	}
	if GuildMemberCount.Valid {
		r.GuildMemberCount = int(GuildMemberCount.Int32)
	}
	if GuildOwnerID.Valid {
		r.GuildOwnerID = GuildOwnerID.String
	}
	if GuildAdminRole.Valid {
		r.GuildAdminRole = GuildAdminRole.String
	}
	if StarUpChannel.Valid {
		r.StarUpChannel = StarUpChannel.String
	}
	if StarDownChannel.Valid {
		r.StarDownChannel = StarDownChannel.String
	}
	if UpdatedDateTime.Valid {
		r.UpdatedDateTime = UpdatedDateTime.Time
	}
	if CreatedDateTime.Valid {
		r.CreatedDateTime = CreatedDateTime.Time
	}
	if IsDevServer.Valid {
		r.IsDevServer = IsDevServer.Bool
	}

	return r, nil
}

func Guild_InsertUpdate(g GuildInfo) (GuildInfo, error) {
	if g.ID == 0 {
		// Insert
		var params []interface{}
		colList := "("
		valList := "("

		colList += "GuildID"
		valList += "?"
		params = append(params, g.GuildID)

		colList += ", GuildName"
		valList += ", ?"
		params = append(params, g.GuildName)

		colList += ", GuildMemberCount"
		valList += ", ?"
		params = append(params, g.GuildMemberCount)

		colList += ", GuildOwnerID"
		valList += ", ?"
		params = append(params, g.GuildOwnerID)

		if g.GuildAdminRole != "" {
			colList += ", GuildAdminRole"
			valList += ", ?"
			params = append(params, g.GuildAdminRole)
		}

		if g.StarUpChannel != "" {
			colList += ", StarUpChannel"
			valList += ", ?"
			params = append(params, g.StarUpChannel)
		}

		if g.StarDownChannel != "" {
			colList += ", StarDownChannel"
			valList += ", ?"
			params = append(params, g.StarDownChannel)
		}

		colList += ", UpdatedDateTime"
		valList += ", NOW()"
		colList += ", CreatedDateTime"
		valList += ", NOW()"

		colList += ", IsDevServer"
		valList += ", ?"
		params = append(params, g.IsDevServer)

		colList += ")"
		valList += ")"

		query := "INSERT INTO Guilds " + colList + " VALUES " + valList
		insertResult, err := Db.ExecContext(context.Background(), query, params...)
		if err != nil {
			logger.Error(g.GuildID, err)
			return g, err
		}
		id, err := insertResult.LastInsertId()
		if err != nil {
			logger.Error(g.GuildID, err)
			return g, err
		} else if id == 0 {
			err = errors.New("starboard insert returned id = 0")
			logger.Error(g.GuildID, err)
			return g, err
		}
		g.ID = int(id)
		return g, nil
	} else {
		// Update
		var params []interface{}
		setList := ""

		setList += "GuildID = ?"
		params = append(params, g.GuildID)

		setList += ", GuildName = ?"
		params = append(params, g.GuildName)

		setList += ", GuildMemberCount = ?"
		params = append(params, g.GuildMemberCount)

		setList += ", GuildOwnerID = ?"
		params = append(params, g.GuildOwnerID)

		if g.GuildAdminRole != "" {
			setList += ", GuildAdminRole = ?"
			params = append(params, g.GuildAdminRole)
		} else {
			setList += ", GuildAdminRole = NULL"
		}

		if g.StarUpChannel != "" {
			setList += ", StarUpChannel = ?"
			params = append(params, g.StarUpChannel)
		} else {
			setList += ", StarUpChannel = NULL"
		}

		if g.StarDownChannel != "" {
			setList += ", StarDownChannel = ?"
			params = append(params, g.StarDownChannel)
		} else {
			setList += ", StarDownChannel = NULL"
		}

		setList += ", UpdatedDateTime = NOW()"

		setList += ", IsDevServer = ?"
		params = append(params, g.IsDevServer)

		query := "UPDATE Guilds SET " + setList + " WHERE ID = ?"
		params = append(params, g.ID)
		_, err := Db.ExecContext(context.Background(), query, params...)
		if err != nil {
			logger.Error(g.GuildID, err)
			return g, err
		}
		return g, nil
	}
}
