package handlers

import (
	"fmt"
	"sort"

	"github.com/bwmarrin/discordgo"
	slash "github.com/hashbat-dev/discgo-bot/Bot/Commands/Slash"
	cache "github.com/hashbat-dev/discgo-bot/Cache"
	config "github.com/hashbat-dev/discgo-bot/Config"
	discord "github.com/hashbat-dev/discgo-bot/Discord"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
	reactions "github.com/hashbat-dev/discgo-bot/Reactions"
)

// ===[Add Slash Commands]=============================================
var slashCommands = []SlashCommand{
	//	/support
	{
		Command: &discordgo.ApplicationCommand{
			Name:        "support",
			Description: "Get Support for Discgo Bot and its features",
		},
		Handler: func(i *discordgo.InteractionCreate, correlationId string) {
			slash.SupportInfo(i, correlationId)
		},
		Complexity: config.TRIVIAL_TASK,
	},
	//	/help
	{
		Command: &discordgo.ApplicationCommand{
			Name:        "help",
			Description: "View all the Bot's commands and their descriptions",
		},
		Handler: func(i *discordgo.InteractionCreate, correlationId string) {
			slash.SendHelp(i, correlationId)
		},
		Complexity: config.TRIVIAL_TASK,
	},
	//	-> Create Meme
	{
		Command: &discordgo.ApplicationCommand{
			Name: "Create Meme",
			Type: discordgo.MessageApplicationCommand,
		},
		Handler: func(i *discordgo.InteractionCreate, correlationId string) {
			slash.MakeMemeInit(i, correlationId)
		},
		Complexity: config.TRIVIAL_TASK,
	},
	//	/admin-role
	{
		Command: &discordgo.ApplicationCommand{
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
		},
		Handler: func(i *discordgo.InteractionCreate, correlationId string) {
			slash.AssignNewAdminRole(i, correlationId)
		},
		Complexity:      config.TRIVIAL_TASK,
		PermissionLevel: config.CommandLevelServerOwner,
	},
	//	/edit-starboard-reactions
	{
		Command: &discordgo.ApplicationCommand{
			Name:        "edit-hall-reactions",
			Description: "[ADMIN] Designate which Emojis constitute up/down votes",
		},
		Handler: func(i *discordgo.InteractionCreate, correlationId string) {
			reactions.EditReactions(i, correlationId)
		},
		Complexity:      config.TRIVIAL_TASK,
		PermissionLevel: config.CommandLevelBotAdmin,
	},
	//	/add-trigger
	{
		Command: &discordgo.ApplicationCommand{
			Name:        "add-trigger",
			Description: "[ADMIN] Add words to be tracked in the Server",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "phrase",
					Description: "The word/phrase to be tracked (case insensitive).",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "notify",
					Description: "Whether to post a message when the word/phrase is detected.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "phrase-only",
					Description: "Whether to ignore if used in a sentence.",
					Required:    true,
				},
			},
		},
		Handler: func(i *discordgo.InteractionCreate, correlationId string) {
			slash.AddTrigger(i, correlationId)
		},
		Complexity:      config.TRIVIAL_TASK,
		PermissionLevel: config.CommandLevelBotAdmin,
	},
	//	/edit-trigger
	{
		Command: &discordgo.ApplicationCommand{
			Name:        "edit-trigger",
			Description: "[ADMIN] Edit a word which is being tracked in the Server",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "phrase",
					Description: "The word/phrase to update",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "notify",
					Description: "If entered, updates whether to post a message when the word/phrase is detected.",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "phrase-only",
					Description: "If entered, updates whether to ignore if used in a sentence.",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "delete",
					Description: "Deletes the phrase. WARNING: Will remove ALL statistics for this phrase in your server.",
					Required:    false,
				},
			},
		},
		Handler: func(i *discordgo.InteractionCreate, correlationId string) {
			slash.EditTrigger(i, correlationId)
		},
		Complexity:      config.TRIVIAL_TASK,
		PermissionLevel: config.CommandLevelBotAdmin,
	},
	//	/phrase-leaderboard
	{
		Command: &discordgo.ApplicationCommand{
			Name:        "phrase-leaderboard",
			Description: "See who's ranking highest for tracked phrases in the server!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "phrase",
					Description: "The word/phrase to see the leaderboard for",
					Required:    true,
				},
			},
		},
		Handler: func(i *discordgo.InteractionCreate, correlationId string) {
			slash.PhraseLeaderboard(i, correlationId)
		},
		Complexity: config.TRIVIAL_TASK,
	},
	//	/wow-stats
	{
		Command: &discordgo.ApplicationCommand{
			Name:        "wow-stats",
			Description: "See yours (or someone else's) Wow stats!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "The user to see the Wow stats for",
					Required:    true,
				},
			},
		},
		Handler: func(i *discordgo.InteractionCreate, correlationId string) {
			slash.WowStats(i, correlationId)
		},
		Complexity: config.TRIVIAL_TASK,
	},
	//	/wow-leaderboard
	{
		Command: &discordgo.ApplicationCommand{
			Name:        "wow-leaderboard",
			Description: "See the server leaderboard for Wows!",
		},
		Handler: func(i *discordgo.InteractionCreate, correlationId string) {
			slash.WowLeaderboard(i, correlationId)
		},
		Complexity: config.TRIVIAL_TASK,
	},
	//	/wow-shop
	{
		Command: &discordgo.ApplicationCommand{
			Name:        "wow-shop",
			Description: "Spend your Wow Bucks!",
		},
		Handler: func(i *discordgo.InteractionCreate, correlationId string) {
			slash.WowShop(i, correlationId)
		},
		Complexity: config.TRIVIAL_TASK,
	},
	//	/wow-inventory
	{
		Command: &discordgo.ApplicationCommand{
			Name:        "wow-inventory",
			Description: "See yours (or someone else's) Wow Inventory, showing all their currently active shop items!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "The user to see the Wow Inventory for",
					Required:    true,
				},
			},
		},
		Handler: func(i *discordgo.InteractionCreate, correlationId string) {
			slash.WowInventory(i, correlationId)
		},
		Complexity: config.TRIVIAL_TASK,
	},
}

