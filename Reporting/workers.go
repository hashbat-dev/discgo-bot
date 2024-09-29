package reporting

import (
	"fmt"
	"time"

	widgets "github.com/hashbat-dev/discgo-bot/Dashboard/Widgets"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
)

var workerData map[int]WorkerData = make(map[int]WorkerData)

type WorkerData struct {
	Name              string
	MaxQueued         int
	Queued            int
	Processing        int
	MaxProcessing     int
	LastQueued        time.Time
	LastProcessStart  time.Time
	LastProcessFinish time.Time
}

var workerColumns []widgets.TableWidgetColumn = []widgets.TableWidgetColumn{
	{Name: "Channel"},
	{Name: "Queued"},
	{Name: "Processing"},
	{Name: "Last Queued"},
	{Name: "Last Started"},
}

var workerOptions widgets.TableWidgetOptions = widgets.TableWidgetOptions{
	Name:  "Worker Channels",
	Width: widgets.WidthQuarter,
}

func CreateWorkerChannel(channelId int, channelName string, maxQueued int, maxProcessing int) {
	workerData[channelId] = WorkerData{
		Name:          channelName,
		MaxQueued:     maxQueued,
		MaxProcessing: maxProcessing,
	}
}

func WorkerSaveWidget(channelId int) {
	var tableRows []widgets.TableWidgetRow
	for _, data := range workerData {
		tableRows = append(tableRows, widgets.TableWidgetRow{
			Values: []widgets.TableWidgetRowValue{
				{Value: data.Name},
				{Value: fmt.Sprintf("%v / %v", data.Queued, data.MaxQueued)},
				{Value: fmt.Sprintf("%v / %v", data.Processing, data.MaxProcessing)},
				{Value: data.LastQueued, TextFormat: widgets.TextFormatTime_TimeWithMs},
				{Value: data.LastProcessStart, TextFormat: widgets.TextFormatTime_TimeWithMs},
			},
		})
	}

	err := widgets.SaveTableWidget(&widgets.TableWidget{
		Options:   workerOptions,
		Columns:   workerColumns,
		Rows:      tableRows,
		RefreshMs: 500,
	})
	if err != nil {
		logger.Error("REPORTING", err)
	}
}

func WorkerQueued(channelId int) {
	worker := workerData[channelId]
	worker.Queued++
	worker.LastQueued = time.Now()
	workerData[channelId] = worker
	WorkerSaveWidget(channelId)
}

func WorkerProcessingStart(channelId int) {
	worker := workerData[channelId]
	worker.Queued--
	if worker.Queued < 0 {
		worker.Queued = 0
	}
	worker.Processing++
	worker.LastProcessStart = time.Now()
	workerData[channelId] = worker
	WorkerSaveWidget(channelId)
}

func WorkerProcessingFinish(channelId int) {
	worker := workerData[channelId]
	worker.Processing--
	if worker.Processing < 0 {
		worker.Processing = 0
	}
	worker.LastProcessFinish = time.Now()
	workerData[channelId] = worker
	WorkerSaveWidget(channelId)
}
