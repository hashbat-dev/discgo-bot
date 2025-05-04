package external

import (
	"encoding/json"
	"io"
	"net/http"

	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

func GetJsonFromUrl(url string) (map[string]interface{}, error) {
	// Make HTTP GET request
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := response.Body.Close()
		if err != nil {
			logger.Error("EXTERNAL", err)
		}
	}()

	// Read response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Parse JSON
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
