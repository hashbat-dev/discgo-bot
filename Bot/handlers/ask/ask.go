package ask

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"regexp"

	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	"github.com/dabi-ngin/discgo-bot/Bot/audit"
	"github.com/dabi-ngin/discgo-bot/Bot/config"
	"github.com/dabi-ngin/discgo-bot/Bot/helpers"
	"github.com/dabi-ngin/discgo-bot/Bot/logging"
)

func Question(interaction *discordgo.InteractionCreate) {

	// Validate Question
	optionMap := helpers.GetOptionMap(interaction)
	inQuestion := helpers.GetOptionStringValue(optionMap, "question")

	isQuestionInt, basicAddOnTrue, basicAddOnFalse := IsQuestion(inQuestion)

	errorText := ""
	switch isQuestionInt {
	case -1:
		errorText = "grug no understand"
	case 0:
		errorText = "Either I'm stupid or that's not a question!"
	}

	if errorText != "" {

		e := embed.NewEmbed()
		e.SetTitle(inQuestion)
		e.SetDescription(errorText)

		err := config.Session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{e.MessageEmbed},
			},
		})

		if err != nil {
			audit.Error(err)
			logging.SendErrorInteraction(interaction)
			return
		}

		return
	}

	// Get random Quote
	url := "https://api.adviceslip.com/advice"

	// Send GET request
	response, err := http.Get(url)
	if err != nil {
		audit.Error(err)
	}
	defer response.Body.Close()

	// Decode JSON
	var data AdviceResp
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		audit.Error(err)
		return
	}

	// DECIDE
	quote := ""
	if data.Slip.Advice != "" {
		quote = data.Slip.Advice
	}

	replyText := ""
	if isQuestionInt == 1 {
		decidedNo := rand.Intn(2) == 0
		if decidedNo {
			replyText += "No"
			if basicAddOnFalse != "" {
				replyText += ", " + basicAddOnFalse + ". "
			} else {
				replyText += ". "
			}
		} else {
			replyText += "Yes"
			if basicAddOnTrue != "" {
				replyText += ", " + basicAddOnTrue + ". "
			} else {
				replyText += ". "
			}
		}
	}

	if quote != "" {
		replyText += quote
	}

	if replyText == "" {
		retText := "Something went wrong!"

		e := embed.NewEmbed()
		e.SetTitle(inQuestion)
		e.SetDescription(retText)

		err := config.Session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{e.MessageEmbed},
			},
		})

		if err != nil {
			audit.Error(err)
			logging.SendErrorInteraction(interaction)
			return
		}
	}

	e := embed.NewEmbed()
	e.SetTitle(inQuestion)
	e.SetDescription(replyText)

	err = config.Session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{e.MessageEmbed},
		},
	})

	if err != nil {
		audit.Error(err)
		logging.SendErrorInteraction(interaction)
		return
	}
}

func IsQuestion(str string) (int, string, string) {

	// -1: Too Complex
	//  0: Not a question
	//  1: Yes
	//	2: Yes, Complex answerable

	tooComplexPattern := `(?i)\b(how|who)\b`
	regex := regexp.MustCompile(tooComplexPattern)
	if regex.MatchString(str) {
		return -1, "", ""
	}

	regex = regexp.MustCompile(`(?i)am\s+i`)
	if regex.MatchString(str) {
		return 1, "you are", "you aren't"
	}

	regex = regexp.MustCompile(`(?i)are you`)
	if regex.MatchString(str) {
		return 1, "I am", "i'm not"
	}

	regex = regexp.MustCompile(`(?i)\bis\b`)
	if regex.MatchString(str) {
		return 1, "", ""
	}

	regex = regexp.MustCompile(`(?i)\b(should|shall)\b`)
	if regex.MatchString(str) {
		return 1, "they should", "they shouldn't"
	}

	regex = regexp.MustCompile(`(?i)\b(do|does)\b`)
	if regex.MatchString(str) {
		return 1, "they do", "they don't"
	}

	regex = regexp.MustCompile(`(?i)\b(did)\b`)
	if regex.MatchString(str) {
		return 1, "they did", "they didn't"
	}

	regex = regexp.MustCompile(`(?i)\bwill\b`)
	if regex.MatchString(str) {
		return 1, "they will", "they won't"
	}

	regex = regexp.MustCompile(`(?i)\bwould\b`)
	if regex.MatchString(str) {
		return 1, "they would", "they wouldn't"
	}

	regex = regexp.MustCompile(`(?i)\bshould\b`)
	if regex.MatchString(str) {
		return 1, "they should", "they shouldn't"
	}

	regex = regexp.MustCompile(`(?i)\bcould\b`)
	if regex.MatchString(str) {
		return 1, "they could", "they couldn't"
	}

	regex = regexp.MustCompile(`(?i)\bare\b`)
	if regex.MatchString(str) {
		return 1, "they are", "they aren't"
	}

	regex = regexp.MustCompile(`(?i)\bcan\b`)
	if regex.MatchString(str) {
		return 1, "they can", "they can't"
	}

	isQuestionPattern := `(?i)\b(what|when|where|why|is|are|can|do|does|did|will|would|should|could|may|might|shall)\b`
	regex = regexp.MustCompile(isQuestionPattern)
	if !regex.MatchString(str) {
		return 0, "", ""
	}

	complexPattern := `(?i)\b(what|when|where|why|how)\b`
	regex = regexp.MustCompile(complexPattern)
	if regex.MatchString(str) {
		return 2, "", ""
	}

	return 1, "", ""
}

type AdviceResp struct {
	Slip AdviceRespInt `json:"slip"`
}

type AdviceRespInt struct {
	Id     int    `json:"id"`
	Advice string `json:"advice"`
}
