package helpers

import (
	"github.com/bwmarrin/discordgo"
	config "github.com/dabi-ngin/discgo-bot/Config"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
)

const (
	AllAccess  = iota
	SuperAdmin = iota
	Admin      = iota
	Moderator  = iota
	NormalUser = iota
)

var commandPriveledges map[string]int

func InitPermissions() {
	commandPriveledges["test"] = 1
}

// DoesUserHavePermissionToUseCommand checks if the user is admin. Later we can change this to command-level permissios.
func DoesUserHavePermissionToUseCommand(message *discordgo.MessageCreate, triggerPhrase string, exclamationCommand string) bool {

	perms, err := config.Session.State.MessagePermissions(message.Message)
	if err != nil {
		logger.Error("", err)
		return false
	}

	if perms&discordgo.PermissionAdministrator == 0 {
		logger.Info("User %v tried to use a command they don't have Permission for :: %v", message.Author.Username, message.Content)
		_, err := config.Session.ChannelMessageSend(message.ChannelID, "You don't have the required permissions for this command!")
		if err != nil {
			logger.Error("", err)
		}
		return false
	} else {
		return true
	}
}
