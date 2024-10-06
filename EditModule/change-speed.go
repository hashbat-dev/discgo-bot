package editmodule

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"io"
	"math"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	cache "github.com/hashbat-dev/discgo-bot/Cache"
	config "github.com/hashbat-dev/discgo-bot/Config"
	discord "github.com/hashbat-dev/discgo-bot/Discord"
	helpers "github.com/hashbat-dev/discgo-bot/Helpers"
	imgwork "github.com/hashbat-dev/discgo-bot/ImgWork"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

type ChangeSpeed struct{}

func (s ChangeSpeed) SelectName() string {
	return "Change Speed"
}

func (s ChangeSpeed) Emoji() *discordgo.ComponentEmoji {
	return &discordgo.ComponentEmoji{Name: "🔁"}
}

func (s ChangeSpeed) PermissionRequirement() int {
	return config.CommandLevelUser
}

func (s ChangeSpeed) Complexity() int {
	return config.TRIVIAL_TASK
}

func (s ChangeSpeed) Execute(i *discordgo.InteractionCreate, correlationId string) {
	// 1. Get the Message object associated with the Interaction request
	messageID, message := discord.GetAssociatedMessageFromInteraction(i)

	// 2. Check there's an associated Image
	imgUrl := helpers.GetImageFromMessage(message, "")
	if imgUrl == "" {
		discord.Interactions_SendError(i, "No image found in message")
		return
	}

	imgExtension := imgwork.GetExtensionFromURL(imgUrl)
	if imgExtension != ".gif" {
		discord.Interactions_SendError(i, fmt.Sprintf("Can't speed up %s's!", imgExtension))
		return
	}

	// => Store these in the Interactions cache for later
	cache.ActiveInteractions[correlationId].Values.String["imgMessageId"] = messageID
	cache.ActiveInteractions[correlationId].Values.String["imgUrl"] = imgUrl
	cache.ActiveInteractions[correlationId].Values.String["imgExtension"] = imgExtension

	// 3. Create the Interaction Objects
	selectMenu := discord.CreateSelectMenu(discordgo.SelectMenu{
		CustomID: "change-speed_speed",
		Options: []discordgo.SelectMenuOption{
			{
				Label: "FASTER",
				Value: "2",
				Emoji: &discordgo.ComponentEmoji{Name: "⏩"},
			},
			{
				Label: "Fast",
				Value: "1",
				Emoji: &discordgo.ComponentEmoji{Name: "➡️"},
			},
			{
				Label: "Slow",
				Value: "-1",
				Emoji: &discordgo.ComponentEmoji{Name: "⬅️"},
			},
			{
				Label: "SLOWER",
				Value: "-2",
				Emoji: &discordgo.ComponentEmoji{Name: "⏪"},
			},
		},
		Placeholder: "Choose Speed...",
	}, correlationId, config.CPU_BOUND_TASK, ChangeSpeedSubmit)

	actionRow := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			selectMenu,
		},
	}

	// 4. Send the Select menu response
	err := config.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Components: []discordgo.MessageComponent{actionRow},
			Flags:      discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		logger.Error(i.GuildID, err)
	}

}

