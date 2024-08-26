package cache

import (
	"time"

	config "github.com/dabi-ngin/discgo-bot/Config"
	helpers "github.com/dabi-ngin/discgo-bot/Helpers"
)

func AddToCommandCache(typeId int, command string, guildId string, userId string, userName string, timeStart time.Time, callDuration time.Duration) {

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
		CommandAverages[avgIndex].AvgDuration = helpers.AverageDuration(CommandAverages[avgIndex].Durations)
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