// Message Commands are not allowed Descriptions, enter User descriptions below for these.
var userDescriptions map[string]string = map[string]string{
	"Create Meme": "Turn an image into a Meme! Either enter an Above Image caption (Option A) or choose Top/Bottom text to add (Option B)",
}

type SlashCommand struct {
	Command         *discordgo.ApplicationCommand
	Handler         func(i *discordgo.InteractionCreate, correlationId string)
	Complexity      int
	PermissionLevel int
}

var SlashCommands map[string]SlashCommand = make(map[string]SlashCommand)

func InitSlashCommands() {
	// Handle adding Slash Commands (these are added below this function)
	for _, cmd := range slashCommands {
		SlashCommands[cmd.Command.Name] = cmd
	}

	writeSlashHelpText()

	config.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		SlashCommandHandler(s, i)
	})
}

func SlashCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	logger.Debug(i.GuildID, "SlashCommandHandler")
	if i.Type == discordgo.InteractionMessageComponent {
		return
	}
	if i.Type == discordgo.InteractionModalSubmit {
		return
	}

	if config.ServiceSettings.ISDEV != cache.ActiveGuilds[i.GuildID].IsDev {
		return
	}

	cmdName := i.ApplicationCommandData().Name
	if cmd, exists := SlashCommands[i.ApplicationCommandData().Name]; exists {

		if !confirmPermissions(i, cmd.PermissionLevel) {
			logger.Event(i.GuildID, "User [ID: %s, UserName: %s] was blocked from using command [%s]", i.Member.User.ID, i.Member.User.Username, cmd.Command.Name)
			discord.SendEmbedFromInteraction(i, "Permission Denied", "You do not have permission to use this command.", 0)
		} else {
			DispatchTask(&Task{
				CommandType: config.CommandTypeSlash,
				Complexity:  cmd.Complexity,
				SlashDetails: &SlashTaskDetails{
					Interaction:  i,
					SlashCommand: cmd,
				},
			})
		}
		return
	}

	logger.ErrorText(i.GuildID, "No Handler found for Slash Command: %v", cmdName)
}

