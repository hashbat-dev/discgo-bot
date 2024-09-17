package slash

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode"

	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	cache "github.com/dabi-ngin/discgo-bot/Cache"
	config "github.com/dabi-ngin/discgo-bot/Config"
	discord "github.com/dabi-ngin/discgo-bot/Discord"
	imgur "github.com/dabi-ngin/discgo-bot/External/Imgur"
	helpers "github.com/dabi-ngin/discgo-bot/Helpers"
	imgwork "github.com/dabi-ngin/discgo-bot/ImgWork"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
	tempfiles "github.com/dabi-ngin/discgo-bot/TempFiles"
)

var (
	MemeGenCustomBase = "https://api.memegen.link/images/custom/"
)

func MakeMemeInit(i *discordgo.InteractionCreate, correlationId string) {
	// 1. Get the Message object associated with the Interaction request
	messageID := i.ApplicationCommandData().TargetID
	if messageID == "" {
		logger.ErrorText(i.GuildID, "MakeMeme: No MessageID provided")
	}

	message := i.ApplicationCommandData().Resolved.Messages[messageID]

	// 2. Check there's an associated Image
	imgUrl := helpers.GetImageFromMessage(message, "")
	if imgUrl == "" {
		discord.SendEmbedFromInteraction(i, "Error", "No image found in message")
		return
	}

	imgExtension := imgwork.GetExtensionFromURL(imgUrl)
	if imgExtension == "" {
		discord.SendEmbedFromInteraction(i, "Error", fmt.Sprintf("Invalid image extension (%s)", imgExtension))
		return
	}

	// => Store these in the Interactions cache for later
	cache.ActiveInteractions[correlationId].Values.String["imgUrl"] = imgUrl
	cache.ActiveInteractions[correlationId].Values.String["imgExtension"] = imgExtension

	// 3. Create the Interaction Objects
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

	// 4. Send the Select menu response
	err := config.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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

	// 2. Is the URL we have Accessible by MemeGen?
	//	  We need a clean URL without any Query strings, Discord Proxy URLs do not work.
	//	  If we DON'T have a clean URL we'll upload it to Imgur and get use that URL.
	var deleteHash string
	var sendImgUrl string
	var sendImgExt string
	var imgSource int // 0: Inbound URL, 1: TempFile, 2: Imgur
	if strings.Contains(cachedInteraction.Values.String["imgUrl"], "?") {
		discord.UpdateInteractionResponse(i, "Creating Meme", "Getting image...")
		if !config.ServiceSettings.ISDEV {
			// Temp file route
			tempFileReader, err := imgwork.DownloadImageToReader(i.GuildID, cachedInteraction.Values.String["imgUrl"], cachedInteraction.Values.String["imgExtension"] == ".gif")
			if err != nil {
				discord.UpdateInteractionResponse(i, "Error", "Couldn't download image.")
				cache.InteractionComplete(correlationId)
				return
			}
			tempFile := tempfiles.AddFile(tempFileReader, cachedInteraction.Values.String["imgExtension"])
			if tempFile == "" {
				discord.UpdateInteractionResponse(i, "Error", "Couldn't download image.")
				cache.InteractionComplete(correlationId)
				return
			}
			imgSource = 1
			sendImgUrl = tempFile
			sendImgExt = cachedInteraction.Values.String["imgExtension"]
		} else {
			// Imgur route
			imgurUrl, imgurDeleteHash, err := getImgurLink(i.GuildID, i.Member.User.ID, cachedInteraction.Values.String["imgUrl"], cachedInteraction.Values.String["imgExtension"])
			if err != nil {
				if strings.Contains(err.Error(), "413") {
					discord.UpdateInteractionResponse(i, "Error", "File size too large.")
				} else {
					discord.UpdateInteractionResponse(i, "Error", "Error getting image.")
				}
				cache.InteractionComplete(correlationId)
				return
			}
			imgSource = 2
			deleteHash = imgurDeleteHash
			sendImgUrl = imgurUrl
			sendImgExt = imgwork.GetExtensionFromURL(imgurUrl)
		}
	} else {
		sendImgUrl = cachedInteraction.Values.String["imgUrl"]
		sendImgExt = cachedInteraction.Values.String["imgExtension"]
	}

	if sendImgUrl == "" || sendImgExt == "" {
		discord.UpdateInteractionResponse(i, "Error", "Couldn't download image.")
		cache.InteractionComplete(correlationId)
		return
	}

	// 3. Generate the Request URL
	discord.UpdateInteractionResponse(i, "Creating Meme", "Building request...")
	url := MemeGenCustomBase
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
		url += "?layout=top&font=notosans&background=" + sendImgUrl
	} else {
		// In Image Caption
		if topText == "" && bottomText == "" {
			discord.UpdateInteractionResponse(i, "Error", "No Captions provided.")
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
		url += sendImgExt + "?font=impact&background=" + sendImgUrl
	}

	// 4. Get the Meme
	discord.UpdateInteractionResponse(i, "Creating Meme", "Getting Meme...")
	newMemeReader, err := getMemeImage(i.GuildID, url)
	if err != nil {
		discord.UpdateInteractionResponse(i, "Error", "Error Getting Meme.")
		cache.InteractionComplete(correlationId)
		return
	}

	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, newMemeReader)
	if err != nil {
		logger.Error(i.GuildID, err)
		discord.UpdateInteractionResponse(i, "Error", "Error Generating Meme.")
		cache.InteractionComplete(correlationId)
		return
	}

	// 5. Send via Discord
	err = config.Session.InteractionResponseDelete(i.Interaction)
	if err != nil {
		logger.Error(i.GuildID, err)
	}
	_, err = discord.SendMessageWithImageBuffer(i.ChannelID, i.GuildID, sendImgExt, &buffer)
	if err != nil {
		logger.Error(i.GuildID, err)
	}

	// 6. Cleanup
	switch imgSource {
	case 1: // => Temp File
		// This is being commented out for Debugging, think we're deleting it too quickly.
		// tempfiles.DeleteFile(i.GuildID, sendImgUrl)
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
		default:
			if unicode.IsLetter(rune(input[i])) || unicode.IsDigit(rune(input[i])) {
				buffer.WriteByte(input[i])
			}
		}
	}

	return buffer.String()
}
