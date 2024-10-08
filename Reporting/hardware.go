package reporting

import (
	"time"

	config "github.com/hashbat-dev/discgo-bot/Config"
	widgets "github.com/hashbat-dev/discgo-bot/Dashboard/Widgets"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

var cpuValuesCache []float64 = []float64{}
var ramValuesCache []float64 = []float64{}
var ramMaxValue int = -1

func Hardware() {
	// CPU -----------------------------------------------------------------------------
	cpuPercentage, err := cpu.Percent(time.Second*time.Duration(config.ServiceSettings.HARDWARESTATINTERVALSECONDS), false)
	if err != nil {
		logger.Error("REPORTING", err)
	} else {
		// CPU Package can return more specific results but we only want the system's overall value so only first value is needed
		cpuValuesCache = append(cpuValuesCache, cpuPercentage[0])

		// Are we over the interval count?
		if len(cpuValuesCache) > config.ServiceSettings.HARDWARESTATMAXINTERVALS {
			cpuValuesCache = cpuValuesCache[1:]
		}
	}

	err = widgets.SaveGraphWidget(widgets.GraphWidget{
		Options: widgets.GraphWidgetOptions{
			Name:                 "CPU Usage",
			Width:                widgets.WidthQuarter,
			GraphWidgetChartType: widgets.GraphWidgetChartTypes[widgets.GraphWidgetChartType_Line],
			Datasets: []widgets.GraphWidgetDataset{
				{
					Label:            "Usage (%)",
					Data:             cpuValuesCache,
					BackgroundColour: []string{config.Colours["blue"].GraphTransparent},
					BorderColour:     []string{config.Colours["blue"].GraphOpaque},
					BorderWidth:      1,
					Fill:             false,
					PointRadius:      0,
				},
			},
			XLabel:   "Time",
			YLabel:   "Usage (%)",
			MaxValue: 100,
		},
	})
	if err != nil {
		logger.Error("REPORTING", err)
	}

	// RAM ------------------------------------------------------------------------------
	ramUsage, err := mem.VirtualMemory()
	if err != nil {
		logger.Error("", err)
	} else {
		availableMb := ramUsage.Available / 1024 / 1024
		ramValuesCache = append(ramValuesCache, float64(availableMb))

		// Are we over the interval count?
		if len(ramValuesCache) > config.ServiceSettings.HARDWARESTATMAXINTERVALS {
			ramValuesCache = ramValuesCache[1:]
		}
	}
	if ramMaxValue < 0 {
		ramMaxValue = int(ramUsage.Total)
	}

	err = widgets.SaveGraphWidget(widgets.GraphWidget{
		Options: widgets.GraphWidgetOptions{
			Name:                 "RAM Usage",
			Width:                widgets.WidthQuarter,
			GraphWidgetChartType: widgets.GraphWidgetChartTypes[widgets.GraphWidgetChartType_Line],
			Datasets: []widgets.GraphWidgetDataset{
				{
					Label:            "Usage (Mb)",
					Data:             ramValuesCache,
					BackgroundColour: []string{config.Colours["green"].GraphTransparent},
					BorderColour:     []string{config.Colours["green"].GraphOpaque},
					BorderWidth:      1,
					Fill:             false,
					PointRadius:      0,
				},
			},
			XLabel: "Time",
			YLabel: "Usage (Mb)",
		},
	})
	if err != nil {
		logger.Error("REPORTING", err)
	}

}
