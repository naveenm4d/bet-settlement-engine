package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	mocks "github.com/naveenm4d/bet-settlement-engine/mocks/adapters"
)

type ServiceTestSuite struct {
	suite.Suite
	asserts          *assert.Assertions
	cacheRepo        *mocks.MockCacheRepository
	resultsValidator *mocks.MockResultsValidator
	accountsService  *mocks.MockAccountsService
}

func (suite *ServiceTestSuite) SetupTest() {
	suite.asserts = assert.New(suite.T())
	suite.cacheRepo = mocks.NewMockCacheRepository(suite.T())
	suite.resultsValidator = mocks.NewMockResultsValidator(suite.T())
	suite.accountsService = mocks.NewMockAccountsService(suite.T())

}
func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
