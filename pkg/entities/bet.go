package entities

import (
	"time"

	"github.com/naveenm4d/bet-settlement-engine/pkg/constants"
)

type Bet struct {
	BetID   string              `json:"bet_id"`
	UserID  string              `json:"user_id"`
	EventID string              `json:"event_id"`
	Odds    int64               `json:"odds"`
	Amount  int64               `json:"amount"`
	Status  constants.BetStatus `json:"status"`
	Payout  *int64              `json:"payout"`

	PlacedAt   time.Time  `json:"placed_at"`
	ResultedAt *time.Time `json:"resulted_at"`
}