func ChangeSpeedSubmit(i *discordgo.InteractionCreate, correlationId string) {
	// 1. Get the Speed information
	cachedInteraction := cache.ActiveInteractions[correlationId]
	speedStringValue := cachedInteraction.Values.String["change-speed_speed"]
	if speedStringValue == "" {
		discord.Interactions_SendError(i, "")
		cache.InteractionComplete(correlationId)
		return
	}

	speedInt, err := strconv.ParseInt(speedStringValue, 10, 32)
	if err != nil {
		logger.Error(i.GuildID, err)
		discord.Interactions_SendError(i, "")
		cache.InteractionComplete(correlationId)
		return
	}

	msgTitle := "Speed Up"
	if speedInt < 0 {
		msgTitle = "Slow Down"
	}
	discord.Interactions_SendMessage(i, msgTitle, "Downloading GIF...")

	// 2. Get the image as an io.Reader object
	imgUrl := cache.ActiveInteractions[correlationId].Values.String["imgUrl"]
	message, err := discord.Message_GetObject(i.GuildID, i.ChannelID, cache.ActiveInteractions[correlationId].Values.String["imgMessageId"])
	if err != nil {
		discord.Interactions_SendError(i, "")
		cache.InteractionComplete(correlationId)
		return
	}

	imageReader, err := imgwork.DownloadImageToReader(i.GuildID, imgUrl, true)
	if err != nil {
		discord.Interactions_SendError(i, "")
		cache.InteractionComplete(correlationId)
		return
	}

	// 4. Change the GIF Speed
	var newImageBuffer bytes.Buffer
	discord.Interactions_SendMessage(i, msgTitle, "Changing Speed...")
	err = changeSpeedGif(i.GuildID, imageReader, &newImageBuffer, int(speedInt))
	if err != nil {
		discord.Interactions_SendError(i, "")
		cache.InteractionComplete(correlationId)
		return
	}

	// 5. Return the reversed Image
	outputImageName := uuid.New().String() + ".gif"
	discord.Interactions_EditText(i, msgTitle+" Completed", "")
	if discord.Message_ReplyWithImage(message, true, outputImageName, &newImageBuffer) != nil {
		logger.Debug(i.GuildID, "Error sending ReplyWithImage")
	}
}

func changeSpeedGif(guildId string, imageReader io.Reader, buffer *bytes.Buffer, newSpeed int) error {
	// 1. Decode the file into a GIF object
	timeStarted := time.Now()
	gifImage, err := gif.DecodeAll(imageReader)
	if err != nil {
		logger.Error(guildId, err)
		return err
	}

	// 2. Lower the Delay of the frames
	outGif := &gif.GIF{}
	var snippedDelay []int
	alreadySlowestCount := 0
	for i := range gifImage.Image {
		frameDelay := gifImage.Delay[i]
		if frameDelay <= 2 {
			alreadySlowestCount++
		}
		var newValue int
		switch newSpeed {
		case 2:
			newValue = int(math.Round(float64(frameDelay) / 4.0))
		case 1:
			newValue = int(math.Round(float64(frameDelay) / 2.0))
		case -1:
			newValue = int(math.Round(float64(frameDelay) * 2.0))
		case -2:
			newValue = int(math.Round(float64(frameDelay) * 4.0))
		default:
			continue
		}

		if newValue < 2 {
			newValue = 2
		}
		snippedDelay = append(snippedDelay, newValue)
	}

	// 3. If it's already as slow as possible due to the delay, cut frames out
	var newPaletted []*image.Paletted
	var newDelay []int
	var newDisposal []byte
	if alreadySlowestCount == len(gifImage.Delay) {
		keptLast := false
		for i, frame := range gifImage.Image {
			if !keptLast {
				newPaletted = append(newPaletted, frame)
				newDelay = append(newDelay, snippedDelay[i])
				newDisposal = append(newDisposal, gifImage.Disposal[i])
				keptLast = true
			} else {
				keptLast = false
			}
		}
	} else {
		newPaletted = gifImage.Image
		newDelay = snippedDelay
		newDisposal = gifImage.Disposal
	}

	outGif.BackgroundIndex = gifImage.BackgroundIndex
	outGif.Config = gifImage.Config
	outGif.Disposal = newDisposal
	outGif.LoopCount = gifImage.LoopCount
	outGif.Image = newPaletted
	outGif.Delay = newDelay

	// 3. Write the new GIF to buffer arg
	err = gif.EncodeAll(buffer, outGif)
	if err != nil {
		logger.Error(guildId, err)
		return err
	}

	direction := "down"
	if newSpeed > 0 {
		direction = "up"
	}

	logger.Info(guildId, "ChangeSpeed [%v] completed after [%v]", direction, time.Since(timeStarted))
	return nil
}
