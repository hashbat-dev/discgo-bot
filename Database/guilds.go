package database

import (
	"database/sql"
	"strings"
	"time"

	logger "github.com/hashbat-dev/discgo-bot/Logger"
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

func Get(guildId string) (GuildInfo, error) {
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
			return GuildInfo{}, err
		} else {
			r.ID = 0
			return r, nil
		}
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

func GuildUpsert(guild GuildInfo) (GuildInfo, error) {
	params := []any{
		"GuildID", guild.GuildID,
		"GuildName", guild.GuildName,
		"GuildMemberCount", guild.GuildMemberCount,
		"GuildOwnerID", guild.GuildOwnerID,
	}
	if guild.GuildAdminRole != "" {
		params = append(params, "GuildAdminRole", guild.GuildAdminRole)
	}
	if guild.StarUpChannel != "" {
		params = append(params, "StarUpChannel", guild.StarUpChannel)
	}
	if guild.StarDownChannel != "" {
		params = append(params, "StarDownChannel", guild.StarDownChannel)
	}

	id, err := Upsert(guild.GuildID, "Guilds", "ID", guild.ID, params...)
	if guild.ID != int(id) {
		guild.ID = int(id)
	}
	return guild, err
}
