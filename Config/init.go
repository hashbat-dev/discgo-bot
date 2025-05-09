package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

var InitComplete bool = false

// These can be swapped around on the go, but pls don't lol. If any are added make sure to also update the map
const (
	LoggingLevelAdmin = iota
	LoggingLevelError
	LoggingLevelWarn
	LoggingLevelEvent
	LoggingLevelInfo
	LoggingLevelDebug
)

var LoggingLevels map[int]LoggingOptions = map[int]LoggingOptions{
	LoggingLevelAdmin: {
		Name:   "Admin",
		Colour: Colours["magenta"],
	},
	LoggingLevelError: {
		Name:   "Error",
		Colour: Colours["red"],
	},
	LoggingLevelWarn: {
		Name:   "Warn",
		Colour: Colours["yellow"],
	},
	LoggingLevelEvent: {
		Name:   "Event",
		Colour: Colours["green"],
	},
	LoggingLevelInfo: {
		Name:   "Info",
		Colour: Colours["white"],
	},
	LoggingLevelDebug: {
		Name:   "Debug",
		Colour: Colours["blue"],
	},
}

type LoggingOptions struct {
	Name   string
	Colour Colour
}

type Colour struct {
	Terminal         string `json:"Terminal,omitempty"`
	Html             string `json:"Html,omitempty"`
	GraphOpaque      string `json:"GraphOpaque,omitempty"`
	GraphTransparent string `json:"GraphTransparent,omitempty"`
}

var Colours map[string]Colour = map[string]Colour{
	"default": {
		Terminal:         "\033[0m",
		Html:             "#000000",
		GraphOpaque:      "rgba(0, 0, 0, 1)",
		GraphTransparent: "rgba(0, 0, 0, 0.2)",
	},
	"white": {
		Terminal:         "\033[97m",
		Html:             "#FFFFFF",
		GraphOpaque:      "rgba(, , , 1)",
		GraphTransparent: "rgba(, , , 0.2)",
	},
	"magenta": {
		Terminal:         "\033[35m",
		Html:             "#C30CC9",
		GraphOpaque:      "rgba(195, 12, 201, 1)",
		GraphTransparent: "rgba(195, 12, 201, 0.2)",
	},
	"yellow": {
		Terminal:         "\033[33m",
		Html:             "#FAF200",
		GraphOpaque:      "rgba(250, 242, 0, 1)",
		GraphTransparent: "rgba(250, 242, 0, 0.2)",
	},
	"green": {
		Terminal:         "\033[32m",
		Html:             "#28F200",
		GraphOpaque:      "rgba(40, 242, 0, 1)",
		GraphTransparent: "rgba(40, 242, 0, 0.2)",
	},
	"red": {
		Terminal:         "\033[31m",
		Html:             "#FF9EA0",
		GraphOpaque:      "rgba(242, 0, 8, 1)",
		GraphTransparent: "rgba(242, 0, 8, 0.2)",
	},
	"blue": {
		Terminal:         "\033[34m",
		Html:             "#25B7FF",
		GraphOpaque:      "rgba(0, 0, 255, 1)",
		GraphTransparent: "rgba(0, 0, 255, 0.2)",
	},
}

const (
	CommandLevelUser = iota
	CommandLevelBotAdmin
	CommandLevelServerOwner
)

// Task categories for channels in message handling
const (
	// TRIVIAL_TASK involves small CPU and no IO waiting
	TRIVIAL_TASK = iota
	// CPU_BOUND_TASK involves intensive operations
	CPU_BOUND_TASK
	// IO_BOUND_TASK involves waiting on API/DB response
	IO_BOUND_TASK
)

const (
	N_TRIVIAL_WORKERS = 50
	N_IO_WORKERS      = 5
)

// Command Types
// This is used to denote types to the Dashboard
// ------------------------------------------------
const ( // A) Populate A + B!
	CommandTypeDefault = iota
	CommandTypeBang
	CommandTypePhrase
	CommandTypeSlash
	CommandTypeSlashResponse
	CommandTypeReactionCheck
)

