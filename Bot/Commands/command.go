package commands

import (
	"github.com/bwmarrin/discordgo"
	bang "github.com/dabi-ngin/discgo-bot/Bot/Commands/Bang"
)

type Command interface {
	Name() string
	Execute(*discordgo.MessageCreate, string) error
	PermissionRequirement() int
	Complexity() int
}

var (
	JumpTable = make(map[string]Command)
)

func init() {
	JumpTable["speech"] = bang.GetImage{}
	JumpTable["addspeech"] = bang.AddImage{}
	JumpTable["delspeech"] = bang.DelImage{}
	JumpTable["makespeech"] = bang.MakeSpeech{}
}
