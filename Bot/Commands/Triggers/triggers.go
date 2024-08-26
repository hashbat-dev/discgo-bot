package triggers

import "time"

type Phrase struct {
	ID                int
	Phrase            string
	NotifyOnDetection bool
	WordOnlyMatch     bool
}

type PhraseLink struct {
	ID            int
	Phrase        Phrase
	GuildID       string
	AddedByUserID string
	AddedDateTime time.Time
}
