package validators

import (
	"go.uber.org/zap"

	"github.com/naveenm4d/bet-settlement-engine/internal/core/adapters"
	"github.com/naveenm4d/bet-settlement-engine/pkg/constants"
	"github.com/naveenm4d/bet-settlement-engine/pkg/entities"
)

var _ = adapters.ResultsValidator(&resultsValidator{})

type resultsValidator struct {
	cacheRepo adapters.CacheRepository
	logger    *zap.SugaredLogger
}

func NewResultsValidator(
	cacheRepo adapters.CacheRepository,
	logger *zap.SugaredLogger,
) adapters.ResultsValidator {
	validator := &resultsValidator{
		cacheRepo: cacheRepo,
		logger:    logger,
	}

	return validator
}

func (rv *resultsValidator) ValidateEvent(event *entities.Event) error {
	events := rv.cacheRepo.GetEvents()

	eventData, exists := events[event.EventID]
	if !exists {
		rv.logger.Errorf("ValidateEvent: Error! invalid event id | EventID - %s", event.EventID)

		return constants.ErrInvalidEventID
	}

	if eventData.Status != constants.Open {
		rv.logger.Errorf("ValidateEvent: Error! invalid event status | EventID - %s | EventStatus", event.EventID, eventData.Status)

		return constants.ErrEventAlreadyResulted
	}

	if event.ResultStatus != constants.Win &&
		event.ResultStatus != constants.Lose {
		rv.logger.Errorf("ValidateEvent: Error! invalid event result status | EventID - %s | EventStatus", event.EventID, eventData.ResultStatus)

		return constants.ErrInvalidEventResultStatus
	}

	return nil
}

func (rv *resultsValidator) ValidateBetForSettlement(bet *entities.Bet) error {
	if bet.Status != constants.Unresulted {
		rv.logger.Errorf("ValidateBetForSettlement: Error! bet already resulted | BetID - %s ", bet.BetID)

		return constants.ErrBetAlreadyResulted
	}

	return nil
}
