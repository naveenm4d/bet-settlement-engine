package adapters

import "github.com/naveenm4d/bet-settlement-engine/pkg/entities"

type CacheRepository interface {
	GetBets() map[string]entities.Bet
	GetEvents() map[string]entities.Event
	GetAccounts() map[string]entities.Account

	UpdateBets(bets map[string]entities.Bet) error
	UpdateAccounts(accounts map[string]entities.Account) error
	UpdateEvents(events map[string]entities.Event) error
}
