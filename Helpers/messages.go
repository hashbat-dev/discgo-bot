package helpers

import (
	"errors"
	"strings"

	"github.com/bwmarrin/discordgo"
	config "github.com/dabi-ngin/discgo-bot/Config"
	external "github.com/dabi-ngin/discgo-bot/External"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
)

// Returns the Image URL if one exists in the message, if not the String is blank. If a specific extension is needed enter it as the second variable, if blank it will accept anything defined in Config.
func GetImageFromMessage(message *discordgo.Message, requiredExtension string) string {
	msgContent := strings.Trim(message.Content, " ")
	msgContentLower := strings.ToLower(msgContent)

	// 1. Check if there's a Replied to Message ====================
	if message.ReferencedMessage != nil {
		repliedToImg := GetImageFromMessage(message.ReferencedMessage, requiredExtension)
		if repliedToImg != "" {
			return repliedToImg
		}
	}

	// 2. Check for a file =========================================
	// A. Are there any Embeds?
	imgLink := ""
	if len(message.Embeds) > 0 {
		for _, embed := range message.Embeds {
			if embed.Type == "video" {
				continue
			}
			// If the Image is a GIF from Tenor then the Thumbnail will be
			// a PNG, this check avoids getting the wrong image.
			if embed.Type != "gif" && embed.Type != "gifv" {
				imgLink = message.Embeds[0].Thumbnail.ProxyURL
			}
		}

	}

	// B. Are there any Attachments?
	if imgLink == "" && len(message.Attachments) > 0 {
		imgLink = message.Attachments[0].ProxyURL
	}

	// C. Is this a Tenor link?
	if imgLink == "" && strings.Contains(msgContentLower, "tenor.com/") {
		tenorLink, err := external.GetImageUrlFromTenor(message.GuildID, msgContent)
		if err != nil {
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

	// 3. Now lets validate the file name (if any) ================================
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

	// 4. Return! =====================================================================
	return imgLink
}
