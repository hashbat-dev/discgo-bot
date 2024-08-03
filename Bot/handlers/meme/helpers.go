package meme

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	"github.com/dabi-ngin/discgo-bot/Bot/audit"
	"github.com/dabi-ngin/discgo-bot/Bot/config"
	"github.com/dabi-ngin/discgo-bot/Bot/constants"
	"github.com/dabi-ngin/discgo-bot/Bot/helpers"
)

func GetMeme(searchTerm string, allowStills bool, fromBoard string) (string, string, string) {
	// CATALOG =========================================================================================
	// Loop each board until we find a thread

	var results []ThreadResults

	if fromBoard != "" {
		url := "https://a.4cdn.org/" + fromBoard + "/catalog.json"

		// Send GET request
		response, err := http.Get(url)
		if err != nil {
			audit.Error(err)
		}
		defer response.Body.Close()

		// Decode JSON
		var threadPages ThreadList
		if err := json.NewDecoder(response.Body).Decode(&threadPages); err != nil {
			audit.Error(err)
		}

		// Search the Threads
		for _, page := range threadPages {
			for _, thread := range page.Thread {
				if thread.Images == 0 {
					continue
				}

				title := thread.Text
				if thread.Title != "" {
					title = thread.Title
				}

				results = append(results, ThreadResults{
					Board:     fromBoard,
					Title:     title,
					OPID:      thread.OPNumber,
					FileCount: thread.Images,
				})
			}
		}
	} else {
		var wg sync.WaitGroup
		for _, board := range BoardsToSearch {
			wg.Add(1)

			go func(board string) {
				defer wg.Done()
				url := "https://a.4cdn.org/" + board + "/catalog.json"

				// Send GET request
				response, err := http.Get(url)
				if err != nil {
					audit.Error(err)
					return
				}
				defer response.Body.Close()
				// Decode JSON
				var threadPages ThreadList
				if err := json.NewDecoder(response.Body).Decode(&threadPages); err != nil {
					audit.Error(err)
					return
				}

				// Search the Threads
				for _, page := range threadPages {
					for _, thread := range page.Thread {
						if thread.Images == 0 {
							continue
						}

						if strings.Contains(strings.ToLower(thread.Title), searchTerm) || strings.Contains(strings.ToLower(thread.Text), searchTerm) {
							title := thread.Text
							if thread.Title != "" {
								title = thread.Title
							}

							results = append(results, ThreadResults{
								Title:     title,
								Board:     board,
								OPID:      thread.OPNumber,
								FileCount: thread.Images,
							})
						}
					}
				}
			}(board)
		}
		wg.Wait()
	}

	if len(results) == 0 {
		return "", "", ""
	}

	// Get a Weighted Thread
	var random []ThreadResults
	logText := "Search Term: " + searchTerm + ", Threads:"
	for _, result := range results {
		logText += " [" + result.Board + ", " + fmt.Sprint(result.FileCount) + "]"
		for i := 0; i < result.FileCount-1; i++ {
			random = append(random, result)
		}
	}
	if fromBoard == "" {
		audit.Log(logText)
	}

	// GET RANDOM FILE =================================================================================
	maxCount := 10
	finalBoard := ""
	fileId := 0
	extension := ""
	thread := 0
	postId := 0
	threadTitle := ""

	for i := 0; i < maxCount; i++ {
		randomThread := random[rand.Intn(len(random))]
		randomFile, board := GetRandomFile(randomThread, allowStills)
		if randomFile.Extension != "" && board != "" {
			threadTitle = randomThread.Title
			finalBoard = board
			fileId = randomFile.FileID
			extension = randomFile.Extension
			thread = randomThread.OPID
			postId = randomFile.PostID
			audit.Log("Selected Board: " + finalBoard + ", File: " + fmt.Sprint(randomFile.FileID) + randomFile.Extension + " on attempt " + fmt.Sprint(i+1))
			break
		}
	}

	memeUrl := "https://i.4cdn.org/" + finalBoard + "/" + fmt.Sprint(fileId) + extension
	threadUrl := "https://boards.4chan.org/" + finalBoard + "/thread/" + fmt.Sprint(thread) + "#p" + fmt.Sprint(postId)
	return memeUrl, threadUrl, threadTitle
}

