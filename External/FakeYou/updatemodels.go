package fakeyou

import (
	"sync"
	"time"

	database "github.com/hashbat-dev/discgo-bot/Database"
	external "github.com/hashbat-dev/discgo-bot/External"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

func UpdateModels() {
	// 1. Check how long it has been since the last update, we'll give a 1 hour lee-way considering this
	// 	  function runs every 12 hours, to avoid it skipping when the last check was 11h 45m ago.
	lastTime, err := database.GetLastCheck("LastFakeYouCheck")
	if err == nil {
		if time.Since(lastTime) <= time.Duration(11*time.Hour) {
			logger.Debug("FAKEYOU", "Skipping Model updates, last update done at: %v", lastTime)
			return
		}
	}

	// 2. Get the Model List, this is a JSON file from FakeYou which lists every TTS Model available to use.
	modelList, err := external.GetJsonFromUrl(FakeYouURLModelList)
	if err != nil {
		logger.Error("FAKEYOU", err)
		return
	}

	// 3. Get a map of all current Models in the Database
	//	  We will delete these as we process them from the JSON list so we'll be left with
	//	  a Map of "orphaned" models which are no longer provided by FakeYou.
	orphanModels, err := database.GetFakeYouModels("")
	if err != nil {
		logger.Error("FAKEYOU", err)
	}

	// 4. Loop through each value and perform our operations.
	//	  Done in a GoRoutine for performance, we want a small pause after each entry.
	var wg sync.WaitGroup

	wg.Add(1)
	go func(list map[string]interface{}) {
		defer wg.Done()
		if models, ok := modelList["models"].([]interface{}); ok {
			for _, model := range models {
				if modelMap, ok := model.(map[string]interface{}); ok {
					title, titleOk := modelMap["title"].(string)
					token, tokenOk := modelMap["model_token"].(string)

					if !titleOk || !tokenOk {
						logger.Debug("FAKEYOU", "Unable to process FakeYou Model Title:[%v] Token:[%v]", titleOk, tokenOk)
						continue
					}

					if database.AddOrUpdateFakeYouModel(title, token) != nil {
						logger.ErrorText("FAKEYOU", "Unable to insert or update model [%v] [%v]", title, token)
					}

					delete(orphanModels, title)
					time.Sleep(time.Millisecond * FAKEYOU_SLEEP_BETWEEN_MODELS_MS)

				} else {
					logger.Error("FAKEYOU", err)
					return
				}
			}
			logger.Info("FAKEYOU", "Completed processing loop of FakeYou TTS Models")
		} else {
			logger.Error("FAKEYOU", err)
			return
		}
	}(modelList)
	wg.Wait()

	// 5. If anything is left in the Orphaned Map, delete it from the Database.
	for _, object := range orphanModels {
		if database.DeleteFakeYouModel(object) != nil {
			logger.Info("FAKEYOU", "Error deleting orphaned model [%v], processing continues", object.Title)
		}
	}

	// 6. Update the Database to say we've completed the check
	database.UpdateLastCheck("LastFakeYouCheck")
}
