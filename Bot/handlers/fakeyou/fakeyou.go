package fakeyou

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ZestHusky/femboy-control/Bot/audit"
	"github.com/ZestHusky/femboy-control/Bot/config"
	"github.com/ZestHusky/femboy-control/Bot/constants"
	dbhelper "github.com/ZestHusky/femboy-control/Bot/dbhelpers"
	"github.com/ZestHusky/femboy-control/Bot/helpers"
	"github.com/ZestHusky/femboy-control/Bot/logging"
	logger "github.com/ZestHusky/femboy-control/Bot/logging"
	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
)

// https://docs.fakeyou.com

// Variables
var FakeYouURLModelList string = "https://api.fakeyou.com/tts/list"
var FakeYouURLRequestTTS string = "https://api.fakeyou.com/tts/inference"
var FakeYouURLCheckRequest string = "https://api.fakeyou.com/tts/job/"
var FakeYouTTSAudioBaseUrl string = "https://storage.googleapis.com/vocodes-public"

// Structs
type TTSRequest struct {
	ModelToken  string `json:"tts_model_token"`
	RequestGuid string `json:"uuid_idempotency_token"`
	RequestText string `json:"inference_text"`
}

type TTSRequestResponse struct {
	Success  bool   `json:"success"`
	JobToken string `json:"inference_job_token"`
}

// Methods
func Search(interaction *discordgo.InteractionCreate) {

	if !helpers.IsAdmin(interaction.Member) {
		audit.Log("Search not from Admin")
		logger.SendErrorMsgInteraction(interaction, "Not an Admin", "This command is only for Admins dummy!", true)
		return
	}

	// Get the Search term
	searchTerm := ""

	options := interaction.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	if opt, ok := optionMap["search"]; ok {
		if opt.Value != "" {
			searchTerm = strings.ToLower(opt.StringValue())
		}
	}

	// Send the FIRST Interaction
	e := embed.NewEmbed()
	e.SetTitle("FakeYou Search - Processing")
	e.SetDescription("Searching TTS Models for '" + searchTerm + "' " + helpers.GetEmote("read", true))

	var embedContent []*discordgo.MessageEmbed
	embedContent = append(embedContent, e.MessageEmbed)
	if embedContent[0].Title == "Error" {
		return
	}

	config.Session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: embedContent,
		},
	})

	modelList, err := helpers.GetJsonFromURL(FakeYouURLModelList)
	if err != nil {
		audit.Error(err)
		logger.SendErrorMsgInteraction(interaction, "Error", "An Error occurred, check logs for more information", true)
		return
	}

	var searchResults []string

	if models, ok := modelList["models"].([]interface{}); ok {
		for _, model := range models {
			if modelMap, ok := model.(map[string]interface{}); ok {
				title, titleOk := modelMap["title"].(string)
				token, tokenOk := modelMap["model_token"].(string)

				if titleOk && tokenOk {
					if !strings.Contains(strings.ToLower(title), searchTerm) {
						continue
					}
					searchResults = append(searchResults, title+" `"+token+"`")
				} else {
					audit.Error(err)
				}
			} else {
				audit.Error(err)
				logger.SendErrorMsgInteraction(interaction, "Error", "An Error occurred, check logs for more information", false)
				return
			}
		}
		audit.Log("Search completed for: " + searchTerm)
	} else {
		audit.Error(err)
		logger.SendErrorMsgInteraction(interaction, "Error", "An Error occurred, check logs for more information", false)
		return
	}

	interactionFinal := fmt.Sprintf("Search completed, found %v results", len(searchResults))
	e.SetDescription(interactionFinal)
	var secondEmbed []*discordgo.MessageEmbed
	secondEmbed = append(secondEmbed, e.MessageEmbed)
	if secondEmbed[0].Title == "Error" {
		audit.Error(err)
	}
	var newResponse *discordgo.WebhookEdit = &discordgo.WebhookEdit{
		Embeds: &secondEmbed,
	}

	_, err = config.Session.InteractionResponseEdit(interaction.Interaction, newResponse)

	if err != nil {
		audit.Error(err)
		logger.SendErrorMsgInteraction(interaction, "Error", "An Error occurred, check logs for more information", false)
	}

	concatMessage := ""
	iterate := 0
	for _, result := range searchResults {

		if len(concatMessage+"\n"+result) > constants.MAX_MESSAGE_LENGTH {
			_, err = config.Session.ChannelMessageSend(interaction.ChannelID, concatMessage)
			if err != nil {
				audit.Error(err)
			}

			concatMessage = ""
			iterate = 0
		}

		if iterate > 0 {
			concatMessage += "\n"
		}

		iterate++
		concatMessage += result

	}

	if concatMessage != "" {
		_, err = config.Session.ChannelMessageSend(interaction.ChannelID, concatMessage)
		if err != nil {
			audit.Error(err)
		}
	}

}

