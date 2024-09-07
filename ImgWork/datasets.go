package imgwork

import "time"

type ImgCategory struct {
	ID       int
	Category string
}

type ImgStorage struct {
	ID          int
	URL         string
	LastChecked time.Time
}

type ImgGuildLink struct {
	ID            int
	Storage       ImgStorage
	Category      ImgCategory
	GuildID       string
	AddedByUserID string
	AddedDateTime time.Time
}
