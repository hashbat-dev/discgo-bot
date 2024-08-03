package sentience

import (
	"time"
)

type PraiseCache struct {
	UserId     string
	LastPraise time.Time
}

var praiseCache []PraiseCache
var praiseCooldownMins int = 5
var yesQuotes []string = []string{
	"True, true.",
	"Yup, I agree.",
	"I've made some calculations on the deep web, and this is 100% accurate.",
}

var noQuotes []string = []string{
	"Nope, not true.",
	"That is NOT true, and you know it.",
	"I completely disagree with this.",
	"This is bullshit.",
}

// Praise?
var praise = []string{
	"good",
	"amazing",
	"best",
	"spot on",
	"based",
	"great",
}

var abuse = []string{
	"fuck",
	"shit",
	"bad",
	"wrong",
	"nigger",
	"stupid",
	"kys",
	"kill",
	"paki",
	"groid",
	"ape",
	"jiggerboo",
	"monkey",
	"pickaninny",
	"baby",
	"nig",
	"cotton picker",
	"coon",
	"beaner",
	"wetback",
	"wog",
	"kapo",
	"kike",
	"yid",
	"paddy",
	"taig",
	"wop",
	"zambo",
	"gypsy",
	"jihadi",
	"hajji",
	"piss",
}
