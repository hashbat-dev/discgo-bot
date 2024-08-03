package helpers

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/bwmarrin/discordgo"
	"github.com/dabi-ngin/discgo-bot/Bot/audit"
	"golang.org/x/net/html"
)

func GetText(command string, forEmbed bool, userFrom string, userTo string) string {
	showText := ""
	switch command {
	case "blush":
		if userTo == "" {
			randText := []string{
				userFrom + " blushes",
				userFrom + " goes full bottom mode",
			}
			showText = GetRandomText(randText)
		} else {
			showText = userFrom + " blushes at " + userTo
		}
	case "bored":
		if userTo == "" {
			randText := []string{
				userFrom + " is bored",
				userFrom + " is as bored as a cucumber sandwich",
				userFrom + " is as bored as the recipient of a Cashew Nuts conversation",
			}
			showText = GetRandomText(randText)
		} else {
			showText = userFrom + " revels in boredom with " + userTo
		}
	case "cry":
		if userTo == "" {
			randText := []string{
				userFrom + " cries " + GetEmote("ae_cry", forEmbed),
				userFrom + " cries like a Femboy outside of Friday " + GetEmote("ae_cry", forEmbed),
			}
			showText = GetRandomText(randText)
		} else {
			showText = userFrom + " cries at " + userTo + " " + GetEmote("ae_cry", forEmbed)
		}
	case "cuddles":
		if userTo == "" {
			showText = userFrom + " cuddles the air " + GetEmote("ae_cry", forEmbed)
		} else {
			showText = userFrom + " cuddles " + userTo
		}
	case "dance":
		randText := []string{
			userFrom + " breaks out the shiggy wiggy",
			userFrom + " busts it down",
			userFrom + " wiggles those elbows",
			userFrom + " gets down with the sickness",
			userFrom + " hits the griddy",
			userFrom + " busts it down sexual style",
			userFrom + " default dances",
		}
		showText = GetRandomText(randText)

		if userTo != "" {
			showText += " at " + userTo
		}
	case "facepalm":
		if userTo == "" {
			showText = userFrom + " facepalms"
		} else {
			showText = userFrom + " facepalms at " + userTo
		}
	case "happy":
		if userTo == "" {
			showText = userFrom + " is happy!"
		} else {
			showText = userTo + " made " + userFrom + " happy!"
		}
	case "highfive":
		if userTo == "" {
			showText = userFrom + " high-fives the air"
		} else {
			showText = userFrom + " high-fives " + userTo
		}
	case "hug":
		if userTo == "" {
			showText = userFrom + " hugs someone"
		} else {
			showText = userFrom + " hugs " + userTo
		}
	case "kick":
		if userTo == "" {
			showText = userFrom + " goes absolutely Jackie Chan"
		} else {
			randText := []string{
				userFrom + " kicks " + userTo,
				userFrom + " absolutely munts " + userTo,
				userFrom + " roundhouses " + userTo,
			}
			showText = GetRandomText(randText)
		}
	case "kiss":
		if userTo == "" {
			showText = userFrom + " kisses the air like a retarded cat"
		} else {
			randText := []string{
				userFrom + " gives " + userTo + " a smooch",
				userFrom + " kisses " + userTo,
				userFrom + " gives lil' love licks to " + userTo,
				userFrom + " samples the saliva of " + userTo,
			}
			showText = GetRandomText(randText)
		}
	case "laugh":
		if userTo == "" {
			randText := []string{
				userFrom + " has a right old giggle",
				userFrom + " laughs",
				userFrom + " is absolutely kekking",
				userFrom + " keks",
			}
			showText = GetRandomText(randText)
		} else {
			showText = userFrom + " laughs at " + userTo
		}
	case "nod":
		if userTo == "" {
			showText = userFrom + " nods"
		} else {
			showText = userFrom + " nods at " + userTo
		}
	case "nom":
		if userTo == "" {
			showText = userFrom + " noms"
		} else {
			showText = userFrom + " noms with " + userTo
		}
	case "nope":
		if userTo == "" {
			showText = userFrom + " gives their biggest nopers"
		} else {
			showText = userFrom + " nopes at " + userTo
		}
	case "pat":
		if userTo == "" {
			showText = userFrom + " sent pats into the void"
		} else {
			showText = userFrom + " pats " + userTo
		}
	case "poke":
		if userTo == "" {
			showText = userFrom + " pokes someone"
		} else {
			showText = userFrom + " pokes " + userTo
		}
	case "pout":
		if userTo == "" {
			showText = userFrom + " gives their biggest pouts"
		} else {
			showText = userFrom + " makes pouty faces at " + userTo
		}
	case "punch":
		if userTo == "" {
			showText = userFrom + " throws a punch"
		} else {
			randText := []string{
				userFrom + " absolutely clackers " + userTo,
				userFrom + " lamps " + userTo,
				userFrom + " knocks out " + userTo,
				userFrom + " knocks seven bells out of " + userTo,
				userFrom + " treats " + userTo + " like a British wife after an England loss",
			}
			showText = GetRandomText(randText)
		}
	case "itsfwiday":
		{
			randText := []string{
				"Get out your finest stockings, IT'S FWIDAY! " + GetEmote("yep", forEmbed),
				"Make sure to give Pog your finest headpats!",
				"Be nice to Pog, the toppiest top in the server " + GetEmote("yep", forEmbed),
				"Bojangles all around!",
				"Bust a groove and hit the griddy! It's mother fucking Fwiday!",
			}
			showText = GetRandomText(randText)
		}
	case "shoot":
		if userTo == "" {
			showText = userFrom + " shoots into the air"
		} else {
			randText := []string{
				userFrom + " shoots " + userTo,
				userFrom + " makes " + userTo + " regret coming to jousting practice",
			}
			showText = GetRandomText(randText)
		}
	case "shrug":
		if userTo == "" {
			showText = userFrom + " shrugs"
		} else {
			showText = userFrom + " shrugs at " + userTo
		}
	case "slap":
		if userTo == "" {
			showText = userFrom + " slaps someone"
		} else {
			showText = userFrom + " slaps " + userTo
		}
	case "sleep":
		if userTo == "" {
			showText = userFrom + " is 'eepy"
		} else {
			showText = userFrom + " sleeps with " + userTo
		}
	case "cancelledfwiday":
		randText := []string{
			"Pog pushed his luck too far.",
			"If Pog confesses the last time he farted we MAY re-instate it.",
			"Have a long hard think about what you've done Pog.",
		}
		showText = GetRandomText(randText)
	case "smile":
		if userTo == "" {
			showText = userFrom + " smiles"
		} else {
			showText = userFrom + " smiles at " + userTo
		}
	case "smug":
		if userTo == "" {
			showText = userFrom + " acts smug"
		} else {
			showText = userFrom + " gives a smug look to " + userTo
		}
	case "stare":
		if userTo == "" {
			showText = userFrom + " stares into the void"
		} else {
			showText = userFrom + " stares at " + userTo
		}
	case "think":
		if userTo == "" {
			showText = userFrom + " has a hecking think"
		} else {
			showText = userFrom + " thinks with " + userTo
		}
	case "thumbsup":
		if userTo == "" {
			showText = userFrom + " gives a hearty thumbs up"
		} else {
			showText = userFrom + " gives " + userTo + " a thumbs up"
		}
	case "wave":
		if userTo == "" {
			showText = userFrom + " waves!"
		} else {
			showText = userFrom + " waves at " + userTo
		}
	case "wink":
		if userTo == "" {
			showText = userFrom + " winks"
		} else {
			showText = userFrom + " winks at " + userTo
		}
	case "yawn":
		if userTo == "" {
			showText = userFrom + " yawns"
		} else {
			showText = userFrom + " yawns at " + userTo
		}
	case "yeet":
		if userTo == "" {
			showText = userFrom + ": YEEEEEEEEEEET"
		} else {
			showText = userFrom + " yeets " + userTo
		}
	case "fart":
		randText := []string{
			userFrom + " executes a sphincter decompression",
			userFrom + " lets one rip",
			userFrom + " cut the cheese",
			userFrom + " farts",
			userFrom + " passed gas",
			userFrom + " provides the room with a warm fragrance",
			userFrom + " filled the air with poo particles",
		}
		showText = GetRandomText(randText)
	case "speen":
		if userTo == "" {
			randText := []string{
				userFrom + " spins!",
				userFrom + " speens",
				userFrom + " speeeeens",
				userFrom + " speeeeeeeeeeens",
				userFrom + " speeeeeeeeeeeeeeeeeeeeeeens",
				userFrom + " TURBO SPEENS",
				userFrom + " TURBO SPEEEEEEEENS",
				userFrom + " TURBO SPEEEEEEEEEEEEEEENS",
				userFrom + " TURBO SPEEEEEEEEEEEEEEEEEEEEEEENS",
			}
			showText = GetRandomText(randText)
		} else {
			randText := []string{
				userFrom + " speens " + userTo,
				userFrom + " speeeeeeeens " + userTo,
				userFrom + " speeeeeeeeeeeeeeeeens " + userTo,
				userFrom + " speeeeeeeeeeeeeeeeeeeeeeeeens " + userTo,
				userFrom + " TURBO SPEENS " + userTo,
				userFrom + " TURBO SPEEEEEEEENS " + userTo,
				userFrom + " TURBO SPEEEEEEEEEEEEEEEENS " + userTo,
			}
			showText = GetRandomText(randText)
		}

	default:
		fmt.Println("GetText returned no value for: [" + command + "]")
	}

	return showText
}

