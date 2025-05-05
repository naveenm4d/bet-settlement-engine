package adapters

import "github.com/naveenm4d/bet-settlement-engine/pkg/constants"

type ResultsService interface {
	SettleBetsForEvent(eventID string, resultStatus constants.EventResultStatus) error
}
