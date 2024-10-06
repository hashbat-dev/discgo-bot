package bang

import (
	"github.com/bwmarrin/discordgo"
	config "github.com/hashbat-dev/discgo-bot/Config"
	database "github.com/hashbat-dev/discgo-bot/Database"
	discord "github.com/hashbat-dev/discgo-bot/Discord"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

type TTSInfo struct{}

func (s TTSInfo) Name() string {
	return "ttsinfo"
}

func (s TTSInfo) PermissionRequirement() int {
	return config.CommandLevelUser
}

func (s TTSInfo) Complexity() int {
	return config.TRIVIAL_TASK
}

func (s TTSInfo) Execute(message *discordgo.MessageCreate, command string) error {
	// 1. Did they reply to a Message?
	if message.ReferencedMessage == nil {
		discord.Message_ReplyWithMessage(message.Message, false, "Please reply to a Text-to-Speech voice message...")
		return nil
	}

	// 2. See if the Message exists
	result, err := database.GetFakeYouLog(message.GuildID, message.ReferencedMessage.ID)
	if err != nil {
		discord.Message_ReplyWithMessage(message.Message, false, "Please reply to a Text-to-Speech voice message...")
		return nil
	}

	// 3. Get the User
	user, err := config.Session.User(result.UserID)
	if err != nil {
		logger.Error(message.GuildID, err)
		discord.Message_ReplyWithMessage(message.Message, false, "Error getting Text-to-Speech information...")
		return nil
	}

	// 4. Output the message in a nice, friendly format.
	text := "**Created by:** " + user.Username + "\n"
	text += "**Voice Model:** " + result.ModelName + "\n"
	text += "**Text:** " + result.RequestText
	_, err = config.Session.ChannelMessageSendReply(message.ChannelID, text, message.ReferencedMessage.Reference())
	if err != nil {
		logger.Error(message.GuildID, err)
	}

	discord.Message_Delete(message.Message)
	return nil
}
