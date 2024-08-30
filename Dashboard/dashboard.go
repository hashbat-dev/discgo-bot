package dashboard

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	logger "github.com/dabi-ngin/discgo-bot/Logger"
)

// CommandLatencies is our growable map which records what the runtime of each command has been as the workers complete running them.
var CommandLatencies map[string][]time.Duration

func init() {
	CommandLatencies = make(map[string][]time.Duration)
}

func Run() {
	http.HandleFunc("/", webhandler)
	err := http.ListenAndServe(":3333", nil)
	if err != nil {
		logger.Error_IgnoreDiscord("DASHBOARD", err)
	}
}

var moduleRoot string = "Dashboard/"

func webhandler(w http.ResponseWriter, r *http.Request) {
	directory := getUrlDirectory(r.URL.Path)
	if directory == "getdata" {
		returnData(w)
	} else {
		getDashboard(w, directory)
	}
}

func returnData(w http.ResponseWriter) {
	// Get the current Packet Data
	packets := PacketCache

	// Send the Data back
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(packets); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error_IgnoreDiscord("DASHBOARD", err)
		return
	}

}

func getDashboard(w http.ResponseWriter, directory string) {
	filePath := moduleRoot
	if strings.Contains(directory, ".") {
		filePath += "Resources/" + directory
	} else {
		filePath += "Pages/" + directory + ".html"
	}
	body, err := os.ReadFile(filePath)

	writeBack := ""
	if err != nil {
		logger.Error_IgnoreDiscord("DASHBOARD", err)
		writeBack = "Uh-oh" // Make sure we write something back, needed for deployment
	} else {
		writeBack = string(body)
	}

	_, err = io.WriteString(w, writeBack)
	if err != nil {
		logger.Error_IgnoreDiscord("DASHBOARD", err)
	}
}

func getUrlDirectory(s string) string {
	if s == "/" {
		return "dashboard"
	}

	firstSlashIndex := strings.Index(s, "/")
	if firstSlashIndex == -1 {
		return ""
	}

	afterFirstSlash := s[firstSlashIndex+1:]

	secondSlashIndex := strings.Index(afterFirstSlash, "/")
	if secondSlashIndex == -1 {
		return strings.ToLower(afterFirstSlash)
	}

	return strings.ToLower(afterFirstSlash[:secondSlashIndex])
}
