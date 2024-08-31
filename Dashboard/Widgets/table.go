package widgets

import (
	"encoding/json"
	"strconv"
	"time"

	config "github.com/dabi-ngin/discgo-bot/Config"
	dashboard "github.com/dabi-ngin/discgo-bot/Dashboard"
)

type TableWidget struct {
	Options TableWidgetOptions
	Columns []TableWidgetColumn
	Rows    []TableWidgetRow
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
	TextFormat int `json:"TextFormat,omitempty"`
	TextColour config.Colour
}

func SaveTableWidget(widget *TableWidget) error {

	for i := range widget.Rows {
		for j := range widget.Rows[i].Values {
			widget.Rows[i].Values[j].Value = formatColumn(widget.Rows[i].Values[j].Value, widget.Rows[i].Values[j].TextFormat)
		}
	}

	jsonData, err := json.Marshal(widget)
	if err != nil {
		return err
	}

	return dashboard.SaveJsonData(widget.Options.Name, jsonData)
}

func formatColumn(value interface{}, format int) string {
	switch v := value.(type) {
	case string:
		return formatStringColumn(v, format)
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case time.Time:
		return formatTimeColumn(v, format)
	case time.Duration:
		return formatDurationColumn(v, format)
	default:
		return "???" // Unknown/unsupported format
	}
}

func formatStringColumn(value string, format int) string {
	switch format {
	case TextFormatString_AbbreviateFromStart:
		if len(value) > AbbrevRemainingChars {
			return "..." + value[len(value)-AbbrevRemainingChars:]
		}
	case TextFormatString_AbbreviateToEnd:
		if len(value) > AbbrevRemainingChars {
			return value[:AbbrevRemainingChars] + "..."
		}
	}
	return value
}

func formatTimeColumn(value time.Time, format int) string {
	switch format {
	case TextFormatTime_DateAndTime:
		return value.Format("2006-01-02 15:04:05")
	case TextFormatTime_DateOnly:
		return value.Format("2006-01-02")
	case TextFormatTime_TimeOnly:
		return value.Format("15:04:05")
	default:
		return value.String()
	}
}

func formatDurationColumn(value time.Duration, format int) string {
	switch format {
	case TextFormatDuration_WithMs:
		return value.String()
	default:
		return value.String()
	}
}
