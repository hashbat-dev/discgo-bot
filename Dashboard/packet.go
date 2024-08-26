package dashboard

import (
	"time"

	cache "github.com/dabi-ngin/discgo-bot/Cache"
	config "github.com/dabi-ngin/discgo-bot/Config"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
	"github.com/google/uuid"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

var PacketCache Packet

type Packet struct {
	PacketInfo   DashboardPacketInfo
	ActiveGuilds []cache.Guild
	HardwareInfo DashboardHardware
	Logging      DashboardLogging
	Commands     DashboardCommands
	Pools        []DashboardChannelInfo
}

type DashboardPacketInfo struct {
	PacketID         string
	TimeStamp        time.Time
	ServerName       string
	MaxLogs          int
	MaxCommands      int
	HardwareInterval int
	HardwareMax      int
}

type DashboardLogging struct {
	LogLevels  map[int]string
	LogEntries []logger.DashboardLog
}

type DashboardHardware struct {
	CpuValues []float64
	RamValues []float64
}

type DashboardCommands struct {
	CommandTypes map[int]string
	Commands     []cache.Command
	CommandInfo  []cache.CmdInfo
}

type CommandTypes struct {
	CommandID   int
	CommandType string
}

func GetPacket() Packet {
	var packet Packet

	packet.PacketInfo = getPacketInfo()
	packet.Logging = getDashboardLogging()
	packet.ActiveGuilds = getDashboardGuilds()
	packet.HardwareInfo = getDashboardHardware()
	packet.Commands = getDashboardCommands()
	packet.Pools = getPoolInfo()

	return packet

}

func getPacketInfo() DashboardPacketInfo {
	var returnStruct DashboardPacketInfo
	returnStruct.PacketID = uuid.New().String()
	returnStruct.TimeStamp = time.Now()
	returnStruct.ServerName = config.HostName
	returnStruct.MaxLogs = config.DashboardMaxLogs
	returnStruct.MaxCommands = config.DashboardMaxCommands
	returnStruct.HardwareInterval = config.HardwareStatIntervalSeconds
	returnStruct.HardwareMax = config.HardwareStatMaxIntervals

	return returnStruct
}

func getDashboardLogging() DashboardLogging {
	var returnStruct DashboardLogging
	returnStruct.LogLevels = config.LoggingLevels
	returnStruct.LogEntries = logger.LogsForDashboard

	return returnStruct
}

func getDashboardGuilds() []cache.Guild {
	return cache.ActiveGuilds
}

var dashboardHardwareCache DashboardHardware

func getDashboardHardware() DashboardHardware {
	var newHardware DashboardHardware
	if dashboardHardwareCache.CpuValues == nil {
		dashboardHardwareCache.CpuValues = []float64{}
		dashboardHardwareCache.RamValues = []float64{}
	} else {
		newHardware = dashboardHardwareCache
	}

	// CPU -----------------------------------------------------------------------------
	cpuPercentage, err := cpu.Percent(time.Second*time.Duration(config.HardwareStatIntervalSeconds), false)
	if err != nil {
		logger.Error("", err)
	} else {
		// CPU Package can return more specific results but we only want the system's overall value so only first value is needed
		newHardware.CpuValues = append(newHardware.CpuValues, cpuPercentage[0])

		// Are we over the interval count?
		if len(newHardware.CpuValues) > config.HardwareStatIntervalSeconds {
			newHardware.CpuValues = newHardware.CpuValues[1:]
		}
	}

	// RAM ------------------------------------------------------------------------------
	ramUsage, err := mem.VirtualMemory()
	if err != nil {
		logger.Error("", err)
	} else {

		availableMb := ramUsage.Available / 1024 / 1024
		newHardware.RamValues = append(newHardware.RamValues, float64(availableMb))

		// Are we over the interval count?
		if len(newHardware.RamValues) > config.HardwareStatMaxIntervals {
			newHardware.RamValues = newHardware.RamValues[1:]
		}
	}

	dashboardHardwareCache = newHardware
	return newHardware
}

func getDashboardCommands() DashboardCommands {
	var returnStruct DashboardCommands
	returnStruct.CommandTypes = config.CommandTypes
	returnStruct.Commands = cache.Commands
	returnStruct.CommandInfo = cache.CommandInfo

	return returnStruct
}

func getPoolInfo() []DashboardChannelInfo {
	var returnStruct []DashboardChannelInfo

	for i, pool := range config.ProcessPools {
		returnStruct = append(returnStruct, DashboardChannelInfo{
			Name:                pool.PoolName,
			ProcessingCount:     PoolProcessing[i],
			ProcessingLastAdded: PoolLastAdded[i],
			QueueCount:          PoolQueue[i],
			QueueLastAdded:      QueueLastAdded[i],
			AverageDuration:     PoolDurations[i].AvgDuration,
		})
	}

	return returnStruct
}
