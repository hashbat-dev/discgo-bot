package wow

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	cache "github.com/hashbat-dev/discgo-bot/Cache"
	database "github.com/hashbat-dev/discgo-bot/Database"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

var (
	dataInit                             = false
	dataLowestWowRank  map[string]string = make(map[string]string)
	dataCurrentWeather WeatherResponse
)

func GetEffectData() {
	dataInit = true
	logger.Info("WOW", "Getting Effect Data")

	getDataLowestWowRanks()
	getDataWeatherData()
}

func getDataLowestWowRanks() {
	for _, g := range cache.ActiveGuilds {
		var UserID sql.NullString

		err := database.Db.QueryRow("SELECT UserID FROM DiscGo.WowStats WHERE GuildID = ? ORDER BY MaxWow LIMIT 1", g.DiscordID).Scan(&UserID)
		if err != nil {
			if err != sql.ErrNoRows {
				logger.Error("WOW", err)
			}
			continue
		}

		if UserID.Valid {
			dataLowestWowRank[g.DiscordID] = UserID.String
		}
	}
}

type WeatherCurrentUnits struct {
	Time          string `json:"time"`
	Interval      string `json:"interval"`
	Temperature2m string `json:"temperature_2m"`
	Rain          string `json:"rain"`
	IsDay         string `json:"is_day"`
	WindSpeed10m  string `json:"wind_speed_10m"`
	CloudCover    string `json:"cloud_cover"`
}

type WeatherCurrent struct {
	Time          string  `json:"time"`
	Interval      int     `json:"interval"`
	Temperature2m float64 `json:"temperature_2m"`
	Rain          float64 `json:"rain"`
	IsDay         int     `json:"is_day"`
	WindSpeed10m  float64 `json:"wind_speed_10m"`
	CloudCover    int     `json:"cloud_cover"`
}

type WeatherResponse struct {
	CurrentUnits WeatherCurrentUnits `json:"current_units"`
	Current      WeatherCurrent      `json:"current"`
}

// Gets Weather data for Manchester
func getDataWeatherData() {
	url := "https://api.open-meteo.com/v1/forecast?latitude=53.4809&longitude=-2.2374&current=temperature_2m,rain,is_day,wind_speed_10m,cloud_cover&wind_speed_unit=mph"

	resp, err := http.Get(url)
	if err != nil {
		logger.Error("WOW", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("WOW", fmt.Errorf("unexpected status code: %d", resp.StatusCode))
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("WOW", err)
		return
	}

	var result WeatherResponse
	if err := json.Unmarshal(body, &result); err != nil {
		logger.Error("WOW", err)
		return
	}

	dataCurrentWeather = result
}
