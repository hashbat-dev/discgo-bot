package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

// These can be swapped around on the go, but pls don't lol. If any are added make sure to also update
const (
	LoggingLevelAdmin = iota
	LoggingLevelError
	LoggingLevelWarn
	LoggingLevelEvent
	LoggingLevelInfo
	LoggingLevelDebug
)

// This is used to denote types to the Dashbaord
var LoggingLevels map[int]string = map[int]string{
	LoggingLevelAdmin: "Admin",
	LoggingLevelError: "Error",
	LoggingLevelWarn:  "Warn",
	LoggingLevelEvent: "Event",
	LoggingLevelInfo:  "Info",
	LoggingLevelDebug: "Debug",
}

const (
	CommandTypeDefault = iota
	CommandTypeBang
	CommandTypePhrase
)

// This is used to denote types to the Dashbaord
var CommandTypes map[int]string = map[int]string{
	CommandTypeDefault: "Default",
	CommandTypeBang:    "Bang",
	CommandTypePhrase:  "Phrase",
}

type Vars struct {
	IsDev       bool
	SuperAdmins []string

	LoggingChannelID   string
	LoggingUsesThreads bool
	VerboseStack       bool
	LogFunctions       bool
	LoggingLevel       int

	DashboardMaxLogs            int
	DashboardMaxCommands        int
	CommandAveragePool          int
	HardwareStatIntervalSeconds int
	HardwareStatMaxIntervals    int

	BotToken string

	DB_NAME       string
	DB_USER       string
	DB_PASSWORD   string
	DB_IP_ADDRESS string
	DB_PORT       string
}

var (
	IsDev       bool
	HostName    string
	SuperAdmins []string

	LoggingChannelID            string
	LoggingUsesThreads          bool
	LoggingVerboseStack         bool
	LoggingLogFunctions         bool
	DashboardMaxLogs            int
	DashboardMaxCommands        int
	CommandAveragePool          int
	HardwareStatIntervalSeconds int
	HardwareStatMaxIntervals    int

	LoggingLevel int

	BotToken string
	Session  *discordgo.Session

	DB_NAME       string
	DB_USER       string
	DB_PASSWORD   string
	DB_IP_ADDRESS string
	DB_PORT       string
)

const (
	MAX_MESSAGE_LENGTH int    = 2000
	ROOT_FOLDER        string = "discgo-bot/"
	BOT_SUB_FOLDER     string = "Bot/"
)

func Init() bool {
	localConfigFile, err := os.ReadFile("config.json")

	if err != nil {
		fmt.Println(fmt.Printf("Config.Init() - Error loading config.json :: %v", err))
		return false
	}

	var configFileVariables Vars
	err = json.Unmarshal([]byte(localConfigFile), &configFileVariables)
	if err != nil {
		fmt.Println(fmt.Printf("Config.Init() - Error unmarshalling config.json :: %v", err))
		return false
	}

	currentHostName, err := os.Hostname()
	if err != nil {
		HostName = "Unknown"
	} else {
		HostName = currentHostName
	}

	IsDev = configFileVariables.IsDev
	SuperAdmins = configFileVariables.SuperAdmins

	LoggingChannelID = configFileVariables.LoggingChannelID
	LoggingUsesThreads = configFileVariables.LoggingUsesThreads
	LoggingVerboseStack = configFileVariables.VerboseStack
	LoggingLogFunctions = configFileVariables.LogFunctions
	LoggingLevel = configFileVariables.LoggingLevel
	DashboardMaxLogs = configFileVariables.DashboardMaxLogs
	DashboardMaxCommands = configFileVariables.DashboardMaxCommands
	CommandAveragePool = configFileVariables.CommandAveragePool
	HardwareStatIntervalSeconds = configFileVariables.HardwareStatIntervalSeconds
	HardwareStatMaxIntervals = configFileVariables.HardwareStatMaxIntervals

	BotToken = configFileVariables.BotToken

	DB_NAME = configFileVariables.DB_NAME
	DB_USER = configFileVariables.DB_USER
	DB_PASSWORD = configFileVariables.DB_PASSWORD
	DB_IP_ADDRESS = configFileVariables.DB_IP_ADDRESS
	DB_PORT = configFileVariables.DB_PORT

	return true
}
