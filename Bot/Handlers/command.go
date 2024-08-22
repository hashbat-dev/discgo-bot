package handlers

import "github.com/bwmarrin/discordgo"

type Command interface {
	Name() string
	Execute(*discordgo.MessageCreate) error
	PermissionRequirement() int
}

func HandleMessage(command Command, message *discordgo.MessageCreate) error {
	command.Execute(message)
}
