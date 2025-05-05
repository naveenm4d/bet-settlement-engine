package services

import (
	"go.uber.org/zap"

	"github.com/naveenm4d/bet-settlement-engine/internal/core/adapters"
	"github.com/naveenm4d/bet-settlement-engine/pkg/constants"
	"github.com/naveenm4d/bet-settlement-engine/pkg/entities"
)

var _ = adapters.AccountsService(&accountsService{})

type accountsService struct {
	cacheRepository adapters.CacheRepository
	logger          *zap.SugaredLogger
}

func NewAccountsService(
	cacheRepository adapters.CacheRepository,
	logger *zap.SugaredLogger,
) adapters.AccountsService {
	svc := &accountsService{
		cacheRepository: cacheRepository,
		logger:          logger,
	}

	return svc
}

func (ac *accountsService) DebitAccount(userID string, betAmount int64) error {
	accounts := ac.cacheRepository.GetAccounts()

	userAccount := accounts[userID]
	userAccount.Balance -= betAmount

	accounts[userID] = userAccount

	return ac.cacheRepository.UpdateAccounts(accounts)
}

func (ac *accountsService) RefundAccount(userID string, betAmount int64) error {
	accounts := ac.cacheRepository.GetAccounts()

	userAccount := accounts[userID]
	userAccount.Balance += betAmount

	accounts[userID] = userAccount

	return ac.cacheRepository.UpdateAccounts(accounts)
}

func (ac *accountsService) GetAccount(userID string) (*entities.Account, error) {
	accounts := ac.cacheRepository.GetAccounts()

	userAccount, exists := accounts[userID]
	if !exists {
		return nil, constants.ErrInvalidUserID
	}

	return &userAccount, nil
}
