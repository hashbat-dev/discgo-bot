package fakeyou

import (
	"fmt"

	external "github.com/hashbat-dev/discgo-bot/External"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

// Returns [Status, .wav Audio Path, Error]
func CheckRequest(guildId string, correlationId string, jobToken string) (string, string, error) {
	reqCheck, err := external.GetJsonFromUrl(FakeYouURLCheckRequest + jobToken)
	if err != nil {
		logger.ErrorText(guildId, "Interaction Request ID: [%v] ERROR: %v", correlationId, err)
		return "", "", err
	}

	if response, ok := reqCheck["state"].(map[string]interface{}); ok {

		status, statusOk := response["status"].(string)
		audioPath, audioPathOk := response["maybe_public_bucket_wav_audio_path"].(string)

		returnStatus := ""
		returnAudioPath := ""

		if statusOk {
			returnStatus = status
		}

		if audioPathOk {
			returnAudioPath = audioPath
		}

		return returnStatus, returnAudioPath, nil

	} else {
		err = fmt.Errorf("Interaction Request ID: [%v] Could not marshal returned response into an interface", correlationId)
		return "", "", err
	}
}
