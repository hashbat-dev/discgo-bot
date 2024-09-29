package reactions

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	cache "github.com/hashbat-dev/discgo-bot/Cache"
	config "github.com/hashbat-dev/discgo-bot/Config"
	database "github.com/hashbat-dev/discgo-bot/Database"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

func AddOrUpdate(message *discordgo.Message, score int, emojiString string) {
	messageEmbed := createMessage(message, score, emojiString)
	currDbObj := database.Starboard_Get(message.GuildID, message.ID)
	isUp := score > 0

	// 1. If we have a DB Object, is it still the same +/- channel?
	if currDbObj.ID > 0 && isUp != currDbObj.IsUpChannel {
		database.Starboard_Delete(message.GuildID, currDbObj.ID)
		currDbObj = database.StarboardMessage{}
	}

	// 2. Get the ChannelID we're working with
	channelId := ""
	if isUp {
		channelId = cache.ActiveGuilds[message.GuildID].StarUpChannel
	} else {
		channelId = cache.ActiveGuilds[message.GuildID].StarDownChannel
	}
	if channelId == "" {
		channelId = CreateChannel(message.GuildID, isUp)
	}
	if channelId == "" {
		return
	}

	// 3. Do we already have a Message?
	if currDbObj.StarboardMessageID != "" {
		_, err := config.Session.ChannelMessageEditEmbed(channelId, currDbObj.StarboardMessageID, messageEmbed)
		if err != nil {
			logger.Error(message.GuildID, err)
			return
		}
	} else {
		newMsg, err := config.Session.ChannelMessageSendEmbed(channelId, messageEmbed)
		if err != nil {
			logger.Error(message.GuildID, err)
			return
		}
		currDbObj.StarboardMessageID = newMsg.ID
	}

	// 4. Update the Object and send back the updates to the Database
	if currDbObj.GuildID == "" {
		currDbObj.GuildID = message.GuildID
	}
	if currDbObj.UserID == "" {
		currDbObj.UserID = message.Author.ID
	}
	if currDbObj.OriginalMessageID == "" {
		currDbObj.OriginalMessageID = message.ID
	}
	currDbObj.IsUpChannel = isUp
	currDbObj.Score = score
	currDbObj.EmojiString = emojiString
	database.Starboard_InsertUpdate(currDbObj)
}

func createMessage(msg *discordgo.Message, score int, emojiString string) *discordgo.MessageEmbed {
	isUp := true
	if score < 0 {
		isUp = false
	}
	color := 0xFFEE00
	if !isUp {
		color = 0xFF0000
	}
	imageURL, newContent := extractImageURL(msg.Content)
	originalMessageLink := fmt.Sprintf("https://discord.com/channels/%s/%s/%s", msg.GuildID, msg.ChannelID, msg.ID)

	embed := &discordgo.MessageEmbed{
		Description: newContent,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    msg.Author.Username,
			IconURL: msg.Author.AvatarURL(""),
		},
		Color: color,
		Footer: &discordgo.MessageEmbedFooter{
			Text: emojiString,
		},

		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "",
				Value: fmt.Sprintf("[Jump to original message](%s)", originalMessageLink),
			},
		},
	}

	if imageURL != "" {
		embed.Image = &discordgo.MessageEmbedImage{
			URL: imageURL,
		}
	}

	if len(msg.Attachments) > 0 {
		attachment := msg.Attachments[0]
		embed.Image = &discordgo.MessageEmbedImage{
			URL: attachment.URL,
		}
	}

	return embed
}

func extractImageURL(content string) (string, string) {
	imageRegex := regexp.MustCompile(`https?://\S+\.(jpg|jpeg|png|gif|webp)(\?\S*)?`)
	imageURLWithQuery := imageRegex.FindString(content)

	if imageURLWithQuery != "" {
		parsedURL, err := url.Parse(imageURLWithQuery)
		if err == nil {
			imageURL := fmt.Sprintf("%s://%s%s", parsedURL.Scheme, parsedURL.Host, parsedURL.Path)
			newContent := strings.Replace(content, imageURLWithQuery, "", 1)
			return imageURL, strings.TrimSpace(newContent)
		}
	}

	return "", content
}
