package fakeyou

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	logger "github.com/dabi-ngin/discgo-bot/Logger"
)

type FakeYouRequestStartOut struct {
	ModelToken  string `json:"tts_model_token"`
	RequestGuid string `json:"uuid_idempotency_token"`
	RequestText string `json:"inference_text"`
}

type FakeYouRequestStartIn struct {
	Success  bool   `json:"success"`
	JobToken string `json:"inference_job_token"`
}

// Create a new Request with FakeYou, returns the FakeYou JobToken or blank if it failed
func CreateRequest(guildId string, correlationId string, voiceModel string, requestText string) string {

	// Build the Request object
	newRequest := FakeYouRequestStartOut{
		ModelToken:  voiceModel,
		RequestGuid: correlationId,
		RequestText: requestText,
	}

	// Marshal the object into JSON
	jsonData, err := json.Marshal(newRequest)
	if err != nil {
		logger.ErrorText(guildId, "Interaction Request ID: [%v] ERROR: %v", correlationId, err)
		return ""
	}

	// Send the request
	byteBuffer := bytes.NewBuffer(jsonData)
	resp, err := http.Post(FakeYouURLRequestTTS, "application/json", byteBuffer)
	if err != nil {
		logger.ErrorText(guildId, "Interaction Request ID: [%v] ERROR: %v", correlationId, err)
		return ""
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("FakeYou HTTP Response upon new request was [%v]", resp.Status)
		logger.ErrorText(guildId, "Interaction Request ID: [%v] ERROR: %v", correlationId, err)
		return ""
	}

	// Get a response struct
	var ttsRequest FakeYouRequestStartIn
	err = json.NewDecoder(resp.Body).Decode(&ttsRequest)
	if err != nil {
		logger.ErrorText(guildId, "Interaction Request ID: [%v] ERROR: %v", correlationId, err)
		return ""
	}

	return ttsRequest.JobToken
}
