package wow

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"

	logger "github.com/hashbat-dev/discgo-bot/Logger"
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

type UselessFact struct {
	Text string `json:"text"`
}

func checkIfRandomFactHasStats() (string, bool) {
	resp, err := http.Get("https://uselessfacts.jsph.pl/api/v2/facts/random")
	if err != nil {
		logger.Error("WOW", err)
		return "", false
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			logger.Error("WOW", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("WOW", err)
		return "", false
	}

	var fact UselessFact
	if err := json.Unmarshal(body, &fact); err != nil {
		logger.Error("WOW", err)
		return "", false
	}

	hasNumbers, _ := regexp.MatchString(`\d`, fact.Text)
	if hasNumbers {
		return fact.Text, true
	}
	return "", false
}
