package config

import (
	"github.com/bwmarrin/discordgo"
)

var (
	// Runtime Environment ===========================================
	IsDev           bool   = false
	ReportedRunning bool   = false
	ServerName      string = ""
	BotToken        string
	BullyTarget     string = ""
	FwidayCancelled bool   = false

	// discordgo (cannot be a constant, we update this on startup, needs to exist here for scope reasons)
	Session *discordgo.Session
)

const (
	// Database Config Settings
	// TODO: Move this to be driven by command line args or env variable we set at bot startup in either dockerfile
	DB_NAME       string = "BotDB"
	DB_USER       string = "bot_access"
	DB_PASSWORD   string = "AVNS_vGqsdK_m6vK2JLODhPI"
	DB_IP_ADDRESS string = "db-mysql-nyc3-09006-do-user-16162189-0.c.db.ondigitalocean.com"
	DB_PORT       string = "25060"
)
