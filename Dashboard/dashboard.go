package dashboard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

	config "github.com/hashbat-dev/discgo-bot/Config"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

var jsonCache map[string]map[string]interface{} = make(map[string]map[string]interface{})
var jsonCacheOrder []JsonCacheOrder

type JsonCacheOrder struct {
	CacheKey  string
	Width     int
	RefreshMs int
}

func Run() {
	http.HandleFunc("/Resources/", resourcesHandler)
	http.HandleFunc("/", webHandler)
	err := http.ListenAndServe(":3333", nil)
	if err != nil {
		logger.Error_IgnoreDiscord("DASHBOARD", err)
	}
}

func webHandler(w http.ResponseWriter, r *http.Request) {
	checkForDashboardMessage(r)
	switch r.URL.Path {
	case "/getData":
		handleGetData(w, r)
	default:
		handleFileRequest(w, r)
	}
}

func resourcesHandler(w http.ResponseWriter, r *http.Request) {
	handleFileRequest(w, r)
}

func checkForDashboardMessage(r *http.Request) {
	if config.ServiceSettings.DASHBOARDURL == "" {
		config.ServiceSettings.DASHBOARDURL = fmt.Sprintf("http://%v/", r.Host)
	}
}

func handleGetData(w http.ResponseWriter, r *http.Request) {
	// What is being requested?
	requestedWidget := ""
	query := r.URL.Query()
	if query != nil {
		requestedWidget = query.Get("widget")
	}

	if requestedWidget == "" {
		// Requesting Widget Overview
		returnWidgetOverview(w)
	} else {
		// Requesting Widget Data
		returnWidgetData(w, requestedWidget)
	}

}

func returnWidgetOverview(w http.ResponseWriter) {
	response := make([]map[string]interface{}, 0, len(jsonCacheOrder))

	// Loop through jsonCacheOrder and add the CacheKey and RefreshMs to the response
	for _, order := range jsonCacheOrder {
		response = append(response, map[string]interface{}{
			"Widget":    order.CacheKey,
			"RefreshMs": order.RefreshMs,
			"SessionID": config.ServiceSettings.SESSIONID,
		})
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		logger.Error("DASHBOARD", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the response headers and write the data
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonData)
	if err != nil {
		logger.Error("DASHBOARD", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func returnWidgetData(w http.ResponseWriter, widget string) {
	w.Header().Set("Content-Type", "application/json")

	if data, exists := jsonCache[widget]; exists {
		// Widget exists, return its data
		if err := json.NewEncoder(w).Encode(data); err != nil {
			logger.Error("DASHBOARD", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// Requested Widget not found
		logger.Debug("DASHBOARD", "requested widget '%s' not found", widget)
		http.Error(w, "Widget not found", http.StatusNotFound)
	}
}

func handleFileRequest(w http.ResponseWriter, r *http.Request) {
	fileRoot := ""
	filePath := r.URL.Path[1:]

	if !strings.Contains(strings.ToLower(r.URL.Path), "/temp/") {
		fileRoot = "Dashboard/"
	}
	if filePath == "" {
		filePath = fileRoot + "Pages/dashboard.html"
	} else {
		filePath = fileRoot + filePath
	}

	file, err := os.Open(filePath)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer func() {
		err := file.Close()
		if err != nil {
			logger.Error("DASHBOARD", err)
		}
	}()

	fileInfo, err := file.Stat()
	if err != nil || fileInfo.IsDir() {
		http.NotFound(w, r)
		return
	}

	http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), file)
}

var saveJsonMutex sync.Mutex

func SaveJsonData(name string, jsonData []byte, width string, refreshMs int) error {
	saveJsonMutex.Lock()
	defer saveJsonMutex.Unlock()

	// Save the JSON Data into the map
	var data map[string]interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return err
	}

	jsonCache[name] = data

	// Remove the entry if it already exists in jsonCacheOrder
	for i, order := range jsonCacheOrder {
		if order.CacheKey == name {
			jsonCacheOrder = append(jsonCacheOrder[:i], jsonCacheOrder[i+1:]...)
			break
		}
	}

	widthInt, err := strconv.Atoi(strings.TrimSuffix(width, "%"))
	if err != nil {
		return err
	}

	jsonCacheOrder = append(jsonCacheOrder, JsonCacheOrder{
		CacheKey:  name,
		Width:     widthInt,
		RefreshMs: refreshMs,
	})

	// Sort jsonCacheOrder by width descending
	sort.SliceStable(jsonCacheOrder, func(i, j int) bool {
		return jsonCacheOrder[i].Width > jsonCacheOrder[j].Width
	})

	// Reorder to attempt to get rows which sum up to 100% in width
	var reordered []JsonCacheOrder
	var remaining []JsonCacheOrder

	for len(jsonCacheOrder) > 0 {
		var currentGroup []JsonCacheOrder
		currentSum := 0
		i := 0

		for i < len(jsonCacheOrder) {
			item := jsonCacheOrder[i]
			if currentSum+item.Width <= 100 {
				currentGroup = append(currentGroup, item)
				currentSum += item.Width
				// Remove the item from the list
				jsonCacheOrder = append(jsonCacheOrder[:i], jsonCacheOrder[i+1:]...)
			} else {
				i++
			}

			// If currentSum is exactly 100%, complete the group
			if currentSum == 100 {
				break
			}
		}

		// If a 100% group was formed, add it to the reordered list
		if currentSum == 100 {
			reordered = append(reordered, currentGroup...)
		} else {
			// If no 100% group could be formed, save the remaining items for later
			remaining = append(remaining, currentGroup...)
		}
	}

	// Add any remaining items that couldn't be grouped to 100% and add to the cache
	reordered = append(reordered, remaining...)
	jsonCacheOrder = reordered

	return nil
}
