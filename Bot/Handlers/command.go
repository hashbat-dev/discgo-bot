package handlers

import (
	"github.com/bwmarrin/discordgo"
	bang "github.com/dabi-ngin/discgo-bot/Bot/Commands/Bang"
	config "github.com/dabi-ngin/discgo-bot/Config"
)

type Command interface {
	Name() string
	Execute(*discordgo.MessageCreate, string) error
	PermissionRequirement() int
	ProcessPool() config.ProcessPool
	LockedByDefault() bool
}

var (
	Commands = make(map[string]Command)
)

func init() {
	Commands["speech"] = bang.GetImage{}
	Commands["addspeech"] = bang.AddImage{}
	Commands["delspeech"] = bang.DelImage{}
	Commands["makespeech"] = bang.MakeSpeech{}
}
