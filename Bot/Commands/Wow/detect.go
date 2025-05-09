package wow

import (
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

func Detect(message *discordgo.MessageCreate) {
	msg := strings.ToLower(strings.TrimSpace(message.Content))
	if msg == "" {
		return
	}

	regCheck := regexp.MustCompile(`(?i)^w+o{1,}w+$`)
	if regCheck.MatchString(msg) {
		logger.Info(message.GuildID, "Wow detected, MessageID: %s", message.ID)
		queueGenerate <- message
	}
}
