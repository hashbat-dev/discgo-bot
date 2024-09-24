package widgets

import (
	"encoding/json"
	"time"

	config "github.com/dabi-ngin/discgo-bot/Config"
	dashboard "github.com/dabi-ngin/discgo-bot/Dashboard"
	"github.com/google/uuid"
)

type TableWidget struct {
	ID        string
	SessionID string
	Timestamp time.Time
	Type      string
	Options   TableWidgetOptions
	Columns   []TableWidgetColumn
	Rows      []TableWidgetRow
	Filters   []TableWidgetFilter `json:"Filters,omitempty"`
	RefreshMs int
}

type TableWidgetOptions struct {
	Name  string `json:"Name,omitempty"`
	Width string `json:"Width,omitempty"`
}

type TableWidgetColumn struct {
	Name     string `json:"Name,omitempty"`
	FontSize int    `json:"FontSize,omitempty"`
	Width    string `json:"Width,omitempty"`
}

type TableWidgetRow struct {
	Values     []TableWidgetRowValue
	TextColour config.Colour
}

type TableWidgetRowValue struct {
	Value      interface{}
	TextFormat int `json:"-"`
	TextColour config.Colour
	HoverText  string `json:"HoverText,omitempty"`
}

type TableWidgetFilter struct {
	Name          string
	FilterType    int
	ColumnNames   []string
	Values        []string `json:"Values,omitempty"`
	FullMatchOnly bool
}

func SaveTableWidget(widget *TableWidget) error {

	for i := range widget.Rows {
		for j := range widget.Rows[i].Values {
			valueText, hoverText := FormatColumn(widget.Rows[i].Values[j].Value, widget.Rows[i].Values[j].TextFormat)
			widget.Rows[i].Values[j].Value = valueText
			if hoverText != "" {
				widget.Rows[i].Values[j].HoverText = hoverText
			}
		}
	}

	if widget.RefreshMs == 0 {
		widget.RefreshMs = DefaultMsTimes[DefaultMsTable]
	}

	widget.ID = uuid.New().String()
	widget.SessionID = config.ServiceSettings.SESSIONID
	widget.Timestamp = time.Now()
	widget.Type = "table"
	jsonData, err := json.Marshal(widget)
	if err != nil {
		return err
	}

	return dashboard.SaveJsonData(widget.Options.Name, jsonData, widget.Options.Width, widget.RefreshMs)
}
