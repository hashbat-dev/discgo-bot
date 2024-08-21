package helpers

import (
	"github.com/bwmarrin/discordgo"
	config "github.com/dabi-ngin/discgo-bot/Config"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
)

func DoesUserHavePermissionToUseCommand(message *discordgo.MessageCreate) bool {

	temp := true
	logger.Remind("Permissions.go :: Temp used - needs to be changed to actually check perms.")

	if !temp {
		logger.Info("User %v tried to use a command they don't have Permission for :: %v", message.Author.Username, message.Content)
		_, err := config.Session.ChannelMessageSend(message.ChannelID, "You don't have the required permissions for this command!")
		if err != nil {
			logger.Error("", err)

		}
	}
	return true
}
