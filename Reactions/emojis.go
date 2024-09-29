package reactions

import database "github.com/hashbat-dev/discgo-bot/Database"

const (
	EmojiCategories = iota
	EmojiCategoryUp
	EmojiCategoryDown
)

func AddGuildEmoji(guildId string, userId string, emoji string, categoryId int) error {
	// 1. Get the EmojiStorage ID
	emojiStorageId, err := database.GetEmojiStorageID(guildId, emoji)
	if err != nil {
		return err
	}

	// 2. Insert the Link
	return database.InsertEmojiGuildLink(emojiStorageId, categoryId, guildId, userId)
}
