package widgets

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
	TextFormatDuration_WithMs
)

var (
	AbbrevRemainingChars = 5
	DateTimeLayout       = "2006-01-02T15:04:05.000"
)
