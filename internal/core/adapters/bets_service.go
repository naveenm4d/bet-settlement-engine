package adapters

import (
	"context"

	"github.com/naveenm4d/bet-settlement-engine/pkg/entities"
)

type BetService interface {
	PlaceBet(ctx context.Context, bet *entities.Bet) (*entities.Bet, error)
}