func GetLettersOnlyCharactersFromString(in string) string {
	var alphabeticChars []rune

	for _, char := range in {
		if unicode.IsLetter(char) {
			alphabeticChars = append(alphabeticChars, char)
		}
	}

	return string(alphabeticChars)
}

func GetNumbersOnlyCharactersFromString(in string) int {
	var numberChars []rune

	for _, char := range in {
		if unicode.IsNumber(char) {
			numberChars = append(numberChars, char)
		}
	}

	retInt, err := strconv.Atoi(string(numberChars))
	if err != nil {
		retInt = 0
		audit.Error(err)
	}

	return retInt
}

var characters string = "¬`!\"£$%^&*()-_+=;:'@#~.>,<?\\|[]{}"
var adjStart []string = []string{}
var adjEnds []string = []string{}

func DoesTextMentionWord(message string, phrase string) bool {
	if !strings.Contains(message, phrase) {
		return false
	}

	if len(adjStart) == 0 {
		for _, c := range characters {
			adjStart = append(adjStart, string(c))
			adjEnds = append(adjEnds, string(c))
		}

		adjStart = append(adjStart, string(""))
		adjEnds = append(adjEnds, string(""))
		adjStart = append(adjStart, string(" "))
		adjEnds = append(adjEnds, string(" "))

	}

	for _, a_s := range adjStart {
		for _, a_e := range adjEnds {

			if a_s == "" && a_e == "" {
				continue
			}

			if strings.Contains(message, a_s+phrase+a_e) {
				return true
			}
		}
	}

	return false
}

func DeleteSourceMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	session.ChannelMessageDelete(message.ChannelID, message.ID)
}

func GetEmote(emojiName string, forEmbed bool) string {
	base := ""
	if emojiName == "ae_cry" {
		base = "ae_cry:1209268603812057118"
	} else if emojiName == "yep" {
		base = "bottom_yup:1223028483475636286"
	} else if emojiName == "nope" {
		base = "bottom_nope:1223028646986645647"
	} else if emojiName == "read" {
		base = "waitwaitwait:1209271418747883580"
	} else if emojiName == "medal" {
		base = "medal:1233399678129668257"
	}

	if base == "" {
		return base
	}

	retString := "<"
	if forEmbed {
		retString += "a"
	}
	retString += ":" + base + ">"
	return retString
}

func GetSuperscriptNumber(num int, addBrackets bool) string {

	retVal := ""
	if addBrackets {
		retVal += "⁽"
	}
	str := strconv.Itoa(num)
	for _, char := range str {
		if char == '0' {
			retVal += "⁰"
		} else if char == '1' {
			retVal += "¹"
		} else if char == '2' {
			retVal += "²"
		} else if char == '3' {
			retVal += "³"
		} else if char == '4' {
			retVal += "⁴"
		} else if char == '5' {
			retVal += "⁵"
		} else if char == '6' {
			retVal += "⁶"
		} else if char == '7' {
			retVal += "⁷"
		} else if char == '8' {
			retVal += "⁸"
		} else if char == '9' {
			retVal += "⁹"
		}
	}
	if addBrackets {
		retVal += "⁾"
	}
	return retVal
}

func ExtractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var sb strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		sb.WriteString(ExtractText(c))
	}
	return sb.String()
}

func StripHTML(input string) (string, error) {
	// Add a temporary outer div tag
	wrappedInput := "<div>" + input + "</div>"
	doc, err := html.Parse(strings.NewReader(wrappedInput))
	if err != nil {
		return "", err
	}
	// Find the first "div" node, which will be the wrapper we added
	var body *html.Node
	var findBody func(*html.Node)
	findBody = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			body = n
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findBody(c)
		}
	}
	findBody(doc)
	if body == nil {
		return "", fmt.Errorf("no div element found")
	}
	text := ExtractText(body)
	unescapedText := html.UnescapeString(text)
	return unescapedText, nil
}
