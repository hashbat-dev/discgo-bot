package module

import "github.com/bwmarrin/discordgo"

type Module interface {
	Command() *discordgo.ApplicationCommand
	Execute(*discordgo.InteractionCreate, string)
	PermissionRequirement() int
	Complexity() int
}

var (
	ModuleList []Module
)

func init() {
	ModuleList = []Module{
		AddImage{"speech"},
		AdminRole{},
		DeleteImage{"speech"},
		EditHallReactions{},
		Support{},
		TTSPlay{},
	}
}