func GetRandomFile(randomThread ThreadResults, allowStills bool) (FoundFiles, string) {
	url := "https://a.4cdn.org/" + randomThread.Board + "/thread/" + fmt.Sprint(randomThread.OPID) + ".json"

	response, err := http.Get(url)
	if err != nil {
		audit.Error(err)
		return FoundFiles{
			Extension: "",
			FileID:    0,
		}, ""
	}
	defer response.Body.Close()

	// Decode JSON
	var threadPosts ThreadPosts
	if err = json.NewDecoder(response.Body).Decode(&threadPosts); err != nil {
		audit.Error(err)
		return FoundFiles{
			PostID:    0,
			Extension: "",
			FileID:    0,
		}, ""
	}

	var fileList []FoundFiles
	for _, post := range threadPosts.Posts {
		if post.Extension == "" || post.FileID == 0 {
			continue
		}

		if allowStills {
			if post.Extension != ".webm" && post.Extension != ".gif" && post.Extension != ".jpg" && post.Extension != ".webp" && post.Extension != ".png" {
				continue
			}
		} else {
			if post.Extension != ".webm" && post.Extension != ".gif" {
				continue
			}
		}

		fileList = append(fileList, FoundFiles{
			PostID:    post.PostID,
			Extension: post.Extension,
			FileID:    post.FileID,
		})
	}

	// Did we find Files?
	if len(fileList) == 0 {
		return FoundFiles{
			PostID:    0,
			Extension: "",
			FileID:    0,
		}, ""
	}

	// Get a random one
	if len(fileList) == 0 {
		return FoundFiles{
			PostID:    0,
			Extension: "",
			FileID:    0,
		}, ""
	} else {
		randomFile := rand.Intn(len(fileList) - 1)
		return fileList[randomFile], randomThread.Board
	}
}

func GetSearchTerm(inSearch string) string {
	searchTerm := ""
	if inSearch == "" {
		randomIndex := rand.Intn(len(GenericSearch))
		searchTerm = strings.ToLower(GenericSearch[randomIndex])
	} else {
		searchTerm = strings.ToLower(inSearch)
	}
	return searchTerm
}

func ErrorEmbed() *discordgo.MessageEmbed {
	errorEmbed := embed.NewEmbed()
	errorEmbed.SetTitle("Oh noes!")
	errorEmbed.SetDescription("An error occurred getting your meme. Pwease don't be mad >w<")
	errorEmbed.SetImage(constants.GIF_AE_CRY)
	return errorEmbed.MessageEmbed
}
func ErrorWebHook() *discordgo.WebhookEdit {
	var webHook discordgo.WebhookEdit
	webHook.Embeds = &[]*discordgo.MessageEmbed{ErrorEmbed()}
	return &webHook
}

func AddToInteractionCache(interaction *discordgo.InteractionCreate) {

	var NewCache []CacheItem
	for _, i := range InteractionCache {
		if time.Since(i.Added).Minutes() <= float64(CacheTimeoutMins) {
			NewCache = append(NewCache, i)
		}
	}

	NewCache = append(InteractionCache, CacheItem{
		Interaction: interaction,
		Added:       time.Now(),
	})

	InteractionCache = NewCache

}

func GetFromInteractionCache(ID string) (discordgo.InteractionCreate, error) {
	for _, i := range InteractionCache {
		if i.Interaction.ID == ID {
			return *i.Interaction, nil
		}
	}

	return discordgo.InteractionCreate{}, errors.New("not found in cache")
}

func FakeInteractionResponse(interaction *discordgo.InteractionCreate) error {
	embedFake := embed.NewEmbed()
	embedFake.SetDescription("Loading...")

	err := config.Session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embedFake.MessageEmbed},
		},
	})
	if err != nil {
		audit.Error(err)
		return err
	}

	err = config.Session.InteractionResponseDelete(interaction.Interaction)
	if err != nil {
		audit.Error(err)
		return err
	}

	return nil
}

func FormatMessage(memeUrl string, threadUrl string, threadTitle string) string {
	title, err := helpers.StripHTML(threadTitle)
	if err != nil {
		title = threadTitle
		audit.Error(err)
	}
	return "[Meme](" + memeUrl + ") | [Thread](<" + threadUrl + ">) | " + title
}