func Add(interaction *discordgo.InteractionCreate) {

	err := AddHandler(interaction)
	if err != nil {
		audit.Error(err)
		logger.SendErrorMsgInteraction(interaction, "Error adding TTS Model", err.Error(), false)
	} else {
		logger.SendMessageInteraction(interaction, "TTS Model Added", "Ready to use! "+helpers.GetEmote("yep", true), "", "", false)
	}
}

func AddHandler(interaction *discordgo.InteractionCreate) error {
	optionMap := helpers.GetOptionMap(interaction)
	inCommand := helpers.GetOptionStringValue(optionMap, "command")
	inModel := helpers.GetOptionStringValue(optionMap, "model")
	inDescription := helpers.GetOptionStringValue(optionMap, "description")

	if inCommand == "" || inModel == "" || inDescription == "" {
		return errors.New("command, model and description are all required")
	}

	ttsResponse := dbhelper.DoesTTSModelExist(inCommand, inModel)
	if ttsResponse != "" {
		return errors.New(ttsResponse)
	}

	err := dbhelper.InsertTTSModel(inCommand, inModel, inDescription, interaction.Member.User.ID)
	if err != nil {
		return err
	}

	return nil
}

func Update(interaction *discordgo.InteractionCreate) {
	err := UpdateHandler(interaction)
	if err != nil {
		audit.Error(err)
		logger.SendErrorMsgInteraction(interaction, "Error adding TTS Model", err.Error(), false)
	} else {
		logger.SendMessageInteraction(interaction, "TTS Model Added", "Ready to use! "+helpers.GetEmote("yep", true), "", "", false)
	}
}

func UpdateHandler(interaction *discordgo.InteractionCreate) error {
	optionMap := helpers.GetOptionMap(interaction)
	inCommand := helpers.GetOptionStringValue(optionMap, "command")
	inNewCommand := helpers.GetOptionStringValue(optionMap, "new-command")
	inNewModel := helpers.GetOptionStringValue(optionMap, "new-model")
	inDescription := helpers.GetOptionStringValue(optionMap, "new-description")

	ttsResponse := dbhelper.DoesTTSModelExist(inCommand, "")
	if ttsResponse != "Command already exists" {
		return errors.New("command to edit does not exist")
	}

	if inNewCommand == "" && inNewModel == "" && inDescription == "" {
		return errors.New("new command or model required")
	}

	err := dbhelper.UpdateTTSModel(inCommand, inNewCommand, inNewModel, inDescription, interaction.Member.User.ID)
	if err != nil {
		return err
	}

	return nil
}

func FakeYouGetTTS(message *discordgo.MessageCreate, voice string, text string) bool {

	delOrig, err := RequestTTS(message, voice, text, "", "")
	if err != nil {
		audit.Error(err)
	}
	return delOrig
}

