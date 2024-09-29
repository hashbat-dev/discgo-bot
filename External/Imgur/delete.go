package imgur

import (
	"bytes"
	"mime/multipart"
	"net/http"

	config "github.com/hashbat-dev/discgo-bot/Config"
	database "github.com/hashbat-dev/discgo-bot/Database"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

func DeleteImgurEntry(guildId string, deleteHash string) error {
	// 1. Send the Delete Request
	url := BaseUrl + deleteHash
	method := "DELETE"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	err := writer.Close()
	if err != nil {
		logger.Error(guildId, err)
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		logger.Error(guildId, err)
		return err
	}
	req.Header.Add("Authorization", "Client-ID "+config.ServiceSettings.IMGURCLIENTID)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		logger.Error(guildId, err)
		return err
	}
	defer res.Body.Close()

	// 2. If successful Delete from the Database
	return database.DeleteImgurLog(guildId, deleteHash)
}
