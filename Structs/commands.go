package structs

import "github.com/bwmarrin/discordgo"

// Needs to be here to avoid Import Cycles
type BangCommand struct {
	Begin       func(message *discordgo.MessageCreate, self BangCommand) error
	ImgCategory string
}
