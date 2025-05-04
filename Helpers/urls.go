package helpers

import (
	"fmt"
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
