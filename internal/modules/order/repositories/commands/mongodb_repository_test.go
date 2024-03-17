package commands_test

import (
	"context"
	"order-service/internal/modules/order"
	"order-service/internal/modules/order/models/request"
	mongoRC "order-service/internal/modules/order/repositories/commands"
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
	repository  order.MongodbRepositoryCommand
	ctx         context.Context
}

func (suite *CommandTestSuite) SetupTest() {
	suite.mockMongodb = new(mocks.Collections)
	suite.mockLogger = &mocklog.Logger{}
	suite.repository = mongoRC.NewCommandMongodbRepository(
		suite.mockMongodb,
		suite.mockLogger,
	)
	suite.ctx = context.Background()
}

func TestCommandTestSuite(t *testing.T) {
	suite.Run(t, new(CommandTestSuite))
}

func (suite *CommandTestSuite) TestUpdateBankTicket() {
	payload := request.UpdateBankTicketReq{
		CountryCode: "code",
		TicketType:  "type",
	}

	// Mock UpsertOne
	expectedResult := make(chan helpers.Result)
	suite.mockMongodb.On("FindOneAndUpdate", mock.Anything, mock.Anything, mock.Anything).Return((<-chan helpers.Result)(expectedResult))

	// Act
	result := suite.repository.UpdateBankTicket(suite.ctx, payload)
	// Asset
	assert.NotNil(suite.T(), result, "Expected a result")

	// Simulate receiving a result from the channel
	go func() {
		expectedResult <- helpers.Result{Data: "result not nil", Error: nil}
		close(expectedResult)
	}()

	// Wait for the goroutine to complete
	<-result

	// Assert UpsertOne
	suite.mockMongodb.AssertCalled(suite.T(), "FindOneAndUpdate", mock.Anything, mock.Anything, mock.Anything)
}
