package reporting

import (
	"sort"

	config "github.com/hashbat-dev/discgo-bot/Config"
	widgets "github.com/hashbat-dev/discgo-bot/Dashboard/Widgets"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

var DashCmdsColumns []widgets.TableWidgetColumn
var DashCmdsInfosColumns []widgets.TableWidgetColumn
var GuildColumns []widgets.TableWidgetColumn
var LogColumns []widgets.TableWidgetColumn
var LogFilters []widgets.TableWidgetFilter

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

	var logLevels []int
	for key := range config.LoggingLevels {
		logLevels = append(logLevels, key)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(logLevels)))
	var typeValues []string
	for _, key := range logLevels {
		typeValues = append(typeValues, config.LoggingLevels[key].Name)
	}

	LogFilters = []widgets.TableWidgetFilter{
		{Name: "Guild ID", FilterType: widgets.TableWidgetType_SelectRegular, ColumnNames: []string{"Guild ID"}, FullMatchOnly: true},
		{Name: "Types", FilterType: widgets.TableWidgetType_SelectCheckbox, ColumnNames: []string{"Type"}, Values: typeValues, FullMatchOnly: true},
		{Name: "Search", FilterType: widgets.TableWidgetType_FreeText, ColumnNames: []string{"Time", "Guild ID", "Source", "Log Text"}},
	}

	err := widgets.SaveInfoWidget(&widgets.InfoWidget{
		Name: "Config Items",
		Items: []widgets.InfoWidgetItem{
			{Name: "HostName", Value: config.ServiceSettings.HOSTNAME, Description: "The name of the machine the instance is running on"},
			{Name: "DashboardMaxLogs", Value: config.ServiceSettings.DASHBOARDMAXLOGS, Description: "The Max number of Logs the reporter will store"},
			{Name: "DashboardMaxCommands", Value: config.ServiceSettings.DASHBOARDMAXCOMMANDS, Description: "The Max number of Commands the reporter will store"},
			{Name: "HardwareStatIntervalSeconds", Value: config.ServiceSettings.HARDWARESTATINTERVALSECONDS, Description: "How often in seconds the reporter will poll Hardware statistics"},
			{Name: "HardwareStatMaxIntervals", Value: config.ServiceSettings.HARDWARESTATMAXINTERVALS, Description: "The maximum number of intervals recording for hardware reporting statistics"},
		},
	})
	if err != nil {
		logger.Error("REPORTING", err)
	}

	// Run retrieval Widgets once on Init so the Maps are available to the Dashboard
	Guilds()
	Logs()
}
