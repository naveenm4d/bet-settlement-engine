package entities

import "github.com/naveenm4d/bet-settlement-engine/pkg/constants"

type Event struct {
	EventID      string                      `json:"event_id"`
	Odds         int64                       `json:"odds"`
	Status       constants.EventStatus       `json:"status"`
	ResultStatus constants.EventResultStatus `json:"result_status"`
}
