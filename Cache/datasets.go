package cache

import (
	"time"

	triggers "github.com/dabi-ngin/discgo-bot/Bot/Commands/Triggers"
)

var Commands []Command
var CommandInfo []CmdInfo
var CommandAverages []CmdAverage

type Command struct {
	TypeID       int
	Command      string
	GuildID      string
	UserID       string
	UserName     string
	CallTime     time.Time
	CallDuration time.Duration
}

type CmdInfo struct {
	TypeID      int
	Command     string
	Count       int
	AvgDuration time.Duration
	LastCall    time.Time
}

type CmdAverage struct {
	TypeID      int
	Command     string
	Durations   []time.Duration
	AvgDuration time.Duration
}

type Guild struct {
	DbID         int
	DiscordID    string
	Name         string
	CommandCount int
	LastCommand  time.Time
	Triggers     []triggers.Phrase
}

type GuildPermissions struct {
	CommandType  int
	RequiredRole string
}