// B) Populate A + B!
var CommandTypes map[int]string = map[int]string{
	CommandTypeDefault:       "Default",
	CommandTypeBang:          "Bang",
	CommandTypePhrase:        "Phrase",
	CommandTypeSlash:         "Slash",
	CommandTypeSlashResponse: "Slash Response",
	CommandTypeReactionCheck: "Reaction",
}

// ------------------------------------------------
// Process Pools
// Used to dispatch BangCommands in the newMessage Handler
// ------------------------------------------------
const (
	ProcessPoolText = iota
	ProcessPoolImages
	ProcessPoolExternal
)

var LastPoolIota int = ProcessPoolExternal

var ProcessPools map[int]ProcessPool = map[int]ProcessPool{
	ProcessPoolText: {
		ProcessPoolIota: ProcessPoolText,
		PoolName:        "Text",
		MaxWorkers:      50,
	},
	ProcessPoolImages: {
		ProcessPoolIota: ProcessPoolImages,
		PoolName:        "Images",
		MaxWorkers:      25,
	},
	ProcessPoolExternal: {
		ProcessPoolIota: ProcessPoolExternal,
		PoolName:        "External",
		MaxWorkers:      10,
	},
}

// -------------------------------------------------
type ProcessPool struct {
	ProcessPoolIota int
	PoolName        string
	MaxWorkers      int
}

type Vars struct {
	SESSIONID   string
	ISDEV       bool
	SUPERADMINS []string

	LOGTODISCORD       bool
	LOGGINGCHANNELID   string
	LOGGINGUSESTHREADS bool
	VERBOSESTACK       bool
	LOGFUNCTIONS       bool
	LOGGINGLEVEL       int

	DASHBOARDMAXDATAPACKETS     int
	DASHBOARDMAXLOGS            int
	DASHBOARDMAXCOMMANDS        int
	DASHBOARDURL                string
	COMMANDAVERAGEPOOL          int
	HARDWARESTATINTERVALSECONDS int
	HARDWARESTATMAXINTERVALS    int
	HOSTNAME                    string

	BOTTOKEN string

	DB_NAME       string
	DB_USER       string
	DB_PASSWORD   string
	DB_IP_ADDRESS string
	DB_PORT       string

	MAXFAKEYOUREQUESTCHECKS int
	MAXFAKEYOUREQUESTERRORS int
	IMGURCLIENTID           string

	TEMPFILEEXPIRYMINS int
	TEMPFILEGRACE      bool
	WOWRETENTIONMINS   int
}

var ServiceSettings Vars

var (
	Logginglevel         string
	Bottoken             string
	Session              *discordgo.Session
	ValidImageExtensions []string
	UserBangHelpText     string
	UserSlashHelpText    string
	UserMsgCmdHelpText   string
)

const (
	MAX_SELECT_LENGTH  int    = 25
	MAX_MESSAGE_LENGTH int    = 2000
	ROOT_FOLDER        string = "discgo-bot/"
	BOT_SUB_FOLDER     string = "Bot/"
	TEMP_FOLDER        string = "temp/"
)

func init() {
	enverr := godotenv.Load(".env")
	if enverr != nil {
		fmt.Println("no .env file found, checking if environment variables already set")
	}

	envVarsErr := parseEnvVariables()
	if envVarsErr != nil {
		errorMsg := fmt.Sprintf("unable to parse environment variables :: %s", envVarsErr.Error())
		panic(errorMsg)
	}
	fmt.Println(ServiceSettings)

	currentHostName, err := os.Hostname()
	if err != nil {
		ServiceSettings.HOSTNAME = "Unknown"
	} else {
		ServiceSettings.HOSTNAME = currentHostName
	}

	ValidImageExtensions = []string{
		".gif",
		".png",
		".jpg",
		".webp",
	}

	ServiceSettings.SESSIONID = uuid.New().String()
}

