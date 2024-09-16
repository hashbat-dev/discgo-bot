package widgets

import (
	"encoding/json"
	"time"

	config "github.com/dabi-ngin/discgo-bot/Config"
	dashboard "github.com/dabi-ngin/discgo-bot/Dashboard"
	"github.com/google/uuid"
)

type InfoWidget struct {
	ID        string
	SessionID string
	Name      string
	Timestamp time.Time
	Type      string
	Items     []InfoWidgetItem
	RefreshMs int
}

type InfoWidgetItem struct {
	Name        string
	Description string
	TextFormat  int
	Value       interface{}
}

func SaveInfoWidget(widget *InfoWidget) error {

	for i := range widget.Items {
		valueText, _ := FormatColumn(widget.Items[i].Value, widget.Items[i].TextFormat)
		widget.Items[i].Value = valueText
	}

	if widget.RefreshMs == 0 {
		widget.RefreshMs = DefaultMsTimes[DefaultMsInfo]
	}

	widget.ID = uuid.New().String()
	widget.SessionID = config.ServiceSettings.SESSIONID
	widget.Timestamp = time.Now()
	widget.Type = "info"
	jsonData, err := json.Marshal(widget)
	if err != nil {
		return err
	}

	return dashboard.SaveJsonData(widget.Name, jsonData, "0%", widget.RefreshMs)
}
