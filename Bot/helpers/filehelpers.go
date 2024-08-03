package helpers

import (
	"bytes"
	"errors"
	"fmt"
	"image/gif"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/ZestHusky/femboy-control/Bot/audit"
	"github.com/ZestHusky/femboy-control/Bot/constants"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

var validExtensions []string = []string{
	".gif",
	".png",
	".jpg",
	".webp",
}

func DoesMessageHaveImage(message *discordgo.Message, requiredExtension string) (string, error) {

	msgContent := strings.Trim(message.Content, " ")
	msgContentLower := strings.ToLower(msgContent)

	// 1. Check for a file =========================================
	// A. Are there any Embeds?
	imgLink := ""
	if len(message.Embeds) > 0 {
		imgLink = message.Embeds[0].Thumbnail.ProxyURL
	}

	// B. Are there any Attachments?
	if imgLink == "" {
		imgLink = message.Attachments[0].ProxyURL
	}

	// C. Is this a Tenor link?
	if strings.Contains(msgContentLower, "tenor.com/") {
		tenorLink, err := GetGifURLFromTenorLink(msgContentLower)
		if err != nil {
			return "", err
		} else if tenorLink != "" {
			imgLink = tenorLink
		}
	}

	// D. Okay last check, what about the body content?
	if imgLink == "" {
		if requiredExtension == "" {
			for _, ext := range validExtensions {
				if strings.Contains(msgContentLower, ext) {
					imgLink = msgContent
					break
				}
			}
		} else {
			if strings.Contains(msgContentLower, requiredExtension) {
				imgLink = msgContent
			}
		}
	}

	// D. No file?
	if imgLink == "" {
		return "", errors.New(ErrorText("no suitable image found"))
	}

	// 2. Now lets validate the file name (if any) ================================
	extValid := false
	if requiredExtension == "" {
		for _, ext := range validExtensions {
			if strings.Contains(imgLink, ext) {
				extValid = true
				break
			}
		}
	} else {
		if strings.Contains(imgLink, requiredExtension) {
			extValid = true
		}
	}

	if !extValid {
		return "", errors.New(ErrorText("image was not a suitable extension type"))
	}

	// 3. Return! =====================================================================
	return imgLink, nil
}

func ErrorText(append string) string {
	errText := ""
	if append != "" {
		errText += append + " (Keep in mind I accept .gif, .png, .jpg or .webp :3)"
	} else {
		errText += "there was an error (Keep in mind I accept .gif, .png, .jpg or .webp :3)"
	}
	errText += "\n\nI also can't accept Tenor links that aren't .gif links, discord's a bit fucky with these >.<."
	errText += "\n\nRight click the gif and click 'Copy Link' to see if it's a .gif link"
	errText += "\n\nYou can get around this by saving the gif to your machine and sending it back into Discord as a message, then adding that"
	return errText
}

func DownloadFileFromURL(url string) (string, error) {
	// Get the file name from the URL
	fileName, extension := GetFileNameAndExtension(url)

	// Create a temporary file to store the downloaded file
	tempFile, err := os.CreateTemp("", fileName+"*"+extension)
	if err != nil {
		audit.Error(err)
		return "", err
	} else {
		audit.Log("TempFile created: [" + url + "] => [" + tempFile.Name() + "]")
	}
	defer tempFile.Close()

	// Download the file from the URL
	resp, err := http.Get(url)
	if err != nil {
		audit.Error(err)
		return "", err
	}
	defer resp.Body.Close()

	// Copy the downloaded file contents to the temporary file
	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		audit.Error(err)
		return "", err
	}
	return tempFile.Name(), nil
}

func GetFileNameAndExtension(file string) (string, string) {
	fullFileName := filepath.Base(file)
	extension := filepath.Ext(fullFileName)
	fileName := strings.Replace(fullFileName, extension, "", -1)
	if strings.Contains(extension, "?") {
		extension = strings.Split(extension, "?")[0]
	}
	return fileName, extension
}

func GetGifURLFromTenorLink(tenorLink string) (string, error) {

	// Step 1: Fetch the HTML from the URL
	res, err := http.Get(tenorLink)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", errors.New("failed to fetch Tenor GIF from URL, response code: " + fmt.Sprint(res.StatusCode))
	}

	// Step 2: Parse the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalf("Failed to parse the HTML document: %v", err)
	}

	// Step 3: Find the GIF on the page
	selection := doc.Find("#single-gif-container").Find("div.Gif").Find("img")
	href, exists := selection.Attr("src")
	if !exists {
		return "", errors.New("couldn't find gif on Tenor page")
	}

	return href, nil
}

func DownloadGif(gifLink string) (*gif.GIF, error) {

	// Download the image from message content to a []byte type
	response, err := http.Get(gifLink)
	if err != nil {
		audit.Error(err)
		return nil, err
	}
	defer response.Body.Close()

	randomName := uuid.New().String()
	tempFile, err := os.CreateTemp(constants.TEMP_DIRECTORY, randomName)
	if err != nil {
		audit.Error(err)
		return nil, err
	}
	defer tempFile.Close()

	if err := tempFile.Close(); err != nil {
		return nil, err
	}

	newName := tempFile.Name() + ".gif"
	err = os.Rename(tempFile.Name(), newName)
	if err != nil {
		audit.Error(err)
		return nil, err
	}

	file, err := os.OpenFile(newName, os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	} else {
		audit.Log("Created Temporary file for GIF: " + file.Name())
	}
	defer os.Remove(file.Name())

	// Copy HTTP response to file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		audit.Error(err)
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	buffer := make([]byte, info.Size())
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, err
	}

	gifData, err := gif.DecodeAll(bytes.NewReader(buffer))
	if err != nil {
		return nil, err
	}
	return gifData, nil
}

func CheckDirectoryExists(dir string) bool {

	// Get the file or directory information
	info, err := os.Stat(dir)

	// If the directory does not exist, create it
	if os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755) // 0755 permissions mean readable and executable by everyone, writable by the owner
		if err != nil {
			audit.Error(err)
			return false
		}

		if os.IsNotExist(err) {
			return false
		} else {
			return true
		}
	} else if err != nil {
		audit.Error(err)
		return false
	} else if !info.IsDir() {
		audit.Error(err)
		return false
	} else {
		return true
	}
}
