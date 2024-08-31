package dashboard

import (
	"encoding/json"
	"net/http"
	"os"

	logger "github.com/dabi-ngin/discgo-bot/Logger"
)

var jsonCache map[string]map[string]interface{} = make(map[string]map[string]interface{})

func Run() {
	http.HandleFunc("/", webHandler)
	err := http.ListenAndServe(":3333", nil)
	if err != nil {
		logger.Error_IgnoreDiscord("DASHBOARD", err)
	}
}

func webHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/getData":
		handleGetData(w)
	default:
		handleFileRequest(w, r)
	}
}

func handleGetData(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(jsonCache); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Error_IgnoreDiscord("DASHBOARD", err)
		return
	}
}

func handleFileRequest(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path[1:]
	if filePath == "" {
		filePath = "Dashboard/Pages/404.html"
	}

	file, err := os.Open(filePath)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil || fileInfo.IsDir() {
		http.NotFound(w, r)
		return
	}

	http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), file)
}

func SaveJsonData(name string, jsonData []byte) error {
	var data map[string]interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return err
	}
	jsonCache[name] = data
	return nil
}
