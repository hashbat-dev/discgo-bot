package handlers

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	commands "github.com/hashbat-dev/discgo-bot/Bot/Commands"
	triggers "github.com/hashbat-dev/discgo-bot/Bot/Commands/Triggers"
	cache "github.com/hashbat-dev/discgo-bot/Cache"
	config "github.com/hashbat-dev/discgo-bot/Config"
	database "github.com/hashbat-dev/discgo-bot/Database"
	discord "github.com/hashbat-dev/discgo-bot/Discord"
	editmodule "github.com/hashbat-dev/discgo-bot/EditModule"
	helpers "github.com/hashbat-dev/discgo-bot/Helpers"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
	module "github.com/hashbat-dev/discgo-bot/Module"
	reactions "github.com/hashbat-dev/discgo-bot/Reactions"
	reporting "github.com/hashbat-dev/discgo-bot/Reporting"
)

type Task struct {
	CommandType           int
	Complexity            int
	BangDetails           *BangTaskDetails
	ModuleDetails         *ModuleDetails
	ModuleResponseDetails *ModuleResponseDetails
	EditModuleDetails     *EditModuleDetails
	PhraseDetails         *PhraseTaskDetails
	MessageObj            *discordgo.Message
}

type BangTaskDetails struct {
	Message       *discordgo.MessageCreate
	Command       commands.Command
	CorrelationId uuid.UUID
}

type ModuleDetails struct {
	Interaction *discordgo.InteractionCreate
	Module      module.Module
}

type EditModuleDetails struct {
	Interaction *discordgo.InteractionCreate
	EditModule  editmodule.EditModule
}

type ModuleResponseDetails struct {
	Interaction   *discordgo.InteractionCreate
	ObjectID      string
	CorrelationID string
}

type PhraseTaskDetails struct {
	Message        *discordgo.MessageCreate
	TriggerPhrases []triggers.Phrase
}

var (
	MAX_QUEUE_IO_TASKS      = 100
	MAX_QUEUE_CPU_TASKS     = 100
	MAX_QUEUE_TRIVIAL_TASKS = 1000
	IO_TASKS                = make(chan *Task, MAX_QUEUE_IO_TASKS)
	CPU_TASKS               = make(chan *Task, MAX_QUEUE_CPU_TASKS)
	TRIVIAL_TASKS           = make(chan *Task, MAX_QUEUE_TRIVIAL_TASKS)
)

func init() {
	reporting.CreateWorkerChannel(config.IO_BOUND_TASK, "IO Bound", MAX_QUEUE_IO_TASKS, config.N_IO_WORKERS)
	reporting.CreateWorkerChannel(config.CPU_BOUND_TASK, "CPU Bound", MAX_QUEUE_CPU_TASKS, runtime.GOMAXPROCS(0))
	reporting.CreateWorkerChannel(config.TRIVIAL_TASK, "Trivial", MAX_QUEUE_TRIVIAL_TASKS, config.N_TRIVIAL_WORKERS)
	for i := 0; i < config.N_TRIVIAL_WORKERS; i++ {
		go worker(i, config.TRIVIAL_TASK, TRIVIAL_TASKS)
	}
	for i := 0; i < config.N_IO_WORKERS; i++ {
		go worker(i, config.IO_BOUND_TASK, IO_TASKS)
	}
	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		go worker(i, config.CPU_BOUND_TASK, CPU_TASKS)
	}

}

func worker(id int, taskId int, ch <-chan *Task) {
	for task := range ch {
		reporting.WorkerProcessingStart(taskId)
		switch task.CommandType {
		case config.CommandTypeBang:
			workerBang(task.BangDetails)
		case config.CommandTypeModule:
			workerModule(task.ModuleDetails)
		case config.CommandTypeModuleResponse:
			workerModuleResponse(task.ModuleResponseDetails)
		case config.CommandTypeEditModule:
			workerEditModule(task.EditModuleDetails)
		case config.CommandTypePhrase:
			workerPhrase(task.PhraseDetails)
		case config.CommandTypeReactionCheck:
			workerReaction(task.MessageObj)
		default:
			logger.ErrorText("WORKER", "Unknown CommandType value [%v]", task.CommandType)
		}
		reporting.WorkerProcessingFinish(taskId)
	}
	logger.Info("WORKER", "worker %d: Channel closed, exiting...", id)
}

