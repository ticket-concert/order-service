package usecases_test

import (
	"context"
	"order-service/internal/modules/order"
	"testing"

	"order-service/internal/modules/order/models/entity"
	"order-service/internal/modules/order/models/request"
	uc "order-service/internal/modules/order/usecases"
	"order-service/internal/pkg/constants"
	"order-service/internal/pkg/errors"
	"order-service/internal/pkg/helpers"
	mockcert "order-service/mocks/modules/order"
	mocklog "order-service/mocks/pkg/log"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type QueryUsecaseTestSuite struct {
	suite.Suite
	mockOrderRepositoryQuery *mockcert.MongodbRepositoryQuery
	mockLogger               *mocklog.Logger
	usecase                  order.UsecaseQuery
	ctx                      context.Context
}

func (suite *QueryUsecaseTestSuite) SetupTest() {
	suite.mockOrderRepositoryQuery = &mockcert.MongodbRepositoryQuery{}
	suite.mockLogger = &mocklog.Logger{}
	suite.ctx = context.Background()
	suite.usecase = uc.NewQueryUsecase(
		suite.mockOrderRepositoryQuery,
		suite.mockLogger,
	)
}
func TestQueryUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(QueryUsecaseTestSuite))
}

func (suite *QueryUsecaseTestSuite) TestFindOrderList() {
	payload := request.OrderList{
		Page:   1,
		Size:   1,
		UserId: "id",
	}
	mockOrderByUser := helpers.Result{
		Data: &[]entity.Order{
			{
				OrderId:      "id",
				PaymentId:    "id",
				FullName:     "name",
				TicketNumber: "111",
			},
		},
		Error: nil,
	}

	suite.mockOrderRepositoryQuery.On("FindOrderByUser", mock.Anything, payload).Return(mockChannel(mockOrderByUser))

	_, err := suite.usecase.FindOrderList(suite.ctx, payload)
	assert.NoError(suite.T(), err)
}

func (suite *QueryUsecaseTestSuite) TestFindOrderListErr() {
	payload := request.OrderList{
		Page:   1,
		Size:   1,
		UserId: "id",
	}
	mockOrderByUser := helpers.Result{
		Data:  nil,
		Error: errors.BadRequest("error"),
	}

	suite.mockOrderRepositoryQuery.On("FindOrderByUser", mock.Anything, payload).Return(mockChannel(mockOrderByUser))

	_, err := suite.usecase.FindOrderList(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *QueryUsecaseTestSuite) TestFindOrderListErrNil() {
	payload := request.OrderList{
		Page:   1,
		Size:   1,
		UserId: "id",
	}
	mockOrderByUser := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockOrderRepositoryQuery.On("FindOrderByUser", mock.Anything, payload).Return(mockChannel(mockOrderByUser))

	_, err := suite.usecase.FindOrderList(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *QueryUsecaseTestSuite) TestFindOrderListErrParse() {
	payload := request.OrderList{
		Page:   1,
		Size:   1,
		UserId: "id",
	}
	mockOrderByUser := helpers.Result{
		Data: &entity.BankTicket{
			TicketNumber: "11",
		},
		Error: nil,
	}

	suite.mockOrderRepositoryQuery.On("FindOrderByUser", mock.Anything, payload).Return(mockChannel(mockOrderByUser))

	_, err := suite.usecase.FindOrderList(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

// Helper function to create a channel
func mockChannel(result helpers.Result) <-chan helpers.Result {
	responseChan := make(chan helpers.Result)

	go func() {
		responseChan <- result
		close(responseChan)
	}()

	return responseChan
}

func (suite *QueryUsecaseTestSuite) TestFindPreOrderList() {
	payload := request.PreOrderList{
		Page:   1,
		Size:   1,
		UserId: "id",
	}
	mockBankTicketByUser := helpers.Result{
		Data: &[]entity.BankTicket{
			{
				TicketNumber:  "111",
				TicketType:    "Gold",
				Price:         50,
				PaymentStatus: constants.Pending,
			},
		},
		Error: nil,
	}

	suite.mockOrderRepositoryQuery.On("FindBankTicketByUser", mock.Anything, payload).Return(mockChannel(mockBankTicketByUser))

	_, err := suite.usecase.FindPreOrderList(suite.ctx, payload)
	assert.NoError(suite.T(), err)
}

func (suite *QueryUsecaseTestSuite) TestFindPreOrderListErr() {
	payload := request.PreOrderList{
		Page:   1,
		Size:   1,
		UserId: "id",
	}
	mockBankTicketByUser := helpers.Result{
		Data: &[]entity.BankTicket{
			{
				TicketNumber:  "111",
				TicketType:    "Gold",
				Price:         50,
				PaymentStatus: constants.Pending,
			},
		},
		Error: errors.BadRequest("error"),
	}

	suite.mockOrderRepositoryQuery.On("FindBankTicketByUser", mock.Anything, payload).Return(mockChannel(mockBankTicketByUser))

	_, err := suite.usecase.FindPreOrderList(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *QueryUsecaseTestSuite) TestFindPreOrderListErrNil() {
	payload := request.PreOrderList{
		Page:   1,
		Size:   1,
		UserId: "id",
	}
	mockBankTicketByUser := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	suite.mockOrderRepositoryQuery.On("FindBankTicketByUser", mock.Anything, payload).Return(mockChannel(mockBankTicketByUser))

	_, err := suite.usecase.FindPreOrderList(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *QueryUsecaseTestSuite) TestFindPreOrderListErrParse() {
	payload := request.PreOrderList{
		Page:   1,
		Size:   1,
		UserId: "id",
	}
	mockBankTicketByUser := helpers.Result{
		Data: &entity.Country{
			Name: "name",
		},
		Error: nil,
	}

	suite.mockOrderRepositoryQuery.On("FindBankTicketByUser", mock.Anything, payload).Return(mockChannel(mockBankTicketByUser))

	_, err := suite.usecase.FindPreOrderList(suite.ctx, payload)
	assert.Error(suite.T(), err)
}
