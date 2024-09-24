package widgets

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	config "github.com/dabi-ngin/discgo-bot/Config"
	dashboard "github.com/dabi-ngin/discgo-bot/Dashboard"
	"github.com/google/uuid"
)

type GraphWidget struct {
	Options   GraphWidgetOptions
	RefreshMs int
}

const (
	GraphWidgetChartType_Bar = iota
	GraphWidgetChartType_Line
	GraphWidgetChartType_Pie
)

var GraphWidgetChartTypes map[int]string = map[int]string{
	GraphWidgetChartType_Bar:  "bar",
	GraphWidgetChartType_Line: "line",
	GraphWidgetChartType_Pie:  "pie",
}

type GraphWidgetOptions struct {
	Name                 string
	Width                string
	GraphWidgetChartType string
	ChartLabels          []string
	Datasets             []GraphWidgetDataset
	MinValue             int `json:"MinValue,omitempty"`
	MaxValue             int `json:"MaxValue,omitempty"`
	XLabel               string
	YLabel               string
}

type GraphWidgetDataset struct {
	Label            string
	Data             interface{}
	BackgroundColour []string
	BorderColour     []string
	BorderWidth      int
	Fill             bool
	PointRadius      int
}

// Writes JSON data for a GraphWidget object
func SaveGraphWidget(widget GraphWidget) error {
	// No Labels?
	if len(widget.Options.ChartLabels) == 0 {
		maxDataset := 0
		for _, dataset := range widget.Options.Datasets {
			switch data := dataset.Data.(type) {
			case []int:
				if len(data) > maxDataset {
					maxDataset = len(data)
				}
			case []float64:
				// Round each float64 value to 3 decimal places
				for i, v := range data {
					data[i] = math.Round(v*1000) / 1000
				}
				if len(data) > maxDataset {
					maxDataset = len(data)
				}
			default:
				return fmt.Errorf("unsupported dataset data type: %T", dataset.Data)
			}
		}

		widget.Options.ChartLabels = make([]string, maxDataset)
		for i := range widget.Options.ChartLabels {
			widget.Options.ChartLabels[i] = fmt.Sprintf("%d", i+1)
		}
	}

	if widget.RefreshMs == 0 {
		widget.RefreshMs = DefaultMsTimes[DefaultMsGraph]
	}

	saveJson := map[string]interface{}{
		"ID":        uuid.New(),
		"SessionID": config.ServiceSettings.SESSIONID,
		"Name":      widget.Options.Name,
		"Options":   widget.Options,
		"Timestamp": time.Now(),
		"Type":      "graph",
	}

	jsonData, err := json.MarshalIndent(saveJson, "", "  ")
	if err != nil {
		return err
	}

	return dashboard.SaveJsonData(widget.Options.Name, jsonData, widget.Options.Width, widget.RefreshMs)
}
