package adapters

import "github.com/naveenm4d/bet-settlement-engine/pkg/entities"

type AccountsService interface {
	DebitAccount(userID string, betAmount int64) error
	RefundAccount(userID string, betAmount int64) error
	GetAccount(userID string) (*entities.Account, error)
}
