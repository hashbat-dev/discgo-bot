package slash

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	cache "github.com/hashbat-dev/discgo-bot/Cache"
	database "github.com/hashbat-dev/discgo-bot/Database"
	discord "github.com/hashbat-dev/discgo-bot/Discord"
)

func EditTrigger(i *discordgo.InteractionCreate, correlationId string) {
	cachedInteraction := cache.ActiveInteractions[correlationId]
	phrase := cachedInteraction.Values.String["phrase"]
	phrase = strings.ToLower(strings.TrimSpace(phrase))

	notify := -1
	if val, ok := cachedInteraction.Values.Bool["notify"]; ok {
		if val {
			notify = 1
		} else {
			notify = 0
		}
	}

	phraseOnly := -1
	if val, ok := cachedInteraction.Values.Bool["phrase-only"]; ok {
		if val {
			phraseOnly = 1
		} else {
			phraseOnly = 0
		}
	}

	delete := -1
	if val, ok := cachedInteraction.Values.Bool["delete"]; ok {
		if val {
			delete = 1
		} else {
			delete = 0
		}
	}

	// 1. Validate
	if phrase == "" {
		discord.SendEmbedFromInteraction(i, "Error", "No Phrase entered!")
		cache.InteractionComplete(correlationId)
		return
	}
	if len(phrase) > 50 {
		discord.SendEmbedFromInteraction(i, "Error", fmt.Sprintf("Phrase too long! Your phrase was %d characters out of the maximum of 50.", len(phrase)))
		cache.InteractionComplete(correlationId)
		return
	}

	triggerPhrase, err := database.GetTriggerPhrase(i.GuildID, phrase)
	if err != nil {
		discord.SendGenericErrorFromInteraction(i)
		cache.InteractionComplete(correlationId)
		return
	}

	linkExists, err := database.DoesPhraseLinkExist(i.GuildID, triggerPhrase.ID)
	if err != nil {
		discord.SendGenericErrorFromInteraction(i)
		cache.InteractionComplete(correlationId)
		return
	}

	if !linkExists {
		discord.SendEmbedFromInteraction(i, "Error", fmt.Sprintf("The phrase '%s' doesn't exist! use /add-phrase to create it", phrase))
		cache.InteractionComplete(correlationId)
		return
	}

	if notify == -1 && phraseOnly == -1 && delete == -1 {
		discord.SendEmbedFromInteraction(i, "Error", "Nothing to update! Our work here is done :)")
		cache.InteractionComplete(correlationId)
		return
	}

	// 2. Perform updates
	err = database.UpdatePhraseGuildLink(triggerPhrase.ID, i.GuildID, notify, phraseOnly, delete)
	if err != nil {
		discord.SendGenericErrorFromInteraction(i)
		cache.InteractionComplete(correlationId)
		return
	}

	// 3. Update the Cache
	newPhrase, err := database.GetGuildPhrases(i.GuildID, triggerPhrase.ID)
	if err != nil {
		discord.SendGenericErrorFromInteraction(i)
		cache.InteractionComplete(correlationId)
		return
	}

	if len(newPhrase) == 0 {
		discord.SendGenericErrorFromInteraction(i)
		cache.InteractionComplete(correlationId)
		return
	}

	cache.AddGuildTrigger(i.GuildID, newPhrase[0].Phrase)

	cache.InteractionComplete(correlationId)
	discord.SendEmbedFromInteraction(i, "Phrase Updated", fmt.Sprintf("The phrase '%s' has been successfully updated!", phrase))
}
