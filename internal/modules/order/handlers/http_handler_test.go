// user_http_handler_test.go

package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"order-service/internal/modules/order/handlers"
	"order-service/internal/modules/order/models/request"
	"order-service/internal/modules/order/models/response"
	"order-service/internal/pkg/constants"
	"order-service/internal/pkg/errors"
	mockcert "order-service/mocks/modules/order"
	mocklog "order-service/mocks/pkg/log"
	mockredis "order-service/mocks/pkg/redis"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
)

type OrderHttpHandlerTestSuite struct {
	suite.Suite

	cUC       *mockcert.UsecaseCommand
	cUQ       *mockcert.UsecaseQuery
	cLog      *mocklog.Logger
	validator *validator.Validate
	handler   *handlers.OrderHttpHandler
	cRedis    *mockredis.Collections
	app       *fiber.App
}

func (suite *OrderHttpHandlerTestSuite) SetupTest() {
	suite.cUC = new(mockcert.UsecaseCommand)
	suite.cUQ = new(mockcert.UsecaseQuery)
	suite.cLog = new(mocklog.Logger)
	suite.validator = validator.New()
	suite.cRedis = new(mockredis.Collections)
	suite.handler = &handlers.OrderHttpHandler{
		OrderUsecaseCommand: suite.cUC,
		OrderUsecaseQuery:   suite.cUQ,
		Logger:              suite.cLog,
		Validator:           suite.validator,
	}
	suite.app = fiber.New()
	handlers.InitOrderHttpHandler(suite.app, suite.cUC, suite.cUQ, suite.cLog, suite.cRedis)
}

func TestUserHttpHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(OrderHttpHandlerTestSuite))
}

