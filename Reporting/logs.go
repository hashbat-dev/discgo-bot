package reporting

import (
	"sort"

	config "github.com/dabi-ngin/discgo-bot/Config"
	widgets "github.com/dabi-ngin/discgo-bot/Dashboard/Widgets"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
)

func Logs() {
	// The process of getting Logs needs to be Reporting getting the information from
	// Logger to avoid import cycles.

	// 1. Get and reset straight away to avoid wiping un-processed logs at the end of the loop
	obtainedLogs := logger.LogsForDashboard

	var newRows []widgets.TableWidgetRow
	for _, log := range obtainedLogs {
		newRows = append(newRows, widgets.TableWidgetRow{
			Values: []widgets.TableWidgetRowValue{
				{Value: log.LogInfo.DateTime, TextFormat: widgets.TextFormatTime_TimeWithMs},
				{Value: config.LoggingLevels[log.LogLevel].Name},
				{Value: log.LogInfo.GuildID, TextFormat: widgets.TextFormatString_AbbreviateToEnd},
				{Value: log.LogInfo.CodeSource},
				{Value: log.LogText},
			},
			TextColour: config.LoggingLevels[log.LogLevel].Colour,
		})
		AddGuildIDToFilter(log.LogInfo.GuildID)
	}
	err := widgets.SaveTableWidget(&widgets.TableWidget{
		Options: widgets.TableWidgetOptions{
			Name:  "Recent Logs",
			Width: widgets.WidthThreeQuarters,
		},
		Columns: LogColumns,
		Rows:    newRows,
		Filters: LogFilters,
	})
	if err != nil {
		logger.Error("REPORTING", err)
	}

}

var guildExists map[string]struct{} = make(map[string]struct{})
var guildValues []string

func AddGuildIDToFilter(guildId string) {
	if guildId == "" {
		return
	}
	if _, exists := guildExists[guildId]; !exists {
		guildExists[guildId] = struct{}{}
		guildValues = append(guildValues, guildId)
		sort.Strings(guildValues)
		LogFilters[0].Values = guildValues
	}
}
