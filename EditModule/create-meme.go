package editmodule

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode"

	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	cache "github.com/hashbat-dev/discgo-bot/Cache"
	config "github.com/hashbat-dev/discgo-bot/Config"
	database "github.com/hashbat-dev/discgo-bot/Database"
	discord "github.com/hashbat-dev/discgo-bot/Discord"
	imgur "github.com/hashbat-dev/discgo-bot/External/Imgur"
	imgwork "github.com/hashbat-dev/discgo-bot/ImgWork"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
	tempfiles "github.com/hashbat-dev/discgo-bot/TempFiles"
)

type CreateMeme struct{}

var (
	MemeGenCustomBase = "https://api.memegen.link/images/custom/"
)

func (s CreateMeme) SelectName() string {
	return "Add Caption"
}

func (s CreateMeme) Emoji() *discordgo.ComponentEmoji {
	return &discordgo.ComponentEmoji{Name: "📝"}
}

func (s CreateMeme) PermissionRequirement() int {
	return config.CommandLevelUser
}

func (s CreateMeme) Complexity() int {
	return config.TRIVIAL_TASK
}

func (s CreateMeme) Execute(i *discordgo.InteractionCreate, correlationId string) {
	response := discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource, // Acknowledges the interaction
	}

	// Respond to the interaction
	err := config.Session.InteractionRespond(cache.ActiveInteractions[correlationId].StartInteraction.Interaction, &response)
	if err != nil {
		fmt.Println("Error acknowledging interaction:", err)
	}

	// 1. Create the Interaction Objects
	captionText := discord.CreateTextInput(discordgo.TextInput{
		CustomID:    "make-meme_caption-text",
		Placeholder: "(A) Above Image: Caption...",
		Style:       discordgo.TextInputShort,
	}, correlationId)
	topText := discord.CreateTextInput(discordgo.TextInput{
		CustomID:    "make-meme_top-text",
		Placeholder: "(B) In Image: Top Text...",
		Style:       discordgo.TextInputShort,
	}, correlationId)
	bottomText := discord.CreateTextInput(discordgo.TextInput{
		CustomID:    "make-meme_bottom-text",
		Placeholder: "(B) In Image: Bottom Text...",
		Style:       discordgo.TextInputShort,
	}, correlationId)

	actionRow2 := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			captionText,
		},
	}
	actionRow4 := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			topText,
		},
	}
	actionRow5 := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			bottomText,
		},
	}

	// 2. Send the Select menu response
	err = config.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: discord.CreateInteractionResponseModal(discordgo.InteractionResponseData{
			CustomID: "make-meme_submit-modal",
			Title:    "Enter either (A) or at least one of (B)",
			Components: []discordgo.MessageComponent{
				actionRow2,
				actionRow4,
				actionRow5,
			},
			Flags: discordgo.MessageFlagsEphemeral,
		}, correlationId, config.IO_BOUND_TASK, MakeMemeStart),
	})
	if err != nil {
		logger.Error(i.GuildID, err)
	}

}

