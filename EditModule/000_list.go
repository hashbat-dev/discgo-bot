package editmodule

import "github.com/bwmarrin/discordgo"

type EditModule interface {
	SelectName() string
	Emoji() *discordgo.ComponentEmoji
	Execute(*discordgo.InteractionCreate, string)
	PermissionRequirement() int
	Complexity() int
}

var (
	EditList []EditModule
)

func init() {

	EditList = []EditModule{
		ChangeSpeed{},
		CreateMeme{},
		DeepFry{},
		FlipImage{}, // Done
		MakeSpeech{},
		Reverse{},
		WidenImage{},
	}
}
