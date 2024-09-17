package imgur

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	config "github.com/dabi-ngin/discgo-bot/Config"
	database "github.com/dabi-ngin/discgo-bot/Database"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
	"github.com/google/uuid"
)

// Returns Imgur Link, Delete Hash, Error
func UploadAndGetUrl(guildId string, userId string, file io.Reader) (string, string, error) {

	// 1. Convert the File to a Base64 encoded string
	data, err := io.ReadAll(file)
	if err != nil {
		logger.Error(guildId, err)
		return "", "", err
	}
	base64String := base64.StdEncoding.EncodeToString(data)

	// 2. Create the Imgur API Request
	payload := &bytes.Buffer{}
	generatedTitle := uuid.New().String()
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("image", base64String)
	_ = writer.WriteField("type", "base64")
	_ = writer.WriteField("title", generatedTitle)
	err = writer.Close()
	if err != nil {
		logger.Error(guildId, err)
		return "", "", err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", BaseUrl, payload)
	if err != nil {
		logger.Error(guildId, err)
		return "", "", err
	}
	req.Header.Add("Authorization", "Client-ID "+config.ServiceSettings.IMGURCLIENTID)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 3. Send the Request
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(guildId, err)
		return "", "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(guildId, err)
		return "", "", err
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http status returned as %s", resp.Status)
		logger.Error(guildId, err)
		return "", "", err
	}

	// 4. If successful, log information
	var imgurResp UploadResponse
	err = json.Unmarshal([]byte(body), &imgurResp)
	if err != nil {
		logger.Error(guildId, err)
		return "", "", err
	}

	if !imgurResp.Success {
		err = fmt.Errorf("imgur provided success=false")
		logger.Error(guildId, err)
		return "", "", err
	}

	err = database.InsertImgurLog(guildId, userId, imgurResp.Data.ID, imgurResp.Data.Type, imgurResp.Data.Title, imgurResp.Data.Link, imgurResp.Data.DeleteHash)
	if err != nil {
		logger.Error(guildId, err)
	}

	return imgurResp.Data.Link, imgurResp.Data.DeleteHash, nil
}