func MakeMemeStart(i *discordgo.InteractionCreate, correlationId string) {
	cachedInteraction := cache.ActiveInteractions[correlationId]

	// 1. Respond to the Modal
	initEmbed := embed.NewEmbed()
	initEmbed.SetDescription("Processing Request...")
	err := config.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{initEmbed.MessageEmbed},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		logger.Error(i.GuildID, err)
	}

	originalUrl := cachedInteraction.Values.String["imgUrl"]
	originalExt := cachedInteraction.Values.String["imgExtension"]

	// 2. Get a Publically accessible URL
	//	  MemeGen's API needs to be able to see the image. We need to factor these points in:
	//		- Localhost is off limits for Development
	//		- The MemeGen API doesn't support .webp
	//		- Discord's CDN actively blocks Proxy URLs from external sites
	//		- The MemeGen API doesn't support links with Query strings
	discord.Interactions_EditText(i, "Creating Meme", "Getting image...")
	var imgSource int
	// 0 - Inbound URL
	// 1 - Temp File
	// 2 - Imgur
	if strings.Contains(cachedInteraction.Values.String["imgUrl"], "?") {
		// Query string exists, can't use Direct URL
		if !config.ServiceSettings.ISDEV {
			imgSource = 1
		} else {
			imgSource = 2
		}
	} else {
		if cachedInteraction.Values.String["imgUrl"] == ".webp" {
			imgSource = 1 // Needed to convert
		} else {
			imgSource = 0
		}
	}

	var (
		sendImgUrl string // URL to send to the API
		sendImgExt string // Extension of the URL going to the API
		deleteHash string // [If using Imgur] Used to delete the image at Imgur after processing
	)
	switch imgSource {
	case 0: // Send inbound URL - No conversion/download required as this can be accessed by the API
		sendImgUrl = cachedInteraction.Values.String["imgUrl"]
		sendImgExt = cachedInteraction.Values.String["imgExtension"]
	case 1: // Download the file to the temp/ directory and provide that URL to the API
		tempFileReader, err := imgwork.DownloadImageToReader(i.GuildID, cachedInteraction.Values.String["imgUrl"], cachedInteraction.Values.String["imgExtension"] == ".gif")
		if err != nil {
			discord.Interactions_EditText(i, "Error", "Couldn't download image.")
			cache.InteractionComplete(correlationId)
			return
		}

		// Convert .webp to .png
		if cachedInteraction.Values.String["imgExtension"] == ".webp" {
			pngReader, err := imgwork.ConvertWebpToPNG(tempFileReader)
			if err != nil {
				discord.Interactions_EditText(i, "Error", "Couldn't convert .webp image.")
				cache.InteractionComplete(correlationId)
				return
			} else {
				tempFileReader = pngReader
			}
		}

		tempFile := tempfiles.AddFile(tempFileReader, cachedInteraction.Values.String["imgExtension"])
		if tempFile == "" {
			discord.Interactions_EditText(i, "Error", "Couldn't download image.")
			cache.InteractionComplete(correlationId)
			return
		}
		imgSource = 1
		sendImgUrl = tempFile
		sendImgExt = cachedInteraction.Values.String["imgExtension"]
	case 2: // Upload to Imgur - Use the Imgur API to upload and provide the Imgur Image URL
		imgurUrl, imgurDeleteHash, err := getImgurLink(i.GuildID, i.Member.User.ID, cachedInteraction.Values.String["imgUrl"], cachedInteraction.Values.String["imgExtension"])
		if err != nil {
			if strings.Contains(err.Error(), "413") {
				discord.Interactions_EditText(i, "Error", "File size too large.")
			} else {
				discord.Interactions_EditText(i, "Error", "Error getting image.")
			}
			cache.InteractionComplete(correlationId)
			return
		}
		imgSource = 2
		deleteHash = imgurDeleteHash
		sendImgUrl = imgurUrl
		sendImgExt = imgwork.GetExtensionFromURL(imgurUrl)
	default:
		logger.ErrorText(i.GuildID, "Unknown imgSource value: %v", imgSource)
		discord.Interactions_EditText(i, "Error", "An Error occured.")
		cache.InteractionComplete(correlationId)
		return
	}

	if sendImgUrl == "" || sendImgExt == "" {
		discord.Interactions_EditText(i, "Error", "Couldn't download image.")
		cache.InteractionComplete(correlationId)
		return
	}

	// 3. Generate the Request URL
	discord.Interactions_EditText(i, "Creating Meme", "Building request...")
	url := MemeGenCustomBase
	logUrl := ""
	captionText := ""
	topText := ""
	bottomText := ""
	if text, exists := cachedInteraction.Values.String["make-meme_caption-text"]; exists {
		captionText = text
	}
	if text, exists := cachedInteraction.Values.String["make-meme_top-text"]; exists {
		topText = text
	}
	if text, exists := cachedInteraction.Values.String["make-meme_bottom-text"]; exists {
		bottomText = text
	}

	if captionText != "" {
		// Top Caption
		url += encodeTextForUrl(captionText) + sendImgExt
		url += "?layout=top&font=notosans&background="
		logUrl = url + "{0}"
		url += sendImgUrl
	} else {
		// In Image Caption
		if topText == "" && bottomText == "" {
			discord.Interactions_EditText(i, "Error", "No Captions provided.")
			cache.InteractionComplete(correlationId)
			return
		}

		if topText == "" {
			url += "_"
		} else {
			url += encodeTextForUrl(topText)
		}
		url += "/"
		if bottomText == "" {
			url += "_"
		} else {
			url += encodeTextForUrl(bottomText)
		}
		url += sendImgExt + "?font=impact&background="
		logUrl = url + "{0}"
		url += sendImgUrl
	}

	// 4. Get the Meme
	discord.Interactions_EditText(i, "Creating Meme", "Getting Meme...")
	newMemeReader, err := getMemeImage(i.GuildID, url)
	if err != nil {
		discord.Interactions_EditText(i, "Error", "Error Getting Meme.")
		cache.InteractionComplete(correlationId)
		return
	}

	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, newMemeReader)
	if err != nil {
		logger.Error(i.GuildID, err)
		discord.Interactions_EditText(i, "Error", "Error Generating Meme.")
		cache.InteractionComplete(correlationId)
		return
	}

	// 5. Send via Discord
	err = config.Session.InteractionResponseDelete(i.Interaction)
	if err != nil {
		logger.Error(i.GuildID, err)
	}
	outMessage, err := discord.Message_SendImage(i.ChannelID, i.GuildID, sendImgExt, &buffer)
	if err != nil {
		logger.Error(i.GuildID, err)
	}

	// 6. Record in the Database
	database.InsertMemeGenLog(i.GuildID, i.Member.User.ID, correlationId, outMessage, originalUrl, originalExt, logUrl)

	// 7. Cleanup
	switch imgSource {
	case 1: // => Temp File
		tempfiles.DeleteFile(i.GuildID, sendImgUrl)
	case 2: // => Imgur
		if deleteHash != "" {
			err = imgur.DeleteImgurEntry(i.GuildID, deleteHash)
			if err != nil {
				logger.Debug(i.GuildID, "Unable to delete temp Imgur image")
			}
		}
	default:
	}
	cache.InteractionComplete(correlationId)
}

