package reporting

import (
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
	}

	err := widgets.SaveTableWidget(&widgets.TableWidget{
		Options: widgets.TableWidgetOptions{
			Name:  "Recent Logs",
			Width: widgets.WidthThreeQuarters,
		},
		Columns:   LogColumns,
		Rows:      newRows,
		RefreshMs: 500,
	})
	if err != nil {
		logger.Error("REPORTING", err)
	}

}
