package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/naveenm4d/bet-settlement-engine/internal/core/adapters"
	"github.com/naveenm4d/bet-settlement-engine/pkg/constants"
	"github.com/naveenm4d/bet-settlement-engine/pkg/entities"
)

var _ = adapters.BetService(&betService{})

type betService struct {
	cacheRepository adapters.CacheRepository
	accountsService adapters.AccountsService
	logger          *zap.SugaredLogger
}

func NewBetService(
	cacheRepository adapters.CacheRepository,
	accountsService adapters.AccountsService,
	logger *zap.SugaredLogger,
) adapters.BetService {
	svc := &betService{
		cacheRepository: cacheRepository,
		accountsService: accountsService,
		logger:          logger,
	}

	return svc
}

func (bs *betService) PlaceBet(
	ctx context.Context,
	bet *entities.Bet,
) (*entities.Bet, error) {
	bs.logger.Info("Starting placing bet...")

	bet.BetID = uuid.NewString()
	bet.Status = constants.Unresulted
	bet.PlacedAt = time.Now()

	existingBets := bs.cacheRepository.GetBets()
	if len(existingBets) < 1 {
		existingBets = make(map[string]entities.Bet)
	}

	existingBets[bet.BetID] = *bet

	if err := bs.accountsService.DebitAccount(bet.UserID, bet.Amount); err != nil {
		return nil, err
	}

	if err := bs.cacheRepository.UpdateBets(existingBets); err != nil {
		if err := bs.accountsService.RefundAccount(bet.UserID, bet.Amount); err != nil {
			return nil, err
		}

		return nil, err
	}

	bs.logger.Info("Bet placed successfully! BetID - %s")

	return bet, nil
}
