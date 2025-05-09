package wow

import (
	"fmt"
	"math/rand"
	"strings"
)

func getRandomNumber(min, max int) int {
	if min > max {
		min, max = max, min // swap if min > max
	}
	return rand.Intn(max-min+1) + min
}

func getOs(n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat("o", n)
}

func intAsSubscript(n int) string {
	subscriptDigits := map[rune]rune{
		'0': '₀',
		'1': '₁',
		'2': '₂',
		'3': '₃',
		'4': '₄',
		'5': '₅',
		'6': '₆',
		'7': '₇',
		'8': '₈',
		'9': '₉',
	}

	str := fmt.Sprintf("%d", n)
	var builder strings.Builder

	for _, ch := range str {
		if sub, ok := subscriptDigits[ch]; ok {
			builder.WriteRune(sub)
		} else {
			builder.WriteRune(ch) // fallback (e.g., minus sign)
		}
	}

	return builder.String()
}
