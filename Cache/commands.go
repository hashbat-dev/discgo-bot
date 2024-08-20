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
	var newAvgCache []CmdAverage
	var cmdAvg CmdAverage
	found := false
	for _, avg := range CommandAverages {
		if typeId == avg.TypeID && command == avg.Command {
			cmdAvg = avg
			found = true
		} else {
			newAvgCache = append(newAvgCache, avg)
		}
	}

	if !found {
		cmdAvg = CmdAverage{
			TypeID:  typeId,
			Command: command,
		}
	}

	cmdAvg.Durations = append(cmdAvg.Durations, callDuration)
	cmdAvg.AvgDuration = AverageDuration(cmdAvg.Durations)

	if len(cmdAvg.Durations) > config.CommandAveragePool {
		cmdAvg.Durations = cmdAvg.Durations[1:]
	}

	newAvgCache = append(newAvgCache, cmdAvg)
	CommandAverages = newAvgCache

	// 2. Update Info
	var newCmdInfo []CmdInfo
	var cmdInfo CmdInfo
	found = false
	for _, cmd := range CommandInfo {
		if typeId == cmd.TypeID && command == cmd.Command {
			cmdInfo = cmd
			found = true
		} else {
			newCmdInfo = append(newCmdInfo, cmd)
		}
	}

	if !found {
		cmdInfo = CmdInfo{
			TypeID:  typeId,
			Command: command,
		}
	}

	cmdInfo.Count++
	cmdInfo.AvgDuration = cmdAvg.AvgDuration
	cmdInfo.LastCall = time.Now()

	newCmdInfo = append(newCmdInfo, cmdInfo)
	CommandInfo = newCmdInfo

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