func (suite *OrderHttpHandlerTestSuite) TestCreateOrderTicket() {

	suite.cUC.On("CreateOrderTicket", mock.Anything, mock.Anything).Return(&response.OrderResp{
		TicketNumber: "111",
	}, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	payload := request.OrderReq{
		UserId:     "id",
		TicketType: "Gold",
		EventId:    "id",
	}

	requestBody, _ := json.Marshal(payload)
	req := httptest.NewRequest(fiber.MethodPost, "/v1/create-order", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Locals("userId", "12345")
	ctx.Request().SetRequestURI("/v1/create-order")
	ctx.Request().Header.SetMethod(fiber.MethodPost)
	ctx.Request().Header.SetContentType("application/json")
	ctx.Request().SetBody(requestBody)

	err := suite.handler.CreateOrder(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *OrderHttpHandlerTestSuite) TestCreateOrderTicketErrBody() {

	suite.cUC.On("CreateOrderTicket", mock.Anything, mock.Anything).Return(&response.OrderResp{
		TicketNumber: "111",
	}, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	payload := request.OrderReq{
		UserId:     "id",
		TicketType: "Gold",
		EventId:    "id",
	}

	requestBody, _ := json.Marshal(payload)
	req := httptest.NewRequest(fiber.MethodPost, "/v1/create-order", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Locals("userId", "12345")
	ctx.Request().SetRequestURI("/v1/create-order")
	ctx.Request().Header.SetMethod(fiber.MethodPost)
	ctx.Request().Header.SetContentType("application/json")
	// ctx.Request().SetBody(requestBody)

	err := suite.handler.CreateOrder(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *OrderHttpHandlerTestSuite) TestCreateOrderTicketErrValidator() {

	suite.cUC.On("CreateOrderTicket", mock.Anything, mock.Anything).Return(&response.OrderResp{
		TicketNumber: "111",
	}, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	payload := request.OrderReq{}

	requestBody, _ := json.Marshal(payload)
	req := httptest.NewRequest(fiber.MethodPost, "/v1/create-order", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Locals("userId", "12345")
	ctx.Request().SetRequestURI("/v1/create-order")
	ctx.Request().Header.SetMethod(fiber.MethodPost)
	ctx.Request().Header.SetContentType("application/json")
	ctx.Request().SetBody(requestBody)

	err := suite.handler.CreateOrder(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *OrderHttpHandlerTestSuite) TestCreateOrderTicketErr() {

	suite.cUC.On("CreateOrderTicket", mock.Anything, mock.Anything).Return(nil, errors.BadRequest("error"))
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	payload := request.OrderReq{
		UserId:     "id",
		TicketType: "Gold",
		EventId:    "id",
	}

	requestBody, _ := json.Marshal(payload)
	req := httptest.NewRequest(fiber.MethodPost, "/v1/create-order", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Locals("userId", "12345")
	ctx.Request().SetRequestURI("/v1/create-order")
	ctx.Request().Header.SetMethod(fiber.MethodPost)
	ctx.Request().Header.SetContentType("application/json")
	ctx.Request().SetBody(requestBody)

	err := suite.handler.CreateOrder(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *OrderHttpHandlerTestSuite) TestGetOrderList() {

	response := &response.OrderListResp{
		CollectionData: []response.OrderList{
			{
				TicketNumber: "id",
				FullName:     "name",
			},
		},
		MetaData: constants.MetaData{},
	}
	suite.cUQ.On("FindOrderList", mock.Anything, mock.Anything).Return(response, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Locals("userId", "12345")
	req := httptest.NewRequest(fiber.MethodGet, "/v1/list?page=1&size=1", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/v1/list?page=1&size=1")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetOrderList(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *OrderHttpHandlerTestSuite) TestGetOrderListErrQuery() {

	response := &response.OrderListResp{
		CollectionData: []response.OrderList{
			{
				TicketNumber: "id",
				FullName:     "name",
			},
		},
		MetaData: constants.MetaData{},
	}
	suite.cUQ.On("FindOrderList", mock.Anything, mock.Anything).Return(response, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Locals("userId", "12345")
	req := httptest.NewRequest(fiber.MethodGet, "/v1/list?page=aa&size=aa", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/v1/list?page=aa&size=aa")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetOrderList(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *OrderHttpHandlerTestSuite) TestGetOrderListErrValidation() {

	response := &response.OrderListResp{
		CollectionData: []response.OrderList{
			{
				TicketNumber: "id",
				FullName:     "name",
			},
		},
		MetaData: constants.MetaData{},
	}
	suite.cUQ.On("FindOrderList", mock.Anything, mock.Anything).Return(response, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Locals("userId", "12345")
	req := httptest.NewRequest(fiber.MethodGet, "/v1/list?page=&size=", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/v1/list?page=&size=")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetOrderList(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *OrderHttpHandlerTestSuite) TestGetOrderListErr() {

	suite.cUQ.On("FindOrderList", mock.Anything, mock.Anything).Return(nil, errors.BadRequest("error"))
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Locals("userId", "12345")
	req := httptest.NewRequest(fiber.MethodGet, "/v1/list?page=1&size=1", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/v1/list?page=1&size=1")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetOrderList(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *OrderHttpHandlerTestSuite) TestGetPreOrderList() {

	response := &response.PreOrderListResp{
		CollectionData: []response.PreOrderList{
			{
				TicketNumber: "id",
				TicketType:   "Gold",
			},
		},
		MetaData: constants.MetaData{},
	}
	suite.cUQ.On("FindPreOrderList", mock.Anything, mock.Anything).Return(response, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Locals("userId", "12345")
	req := httptest.NewRequest(fiber.MethodGet, "/v1/preorder-list?page=1&size=1", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/v1/preorder-list?page=1&size=1")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetPreOrderList(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *OrderHttpHandlerTestSuite) TestGetPreOrderListErrQuery() {

	response := &response.PreOrderListResp{
		CollectionData: []response.PreOrderList{
			{
				TicketNumber: "id",
				TicketType:   "Gold",
			},
		},
		MetaData: constants.MetaData{},
	}
	suite.cUQ.On("FindPreOrderList", mock.Anything, mock.Anything).Return(response, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Locals("userId", "12345")
	req := httptest.NewRequest(fiber.MethodGet, "/v1/preorder-list?page=aa&size=aa", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/v1/preorder-list?page=aa&size=aa")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetPreOrderList(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *OrderHttpHandlerTestSuite) TestGetPreOrderListErrValidation() {

	response := &response.PreOrderListResp{
		CollectionData: []response.PreOrderList{
			{
				TicketNumber: "id",
				TicketType:   "Gold",
			},
		},
		MetaData: constants.MetaData{},
	}
	suite.cUQ.On("FindPreOrderList", mock.Anything, mock.Anything).Return(response, nil)
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Locals("userId", "12345")
	req := httptest.NewRequest(fiber.MethodGet, "/v1/preorder-list?page=&size=", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/v1/preorder-list?page=&size=")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetPreOrderList(ctx)
	assert.Nil(suite.T(), err)
}

func (suite *OrderHttpHandlerTestSuite) TestGetPreOrderListErr() {

	suite.cUQ.On("FindPreOrderList", mock.Anything, mock.Anything).Return(nil, errors.BadRequest("error"))
	suite.cLog.On("Info", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ctx := suite.app.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Locals("userId", "12345")
	req := httptest.NewRequest(fiber.MethodGet, "/v1/preorder-list?page=1&size=1", nil)
	req.Header.Set("Content-Type", "application/json")
	ctx.Request().SetRequestURI("/v1/preorder-list?page=1&size=1")
	ctx.Request().Header.SetMethod(fiber.MethodGet)
	ctx.Request().Header.SetContentType("application/json")

	err := suite.handler.GetPreOrderList(ctx)
	assert.Nil(suite.T(), err)
}