func confirmPermissions(i *discordgo.InteractionCreate, permLevel int) bool {
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

func RefreshSlashCommands(guildId string) {
	if len(SlashCommands) == 0 {
		InitSlashCommands()
	}

	// Fetch all registered commands assigned from the Bot -> the Guild
	registeredCommands, err := config.Session.ApplicationCommands(config.Session.State.User.ID, guildId)
	if err != nil {
		logger.Error(guildId, err)
		return
	}

	// Check Existing Commands
	var validated = make(map[string]interface{})
	for _, cmd := range registeredCommands {
		deleteCmd := false
		if local, exists := SlashCommands[cmd.Name]; exists {
			// Slash Command already registered, are the options the same?
			var liveOpts = make(map[string]interface{})
			for _, opt := range cmd.Options {
				liveOpts[fmt.Sprintf("%s|%s", opt.Name, opt.Type.String())] = struct{}{}
			}
			var currOpts = make(map[string]interface{})
			for _, opt := range local.Command.Options {
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
		} else {
			// Slash Command exists externally but not in our map, delete it
			deleteCmd = true
		}
		if deleteCmd {
			err := config.Session.ApplicationCommandDelete(config.Session.State.User.ID, guildId, cmd.ID)
			if err != nil {
				logger.ErrorText(guildId, "Error deleting Command: %s, Error: %e", cmd.Name, err)
			}
		} else {
			validated[cmd.Name] = struct{}{}
		}
	}

	// Register the Commands again
	for _, cmd := range SlashCommands {
		if _, validated := validated[cmd.Command.Name]; validated {
			// Already registered
			logger.Debug(guildId, "Command already registered: %v", cmd.Command.Name)
			continue
		}

		_, err := config.Session.ApplicationCommandCreate(config.Session.State.User.ID, guildId, cmd.Command)
		if err != nil {
			logger.ErrorText(guildId, "Error registering Command: %v, Error: %v", cmd.Command.Name, err)
		} else {
			logger.Debug(guildId, "Successfully registered command: %v", cmd.Command.Name)
		}
	}
}

func writeSlashHelpText() {
	// 1. Sort the slice Alphabetically
	slashCmds := slashCommands
	sort.Slice(slashCmds, func(i, j int) bool {
		return slashCmds[i].Command.Name < slashCmds[j].Command.Name
	})

	// 2. Generate the Help Text
	slashText := ""
	msgCmdText := ""

	for _, cmd := range slashCmds {
		if cmd.Command.Type == discordgo.MessageApplicationCommand {
			// Message Commands
			if len(msgCmdText) > 0 {
				msgCmdText += "\n"
			}
			msgCmdText += "**" + cmd.Command.Name + "**: " + userDescriptions[cmd.Command.Name]
		} else {
			// Regular Slash Commands
			cmdText := "**/" + cmd.Command.Name + "**: " + cmd.Command.Description
			if len(cmd.Command.Options) > 0 {
				cmdText += " It accepts these parameters: "
			}
			for i, opt := range cmd.Command.Options {
				if i > 0 {
					cmdText += ", "
				}
				cmdText += opt.Name + ": *" + opt.Description + "*"
			}
			if len(slashText) > 0 {
				slashText += "\n"
			}
			slashText += cmdText
		}

	}

	// 3. Set it in the Config
	config.UserSlashHelpText = slashText
	config.UserMsgCmdHelpText = msgCmdText
}
