package fakeyou

import (
	"sync"
	"time"

	database "github.com/dabi-ngin/discgo-bot/Database"
	external "github.com/dabi-ngin/discgo-bot/External"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
)

func UpdateModels() {
	// 1. Get the Model List, this is a JSON file from FakeYou which lists every TTS Model available to use.
	modelList, err := external.GetJsonFromUrl(FakeYouURLModelList)
	if err != nil {
		logger.Error("FAKEYOU", err)
		return
	}

	// 2. Get a map of all current Models in the Database
	//	  We will delete these as we process them from the JSON list so we'll be left with
	//	  a Map of "orphaned" models which are no longer provided by FakeYou.
	orphanModels, err := database.GetFakeYouModels("")
	if err != nil {
		logger.Error("FAKEYOU", err)
	}

	// 3. Loop through each value and perform our operations.
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

	// 4. If anything is left in the Orphaned Map, delete it from the Database.
	for _, object := range orphanModels {
		if database.DeleteFakeYouModel(object) != nil {
			logger.Info("FAKEYOU", "Error deleting orphaned model [%v], processing continues", object.Title)
		}
	}

}
