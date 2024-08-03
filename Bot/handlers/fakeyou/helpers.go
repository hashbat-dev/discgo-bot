package fakeyou

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dabi-ngin/discgo-bot/Bot/audit"
	"github.com/dabi-ngin/discgo-bot/Bot/config"
	"github.com/dabi-ngin/discgo-bot/Bot/constants"
	"github.com/dabi-ngin/discgo-bot/Bot/helpers"
	"github.com/google/uuid"
)

func RequestAudioFromFakeYou(voiceModelID string, textToConvert string, message *discordgo.Message, ttsTypeName string) (string, error) {
	// Build Request
	reqID := uuid.New()
	newRequest := TTSRequest{
		ModelToken:  voiceModelID,
		RequestGuid: reqID.String(),
		RequestText: textToConvert,
	}

	jsonData, err := json.Marshal(newRequest)
	if err != nil {

		return "", err
	}

	byteBuffer := bytes.NewBuffer(jsonData)
	resp, err := http.Post(FakeYouURLRequestTTS, "application/json", byteBuffer)
	if err != nil {
		audit.Error(err)
		return "", err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		config.Session.ChannelMessageEdit(message.ChannelID, message.ID, "Uh oh...")
		return "", fmt.Errorf("HTTP Response was %v", resp.Status)
	}

	// Get a response struct
	var ttsRequest TTSRequestResponse
	err = json.NewDecoder(resp.Body).Decode(&ttsRequest)
	if err != nil {
		config.Session.ChannelMessageEdit(message.ChannelID, message.ID, "Uh oh...")
		return "", err
	}

	// Start checking for Responses
	ttsAudioPath := ""
	messageSuffix := ""

	// Loop a certain number of times
	for i := 0; i < constants.TTS_CHECK_ATTEMPTS; i++ {

		reqCheck, err := helpers.GetJsonFromURL(FakeYouURLCheckRequest + ttsRequest.JobToken)
		if err != nil {
			config.Session.ChannelMessageEdit(message.ChannelID, message.ID, "Uh oh...")
			return "", err
		}

		lastResponse := ""
		if response, ok := reqCheck["state"].(map[string]interface{}); ok {

			status, statusOk := response["status"].(string)
			audioPath, audioPathOk := response["maybe_public_bucket_wav_audio_path"].(string)

			if statusOk && lastResponse != status {
				lastResponse = status
			} else if !statusOk {
				audit.Error(err)
			}

			if audioPathOk && ttsAudioPath != audioPath {
				ttsAudioPath = audioPath
			} else if !statusOk {
				audit.Error(err)
			}

			config.Session.ChannelMessageEdit(message.ChannelID, message.ID, fmt.Sprintf(ttsTypeName+" Request Status: %v"+messageSuffix, lastResponse))
			if len(messageSuffix) >= 5 {
				messageSuffix = ""
			} else {
				messageSuffix += "."
			}

			if status == "complete_failure" || status == "attempt_failed" || status == "dead" {
				config.Session.ChannelMessageDelete(message.ChannelID, message.ID)
				return "", err
			}

			if ttsAudioPath != "" {
				audit.Log(fmt.Sprintf("Got Audio path after %v attempts: %v", i, ttsAudioPath))
				break
			}

		} else {
			audit.Error(err)
		}

		time.Sleep(time.Duration(constants.TTS_CHECK_DELAY) * time.Millisecond)
	}

	config.Session.ChannelMessageDelete(message.ChannelID, message.ID)

	if ttsAudioPath == "" {
		err = errors.New("no tts audio path returned after " + fmt.Sprint(constants.TTS_CHECK_ATTEMPTS) + " attempts")
		audit.Error(err)
		return "", err
	}

	return ttsAudioPath, nil
}
