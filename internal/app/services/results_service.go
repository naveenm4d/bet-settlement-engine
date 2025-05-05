package services

import (
	"time"

	"go.uber.org/zap"

	"github.com/naveenm4d/bet-settlement-engine/internal/core/adapters"
	"github.com/naveenm4d/bet-settlement-engine/pkg/constants"
)

var _ = adapters.ResultsService(&resultsService{})

type resultsService struct {
	validator       adapters.ResultsValidator
	cacheRepository adapters.CacheRepository
	logger          *zap.SugaredLogger
}

func NewResultsService(
	validator adapters.ResultsValidator,
	cacheRepository adapters.CacheRepository,
	logger *zap.SugaredLogger,
) adapters.ResultsService {
	svc := &resultsService{
		validator:       validator,
		cacheRepository: cacheRepository,
		logger:          logger,
	}

	return svc
}

func (rs *resultsService) SettleBetsForEvent(
	eventID string,
	resultStatus constants.EventResultStatus,
) error {
	bets := rs.cacheRepository.GetBets()

	for betID, bet := range bets {
		if bet.EventID != eventID {
			continue
		}

		if err := rs.validator.ValidateBetForSettlement(&bet); err != nil {
			rs.logger.Errorf("SettleBetsForEvent failed! BetID - %s | Err - %s", betID, err)

			continue
		}

		betStatus := constants.ResultedWin
		if resultStatus == constants.Lose {
			betStatus = constants.ResultedLoss
		}

		betPayout := int64(0)
		if betStatus == constants.ResultedWin {
			betPayout = bet.Odds * bet.Amount
		}

		resultedAt := time.Now()

		bet.Status = betStatus
		bet.Payout = &betPayout
		bet.ResultedAt = &resultedAt

		bets[betID] = bet
	}

	events := rs.cacheRepository.GetEvents()

	eventData := events[eventID]
	eventData.Status = constants.Resulted
	eventData.ResultStatus = resultStatus

	events[eventID] = eventData

	if err := rs.cacheRepository.UpdateEvents(events); err != nil {
		return err
	}

	return rs.cacheRepository.UpdateBets(bets)
}
