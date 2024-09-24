package data

import "time"

type GuildEmoji struct {
	ID            int
	EmojiID       int
	CategoryID    int
	GuildID       string
	AddedByUserID string
	AddedDateTime time.Time
	EmojiCategory string
	Emoji         string
}
