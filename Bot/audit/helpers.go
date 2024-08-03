package audit

import (
	"runtime"
	"strings"
)

func GetStack() string {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			return ParseStackTrace(string(buf[:n]))
		}
		buf = make([]byte, len(buf)*2)
	}
}

func ParseStackTrace(stack string) string {

	retVal := ""
	lines := strings.Split(stack, "\n")

	isFirst := true

	for i := 1; i < len(lines)-1; i++ {

		line := lines[i]
		if strings.Contains(line, "audit") || strings.Contains(line, "created by") || strings.Contains(line, "main.go") {
			continue
		}

		botIndex := strings.Index(line, "Bot/")
		botIndexLength := 4
		if botIndex == -1 {
			botIndex = strings.Index(line, "femboy-control/")
			botIndexLength = 15
		}
		lastIndex := 0

		isFileLine := strings.Contains(line, " +")
		if isFileLine {
			lastIndex = strings.LastIndex(line, " +") - 1
		} else {
			lastIndex = strings.LastIndex(line, ")")
		}

		if !isFileLine && !logFunctions {
			continue
		}

		if botIndex != -1 && lastIndex != -1 {

			appended := false
			if isFirst {
				isFirst = false
			} else {
				if !verboseStack {
					break
				}
				retVal += " <= "
				appended = true
			}

			retVal += line[botIndex+botIndexLength : lastIndex+1]

			if appended && !verboseStack {
				break
			}

		}
	}

	return RemoveTextInParentheses(retVal)
}

func RemoveTextInParentheses(input string) string {
	var result strings.Builder
	inParens := false

	for _, char := range input {
		if char == '(' {
			inParens = true
			result.WriteRune(char)
		} else if char == ')' {
			inParens = false
			result.WriteRune(char)
		} else if !inParens {
			result.WriteRune(char)
		}
	}

	return result.String()
}