func getMemeImage(guildId string, url string) (io.Reader, error) {
	logger.Debug(guildId, "Requesting Meme generation: [%s]", url)
	resp, err := http.Get(url)
	if err != nil {
		logger.Error(guildId, err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("get meme image returned http status: %s", resp.Status)
		logger.Error(guildId, err)
		return nil, err
	}

	return resp.Body, nil
}

// Returns Imgur Link, Delete Hash, Error
func getImgurLink(guildId string, userId string, imgUrl string, imgExtension string) (string, string, error) {
	imageReader, err := imgwork.DownloadImageToReader(guildId, imgUrl, imgExtension == ".gif")
	if err != nil {
		return "", "", err
	}
	return imgur.UploadAndGetUrl(guildId, userId, imageReader)
}

func encodeTextForUrl(input string) string {
	var buffer bytes.Buffer

	for i := 0; i < len(input); i++ {
		switch input[i] {
		case ' ':
			buffer.WriteByte('_')
		case '-':
			buffer.WriteString("--")
		case '_':
			buffer.WriteString("__")
		case '\n':
			buffer.WriteString("~n")
		case '?':
			buffer.WriteString("~q")
		case '&':
			buffer.WriteString("~a")
		case '%':
			buffer.WriteString("~p")
		case '#':
			buffer.WriteString("~h")
		case '/':
			buffer.WriteString("~s")
		case '\\':
			buffer.WriteString("~b")
		case '<':
			buffer.WriteString("~l")
		case '>':
			buffer.WriteString("~g")
		case '"':
			buffer.WriteString("''")
		case '\'':
			buffer.WriteString("'")
		default:
			if unicode.IsLetter(rune(input[i])) || unicode.IsDigit(rune(input[i])) {
				buffer.WriteByte(input[i])
			}
		}
	}

	return buffer.String()
}
