package slash

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	cache "github.com/hashbat-dev/discgo-bot/Cache"
	database "github.com/hashbat-dev/discgo-bot/Database"
	discord "github.com/hashbat-dev/discgo-bot/Discord"
)

func AddTrigger(i *discordgo.InteractionCreate, correlationId string) {
	cachedInteraction := cache.ActiveInteractions[correlationId]
	phrase := cachedInteraction.Values.String["phrase"]
	notify := cachedInteraction.Values.Bool["notify"]
	phraseOnly := cachedInteraction.Values.Bool["phrase-only"]
	phrase = strings.ToLower(strings.TrimSpace(phrase))

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

	// 2. Get the TriggerPhrase
	triggerPhrase, err := database.GetTriggerPhrase(i.GuildID, phrase)
	if err != nil {
		discord.SendGenericErrorFromInteraction(i)
		cache.InteractionComplete(correlationId)
		return
	}

	// 3. Does a Phrase Link already exist?
	linkExists, err := database.DoesPhraseLinkExist(i.GuildID, triggerPhrase.ID)
	if err != nil {
		discord.SendGenericErrorFromInteraction(i)
		cache.InteractionComplete(correlationId)
		return
	}

	if linkExists {
		discord.SendEmbedFromInteraction(i, "Error", "This Phrase already exists! Use the /edit-phrase command instead to update it.")
		cache.InteractionComplete(correlationId)
		return
	}

	// 4. Insert the Phrase Link
	err = database.InsertPhraseGuildLink(triggerPhrase.ID, i.GuildID, i.Member.User.ID, notify, phraseOnly)
	if err != nil {
		discord.SendGenericErrorFromInteraction(i)
		cache.InteractionComplete(correlationId)
		return
	}

	// 5. Add to the Cache
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

	// Done!
	cache.InteractionComplete(correlationId)
	discord.SendEmbedFromInteraction(i, "Phrase Added", fmt.Sprintf("The phrase '%s' is now being tracked in the server!", phrase))
}
