package services

import (
	"errors"

	"github.com/naveenm4d/bet-settlement-engine/pkg/constants"
	"github.com/naveenm4d/bet-settlement-engine/pkg/entities"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func (suite *ServiceTestSuite) Test_SettleBetsForEvent_Success() {
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()

	logger := zapLogger.Sugar()

	svc := NewResultsService(suite.resultsValidator, suite.cacheRepo, logger)

	suite.cacheRepo.
		On("GetBets").
		Return(map[string]entities.Bet{
			"1234": entities.Bet{UserID: "1234", EventID: "2345", Odds: 15, Status: constants.Unresulted},
		})

	suite.cacheRepo.
		On("GetEvents").
		Return(map[string]entities.Event{
			"2345": entities.Event{EventID: "2345", Status: constants.Open},
		})

	suite.resultsValidator.
		On("ValidateBetForSettlement", mock.Anything).
		Return(nil)

	suite.cacheRepo.
		On("UpdateEvents", mock.Anything).
		Return(nil)

	suite.cacheRepo.
		On("UpdateBets", mock.Anything).
		Return(nil)

	err := svc.SettleBetsForEvent("2345", constants.Win)

	suite.asserts.NoError(err)
}

func (suite *ServiceTestSuite) Test_SettleBetsForEvent_Error() {
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()

	logger := zapLogger.Sugar()
	testError := errors.New("test error")

	svc := NewResultsService(suite.resultsValidator, suite.cacheRepo, logger)

	suite.Run("UpdateBets Error", func() {
		suite.cacheRepo.
			On("GetBets").
			Return(map[string]entities.Bet{
				"1234": entities.Bet{UserID: "1234", EventID: "2345", Odds: 15, Status: constants.Unresulted},
			})

		suite.cacheRepo.
			On("GetEvents").
			Return(map[string]entities.Event{
				"2345": entities.Event{EventID: "2345", Status: constants.Open},
			})

		suite.resultsValidator.
			On("ValidateBetForSettlement", mock.Anything).
			Return(nil)

		suite.cacheRepo.
			On("UpdateEvents", mock.Anything).
			Return(nil)

		suite.cacheRepo.
			On("UpdateBets", mock.Anything).
			Return(testError)

		err := svc.SettleBetsForEvent("2345", constants.Win)

		suite.asserts.Error(err)
		suite.asserts.ErrorIs(err, testError)
	})

	suite.Run("UpdateEvents Error", func() {
		suite.cacheRepo.
			On("GetBets").
			Return(map[string]entities.Bet{
				"1234": entities.Bet{UserID: "1234", EventID: "2345", Odds: 15, Status: constants.Unresulted},
			})

		suite.cacheRepo.
			On("GetEvents").
			Return(map[string]entities.Event{
				"2345": entities.Event{EventID: "2345", Status: constants.Open},
			})

		suite.resultsValidator.
			On("ValidateBetForSettlement", mock.Anything).
			Return(nil)

		suite.cacheRepo.
			On("UpdateEvents", mock.Anything).
			Return(testError)

		err := svc.SettleBetsForEvent("2345", constants.Win)

		suite.asserts.Error(err)
		suite.asserts.ErrorIs(err, testError)
	})
}
