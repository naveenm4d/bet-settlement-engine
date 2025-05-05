package repositories

import (
	"sync"

	"go.uber.org/zap"

	"github.com/naveenm4d/bet-settlement-engine/internal/core/adapters"
	"github.com/naveenm4d/bet-settlement-engine/pkg/cache"
	"github.com/naveenm4d/bet-settlement-engine/pkg/constants"
	"github.com/naveenm4d/bet-settlement-engine/pkg/entities"
)

var _ = adapters.CacheRepository(&cacheRepo{})

type cacheRepo struct {
	cache  cache.Cache
	mu     *sync.Mutex
	logger *zap.SugaredLogger
}

func NewCacheRepository(
	cache cache.Cache,
	logger *zap.SugaredLogger,
) adapters.CacheRepository {
	repo := &cacheRepo{
		cache:  cache,
		mu:     &sync.Mutex{},
		logger: logger,
	}

	return repo
}

func (c *cacheRepo) GetBets() map[string]entities.Bet {
	var bets map[string]entities.Bet

	val, err := c.cache.Get(constants.CacheKeyBets)
	if err == nil {
		betsData, ok := val.(map[string]entities.Bet)
		if ok {
			bets = betsData
		}
	}

	c.logger.Infof("GetBets(): bets - %+v", bets)

	return bets
}

func (c *cacheRepo) GetEvents() map[string]entities.Event {
	var events map[string]entities.Event

	val, err := c.cache.Get(constants.CacheKeyEvents)
	if err == nil {
		eventsData, ok := val.(map[string]entities.Event)
		if ok {
			events = eventsData
		}
	}

	c.logger.Infof("GetEvents(): events - %+v", events)

	return events
}

func (c *cacheRepo) GetAccounts() map[string]entities.Account {
	var accounts map[string]entities.Account

	val, err := c.cache.Get(constants.CacheKeyAccounts)
	if err == nil {
		accountsData, ok := val.(map[string]entities.Account)
		if ok {
			accounts = accountsData
		}
	}

	c.logger.Infof("GetAccounts(): accounts - %+v", accounts)

	return accounts
}

func (c *cacheRepo) UpdateBets(bets map[string]entities.Bet) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.cache.Set(constants.CacheKeyBets, bets)
}

func (c *cacheRepo) UpdateAccounts(accounts map[string]entities.Account) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.cache.Set(constants.CacheKeyAccounts, accounts)
}

func (c *cacheRepo) UpdateEvents(events map[string]entities.Event) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.cache.Set(constants.CacheKeyEvents, events)
}
