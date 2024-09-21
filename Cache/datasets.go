package cache

import (
	"time"

	triggers "github.com/dabi-ngin/discgo-bot/Bot/Commands/Triggers"
)

type Guild struct {
	DbID            int
	DiscordID       string
	Name            string
	CommandCount    int
	IsDev           bool
	LastCommand     time.Time
	Triggers        []triggers.Phrase
	StarUpChannel   string
	StarDownChannel string
}

type GuildPermissions struct {
	CommandType  int
	RequiredRole string
}
