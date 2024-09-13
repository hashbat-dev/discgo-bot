package slash

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	cache "github.com/dabi-ngin/discgo-bot/Cache"
	config "github.com/dabi-ngin/discgo-bot/Config"
	database "github.com/dabi-ngin/discgo-bot/Database"
	discord "github.com/dabi-ngin/discgo-bot/Discord"
	fakeyou "github.com/dabi-ngin/discgo-bot/External/FakeYou"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var resultsEmoji discordgo.ComponentEmoji = discordgo.ComponentEmoji{Name: "🎤"}

func TtsPlay(i *discordgo.InteractionCreate, correlationId string) {
	cachedInteraction := cache.ActiveInteractions[correlationId]
	searchTerm := cachedInteraction.Values.String["voice"]
	ttsText := cachedInteraction.Values.String["text"]

	// 1. Search the Database for TTS Models
	modelResults, err := database.GetFakeYouModels(searchTerm)
	if err != nil {
		logger.Error(i.GuildID, err)
		discord.SendGenericErrorFromInteraction(i)
		return
	}

	if len(modelResults) == 0 {
		// No Results
		retText := fmt.Sprintf("No voice models were found for the search term: '%v'", searchTerm)
		retText += "\nPlease try again, your requested text is below:\n\n" + ttsText
		discord.SendEmbedFromInteraction(i, "No Results", retText)
		cache.InteractionComplete(correlationId)
		return
	} else if len(modelResults) == 1 {
		// 1 Result, straight to the queue!
		for key, value := range modelResults {
			TtsPlayCreateRequest(correlationId, key, value.ModelToken, cachedInteraction.Values.String["text"])
			return
		}
	}

	// 2. Put the first 25 Results into the Select menu
	var resultList []discordgo.SelectMenuOption
	j := 0
	for _, result := range modelResults {
		j++
		if j > config.MAX_SELECT_LENGTH {
			break
		}
		resultList = append(resultList, discordgo.SelectMenuOption{
			Label: result.Title,
			Value: result.ModelToken,
			Emoji: resultsEmoji,
		})
	}

	// 3. Create an action row and add the select menu to it
	selectMenu := discord.CreateSelectMenu(discordgo.SelectMenu{
		CustomID:    "tts-play_select-model",
		Options:     resultList,
		Placeholder: "Choose Voice Model to use...",
	}, correlationId, config.IO_BOUND_TASK, TtsPlaySelectModel)

	actionRow := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			selectMenu,
		},
	}

	// 4. Generate the text to show the user
	showText := fmt.Sprintf("Found %v results for '%v'.", j, searchTerm)
	if len(modelResults) > config.MAX_SELECT_LENGTH {
		showText = fmt.Sprintf("Showing first %v results of %v for '%v'.", j, len(modelResults), searchTerm)
		showText += "\nPlease refine your search term for more accurate results."
	}

	// 5. Send the Select menu response
	err = config.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    showText,
			Components: []discordgo.MessageComponent{actionRow},
			Flags:      discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		logger.Error(i.GuildID, err)
	}

}

func TtsPlaySelectModel(i *discordgo.InteractionCreate, correlationId string) {
	cachedInteraction := cache.ActiveInteractions[correlationId]

	// 1. Can we get the selected Voice Model?
	voiceModel := cachedInteraction.Values.String["tts-play_select-model"]
	if voiceModel == "" {
		logger.ErrorText(i.GuildID, "VoiceModel property was blank, cannot process request")
		discord.UpdateInteractionResponseWithGenericError(&cachedInteraction.StartInteraction)
		cache.InteractionComplete(correlationId)
		return
	}

	// 2. If so, create a new Request
	TtsPlayCreateRequest(correlationId, "", voiceModel, cachedInteraction.Values.String["text"])
}

type ActiveRequest struct {
	CorrelationId    string
	Interaction      discordgo.InteractionCreate
	RequestModelName string
	RequestModelID   string
	RequestText      string
	LoadingSpinner   string
	Started          time.Time
	FakeYouJobToken  string
	LastStatus       string
	UpdateCount      int
	ErrorCount       int
}

var ActiveRequests map[string]ActiveRequest = make(map[string]ActiveRequest)
var RequestEmbedTitle string = "Text-to-Speech Request"

func requestEmbedText(voiceModelName string, requestedText string, loadingText string, loadingSpinner string) (string, string) {
	desc := fmt.Sprintf("**Voice**: %s\n**Text**: %s\n\n%s%s", voiceModelName, requestedText, loadingText, loadingSpinner)
	return RequestEmbedTitle, desc
}

func incrementLoadingSpinner(currentSpinner string) string {
	if len(currentSpinner) > 4 {
		return "..."
	} else {
		return currentSpinner + "."
	}
}

func getUserFriendlyStatus(status string) string {
	status = strings.ReplaceAll(status, "_", " ")
	words := strings.Fields(status)
	caser := cases.Title(language.English)
	for i, word := range words {
		words[i] = caser.String(word)
	}

	return strings.Join(words, " ")
}

func TtsPlayCreateRequest(correlationId string, requestModelName string, requestModelId string, requestText string) {

	cachedInteraction := cache.ActiveInteractions[correlationId]

	// 1. Get the Model Name from the Token if blank
	if requestModelName == "" {
		modelName, err := database.GetModelNameFromToken(cachedInteraction.StartInteraction.GuildID, requestModelId)
		if err != nil {
			modelName = "Unknown"
		}
		requestModelName = modelName
	}

	// 2. Insert the Request
	ActiveRequests[correlationId] = ActiveRequest{
		CorrelationId:    correlationId,
		Interaction:      cachedInteraction.StartInteraction,
		RequestModelName: requestModelName,
		RequestModelID:   requestModelId,
		RequestText:      requestText,
		LastStatus:       "Creating Request",
		LoadingSpinner:   "...",
		Started:          time.Now(),
	}

	// 3. Update with the Initial Embed
	title, desc := requestEmbedText(requestModelName, requestText, "Creating Request", "...")
	discord.UpdateInteractionResponse(&cachedInteraction.StartInteraction, title, desc)

	// From here, the Scheduler will begin processing queued items in the function below
}

