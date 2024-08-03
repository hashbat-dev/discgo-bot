package wow

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type WowStatItem struct {
	MessageID   string
	Rolls       []DiceRoll
	MiddleCount int
	UserID      string
	Added       time.Time
}

type DiceRoll struct {
	RollContinue  int
	RollLength    int
	RollSpecial   []int
	Special       string
	SpecialString string
}

var WowStatCache []WowStatItem

var wowSpamCache []WowSpamCache

type WowSpamCache struct {
	UserID         string
	LastUsed       time.Time
	SessionCount   int
	SessionStart   time.Time
	SessionApplied []string
}

type WowEffect struct {
	UserID               string
	EffectName           string
	EffectDescription    string
	ActiveUntil          time.Time
	ContinueModifier     int
	ContinueRollModifier int
	WowRollModifier      int
	FreeRolls            int
	TempFreeRolls        int
	SessionBased         bool
}

var wowEffects []WowEffect

type WowSendCache struct {
	ChannelID   string
	ReplyText   string
	MessageRef  discordgo.MessageReference
	WowRolls    []DiceRoll
	WowCount    int
	UserID      string
	MessageSent bool
	ErrorCount  int
}

var SendCache []WowSendCache
