package interactions

import (
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

var Active map[string]ActiveInteraction = make(map[string]ActiveInteraction)   // Key: CorrelationID
var Handlers map[string]HandleInteraction = make(map[string]HandleInteraction) // Key: ObjectID

var muActive sync.Mutex

type ActiveInteraction struct {
	Interactions []discordgo.InteractionCreate
	Values       InteractionValues
	Timestamps   []time.Time
}

type HandleInteraction struct {
	Complexity int
	Execute    func(i *discordgo.InteractionCreate, correlationId string)
}

type InteractionValues struct { // Keys: Source ObjectID
	String  map[string]string
	Integer map[string]int64
	Bool    map[string]bool
	User    map[string]*discordgo.User
	Channel map[string]*discordgo.Channel
	Role    map[string]*discordgo.Role
	Number  map[string]float64
}

// Adds the Interaction to the local cache and retrieves Values from any submitted form elements
func UpdateCache(correlationId string, i *discordgo.InteractionCreate) {
	if _, exists := Active[correlationId]; exists {
		updateInteraction(correlationId, i)
	} else {
		insertInteraction(correlationId, i)
	}
	ExtractValues(correlationId, i)
}

func Complete(correlationId string) {
	delete(Active, correlationId)
}

func insertInteraction(correlationId string, i *discordgo.InteractionCreate) {
	muActive.Lock()
	Active[correlationId] = ActiveInteraction{
		Interactions: []discordgo.InteractionCreate{*i},
		Values: InteractionValues{
			String:  make(map[string]string),
			Integer: make(map[string]int64),
			Bool:    make(map[string]bool),
			User:    make(map[string]*discordgo.User),
			Channel: make(map[string]*discordgo.Channel),
			Role:    make(map[string]*discordgo.Role),
			Number:  make(map[string]float64),
		},
		Timestamps: []time.Time{time.Now()},
	}
	muActive.Unlock()
}

func updateInteraction(correlationId string, i *discordgo.InteractionCreate) {
	muActive.Lock()
	cache := Active[correlationId]
	cache.Interactions = append(cache.Interactions, *i)
	cache.Timestamps = append(cache.Timestamps, time.Now())
	Active[correlationId] = cache
	muActive.Unlock()
}
