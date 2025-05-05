package validators

import (
	"context"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/naveenm4d/bet-settlement-engine/internal/core/adapters"
	"github.com/naveenm4d/bet-settlement-engine/pkg/constants"
	"github.com/naveenm4d/bet-settlement-engine/pkg/entities"
)

var _ = adapters.BetsValidator(&betsValidator{})

type betsValidator struct {
	cacheRepo adapters.CacheRepository
	logger    *zap.SugaredLogger
}

func NewBetsValidator(
	cacheRepo adapters.CacheRepository,
	logger *zap.SugaredLogger,
) adapters.BetsValidator {
	validator := &betsValidator{
		cacheRepo: cacheRepo,
		logger:    logger,
	}

	return validator
}

func (v *betsValidator) ValidateBet(ctx context.Context, bet *entities.Bet) error {
	eGroup, _ := errgroup.WithContext(ctx)

	eGroup.Go(func() error {
		return v.validateEvent(bet.EventID, bet.Odds)
	})

	eGroup.Go(func() error {
		return v.validateBetAmount(bet.UserID, bet.Amount)
	})

	return eGroup.Wait()
}

func (v *betsValidator) validateEvent(eventID string, betOdds int64) error {
	events := v.cacheRepo.GetEvents()

	betEvent, exists := events[eventID]
	if !exists {
		v.logger.Errorf("ValidateBet: Error! invalid event id | EventID - %s", eventID)

		return constants.ErrInvalidEventID
	}

	if betEvent.Status != constants.Open {
		v.logger.Errorf("ValidateBet: Error! invalid event status | EventID - %s | EventStatus", eventID, betEvent.Status)

		return constants.ErrInvalidEventStatus
	}

	if betEvent.Odds != betOdds {
		v.logger.Errorf("ValidateBet: Error! invalid bet odd | EventID - %s | BetOdds - %v | EventOdds - %v", eventID, betOdds, betEvent.Odds)

		return constants.ErrInvalidOdds
	}

	return nil
}

func (v *betsValidator) validateBetAmount(userID string, amount int64) error {
	accounts := v.cacheRepo.GetAccounts()

	userAccount, exists := accounts[userID]
	if !exists {
		v.logger.Errorf("ValidateBet: Error! invalid user id | UserID - %s", userID)

		return constants.ErrInvalidUserID
	}

	if amount <= 0 {
		v.logger.Errorf("ValidateBet: Error! invalid bet amount | UserID - %s | Amount - %v", userID, amount)

		return constants.ErrInvalidAmount
	}

	if userAccount.Balance < amount {
		v.logger.Errorf("ValidateBet: Error! insuffiecient balance | UserID - %s | Balance - %v", userID, userAccount.Balance)

		return constants.ErrInsuffiecientBalance
	}

	return nil
}
