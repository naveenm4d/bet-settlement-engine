package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/allegro/bigcache"
	"go.uber.org/zap"

	"github.com/naveenm4d/bet-settlement-engine/internal/app/repositories"
	"github.com/naveenm4d/bet-settlement-engine/internal/app/services"
	"github.com/naveenm4d/bet-settlement-engine/internal/app/validators"
	"github.com/naveenm4d/bet-settlement-engine/internal/config"
	"github.com/naveenm4d/bet-settlement-engine/internal/handlers"

	"github.com/naveenm4d/bet-settlement-engine/pkg/cache"
	"github.com/naveenm4d/bet-settlement-engine/pkg/constants"
	"github.com/naveenm4d/bet-settlement-engine/pkg/entities"
)

func initCache(ctx context.Context) (cache.Cache, error) {
	cacheConfig := bigcache.Config{
		LifeWindow:       3600 * time.Second,
		CleanWindow:      3600 * time.Second,
		HardMaxCacheSize: 4,
		Shards:           2,
	}

	cache, err := cache.NewCache(ctx, cacheConfig)
	if err != nil {
		return nil, err
	}

	return cache, nil
}

func main() {
	appName := "bet-settlement-engine"

	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()

	logger := zapLogger.Sugar()

	logger.Infof("Starting app %s...", appName)

	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-c
		logger.Infof("Got %s signal. Cancelling", sig)

		cancel()
	}()

	memCache, errM := initCache(ctx)
	if errM != nil {
		logger.Errorf("could not create tenant cache: error - %s", errM)
	}

	//////////
	logger.Info("populating cache...")

	memCache.Set(constants.CacheKeyEvents, map[string]entities.Event{
		"1234": {
			EventID: "1234",
			Odds:    15,
			Status:  constants.Open,
		},
	})

	memCache.Set(constants.CacheKeyAccounts, map[string]entities.Account{
		"2345": {
			UserID:  "2345",
			Balance: 100000,
		},
	})

	logger.Info("populating cache completed")
	//////////

	cacheRespository := repositories.NewCacheRepository(memCache, logger)

	accountsService := services.NewAccountsService(cacheRespository, logger)

	betService := services.NewBetService(cacheRespository, accountsService, logger)

	betValidator := validators.NewBetsValidator(cacheRespository, logger)

	resultsValidator := validators.NewResultsValidator(cacheRespository, logger)

	resultsService := services.NewResultsService(resultsValidator, cacheRespository, logger)

	handler := handlers.NewHandler(betValidator, resultsValidator, betService, resultsService, accountsService, logger)

	router := handlers.NewRouter(handler)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%v", *config.Config.HTTPPort),
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Infof("Initializing server on port : %v", *config.Config.HTTPPort)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("ListenAndServe error: %s", err)
		}

		cancel()
	}()

	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Errorf("Shutdown error: %s", err)
	}

	logger.Infof("Exit!")
}
