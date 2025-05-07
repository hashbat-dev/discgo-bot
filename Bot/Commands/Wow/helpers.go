package wow

import (
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