func RequestTTS(message *discordgo.MessageCreate, voice string, text string, staticResponse string, ttsObjectName string) (bool, error) {

	// Build a Message Response first
	if ttsObjectName == "" {
		ttsObjectName = "TTS"
	}

	msg, err := config.Session.ChannelMessageSend(message.ChannelID, "Requesting "+ttsObjectName+"...")

	if err != nil {
		audit.Error(err)
	}

	// Can we get a Model?
	token, err := dbhelper.GetTTSModel(voice)
	if err != nil {
		config.Session.ChannelMessageEdit(message.ChannelID, message.ID, "Uh oh...")
		return false, err
	} else if token == "" {
		config.Session.ChannelMessageEdit(message.ChannelID, message.ID, "Couldn't find a "+ttsObjectName+" Model for !tts"+voice)
		return false, errors.New("did not find TTS model for command: " + voice)
	}

	ttsText := strings.TrimSpace(text)
	replyMessage := ""
	if message.ReferencedMessage != nil {
		replyMessage = message.ReferencedMessage.ID
		if ttsText == "" {
			ttsText = strings.TrimSpace(message.ReferencedMessage.Content)
		}
	}

	delOrig := replyMessage != ""

	if ttsText == "" {
		config.Session.ChannelMessageEdit(message.ChannelID, message.ID, "Uh oh...")
		return false, errors.New("no text found! " + helpers.GetEmote("ae_cry", true))
	}

	ttsAudioPath, err := RequestAudioFromFakeYou(token, ttsText, msg, ttsObjectName)
	if err != nil {
		config.Session.ChannelMessageEdit(message.ChannelID, message.ID, "Error getting "+ttsObjectName+" >w<...\n"+err.Error())
	} else if ttsAudioPath == "" {
		config.Session.ChannelMessageEdit(message.ChannelID, message.ID, "Error getting "+ttsObjectName+" >w<...\nThe API didn't give me a URL!")
	} else {
		// Delete temp Message
		config.Session.ChannelMessageDelete(message.ChannelID, message.ID)

		// Download the .wav file and Send it
		fileURL := FakeYouTTSAudioBaseUrl + ttsAudioPath

		// Download the .wav file from the URL
		filePath, err := helpers.DownloadFileFromURL(fileURL)
		if err != nil {
			return false, err
		}
		defer os.Remove(filePath) // Delete the temporary file once done

		// Create a File object for the downloaded .wav file
		file, err := os.Open(filePath)
		if err != nil {
			return false, err
		}
		defer file.Close()

		// Seek to the beginning of the file
		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			return false, err
		}

		// Create a File object with the temporary file
		fileName := staticResponse
		if fileName == "" {
			fileName = ttsText
		}
		fileObj := &discordgo.File{
			Name:   helpers.ConvertTextToFilename(fileName, 150) + ".wav",
			Reader: file,
		}

		if message.ReferencedMessage == nil {
			_, err = config.Session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
				Files: []*discordgo.File{fileObj},
			})

			if err != nil {
				return false, err
			}
		} else {
			msgRef := logger.MessageRefObj(message.ReferencedMessage)
			_, err = config.Session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
				Files:     []*discordgo.File{fileObj},
				Reference: &msgRef,
			})

			if err != nil {
				return false, err
			}
		}
	}

	return delOrig, nil
}

func GiveTTSHandler(interaction *discordgo.InteractionCreate) {
	fmt.Println("[GIVETTS] Received Give TTS List Handler command")

	ttsList, err := dbhelper.GetTTSList()
	if err != nil {
		logging.SendErrorInteraction(interaction)
	}

	e := embed.NewEmbed()
	e.SetTitle("TTS Voices")
	descText := "Type '!ttsexample Your Text' or reply to a Message with !ttsexample to generate AI Text to Speech\n"
	for _, tts := range ttsList {
		descText += "\n**!tts" + tts.Command + "** - " + tts.Description
	}
	e.SetDescription(descText)

	var embedContent []*discordgo.MessageEmbed
	embedContent = append(embedContent, e.MessageEmbed)
	if embedContent[0].Title == "Error" {
		return
	}

	config.Session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: embedContent,
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})
}
