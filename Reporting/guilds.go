package reporting

import (
	cache "github.com/dabi-ngin/discgo-bot/Cache"
	config "github.com/dabi-ngin/discgo-bot/Config"
	widgets "github.com/dabi-ngin/discgo-bot/Dashboard/Widgets"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
)

func Guilds() {

	var GuildRows []widgets.TableWidgetRow
	for _, value := range cache.ActiveGuilds {
		if config.ServiceSettings.ISDEV != value.IsDev {
			continue
		}
		GuildRows = append(GuildRows, widgets.TableWidgetRow{
			Values: []widgets.TableWidgetRowValue{
				{Value: value.Name},
				{Value: value.DiscordID, TextFormat: widgets.TextFormatString_AbbreviateToEnd},
				{Value: value.DbID},
				{Value: value.CommandCount},
				{Value: value.LastCommand, TextFormat: widgets.TextFormatTime_TimeOnly},
			},
		})
	}

	err := widgets.SaveTableWidget(&widgets.TableWidget{
		Options: widgets.TableWidgetOptions{
			Name:  "Active Guilds",
			Width: widgets.WidthQuarter,
		},
		Columns: GuildColumns,
		Rows:    GuildRows,
	})
	if err != nil {
		logger.Error("REPORTING", err)
	}

}
