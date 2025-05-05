package constants

type BetStatus string
type EventStatus string
type EventResultStatus string

var (
	Unresulted   BetStatus = "UNRESULTED"
	ResultedWin  BetStatus = "RESULTED_WIN"
	ResultedLoss BetStatus = "RESULTED_LOSS"
)

var (
	Open     EventStatus = "OPEN"
	Resulted EventStatus = "RESULTED"
)

var (
	Win  EventResultStatus = "WIN"
	Lose EventResultStatus = "LOSS"
)
