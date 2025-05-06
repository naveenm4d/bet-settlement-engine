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
	logger.Info("populating test cache data...")

	memCache.Set(constants.CacheKeyEvents, map[string]entities.Event{
		"91bb5494-ca38-42e9-a20d-cfa9a07900e6": {
			EventID: "91bb5494-ca38-42e9-a20d-cfa9a07900e6",
			Odds:    150,
			Status:  constants.Open,
		},
		"5b1043d8-151f-441d-868b-d8227961d54f": {
			EventID: "5b1043d8-151f-441d-868b-d8227961d54f",
			Odds:    160,
			Status:  constants.Open,
		},
		"38c86515-3eb8-458b-9f11-5cb677dbbe6f": {
			EventID: "38c86515-3eb8-458b-9f11-5cb677dbbe6f",
			Odds:    170,
			Status:  constants.Resulted,
		},
	})

	memCache.Set(constants.CacheKeyAccounts, map[string]entities.Account{
		"6209cf4b-92f0-4c0e-8f8e-e4a518cd2430": {
			UserID:  "6209cf4b-92f0-4c0e-8f8e-e4a518cd2430",
			Balance: 100000,
		},
		"774ab462-b536-458b-ab46-30903f45001f": {
			UserID:  "774ab462-b536-458b-ab46-30903f45001f",
			Balance: 0,
		},
	})

	logger.Info("populating test cache data completed")
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
