package reporting

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	commands "github.com/dabi-ngin/discgo-bot/Bot/Commands"
	cache "github.com/dabi-ngin/discgo-bot/Cache"
	config "github.com/dabi-ngin/discgo-bot/Config"
	widgets "github.com/dabi-ngin/discgo-bot/Dashboard/Widgets"
	database "github.com/dabi-ngin/discgo-bot/Database"
	helpers "github.com/dabi-ngin/discgo-bot/Helpers"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
	"github.com/google/uuid"
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

func Command(commandTypeId int, Message *discordgo.MessageCreate, Command commands.Command, CorrelationId uuid.UUID, timeStarted time.Time) {

	// 1. Calculate the time taken straight away and log the Event
	timeTaken := time.Since(timeStarted)
	commandName := Command.Name()
	logger.Event(Message.GuildID, fmt.Sprintf("[%v] Command completed successfully after %v [%v]", CorrelationId, timeTaken, commandName))
	database.LogCommandUsage(Message.GuildID, Message.Author.ID, commandTypeId, commandName)

	// Command Log (individual)
	newCmd := widgets.TableWidgetRow{
		Values: []widgets.TableWidgetRowValue{
			{Value: config.CommandTypes[commandTypeId]},
			{Value: commandName},
			{Value: Message.GuildID, TextFormat: widgets.TextFormatString_AbbreviateToEnd},
			{Value: Message.Author.ID, TextFormat: widgets.TextFormatString_AbbreviateToEnd},
			{Value: Message.Author.Username},
			{Value: timeStarted, TextFormat: widgets.TextFormatTime_TimeOnly},
			{Value: timeTaken, TextFormat: widgets.TextFormatDuration_WithMs},
		},
	}
	newDashCmds := append([]widgets.TableWidgetRow{newCmd}, DashCmdRows...)

	if len(newDashCmds) > config.DashboardMaxCommands {
		newDashCmds = newDashCmds[1:]
	}

	DashCmdRows = newDashCmds

	// Command Info (grouped)
	// 1. Work out the new averages
	CmdKey := fmt.Sprintf("%v:%v", commandTypeId, commandName)
	if avg, ok := DashCmdAvgMap[CmdKey]; ok {
		avg.Durations = append(avg.Durations, timeTaken)
		avg.AvgDuration = helpers.AverageDuration(avg.Durations)
		DashCmdAvgMap[CmdKey] = avg
	} else {
		DashCmdAvgMap[CmdKey] = DashCmdAvg{
			TypeID:      commandTypeId,
			Command:     commandName,
			Durations:   []time.Duration{timeTaken},
			AvgDuration: timeTaken,
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
			TypeID:      commandTypeId,
			Type:        config.CommandTypes[commandTypeId],
			Command:     commandName,
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

	err := widgets.SaveTableWidget(&widgets.TableWidget{
		Options: widgets.TableWidgetOptions{
			Name:  "Command Log",
			Width: widgets.WidthHalf,
		},
		Columns: DashCmdsColumns,
		Rows:    DashCmdRows,
	})

	if err != nil {
		logger.Error("REPORTING", err)
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
		logger.Error("REPORTING", err)
	}

	cache.UpdateLastGuildCommand(Message.GuildID)
}
