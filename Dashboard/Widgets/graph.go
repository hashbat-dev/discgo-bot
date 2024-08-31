package widgets

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	dashboard "github.com/dabi-ngin/discgo-bot/Dashboard"
)

type GraphWidget struct {
	Name    string
	Data    interface{}
	Options GraphWidgetOptions
}

type GraphWidgetOptions struct {
	Width      string
	LineColour string
	MinValue   int
	MaxValue   int
}

// Writes JSON data for a GraphWidget object
func SaveGraphWidget(widget GraphWidget) error {
	dataValue := reflect.ValueOf(widget.Data)

	// Check that Data is a slice and []int or []float64
	if dataValue.Kind() != reflect.Slice {
		return fmt.Errorf("Data should be a slice")
	}

	elemKind := dataValue.Type().Elem().Kind()
	if elemKind != reflect.Int && elemKind != reflect.Float64 {
		return fmt.Errorf("Data slice not an accepted format for a graph")
	}

	dataSlice := make([]interface{}, dataValue.Len())
	for i := 0; i < dataValue.Len(); i++ {
		dataSlice[i] = dataValue.Index(i).Interface()
	}

	saveJson := map[string]interface{}{
		"Data":      dataSlice,
		"Options":   widget.Options,
		"Timestamp": time.Now(),
	}

	jsonData, err := json.MarshalIndent(saveJson, "", "  ")
	if err != nil {
		return err
	}

	dashboard.SaveJsonData(widget.Name, jsonData)
	return nil
}
