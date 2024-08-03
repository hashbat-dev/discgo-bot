package meme

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

var InteractionCache []CacheItem
var CacheTimeoutMins int = 60

type CacheItem struct {
	Interaction *discordgo.InteractionCreate
	Added       time.Time
}
type ThreadResults struct {
	Title     string
	Board     string
	OPID      int
	FileCount int
}

var BoardsToSearch []string = []string{
	"wsg",
	"vg",
	"v",
	"vm",
	"c",
	"b",
	"gif",
	"mlp",
	"pol",
	"x",
	"vp",
	"vt",
}
var GenericSearch []string = []string{
	"neco",
	"gif caption",
	"ylyl",
	"family guy",
	"simpsons",
	"cat ",
	"brainrot",
}

type FoundFiles struct {
	PostID    int
	FileID    int
	Extension string
}

type ThreadPosts struct {
	Posts []ThreadPost `json:"posts"`
}

type ThreadPost struct {
	PostID    int    `json:"no,omitempty"`
	Extension string `json:"ext,omitempty"`
	FileID    int    `json:"tim,omitempty"`
}

type ThreadList []struct {
	Page   int            `json:"page"`
	Thread []ThreadObject `json:"threads"`
}

type ThreadObject struct {
	OPNumber int    `json:"no"`
	Title    string `json:"sub"`
	Text     string `json:"com"`
	Images   int    `json:"images"`
}
