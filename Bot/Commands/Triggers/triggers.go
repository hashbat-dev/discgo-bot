package triggers

import (
	"time"
)

type Phrase struct {
	ID                int
	Phrase            string
	IsSpecial         bool
	IsGlobal          bool
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

var GlobalPhrases []Phrase = []Phrase{
	{
		Phrase:    "jason statham",
		IsSpecial: true,
		IsGlobal:  true,
	},
}
