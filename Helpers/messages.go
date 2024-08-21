package helpers

import (
	"errors"
	"strings"

	"github.com/bwmarrin/discordgo"
	config "github.com/dabi-ngin/discgo-bot/Config"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
)

// Checks for, and returns if exists a !command
func CheckForBangCommand(messageContent string) string {
	if len(messageContent) == 0 {
		return ""
	}

	if string([]rune(messageContent)[0]) == "!" {
		spaceIndex := strings.Index(messageContent, " ")
		if spaceIndex == -1 {
			// No spaces in the Content, we assume the whole message is the ! command
			return messageContent[1:]
		} else {
			return strings.Split(messageContent, " ")[0]
		}
	}
	return ""
}

// Returns the Image URL if one exists in the message, if not the String is blank. If a specific extension is needed enter it as the second variable, if blank it will accept anything defined in Config.
func GetImageFromMessage(message *discordgo.Message, requiredExtension string) string {

	msgContent := strings.Trim(message.Content, " ")
	msgContentLower := strings.ToLower(msgContent)

	// 1. Check for a file =========================================
	// A. Are there any Embeds?
	imgLink := ""
	if len(message.Embeds) > 0 {
		imgLink = message.Embeds[0].Thumbnail.ProxyURL
	}

	// B. Are there any Attachments?
	if imgLink == "" && len(message.Attachments) > 0 {
		imgLink = message.Attachments[0].ProxyURL
	}

	// C. Is this a Tenor link?
	if strings.Contains(msgContentLower, "tenor.com/") {
		tenorLink, err := GetImageUrlFromTenor(msgContentLower)
		if err != nil {
			logger.Error(message.GuildID, err)
			return ""
		} else if tenorLink != "" {
			imgLink = tenorLink
		}
	}

	// D. Okay last check, what about the body content?
	if imgLink == "" {
		if requiredExtension == "" {
			for _, ext := range config.ValidImageExtensions {
				if strings.Contains(msgContentLower, ext) {
					imgLink = msgContent
					break
				}
			}
		} else {
			if strings.Contains(msgContentLower, requiredExtension) {
				imgLink = msgContent
			}
		}
	}

	// D. No file?
	if imgLink == "" {
		logger.Error(message.GuildID, errors.New("no suitable image found"))
		return ""
	}

	// 2. Now lets validate the file name (if any) ================================
	extValid := false
	if requiredExtension == "" {
		for _, ext := range config.ValidImageExtensions {
			if strings.Contains(imgLink, ext) {
				extValid = true
				break
			}
		}
	} else {
		if strings.Contains(imgLink, requiredExtension) {
			extValid = true
		}
	}

	if !extValid {
		logger.Error(message.GuildID, errors.New("image was not a suitable extension type"))
		return ""
	}

	// 3. Return! =====================================================================
	return imgLink
}
