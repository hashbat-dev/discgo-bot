package widgets

import (
	"strconv"
	"time"
)

const (
	WidthQuarter       = "25%"
	WidthHalf          = "50%"
	WidthThreeQuarters = "75%"
	WidthFull          = "100%"
)

const (
	TextFormatString = iota
	TextFormatString_AbbreviateFromStart
	TextFormatString_AbbreviateToEnd
	TextFormatTime_DateAndTime
	TextFormatTime_DateOnly
	TextFormatTime_TimeOnly
	TextFormatTime_TimeWithMs
	TextFormatDuration_WithMs
)

var (
	AbbrevRemainingChars = 5
	DateTimeLayout       = "2006-01-02T15:04:05.000"
)

const (
	DefaultMsGraph = iota
	DefaultMsInfo
	DefaultMsTable
)

var DefaultMsTimes map[int]int = map[int]int{
	DefaultMsGraph: 2000,
	DefaultMsInfo:  10000,
	DefaultMsTable: 1000,
}

// Returns (Text, HoverText)
func FormatColumn(value interface{}, format int) (string, string) {
	switch v := value.(type) {
	case string:
		return formatStringColumn(v, format)
	case int:
		return strconv.Itoa(v), ""
	case int64:
		return strconv.FormatInt(v, 10), ""
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), ""
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32), ""
	case time.Time:
		return formatTimeColumn(v, format), ""
	case time.Duration:
		return formatDurationColumn(v, format), ""
	default:
		return "???", "" // Unknown/unsupported format
	}
}

func formatStringColumn(value string, format int) (string, string) {
	switch format {
	case TextFormatString_AbbreviateFromStart:
		if len(value) > AbbrevRemainingChars {
			return "..." + value[len(value)-AbbrevRemainingChars:], value
		}
	case TextFormatString_AbbreviateToEnd:
		if len(value) > AbbrevRemainingChars {
			return value[:AbbrevRemainingChars] + "...", value
		}
	}
	return value, ""
}

func formatTimeColumn(value time.Time, format int) string {
	switch format {
	case TextFormatTime_DateAndTime:
		return value.Format("2006-01-02 15:04:05")
	case TextFormatTime_DateOnly:
		return value.Format("2006-01-02")
	case TextFormatTime_TimeOnly:
		return value.Format("15:04:05")
	case TextFormatTime_TimeWithMs:
		return value.Format("15:04:05.000")
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
