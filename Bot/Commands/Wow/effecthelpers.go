package wow

import (
	"fmt"
	"strconv"
)

func countMatchingLastDigits(s string) int {
	if len(s) == 0 {
		return 0
	}

	lastChar := s[len(s)-1]
	count := 1

	for i := len(s) - 2; i >= 0; i-- {
		if s[i] == lastChar {
			count++
		} else {
			break
		}
	}

	return count
}

func sumDigits(s string) (int, error) {
	sum := 0
	for _, ch := range s {
		digit, err := strconv.Atoi(string(ch))
		if err != nil {
			return 0, fmt.Errorf("invalid character '%c' in input string", ch)
		}
		sum += digit
	}
	return sum, nil
}

func endsInZero(n int) bool {
	return n%10 == 0
}
