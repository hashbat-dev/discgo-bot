package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/ZestHusky/femboy-control/Bot/audit"
	"github.com/ZestHusky/femboy-control/Bot/config"
	"github.com/ZestHusky/femboy-control/Bot/constants"
	"github.com/ZestHusky/femboy-control/Bot/logging"
	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
)

func GetRandomText(inputSlice []string) string {
	return inputSlice[RandInt(0, len(inputSlice))]
}

func GetRandomInt(inputSlice []int) int {
	return inputSlice[RandInt(0, len(inputSlice))]
}

func RandInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func GetReactionSourceUsers(message *discordgo.MessageCreate) (string, string) {

	userFrom := "<@" + message.Author.ID + ">"
	userTo := ""

	if message.ReferencedMessage != nil {
		// Replied to message?
		userTo = "<@" + message.ReferencedMessage.Author.ID + ">"
		fmt.Println("[ANIMEGIF] message.ReferencedMessage != nil, userTo: " + userTo)
	} else if message.Mentions != nil {
		// @'d user(s)?
		fmt.Println("[ANIMEGIF] message.Mentions != nil")
		var userList []string
		for _, user := range message.Mentions {
			userList = append(userList, user.ID)
		}

		userString := ""
		userCount := len(userList)
		for i, user := range userList {
			if userCount > 0 {
				if i == userCount-1 && userCount > 1 {
					userString += " and "
				} else if i > 0 {
					userString += ", "
				}
			}
			userString += "<@" + user + ">"
			fmt.Println("[ANIMEGIF] Context is: @'d users, string: " + userString)
		}
		userTo = userString
	}

	return userFrom, userTo
}

func GetReactionSourceUsers_FromEdit(message *discordgo.MessageUpdate) (string, string) {

	userFrom := "<@" + message.Author.ID + ">"
	userTo := ""

	if message.ReferencedMessage != nil {
		// Replied to message?
		userTo = "<@" + message.ReferencedMessage.Author.ID + ">"
		fmt.Println("[ANIMEGIF] message.ReferencedMessage != nil, userTo: " + userTo)
	} else if message.Mentions != nil {
		// @'d user(s)?
		fmt.Println("[ANIMEGIF] message.Mentions != nil")
		var userList []string
		for _, user := range message.Mentions {
			userList = append(userList, user.ID)
		}

		userString := ""
		userCount := len(userList)
		for i, user := range userList {
			if userCount > 0 {
				if i == userCount-1 && userCount > 1 {
					userString += " and "
				} else if i > 0 {
					userString += ", "
				}
			}
			userString += "<@" + user + ">"
			fmt.Println("[ANIMEGIF] Context is: @'d users, string: " + userString)
		}
		userTo = userString
	}

	return userFrom, userTo
}

func GetJsonFromURL(url string) (map[string]interface{}, error) {

	// Make HTTP GET request
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Read response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Parse JSON
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, nil

}

func ConvertTextToFilename(text string, maxLength int) string {
	// Define a regular expression pattern to match non-letter characters
	reg := regexp.MustCompile("[^a-zA-Z0-9 ]")

	// Replace non-letter characters with an empty string
	removed := reg.ReplaceAllString(text, "")

	// Replace spaces with underscores
	replaced := strings.ReplaceAll(removed, " ", "_")

	output := replaced
	if len(replaced) > maxLength {
		runes := []rune(replaced)
		output = string(runes[:maxLength])
	}

	return output
}

func ReactList() *discordgo.MessageEmbed {
	e := embed.NewEmbed()
	e.SetTitle("Anime Reactions")
	e.SetFooter("Reply to a message or @mention 1+ users to direct the reaction")

	s := "**!baka**, **!bite**, **!blush**, **!bored**,\n**!cry**, **!cuddle**,\n**!dance**,\n**!facepalm**, **!feed**,\n**!handhold**, **!handshake**, **!happy**, **!highfive**, **!hug**,\n"
	s += "**!kick**, **!kiss**,\n**!laugh**, **!lurk**,\n**!nod**, **!nom**, **!nope**,\n**!pat**, **!peck**, **!poke**, **!pout**, **!punch**,\n**!shoot**, "
	s += "**!shrug**, **!slap**, **!sleep**, **!smile**, **!smug**, **!stare**,\n**!think**, **!thumbsup**, **!tickle**,\n**!wave**, **!wink**,\n**!yawn**, **!yeet**"
	e.SetDescription(s)

	return e.MessageEmbed
}

func IsAdmin(member *discordgo.Member) bool {
	for _, role := range member.Roles {
		if role == constants.ROLE_BOT_DEVELOPER {
			return true
		}
	}
	return false
}

func GetNicknameFromID(guildId string, userId string) string {
	member, err := config.Session.GuildMember(guildId, userId)
	if err != nil {
		audit.Error(err)
		return ""
	}

	nickname := member.Nick

	// If the member doesn't have a nickname, use their username
	if nickname == "" {
		nickname = member.User.Username
	}

	if nickname == "" {
		nickname = "Unknown"
	}
	return nickname
}

func GetAvatarURLFromID(guildId string, userId string) string {
	member, err := config.Session.GuildMember(guildId, userId)
	if err != nil {
		audit.Error(err)
		return ""
	}

	return member.AvatarURL("")
}

func CommandInDevelopment(interaction *discordgo.InteractionCreate, userText string) bool {
	if config.IsDev {
		return false
	} else {
		logging.SendMessageInteraction(interaction, "Feature In Development", userText, "", "", false)
		return true
	}
}

func GenericErrorEmbed(title string, desc string) *discordgo.MessageEmbed {
	errorEmbed := embed.NewEmbed()
	errorEmbed.SetTitle(title)
	errorEmbed.SetDescription(desc)
	errorEmbed.SetImage(constants.GIF_AE_CRY)
	return errorEmbed.MessageEmbed
}

func GenericErrorWebHook(title string, desc string) *discordgo.WebhookEdit {
	var webHook discordgo.WebhookEdit
	webHook.Embeds = &[]*discordgo.MessageEmbed{GenericErrorEmbed(title, desc)}
	return &webHook
}

func GenericEmbed(title string, desc string) *discordgo.MessageEmbed {
	errorEmbed := embed.NewEmbed()
	errorEmbed.SetTitle(title)
	errorEmbed.SetDescription(desc)
	return errorEmbed.MessageEmbed
}

func GenericWebHook(title string, desc string) *discordgo.WebhookEdit {
	var webHook discordgo.WebhookEdit
	webHook.Embeds = &[]*discordgo.MessageEmbed{GenericEmbed(title, desc)}
	return &webHook
}

func GetRandomNumber(min int, max int) int {
	return rand.Intn(max) + min
}

func IsStringInArray(inArray []string, checkFor string) bool {
	for _, s := range inArray {
		if s == checkFor {
			return true
		}
	}
	return false
}

func ReverseIntArray(arr []int) {
	n := len(arr)
	for i := 0; i < n/2; i++ {
		arr[i], arr[n-1-i] = arr[n-1-i], arr[i]
	}
}

func GetIndividualInts(n int) []int {
	str := strconv.Itoa(n)
	digits := make([]int, len(str))
	for i, char := range str {
		digit, err := strconv.Atoi(string(char))
		if err != nil {
			return nil
		}
		digits[i] = digit
	}
	return digits
}

func CurrentFunctionName() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return ""
	}
	return fn.Name()
}

func IsIntBetweenXandY(num int, lower int, upper int) bool {
	return num >= lower && num <= upper
}
