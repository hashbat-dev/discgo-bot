package bangCommands

import (
	"github.com/bwmarrin/discordgo"
	config "github.com/dabi-ngin/discgo-bot/Config"
)

type Reverse struct {
}

func (r Reverse) Name() string {
	return "reverse"
}

func (r Reverse) Execute(message *discordgo.MessageCreate) error {
	// do the funky image work
	// run some db query
	return nil
}

func (r Reverse) PermissionRequirement() int {
	return config.CommandLevelAdmin
}

func (r Reverse) IsComplex() bool {
	return true
}
