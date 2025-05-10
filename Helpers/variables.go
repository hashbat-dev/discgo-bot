package helpers

import (
	"fmt"
	"math/rand"
	"reflect"
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

func IsZero(v any) bool {
	switch val := v.(type) {
	case int:
		return val == 0
	case int64:
		return val == 0
	case uint:
		return val == 0
	case uint64:
		return val == 0
	default:
		return reflect.ValueOf(v).IsZero()
	}
}

func NiceDateFormat(t time.Time) string {
	now := time.Now()
	loc := now.Location()
	t = t.In(loc)

	// Strip time to midnight for date comparison
	y, m, d := now.Date()
	today := time.Date(y, m, d, 0, 0, 0, 0, loc)

	y2, m2, d2 := t.Date()
	inputDate := time.Date(y2, m2, d2, 0, 0, 0, 0, loc)

	diff := today.Sub(inputDate).Hours() / 24
	timePart := t.Format("15:04")

	switch {
	case diff == 0:
		return "Today " + timePart
	case diff == 1:
		return "Yesterday " + timePart
	case diff > 1 && diff <= 7:
		return fmt.Sprintf("%.0f days ago %s", diff, timePart)
	default:
		day := t.Day()
		month := t.Month().String()
		return fmt.Sprintf("%d%s %s %s", day, ordinal(day), month, timePart)
	}
}

// ordinal returns the English ordinal suffix for a given day
func ordinal(n int) string {
	if n >= 11 && n <= 13 {
		return "th"
	}
	switch n % 10 {
	case 1:
		return "st"
	case 2:
		return "nd"
	case 3:
		return "rd"
	default:
		return "th"
	}
}

func GetRandomNumber(min, max int) int {
	if min > max {
		min, max = max, min // swap if min > max
	}
	return rand.Intn(max-min+1) + min
}
