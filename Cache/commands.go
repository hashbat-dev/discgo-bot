package cache

import (
	"time"

	config "github.com/dabi-ngin/discgo-bot/Config"
)

type Command struct {
	TypeID       int
	Command      string
	GuildID      string
	UserID       string
	UserName     string
	CallTime     time.Time
	CallDuration time.Duration
}

type CmdInfo struct {
	TypeID      int
	Command     string
	Count       int
	AvgDuration time.Duration
	LastCall    time.Time
}

type CmdAverage struct {
	TypeID      int
	Command     string
	Durations   []time.Duration
	AvgDuration time.Duration
}

var Commands []Command
var CommandInfo []CmdInfo
var CommandAverages []CmdAverage

func AddToCommandCache(typeId int, command string, guildId string, userId string, userName string, timeStart time.Time, timeFinish time.Time) {
	callDuration := timeFinish.Sub(timeStart)

	// Commands
	newCmd := Command{
		TypeID:       typeId,
		Command:      command,
		GuildID:      guildId,
		UserID:       userId,
		UserName:     userName,
		CallTime:     timeStart,
		CallDuration: callDuration,
	}
	NewCommandCache := append([]Command{newCmd}, Commands...)

	if len(Commands) > config.DashboardMaxCommands {
		NewCommandCache = NewCommandCache[1:]
	}

	Commands = NewCommandCache

	// CommandInfo
	// 1. Work out Averages
	avgIndex := -1
	for i, avg := range CommandAverages {
		if typeId == avg.TypeID && command == avg.Command {
			avgIndex = i
			break
		}
	}

	if avgIndex == -1 {
		CommandAverages = append(CommandAverages, CmdAverage{
			TypeID:      typeId,
			Command:     command,
			Durations:   []time.Duration{callDuration},
			AvgDuration: callDuration,
		})
		avgIndex = len(CommandAverages) - 1
	} else {
		CommandAverages[avgIndex].Durations = append(CommandAverages[avgIndex].Durations, callDuration)
		CommandAverages[avgIndex].AvgDuration = AverageDuration(CommandAverages[avgIndex].Durations)
	}

	if len(CommandAverages[avgIndex].Durations) > config.CommandAveragePool {
		CommandAverages[avgIndex].Durations = CommandAverages[avgIndex].Durations[1:]
	}

	// 2. Update Info
	infoIndex := -1
	for i, cmd := range CommandInfo {
		if typeId == cmd.TypeID && command == cmd.Command {
			infoIndex = i
			break
		}
	}

	if infoIndex == -1 {
		CommandInfo = append(CommandInfo, CmdInfo{
			TypeID:      typeId,
			Command:     command,
			Count:       1,
			AvgDuration: CommandAverages[avgIndex].AvgDuration,
			LastCall:    time.Now(),
		})
	} else {
		CommandInfo[infoIndex].Count++
		CommandInfo[infoIndex].AvgDuration = CommandAverages[avgIndex].AvgDuration
		CommandInfo[infoIndex].LastCall = time.Now()
	}

}

func AverageDuration(durations []time.Duration) time.Duration {
	var total time.Duration

	for _, duration := range durations {
		total += duration
	}

	if len(durations) == 0 {
		return 0
	}

	average := total / time.Duration(len(durations))
	return average
}
