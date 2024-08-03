package imgwork

import (
	"strings"

	config "github.com/dabi-ngin/discgo-bot/Config"
)

func GetExtensionFromURL(url string) string {
	if qIndex := strings.Index(url, "?"); qIndex != -1 {
		url = url[:qIndex]
	}

	lastSlashIndex := strings.LastIndex(url, "/")
	if lastSlashIndex != -1 {
		url = url[lastSlashIndex+1:]
	}

	lastDotIndex := strings.LastIndex(url, ".")
	if lastDotIndex != -1 {
		extension := url[lastDotIndex:]
		if IsExtensionAccepted(extension) {
			return extension
		} else {
			return ""
		}
	}

	return ""
}

func IsExtensionAccepted(extension string) bool {
	for _, ext := range config.ValidImageExtensions {
		if extension == ext {
			return true
		}
	}
	return false
}
