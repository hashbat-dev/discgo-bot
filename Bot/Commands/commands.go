package commands

import "github.com/bwmarrin/discordgo"

type Command interface {
	Name() string
	Execute(*discordgo.MessageCreate) error
	PermissionRequirement() int
	IsComplex() bool
}
