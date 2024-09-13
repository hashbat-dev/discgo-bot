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
	JumpTable["ttsinfo"] = bang.TTSInfo{}
	JumpTable["flipleft"] = bang.FlipImage{FlipDirection: "left"}
	JumpTable["flipright"] = bang.FlipImage{FlipDirection: "right"}
	JumpTable["flipup"] = bang.FlipImage{FlipDirection: "up"}
	JumpTable["flipdown"] = bang.FlipImage{FlipDirection: "down"}
	JumpTable["flipboth"] = bang.FlipImage{FlipDirection: "both"}
	JumpTable["flipall"] = bang.FlipImage{FlipDirection: "all"}
	JumpTable["reverse"] = bang.Reverse{}
	JumpTable["speedup"] = bang.ChangeSpeed{SpeedUp: true}
	JumpTable["slowdown"] = bang.ChangeSpeed{SpeedUp: false}
}
