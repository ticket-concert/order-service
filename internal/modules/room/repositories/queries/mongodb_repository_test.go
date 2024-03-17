package queries_test

import (
	"context"
	"order-service/internal/modules/room"
	mongoRQ "order-service/internal/modules/room/repositories/queries"
	"order-service/internal/pkg/helpers"
	mocks "order-service/mocks/pkg/databases/mongodb"
	mocklog "order-service/mocks/pkg/log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CommandTestSuite struct {
	suite.Suite
	mockMongodb *mocks.Collections
	mockLogger  *mocklog.Logger
	repository  room.MongodbRepositoryQuery
	ctx         context.Context
}

func (suite *CommandTestSuite) SetupTest() {
	suite.mockMongodb = new(mocks.Collections)
	suite.mockLogger = &mocklog.Logger{}
	suite.repository = mongoRQ.NewQueryMongodbRepository(
		suite.mockMongodb,
		suite.mockLogger,
	)
	suite.ctx = context.Background()
}

func TestCommandTestSuite(t *testing.T) {
	suite.Run(t, new(CommandTestSuite))
}

func (suite *CommandTestSuite) TestFindOneLastQueue() {

	// Mock FindOne
	expectedResult := make(chan helpers.Result)
	suite.mockMongodb.On("FindOne", mock.Anything, mock.Anything).Return((<-chan helpers.Result)(expectedResult))

	// Act
	result := suite.repository.FindOneLastQueue(suite.ctx, mock.Anything)
	// Asset
	assert.NotNil(suite.T(), result, "Expected a result")

	// Simulate receiving a result from the channel
	go func() {
		expectedResult <- helpers.Result{Data: "result not nil", Error: nil}
		close(expectedResult)
	}()

	// Wait for the goroutine to complete
	<-result

	// Assert FindOne
	suite.mockMongodb.AssertCalled(suite.T(), "FindOne", mock.Anything, mock.Anything)
}

func (suite *CommandTestSuite) TestFindOneQueueByUserId() {

	// Mock FindOne
	expectedResult := make(chan helpers.Result)
	suite.mockMongodb.On("FindOne", mock.Anything, mock.Anything).Return((<-chan helpers.Result)(expectedResult))

	// Act
	result := suite.repository.FindOneQueueByUserId(suite.ctx, mock.Anything, mock.Anything)
	// Asset
	assert.NotNil(suite.T(), result, "Expected a result")

	// Simulate receiving a result from the channel
	go func() {
		expectedResult <- helpers.Result{Data: "result not nil", Error: nil}
		close(expectedResult)
	}()

	// Wait for the goroutine to complete
	<-result

	// Assert FindOne
	suite.mockMongodb.AssertCalled(suite.T(), "FindOne", mock.Anything, mock.Anything)
}