func parseEnvVariables() error {
	// setting a base convError value we'll reuse across each attempted conversion.
	// shouldn't need to worry about shadowing variables as this will return out at first instance of nil
	var convErr error
	isdev := os.Getenv("ISDEV")
	if isdev == "" {
		return errors.New("could not find value for ISDEV in environment variables")
	}
	ServiceSettings.ISDEV, convErr = strconv.ParseBool(isdev)
	if convErr != nil {
		return convErr
	}

	if ServiceSettings.ISDEV {
		ServiceSettings.DASHBOARDURL = "http://localhost:3333/"
	}

	superadminstr := os.Getenv("SUPERADMINS")
	if superadminstr == "" {
		return errors.New("could not find value for SUPERADMINS in environment variables")
	}
	ServiceSettings.SUPERADMINS = strings.Split(superadminstr, ",")
	if len(ServiceSettings.SUPERADMINS) == 0 {
		return errors.New("could not set any values for SUPERADMINS")
	}

	ServiceSettings.LOGGINGCHANNELID = os.Getenv("LOGGINGCHANNELID")
	if ServiceSettings.LOGGINGCHANNELID == "" {
		return errors.New("could not find value for LOGGINGCHANNELID in environment variables")
	}

	loggingusesthreads := os.Getenv("LOGGINGUSESTHREADS")
	if loggingusesthreads == "" {
		return errors.New("could not find value for LOGGINGUSESTHREADS in environment variables")
	}
	ServiceSettings.LOGGINGUSESTHREADS, convErr = strconv.ParseBool(loggingusesthreads)
	if convErr != nil {
		return convErr
	}

	verbosestack := os.Getenv("VERBOSESTACK")
	if verbosestack == "" {
		return errors.New("could not find value for VERBOSESTACK in environment variables")
	}
	ServiceSettings.VERBOSESTACK, convErr = strconv.ParseBool(verbosestack)
	if convErr != nil {
		return convErr
	}

	logfunctions := os.Getenv("LOGFUNCTIONS")
	if logfunctions == "" {
		return errors.New("could not find value for LOGFUNCTIONS in environment variables")
	}
	ServiceSettings.LOGFUNCTIONS, convErr = strconv.ParseBool(logfunctions)
	if convErr != nil {
		return convErr
	}

	dashboardmaxlogs := os.Getenv("DASHBOARDMAXLOGS")
	if dashboardmaxlogs == "" {
		return errors.New("could not find value for DASHBOARDMAXLOGS in environment variables")
	}
	ServiceSettings.DASHBOARDMAXLOGS, convErr = strconv.Atoi(dashboardmaxlogs)
	if convErr != nil {
		return convErr
	}

	dashboardmaxcommands := os.Getenv("DASHBOARDMAXCOMMANDS")
	if dashboardmaxcommands == "" {
		return errors.New("could not find value for DASHBOARDMAXCOMMANDS in environment variables")
	}
	ServiceSettings.DASHBOARDMAXCOMMANDS, convErr = strconv.Atoi(dashboardmaxcommands)
	if convErr != nil {
		return convErr
	}

	commandaveragepool := os.Getenv("COMMANDAVERAGEPOOL")
	if commandaveragepool == "" {
		return errors.New("could not find value for COMMANDAVERAGEPOOL in environment variables")
	}
	ServiceSettings.COMMANDAVERAGEPOOL, convErr = strconv.Atoi(commandaveragepool)
	if convErr != nil {
		return convErr
	}

	hardwarestatintervalseconds := os.Getenv("HARDWARESTATINTERVALSECONDS")
	if hardwarestatintervalseconds == "" {
		return errors.New("could not find value for HARDWARESTATINTERVALSECONDS in environment variables")
	}
	ServiceSettings.HARDWARESTATINTERVALSECONDS, convErr = strconv.Atoi(hardwarestatintervalseconds)
	if convErr != nil {
		return convErr
	}

	hardwarestatmaxintervals := os.Getenv("HARDWARESTATMAXINTERVALS")
	if hardwarestatmaxintervals == "" {
		return errors.New("could not find value for HARDWARESTATMAXINTERVALS in environment variables")
	}
	ServiceSettings.HARDWARESTATMAXINTERVALS, convErr = strconv.Atoi(hardwarestatmaxintervals)
	if convErr != nil {
		return convErr
	}

	ServiceSettings.BOTTOKEN = os.Getenv("BOTTOKEN")
	if ServiceSettings.BOTTOKEN == "" {
		return errors.New("could not find value for BOTTOKEN in environment variables")
	}

	ServiceSettings.DB_NAME = os.Getenv("DB_NAME")
	if ServiceSettings.DB_NAME == "" {
		return errors.New("could not find value for DB_NAME in environment variables")
	}

	ServiceSettings.DB_USER = os.Getenv("DB_USER")
	if ServiceSettings.DB_USER == "" {
		return errors.New("could not find value for DB_USER in environment variables")
	}

	ServiceSettings.DB_PASSWORD = os.Getenv("DB_PASSWORD")
	if ServiceSettings.DB_PASSWORD == "" {
		return errors.New("could not find value for DB_PASSWORD in environment variables")
	}

	ServiceSettings.DB_IP_ADDRESS = os.Getenv("DB_IP_ADDRESS")
	if ServiceSettings.DB_IP_ADDRESS == "" {
		return errors.New("could not find value for DB_IP_ADDRESS in environment variables")
	}

	ServiceSettings.DB_PORT = os.Getenv("DB_PORT")
	if ServiceSettings.DB_PORT == "" {
		return errors.New("could not find value for DB_PORT in environment variables")
	}

	maxfakeyourequestchecks := os.Getenv("MAXFAKEYOUREQUESTCHECKS")
	if maxfakeyourequestchecks == "" {
		return errors.New("could not find value for MAXFAKEYOUREQUESTCHECKS in environment variables")
	}
	ServiceSettings.MAXFAKEYOUREQUESTCHECKS, convErr = strconv.Atoi(maxfakeyourequestchecks)
	if convErr != nil {
		return convErr
	}

	maxfakeyourequesterrors := os.Getenv("MAXFAKEYOUREQUESTERRORS")
	if maxfakeyourequesterrors == "" {
		return errors.New("could not find value for MAXFAKEYOUREQUESTERRORS in environment variables")
	}
	ServiceSettings.MAXFAKEYOUREQUESTERRORS, convErr = strconv.Atoi(maxfakeyourequesterrors)
	if convErr != nil {
		return convErr
	}

	ServiceSettings.IMGURCLIENTID = os.Getenv("IMGURCLIENTID")
	if ServiceSettings.IMGURCLIENTID == "" {
		return errors.New("could not find value for IMGURCLIENTID in environment variables")
	}

	tempFileExpiryMins := os.Getenv("TEMPFILEEXPIRYMINS")
	if tempFileExpiryMins == "" {
		return errors.New("could not find value for TEMPFILEEXPIRYMINS in environment variables")
	}
	ServiceSettings.TEMPFILEEXPIRYMINS, convErr = strconv.Atoi(tempFileExpiryMins)
	if convErr != nil {
		return convErr
	}

	tempFileGrace := os.Getenv("TEMPFILEGRACE")
	if tempFileGrace == "" {
		return errors.New("could not find value for TEMPFILEGRACE in environment variables")
	}
	ServiceSettings.TEMPFILEGRACE, convErr = strconv.ParseBool(tempFileGrace)
	if convErr != nil {
		return convErr
	}

	tempWowRetentionMins := os.Getenv("WOWRETENTIONMINS")
	if tempWowRetentionMins == "" {
		fmt.Println("Applied default WOWRETENTIONMINS value of 60")
		tempWowRetentionMins = "60"
	}
	ServiceSettings.WOWRETENTIONMINS, convErr = strconv.Atoi(tempWowRetentionMins)
	if convErr != nil {
		return convErr
	}

	return nil
}
