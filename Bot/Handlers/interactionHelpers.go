package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	cache "github.com/hashbat-dev/discgo-bot/Cache"
	config "github.com/hashbat-dev/discgo-bot/Config"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
	module "github.com/hashbat-dev/discgo-bot/Module"
)

func RegisterModules() {

	// Fetch all registered commands assigned globally from the Bot
	registeredCommands, err := config.Session.ApplicationCommands(config.Session.State.User.ID, "")
	if err != nil {
		logger.Error("MODULES", err)
		return
	}

	// Check Existing Commands
	var validated map[string]interface{} = make(map[string]interface{})
	for _, cmd := range registeredCommands {
		deleteCmd := false
		foundCmd := false
		for _, module := range module.ModuleList {
			if module.Command().Name == cmd.Name {
				foundCmd = true
				// Slash Command already registered, are the options the same?
				var liveOpts map[string]interface{} = make(map[string]interface{})
				for _, opt := range cmd.Options {
					liveOpts[fmt.Sprintf("%s|%s", opt.Name, opt.Type.String())] = struct{}{}
				}
				var currOpts map[string]interface{} = make(map[string]interface{})
				for _, opt := range module.Command().Options {
					currOpts[fmt.Sprintf("%s|%s", opt.Name, opt.Type.String())] = struct{}{}
				}

				// Different option lengths?
				if len(liveOpts) != len(currOpts) {
					deleteCmd = true
				}

				// Matching options?
				if !deleteCmd {
					for item := range liveOpts {
						if _, found := currOpts[item]; !found {
							deleteCmd = true
						}
					}
				}

				if !deleteCmd {
					validated[cmd.Name] = struct{}{}
				}
				break
			}

			if foundCmd {
				break
			}
		}

		if !foundCmd {
			deleteCmd = true
		}

		if deleteCmd {
			err := config.Session.ApplicationCommandDelete(config.Session.State.User.ID, "", cmd.ID)
			if err != nil {
				logger.ErrorText("MODULE", "Error deleting Command: %s, Error: %e", cmd.Name, err)
			}
		} else {
			validated[cmd.Name] = struct{}{}
		}
	}

	// Register the Commands again
	for _, cmd := range module.ModuleList {

		if _, validated := validated[cmd.Command().Name]; validated {
			// Already registered
			logger.Debug("MODULE", "Command already registered: %v", cmd.Command().Name)
			continue
		}

		_, err := config.Session.ApplicationCommandCreate(config.Session.State.User.ID, "", cmd.Command())
		if err != nil {
			logger.ErrorText("MODULE", "Error registering Command: %v, Error: %v", cmd.Command().Name, err)
		} else {
			logger.Debug("MODULE", "Successfully registered command: %v", cmd.Command().Name)
		}
	}
}

func ConfirmPermissions(i *discordgo.InteractionCreate, permLevel int) bool {
	switch permLevel {
	case config.CommandLevelUser:
		return true
	case config.CommandLevelBotAdmin:
		guildAdminRole := cache.ActiveGuilds[i.GuildID].BotAdminRole
		if guildAdminRole == "" {
			return i.Member.User.ID == cache.ActiveGuilds[i.GuildID].ServerOwner
		} else {
			for _, r := range i.Member.Roles {
				if r == guildAdminRole {
					return true
				}
			}
			return false
		}
	case config.CommandLevelServerOwner:
		return i.Member.User.ID == cache.ActiveGuilds[i.GuildID].ServerOwner
	default:
	}
	return false
}
