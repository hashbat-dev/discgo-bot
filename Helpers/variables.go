package helpers

import (
	"bytes"
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

func LettersNumbersAndDashesOnly(input string) string {
	var buffer bytes.Buffer

	for i := 0; i < len(input); i++ {
		switch input[i] {
		case ' ':
			buffer.WriteByte('_')
		case '-':
			buffer.WriteString("--")
		default:
			if unicode.IsLetter(rune(input[i])) || unicode.IsDigit(rune(input[i])) {
				buffer.WriteByte(input[i])
			}
		}
	}

	return buffer.String()
}
