package adapters

import "github.com/naveenm4d/bet-settlement-engine/pkg/entities"

type ResultsValidator interface {
	ValidateEvent(event *entities.Event) error
	ValidateBetForSettlement(bet *entities.Bet) error
}
