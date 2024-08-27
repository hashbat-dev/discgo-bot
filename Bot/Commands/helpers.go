package commands

import "strings"

func GetBangCommand(messageContent string) string {
	if len(messageContent) == 0 {
		return ""
	}

	if strings.HasPrefix(messageContent, "!") {
		spaceIndex := strings.Index(messageContent, " ")
		if spaceIndex == -1 {
			// No spaces in the Content, we assume the whole message is the ! command
			return messageContent[1:]
		} else {
			return strings.Split(messageContent, " ")[0]
		}
	}
	return ""
}