func workerBang(msg *BangTaskDetails) {
	logger.Info(msg.Message.GuildID, "worker :: processing command [%v] correlation-id :: %v", msg.Command.Name(), msg.CorrelationId)
	timeStart := time.Now()

	execErr := msg.Command.Execute(msg.Message, msg.Command.Name())
	if execErr != nil {
		logger.ErrorText(msg.Message.GuildID, "worker :: [%v] error :: %v :: correlation-id :: %v", msg.Command.Name(), execErr.Error(), msg.CorrelationId)
		return // Failed to execute, skip loop iteration
	}

	reporting.Command(config.CommandTypeBang, msg.Message.GuildID, msg.Message.Author.ID, msg.Message.Author.Username, msg.Command.Name(), msg.CorrelationId.String(), timeStart)
}

func workerModule(msg *ModuleDetails) {
	correlationId := cache.AddInteraction(msg.Interaction, msg.Module.Command().Name)
	timeStarted := time.Now()
	msg.Module.Execute(msg.Interaction, correlationId)
	reporting.Command(config.CommandTypeModule, msg.Interaction.GuildID, msg.Interaction.Member.User.ID, msg.Interaction.Member.User.Username, msg.Module.Command().Name, correlationId, timeStarted)
}

func workerEditModule(msg *EditModuleDetails) {
	correlationId := cache.AddInteraction(msg.Interaction, msg.EditModule.SelectName())
	timeStarted := time.Now()
	msg.EditModule.Execute(msg.Interaction, correlationId)
	reporting.Command(config.CommandTypeModule, msg.Interaction.GuildID, msg.Interaction.Member.User.ID, msg.Interaction.Member.User.Username, msg.EditModule.SelectName(), correlationId, timeStarted)
}

func workerModuleResponse(msg *ModuleResponseDetails) {
	logger.Info(msg.Interaction.GuildID, "Interaction ID: [%v] Processing Response, Object: [%v]", msg.CorrelationID, msg.ObjectID)
	timeStarted := time.Now()

	// Update the Interaction Cache with the provided options
	cache.UpdateInteraction(msg.CorrelationID, msg.Interaction)

	// Execute the Command Response
	if msg.ObjectID == "edit-image_select" {
		// => Edit Image response needs to be handled here due to the modular nature
		editFound := false
		editMod := ""
		for _, edit := range editmodule.EditList {
			if _, exists := cache.ActiveInteractions[msg.CorrelationID].Values.String["edit-image_select"]; !exists {
				logger.ErrorText(msg.Interaction.GuildID, "Edit Module Select value not present")
				break
			}

			if cache.ActiveInteractions[msg.CorrelationID].Values.String["edit-image_select"] == helpers.LettersNumbersAndDashesOnly(edit.SelectName()) {
				editMod = cache.ActiveInteractions[msg.CorrelationID].Values.String["edit-image_select"]
				editFound = true
				edit.Execute(msg.Interaction, msg.CorrelationID)
			}
		}

		if !editFound {
			logger.ErrorText(msg.Interaction.GuildID, "Unable to find Edit Module: [%s]", editMod)
		}
	} else {
		discord.InteractionResponseHandlers[msg.ObjectID].Execute(msg.Interaction, msg.CorrelationID)
	}
	reporting.Command(config.CommandTypeModuleResponse, msg.Interaction.GuildID, msg.Interaction.Member.User.ID, msg.Interaction.Member.User.Username, msg.ObjectID, msg.CorrelationID, timeStarted)
}

