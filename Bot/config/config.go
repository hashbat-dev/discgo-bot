package config

import (
	"encoding/json"
	"os"

	"github.com/bwmarrin/discordgo"
)

type Vars struct {
	IsDev           bool
	ReportedRunning bool
	ServerName      string
	BotToken        string
	BullyTarget     string
	FwidayCancelled bool

	// discordgo (cannot be a constant, we update this on startup, needs to exist here for scope reasons)
	Session *discordgo.Session

	DB_NAME       string
	DB_USER       string
	DB_PASSWORD   string
	DB_IP_ADDRESS string
	DB_PORT       string
}

var (
	// Runtime Environment ===========================================
	IsDev           bool   = false
	ReportedRunning bool   = false
	ServerName      string = ""
	BotToken        string = ""
	BullyTarget     string = ""
	FwidayCancelled bool   = false

	// discordgo (cannot be a constant, we update this on startup, needs to exist here for scope reasons)
	Session *discordgo.Session

	// THESE DON'T CHANGE AT RUNTIME
	// Database Config Settings
	DB_NAME       string = ""
	DB_USER       string = ""
	DB_PASSWORD   string = ""
	DB_IP_ADDRESS string = ""
	DB_PORT       string = ""
)

func SetVars() {

	fileText, err := os.ReadFile("Bot/config/config.JSON")

	if err != nil {
		panic("oops")
		//using audit here causes import cycle
	}

	var vars Vars

	jsonerr := json.Unmarshal([]byte(fileText), &vars)
	if jsonerr != nil {
		panic("oops")
		//using audit here causes import cycle
	}

	// Runtime Environment ===========================================
	IsDev = vars.IsDev
	ReportedRunning = vars.ReportedRunning
	ServerName = vars.ServerName
	BotToken = vars.BotToken
	BullyTarget = vars.BullyTarget
	FwidayCancelled = vars.FwidayCancelled

	// discordgo (cannot be a constant, we update this on startup, needs to exist here for scope reasons)
	Session = vars.Session

	// THESE DON'T CHANGE AT RUNTIME
	// Database Config Settings
	DB_NAME = vars.ServerName
	DB_USER = vars.DB_USER
	DB_PASSWORD = vars.DB_PASSWORD
	DB_IP_ADDRESS = vars.DB_IP_ADDRESS
	DB_PORT = vars.DB_PORT

}
