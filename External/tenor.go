package external

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
)

func GetImageUrlFromTenor(guildId string, tenorLink string) (string, error) {
	// Step 1: Fetch the HTML from the URL
	res, err := http.Get(tenorLink)
	if err != nil {
		logger.Error(guildId, err)
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		err = fmt.Errorf("failed to fetch Tenor GIF from URL, Tenor responded with Status Code: %d", res.StatusCode)
		logger.Error(guildId, err)
		return "", err
	}

	// Step 2: Parse the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		logger.Error(guildId, err)
		return "", err
	}

	// Step 3: Get the Tenor ID from the URL
	lastSlashIndex := strings.LastIndex(tenorLink, "/")
	secondLastSlashIndex := strings.LastIndex(tenorLink[:lastSlashIndex], "/")
	var tenorId string
	if secondLastSlashIndex != -1 && lastSlashIndex != -1 && secondLastSlashIndex < lastSlashIndex {
		tenorId = tenorLink[secondLastSlashIndex+1 : lastSlashIndex]
	}

	if tenorId == "" {
		err = errors.New("could not find tenor id in url")
		logger.Error(guildId, err)
		return "", err
	}

	// Step 4: Find the GIF on the page
	gifUrl := ""

	// => Attempt 1: Does the link take us to the full website?
	selection := doc.Find("#single-gif-container").Find("div.Gif").Find("img")
	href, exists := selection.Attr("src")
	if exists {
		gifUrl = href
	}

	// => Attempt 2: Does the link take us to the image on a lightbox?
	if gifUrl == "" {
		doc.Find("img").Each(func(index int, img *goquery.Selection) {
			src, exists := img.Attr("src")
			if exists && strings.Contains(src, tenorId) {
				gifUrl = src
				return
			}
		})
	}

	if gifUrl == "" {
		err = errors.New("could not find gif on tenor page")
		logger.Error(guildId, err)
		return "", err
	}

	return gifUrl, nil
}
