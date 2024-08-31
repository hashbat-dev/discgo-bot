package handlers

import (
	"fmt"
	"time"

	config "github.com/dabi-ngin/discgo-bot/Config"
	widgets "github.com/dabi-ngin/discgo-bot/Dashboard/Widgets"
	helpers "github.com/dabi-ngin/discgo-bot/Helpers"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
)

type DashCmdAvg struct {
	TypeID      int
	Command     string
	Durations   []time.Duration
	AvgDuration time.Duration
}
type DashCmdInfo struct {
	TypeID      int
	Type        string
	Command     string
	Count       int
	AvgDuration time.Duration
	LastCall    time.Time
}

var DashCmdRows []widgets.TableWidgetRow
var DashCmdInfoRows []widgets.TableWidgetRow

var DashCmdInfoMap map[string]DashCmdInfo = make(map[string]DashCmdInfo)
var DashCmdAvgMap map[string]DashCmdAvg = make(map[string]DashCmdAvg)

var DashCmdsColumns []widgets.TableWidgetColumn
var DashCmdsInfosColumns []widgets.TableWidgetColumn

func ReportCommand(typeId int, command string, guildId string, userId string, userName string, timeStart time.Time, callDuration time.Duration) {

	logger.Event(guildId, fmt.Sprintf("Command completed successfully after %v [%v]", callDuration, command))

	// Command Log (individual)
	newCmd := widgets.TableWidgetRow{
		Values: []widgets.TableWidgetRowValue{
			{Value: config.CommandTypes[typeId]},
			{Value: guildId, TextFormat: widgets.TextFormatString_AbbreviateToEnd},
			{Value: userId, TextFormat: widgets.TextFormatString_AbbreviateToEnd},
			{Value: userName},
			{Value: timeStart, TextFormat: widgets.TextFormatTime_TimeOnly},
			{Value: callDuration, TextFormat: widgets.TextFormatDuration_WithMs},
		},
	}
	newDashCmds := append([]widgets.TableWidgetRow{newCmd}, DashCmdRows...)

	if len(newDashCmds) > config.DashboardMaxCommands {
		newDashCmds = newDashCmds[1:]
	}

	DashCmdRows = newDashCmds

	// Command Info (grouped)
	// 1. Work out the new averages
	CmdKey := fmt.Sprintf("%v:%v", typeId, command)
	if avg, ok := DashCmdAvgMap[CmdKey]; ok {
		avg.Durations = append(avg.Durations, callDuration)
		avg.AvgDuration = helpers.AverageDuration(avg.Durations)
		DashCmdAvgMap[CmdKey] = avg
	} else {
		DashCmdAvgMap[CmdKey] = DashCmdAvg{
			TypeID:      typeId,
			Command:     command,
			Durations:   []time.Duration{callDuration},
			AvgDuration: callDuration,
		}
	}

	// 2. Update Info
	if info, ok := DashCmdInfoMap[CmdKey]; ok {
		info.Count++
		info.AvgDuration = DashCmdAvgMap[CmdKey].AvgDuration
		info.LastCall = time.Now()
		DashCmdInfoMap[CmdKey] = info
	} else {
		DashCmdInfoMap[CmdKey] = DashCmdInfo{
			TypeID:      typeId,
			Type:        config.CommandTypes[typeId],
			Command:     command,
			Count:       1,
			AvgDuration: DashCmdAvgMap[CmdKey].AvgDuration,
			LastCall:    time.Now(),
		}
	}

	// 3. Generate Widget
	var newCmdInfoRows []widgets.TableWidgetRow
	for _, value := range DashCmdInfoMap {
		newCmdInfoRows = append(newCmdInfoRows, widgets.TableWidgetRow{
			Values: []widgets.TableWidgetRowValue{
				{Value: value.Type},
				{Value: value.Command},
				{Value: value.Count},
				{Value: value.AvgDuration, TextFormat: widgets.TextFormatDuration_WithMs},
				{Value: value.LastCall, TextFormat: widgets.TextFormatTime_TimeOnly},
			},
		})
	}

	DashCmdInfoRows = newCmdInfoRows

	checkWidgetColumnsExist()

	err := widgets.SaveTableWidget(&widgets.TableWidget{
		Options: widgets.TableWidgetOptions{
			Name:  "Command Log",
			Width: widgets.WidthHalf,
		},
		Columns: DashCmdsColumns,
		Rows:    DashCmdRows,
	})

	if err != nil {
		logger.Error("DASHBOARD", err)
	}

	err = widgets.SaveTableWidget(&widgets.TableWidget{
		Options: widgets.TableWidgetOptions{
			Name:  "Command Info",
			Width: widgets.WidthQuarter,
		},
		Columns: DashCmdsInfosColumns,
		Rows:    DashCmdInfoRows,
	})

	if err != nil {
		logger.Error("DASHBOARD", err)
	}
}

func checkWidgetColumnsExist() {
	if len(DashCmdsColumns) == 0 {
		DashCmdsColumns = []widgets.TableWidgetColumn{
			{Name: "Type"},
			{Name: "Name"},
			{Name: "Guild ID"},
			{Name: "User ID"},
			{Name: "User Name"},
			{Name: "Called"},
			{Name: "Duration"},
		}
	}

	if len(DashCmdsInfosColumns) == 0 {
		DashCmdsInfosColumns = []widgets.TableWidgetColumn{
			{Name: "Type"},
			{Name: "Command"},
			{Name: "Count"},
			{Name: "Avg. Time"},
			{Name: "Last Call"},
		}
	}
}
