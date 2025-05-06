package services

import (
	"errors"

	"github.com/naveenm4d/bet-settlement-engine/pkg/entities"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func (suite *ServiceTestSuite) Test_DebitAccount_Success() {
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()

	logger := zapLogger.Sugar()

	svc := NewAccountsService(suite.cacheRepo, logger)

	suite.cacheRepo.
		On("GetAccounts").
		Return(map[string]entities.Account{
			"1234": entities.Account{UserID: "1234", Balance: int64(1000)},
		})

	suite.cacheRepo.
		On("UpdateAccounts", mock.Anything).
		Return(nil)

	err := svc.DebitAccount("1234", int64(100))

	suite.asserts.NoError(err)
}

func (suite *ServiceTestSuite) Test_DebitAccount_Error() {
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()

	logger := zapLogger.Sugar()
	testError := errors.New("test error")

	svc := NewAccountsService(suite.cacheRepo, logger)

	suite.cacheRepo.
		On("GetAccounts").
		Return(map[string]entities.Account{
			"1234": entities.Account{UserID: "1234", Balance: int64(1000)},
		})

	suite.cacheRepo.
		On("UpdateAccounts", mock.Anything).
		Return(testError)

	err := svc.DebitAccount("1234", int64(100))

	suite.asserts.Error(err)
	suite.asserts.ErrorIs(err, testError)
}
