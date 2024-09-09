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

	// ===[Add Bang Commands]===========================================
	JumpTable["speech"] = bang.GetImage{ImageCategory: "speech"}
	JumpTable["addspeech"] = bang.AddImage{ImageCategory: "speech"}
	JumpTable["delspeech"] = bang.DelImage{ImageCategory: "speech"}
	JumpTable["makespeech"] = bang.MakeSpeech{}
}
