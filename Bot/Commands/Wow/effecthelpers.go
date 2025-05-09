package wow

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
