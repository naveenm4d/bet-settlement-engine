package adapters

import (
	"context"

	"github.com/naveenm4d/bet-settlement-engine/pkg/entities"
)

type BetsValidator interface {
	ValidateBet(ctx context.Context, bet *entities.Bet) error
}