func workerPhrase(msg *PhraseTaskDetails) {
	var notifyPhrases []string
	timeStarted := time.Now()

	// Process all generic Trigger Phrases first
	for _, phrase := range msg.TriggerPhrases {
		if phrase.IsSpecial {
			continue
		}
		reporting.Command(config.CommandTypePhrase, msg.Message.GuildID, msg.Message.Author.ID, msg.Message.Author.Username, phrase.Phrase, uuid.New().String(), timeStarted)
		if phrase.NotifyOnDetection {
			notifyPhrases = append(notifyPhrases, phrase.Phrase)
		}
	}

	if len(notifyPhrases) > 0 {
		showText := strings.ToUpper(helpers.ConcatStringWithAnd(notifyPhrases)) + " MENTIONED"
		_, err := config.Session.ChannelMessageSend(msg.Message.ChannelID, showText)
		if err != nil {
			logger.Error(msg.Message.GuildID, err)
		}
	}

	// Now process any Special ones
	for _, phrase := range msg.TriggerPhrases {
		if !phrase.IsSpecial {
			continue
		}
		switch phrase.Phrase {
		case "jason statham":
			webm, err := database.GetRandomResource(msg.Message.GuildID, 1)
			if err != nil {
				continue
			}
			discord.Message_ReplyWithMessage(msg.Message.Message, false, fmt.Sprintf("[%s](%s)", "Jason Statham?", webm))
		default:
			logger.ErrorText(msg.Message.GuildID, "Unhandled Special Phrase [%v]", phrase.Phrase)
		}
	}
}

func workerReaction(msg *discordgo.Message) {
	// 1. Get the Reactions
	upCount := 0
	upString := ""
	downCount := 0
	downString := ""
	var userIds map[string]interface{} = make(map[string]interface{})

	for _, reaction := range msg.Reactions {
		emojiIdentifier := reaction.Emoji.Name
		if reaction.Emoji.ID != "" {
			emojiIdentifier = fmt.Sprintf("%s:%s", reaction.Emoji.Name, reaction.Emoji.ID)
		}
		users, err := config.Session.MessageReactions(msg.ChannelID, msg.ID, emojiIdentifier, 100, "", "")
		if err != nil {
			logger.Error(msg.GuildID, err)
			continue
		}
		found := false
		// Check for upvote emojis
		isUpVote := false
		for _, emoji := range reactions.UpvoteEmojis {
			if emojiIdentifier == emoji {
				isUpVote = true
				for _, user := range users {
					if _, exists := userIds[user.ID]; !exists {
						userIds[user.ID] = struct{}{}
						upCount += reaction.Count
						if len(upString) > 0 {
							upString += " "
						}
						upString += addEmojis(emoji, reaction.Count)
						found = true
						break
					}
				}
				break
			}
		}
		if found {
			continue
		}
		// Check for downvote emojis
		isDownVote := false
		for _, emoji := range reactions.DownvoteEmojis {
			if emojiIdentifier == emoji {
				isDownVote = true
				for _, user := range users {
					if _, exists := userIds[user.ID]; !exists {
						userIds[user.ID] = struct{}{}
						downCount += reaction.Count
						if len(downString) > 0 {
							downString += " "
						}
						downString += addEmojis(emoji, reaction.Count)
						found = true
						break
					}
				}
				break
			}
		}
		if !isUpVote && !isDownVote {
			logger.Debug(msg.GuildID, "Unclassified Up/Down Emoji: %s", emojiIdentifier)
		}
	}

	// 2. Check whether the Reaction count passes the threshold for saving
	score := upCount - downCount
	threshold := 3
	if score >= threshold || score <= -threshold {
		emojiString := upString
		if score < 0 {
			emojiString = downString
		}
		reactions.AddOrUpdate(msg, score, emojiString)
	} else {
		reactions.DeleteIfExists(msg)
	}

	logger.Info(msg.GuildID, "Finished Reaction processing for Message ID: %s", msg.ID)
}

func addEmojis(emoji string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		if i > 0 {
			result += " "
		}
		result += emoji
	}
	return result
}
