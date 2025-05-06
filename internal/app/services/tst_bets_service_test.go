package services

import (
	"context"
	"errors"

	"github.com/naveenm4d/bet-settlement-engine/pkg/constants"
	"github.com/naveenm4d/bet-settlement-engine/pkg/entities"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func (suite *ServiceTestSuite) Test_PlaceBet_Success() {
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()

	logger := zapLogger.Sugar()
	ctx := context.TODO()

	bet := entities.Bet{UserID: "1234", EventID: "2345", Odds: 15, Status: constants.Unresulted}

	svc := NewBetService(suite.cacheRepo, suite.accountsService, logger)

	suite.cacheRepo.
		On("GetBets").
		Return(map[string]entities.Bet{})

	suite.accountsService.
		On("DebitAccount", mock.Anything, mock.Anything).
		Return(nil)

	suite.cacheRepo.
		On("UpdateBets", mock.Anything).
		Return(nil)

	resp, err := svc.PlaceBet(ctx, &bet)

	suite.asserts.NoError(err)
	suite.asserts.NotNil(resp)
	suite.asserts.Equal(bet.UserID, "1234")
	suite.asserts.Equal(bet.EventID, "2345")
	suite.asserts.Equal(bet.Status, constants.Unresulted)
}

func (suite *ServiceTestSuite) Test_PlaceBet_Error() {
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()

	logger := zapLogger.Sugar()
	ctx := context.TODO()
	testError := errors.New("test error")

	bet := entities.Bet{UserID: "1234", EventID: "2345", Odds: 15, Status: constants.Unresulted}

	svc := NewBetService(suite.cacheRepo, suite.accountsService, logger)

	suite.Run("UpdateBets error", func() {
		suite.cacheRepo.
			On("GetBets").
			Return(map[string]entities.Bet{})

		suite.accountsService.
			On("DebitAccount", mock.Anything, mock.Anything).
			Return(nil)

		suite.cacheRepo.
			On("UpdateBets", mock.Anything).
			Return(testError)

		suite.accountsService.
			On("RefundAccount", mock.Anything, mock.Anything).
			Return(nil)

		resp, err := svc.PlaceBet(ctx, &bet)

		suite.asserts.Error(err)
		suite.asserts.ErrorIs(err, testError)
		suite.asserts.Nil(resp)
	})

	suite.Run("DebitAccount error", func() {
		suite.cacheRepo.
			On("GetBets").
			Return(map[string]entities.Bet{})

		suite.accountsService.
			On("DebitAccount", mock.Anything, mock.Anything).
			Return(testError)

		resp, err := svc.PlaceBet(ctx, &bet)

		suite.asserts.Error(err)
		suite.asserts.ErrorIs(err, testError)
		suite.asserts.Nil(resp)
	})
}
