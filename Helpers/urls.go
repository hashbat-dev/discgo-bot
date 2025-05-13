package helpers

import (
	"fmt"
	"io"
	"net/http"

	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

func DoesLinkWork(link string) (bool, error) {
	resp, err := http.Head(link)
	if err != nil {
		return false, fmt.Errorf("error making request: %v", err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			logger.Error("EXTERNAL", err)
		}
	}()

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}
	return true, nil
}

func GetBytesFromURL(guildId string, url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		logger.Error(guildId, err)
		return nil, err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			logger.Error(guildId, err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(guildId, err)
		return nil, err
	}

	return body, nil
}
