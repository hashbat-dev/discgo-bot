package helpers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func GetImageUrlFromTenor(tenorLink string) (string, error) {

	// Step 1: Fetch the HTML from the URL
	res, err := http.Get(tenorLink)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", fmt.Errorf("failed to fetch Tenor GIF from URL, Tenor responded with Status Code: %d", res.StatusCode)
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
