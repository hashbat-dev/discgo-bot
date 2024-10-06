package module

import (
	"github.com/bwmarrin/discordgo"
	cache "github.com/hashbat-dev/discgo-bot/Cache"
	config "github.com/hashbat-dev/discgo-bot/Config"
	database "github.com/hashbat-dev/discgo-bot/Database"
	discord "github.com/hashbat-dev/discgo-bot/Discord"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

type AdminRole struct{}

func (s AdminRole) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "admin-role",
		Description: "[SERVER OWNER ONLY] Designate role with access to all Bot options",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionRole,
				Name:        "role",
				Description: "The Role you want to assign as the Bot Administrator.",
				Required:    true,
			},
		},
	}
}

func (s AdminRole) PermissionRequirement() int {
	return config.CommandLevelServerOwner
}

func (s AdminRole) Complexity() int {
	return config.TRIVIAL_TASK
}

func (s AdminRole) Execute(i *discordgo.InteractionCreate, correlationId string) {
	cachedInteraction := cache.ActiveInteractions[correlationId]
	newRoleId := ""
	newRoleName := ""
	if _, exists := cachedInteraction.Values.Role["role"]; exists {
		newRoleId = cachedInteraction.Values.Role["role"].ID
		newRoleName = cachedInteraction.Values.Role["role"].Name
	}

	if newRoleId == "" {
		logger.ErrorText(i.GuildID, "Couldn't find input role id")
		discord.Interactions_SendError(i, "")
		return
	}

	guildDbObj, err := database.Get(i.GuildID)
	oldRoleId := guildDbObj.GuildAdminRole
	if err != nil {
		discord.Interactions_SendError(i, "")
		return
	}

	guildDbObj.GuildAdminRole = newRoleId
	_, err = database.Upsert(guildDbObj)
	if err != nil {
		discord.Interactions_SendError(i, "")
		return
	}

	logger.Event(i.GuildID, "Bot Administrator role changed from [%s] to [%s]", oldRoleId, newRoleId)
	discord.Interactions_SendMessage(i, "Admin Role Update", "Bot Administrator role is now: ["+newRoleName+"]")
}
