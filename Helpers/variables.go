package helpers

import (
	"strings"
	"time"
	"unicode"
)

func RemoveStartingXCharacters(inMsg string, removeLength int) string {
	if len(inMsg) < removeLength {
		return inMsg
	} else {
		return inMsg[3:]
	}
}

func GetNullDateTime() time.Time {
	return time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
}

func ConcatStringWithAnd(words []string) string {
	switch len(words) {
	case 0:
		return ""
	case 1:
		return words[0]
	case 2:
		return words[0] + " and " + words[1]
	default:
		return strings.Join(words[:len(words)-1], ", ") + " and " + words[len(words)-1]
	}
}

func CapitaliseWords(s string) string {
	words := strings.Fields(s)
	for i, word := range words {
		if len(word) > 0 {
			r := []rune(word)
			r[0] = unicode.ToUpper(r[0])
			words[i] = string(r)
		}
	}
	return strings.Join(words, " ")
}
