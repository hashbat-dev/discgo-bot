package slash

import (
	"github.com/bwmarrin/discordgo"
	cache "github.com/hashbat-dev/discgo-bot/Cache"
	database "github.com/hashbat-dev/discgo-bot/Database"
	discord "github.com/hashbat-dev/discgo-bot/Discord"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

func AssignNewAdminRole(i *discordgo.InteractionCreate, correlationId string) {
	cachedInteraction := cache.ActiveInteractions[correlationId]
	newRoleId := ""
	newRoleName := ""
	if _, exists := cachedInteraction.Values.Role["role"]; exists {
		newRoleId = cachedInteraction.Values.Role["role"].ID
		newRoleName = cachedInteraction.Values.Role["role"].Name
	}

	if newRoleId == "" {
		logger.ErrorText(i.GuildID, "Couldn't find input role id")
		discord.SendGenericErrorFromInteraction(i)
		return
	}

	guildDbObj, err := database.Get(i.GuildID)
	oldRoleId := guildDbObj.GuildAdminRole
	if err != nil {
		discord.SendGenericErrorFromInteraction(i)
		return
	}

	guildDbObj.GuildAdminRole = newRoleId
	_, err = database.Upsert(guildDbObj)
	if err != nil {
		discord.SendGenericErrorFromInteraction(i)
		return
	}

	logger.Event(i.GuildID, "Bot Administrator role changed from [%s] to [%s]", oldRoleId, newRoleId)
	discord.SendEmbedFromInteraction(i, "Admin Role Update", "Bot Administrator role is now: ["+newRoleName+"]")
}
