package module

import (
	"github.com/bwmarrin/discordgo"
	config "github.com/hashbat-dev/discgo-bot/Config"
	discord "github.com/hashbat-dev/discgo-bot/Discord"
)

type EditImage struct{}

func (s EditImage) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "support",
		Description: "Get Support for Discgo Bot and its features",
	}
}

func (s EditImage) PermissionRequirement() int {
	return config.CommandLevelUser
}

func (s EditImage) Complexity() int {
	return config.TRIVIAL_TASK
}

func (s EditImage) Execute(i *discordgo.InteractionCreate, correlationId string) {
	discord.Interactions_SendMessage(i, "Help & Support", "For help and support on Discgo Bot please contact us.")
}
