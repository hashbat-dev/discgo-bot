package fakeyou

import (
	"fmt"
	"io"
	"net/http"

	logger "github.com/dabi-ngin/discgo-bot/Logger"
)

func DownloadFile(guildId string, correlationId string, audioPath string) (io.Reader, error) {
	url := FakeYouTTSAudioBaseUrl + audioPath

	resp, err := http.Get(url)
	if err != nil {
		logger.ErrorText(guildId, "Interaction Request ID: [%v] ERROR: %v", correlationId, err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("HTTP Response code returned as [%v]", resp.Status)
		logger.ErrorText(guildId, "Interaction Request ID: [%v] ERROR: %v", correlationId, err)
		return nil, err
	}

	return resp.Body, nil
}