// Called every 2 seconds from the Scheduler
func ProcessQueue() {
	for k, v := range ActiveRequests {
		go func(key string, value ActiveRequest) {
			value.UpdateCount++
			value.LoadingSpinner = incrementLoadingSpinner(value.LoadingSpinner)

			// 1. Has the Request exceed the Maximum update attempts?
			if value.UpdateCount > config.MaxFakeYouRequestChecks {
				sendText := fmt.Sprintf("Text-to-Speech request hit the Maximum number of attempts (%v), please try again later.", value.UpdateCount)
				discord.UpdateInteractionResponse(&value.Interaction, RequestEmbedTitle, sendText)
				cache.InteractionComplete(value.CorrelationId)
				delete(ActiveRequests, key)
				return
			}

			// 2. Has the Request been started?
			if value.FakeYouJobToken == "" {
				jobToken := fakeyou.CreateRequest(value.Interaction.GuildID, value.CorrelationId, value.RequestModelID, value.RequestText)
				if jobToken == "" {
					// Request failed, provide error and complete the request
					discord.UpdateInteractionResponseWithGenericError(&value.Interaction)
					cache.InteractionComplete(value.CorrelationId)
					delete(ActiveRequests, key)
				} else {
					// Update the Request
					value.FakeYouJobToken = jobToken
					value.LastStatus = "Request Sent"
					title, desc := requestEmbedText(value.RequestModelName, value.RequestText, value.LastStatus, value.LoadingSpinner)
					discord.UpdateInteractionResponse(&value.Interaction, title, desc)
					ActiveRequests[key] = value
				}
				return
			}

			// 3. Request already Started, get its status
			status, audioPath, err := fakeyou.CheckRequest(value.Interaction.GuildID, value.CorrelationId, value.FakeYouJobToken)
			if err != nil {
				value.ErrorCount++
				if value.ErrorCount >= config.MaxFakeYouRequestErrors {
					// Too many errors, provide error and complete the request
					discord.UpdateInteractionResponseWithGenericError(&value.Interaction)
					cache.InteractionComplete(value.CorrelationId)
					delete(ActiveRequests, key)
					return
				} else {
					// Try again
					title, desc := requestEmbedText(value.RequestModelName, value.RequestText, value.LastStatus, value.LoadingSpinner)
					discord.UpdateInteractionResponse(&value.Interaction, title, desc)
					ActiveRequests[key] = value
					return
				}
			}

			// 4. Not completed yet, post update and continue
			if audioPath == "" {
				if status == "complete_failure" || status == "attempt_failed" || status == "dead" {
					// Has the Request errored at FakeYou's end?
					sendText := fmt.Sprintf("Text-to-Speech request failed after %v attempts with the status '%v'", value.UpdateCount, getUserFriendlyStatus(status))
					discord.UpdateInteractionResponse(&value.Interaction, RequestEmbedTitle, sendText)
					cache.InteractionComplete(value.CorrelationId)
					delete(ActiveRequests, key)
				} else {
					// Update user with the current status
					value.LastStatus = getUserFriendlyStatus(status)
					title, desc := requestEmbedText(value.RequestModelName, value.RequestText, value.LastStatus, value.LoadingSpinner)
					discord.UpdateInteractionResponse(&value.Interaction, title, desc)
					ActiveRequests[key] = value
				}
				return
			}

			// 5. Request was Completed!
			// => Download the .wav file as an io.Reader
			wavReader, err := fakeyou.DownloadFile(value.Interaction.GuildID, value.CorrelationId, audioPath)
			if err != nil {
				sendText := "Text-to-Speech request failed, could not download the result"
				discord.UpdateInteractionResponse(&value.Interaction, RequestEmbedTitle, sendText)
				cache.InteractionComplete(value.CorrelationId)
				delete(ActiveRequests, key)
				return
			}

			// => Send the Result
			fileObj := &discordgo.File{
				Name:   value.CorrelationId + ".wav",
				Reader: wavReader,
			}

			message, err := config.Session.ChannelMessageSendComplex(value.Interaction.ChannelID, &discordgo.MessageSend{
				Files: []*discordgo.File{fileObj},
			})

			if err != nil {
				logger.Error(value.Interaction.GuildID, err)
				sendText := "Text-to-Speech request failed, could not send the downloaded result"
				discord.UpdateInteractionResponse(&value.Interaction, RequestEmbedTitle, sendText)
				cache.InteractionComplete(value.CorrelationId)
				delete(ActiveRequests, key)
				return
			}

			// => Delete the Requesting Interaction
			err = config.Session.InteractionResponseDelete(value.Interaction.Interaction)
			if err != nil {
				logger.Error(value.Interaction.GuildID, err)
			}

			// => Log the Final Result in the Database
			database.InsertFakeYouLog(value.Interaction.GuildID, value.Interaction.Member.User.ID, value.CorrelationId, message.ID, value.RequestModelName,
				value.RequestModelID, value.RequestText, value.UpdateCount)

			// 6. Complete the Request
			logger.Event(value.Interaction.GuildID, "[%v] Completed /tts-play interaction after %v attempts", value.CorrelationId, value.UpdateCount)
			cache.InteractionComplete(value.CorrelationId)
			delete(ActiveRequests, key)
		}(k, v)
	}
}
