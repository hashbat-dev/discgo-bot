package reporting

import (
	config "github.com/dabi-ngin/discgo-bot/Config"
	widgets "github.com/dabi-ngin/discgo-bot/Dashboard/Widgets"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
)

var DashCmdsColumns []widgets.TableWidgetColumn
var DashCmdsInfosColumns []widgets.TableWidgetColumn
var GuildColumns []widgets.TableWidgetColumn
var LogColumns []widgets.TableWidgetColumn

func init() {

	DashCmdsColumns = []widgets.TableWidgetColumn{
		{Name: "Type"},
		{Name: "Name"},
		{Name: "Guild ID"},
		{Name: "User ID"},
		{Name: "User Name"},
		{Name: "Called"},
		{Name: "Duration"},
	}

	DashCmdsInfosColumns = []widgets.TableWidgetColumn{
		{Name: "Type"},
		{Name: "Command"},
		{Name: "Count"},
		{Name: "Avg. Time"},
		{Name: "Last Call"},
	}

	GuildColumns = []widgets.TableWidgetColumn{
		{Name: "Name"},
		{Name: "Guild ID"},
		{Name: "Db ID"},
		{Name: "Calls"},
		{Name: "Last Command"},
	}

	LogColumns = []widgets.TableWidgetColumn{
		{Name: "Time"},
		{Name: "Type"},
		{Name: "Guild ID"},
		{Name: "Source"},
		{Name: "Log Text"},
	}

	err := widgets.SaveInfoWidget(&widgets.InfoWidget{
		Name: "Config Items",
		Items: []widgets.InfoWidgetItem{
			{Name: "HostName", Value: config.HostName, Description: "The name of the machine the instance is running on"},
			{Name: "DashboardMaxLogs", Value: config.DashboardMaxLogs, Description: "The Max number of Logs the reporter will store"},
			{Name: "DashboardMaxCommands", Value: config.DashboardMaxCommands, Description: "The Max number of Commands the reporter will store"},
			{Name: "HardwareStatIntervalSeconds", Value: config.HardwareStatIntervalSeconds, Description: "How often in seconds the reporter will poll Hardware statistics"},
			{Name: "HardwareStatMaxIntervals", Value: config.HardwareStatMaxIntervals, Description: "The maximum number of intervals recording for hardware reporting statistics"},
		},
	})
	if err != nil {
		logger.Error("REPORTING", err)
	}

	// Run retrieval Widgets once on Init so the Maps are available to the Dashboard
	Guilds()
	Logs()

}
