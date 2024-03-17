package usecases_test

import (
	"context"
	"order-service/internal/modules/order"
	"order-service/internal/pkg/constants"
	"order-service/internal/pkg/errors"
	"order-service/internal/pkg/helpers"
	"testing"

	eventEntity "order-service/internal/modules/event/models/entity"
	"order-service/internal/modules/order/models/entity"
	"order-service/internal/modules/order/models/request"
	uc "order-service/internal/modules/order/usecases"
	roomEntity "order-service/internal/modules/room/models/entity"
	ticketEntity "order-service/internal/modules/ticket/models/entity"
	userEntity "order-service/internal/modules/user/models/entity"
	mockcertEvent "order-service/mocks/modules/event"
	mockcert "order-service/mocks/modules/order"
	mockcertRoom "order-service/mocks/modules/room"
	mockcertTicket "order-service/mocks/modules/ticket"
	mockcertUser "order-service/mocks/modules/user"
	mocklog "order-service/mocks/pkg/log"
	mockredis "order-service/mocks/pkg/redis"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CommandUsecaseTestSuite struct {
	suite.Suite
	mockOrderRepositoryQuery    *mockcert.MongodbRepositoryQuery
	mockOrderRepositoryCommand  *mockcert.MongodbRepositoryCommand
	mockRoomRepositoryQuery     *mockcertRoom.MongodbRepositoryQuery
	mockTicketRepositoryQuery   *mockcertTicket.MongodbRepositoryQuery
	mockTicketRepositoryCommand *mockcertTicket.MongodbRepositoryCommand
	mockEventRepositoryQuery    *mockcertEvent.MongodbRepositoryQuery
	mockUserRepositoryQuery     *mockcertUser.MongodbRepositoryQuery
	mockLogger                  *mocklog.Logger
	mockRedis                   *mockredis.Collections
	usecase                     order.UsecaseCommand
	ctx                         context.Context
}

func (suite *CommandUsecaseTestSuite) SetupTest() {
	suite.mockOrderRepositoryQuery = &mockcert.MongodbRepositoryQuery{}
	suite.mockOrderRepositoryCommand = &mockcert.MongodbRepositoryCommand{}
	suite.mockRoomRepositoryQuery = &mockcertRoom.MongodbRepositoryQuery{}
	suite.mockTicketRepositoryQuery = &mockcertTicket.MongodbRepositoryQuery{}
	suite.mockTicketRepositoryCommand = &mockcertTicket.MongodbRepositoryCommand{}
	suite.mockUserRepositoryQuery = &mockcertUser.MongodbRepositoryQuery{}
	suite.mockEventRepositoryQuery = &mockcertEvent.MongodbRepositoryQuery{}
	suite.mockLogger = &mocklog.Logger{}
	suite.mockRedis = &mockredis.Collections{}
	suite.ctx = context.Background()
	suite.usecase = uc.NewCommandUsecase(
		suite.mockOrderRepositoryCommand,
		suite.mockOrderRepositoryQuery,
		suite.mockRoomRepositoryQuery,
		suite.mockTicketRepositoryQuery,
		suite.mockTicketRepositoryCommand,
		suite.mockEventRepositoryQuery,
		suite.mockUserRepositoryQuery,
		suite.mockLogger,
		suite.mockRedis,
	)
}

func TestCommandUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(CommandUsecaseTestSuite))
}

func (suite *CommandUsecaseTestSuite) TestCreateOrderTicket() {
	payload := request.OrderReq{
		UserId:     "id",
		TicketType: "type",
		EventId:    "id",
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag: "tag",
		},
		Error: nil,
	}
	mockQueueByUser := helpers.Result{
		Data: &roomEntity.QueueRoom{
			QueueId: "id",
		},
		Error: nil,
	}
	mockBankTicketByParam := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	mockTicketByEvent := helpers.Result{
		Data: &ticketEntity.Ticket{
			TicketId:       "id",
			TicketPrice:    50,
			TotalRemaining: 10,
		},
		Error: nil,
	}
	mockUserById := helpers.Result{
		Data: &userEntity.User{
			Country: userEntity.Country{
				Code: "ID",
			},
		},
		Error: nil,
	}
	mockUpdateBankTicket := helpers.Result{
		Data: &entity.BankTicket{
			TicketId:     "id",
			EventId:      "id",
			TicketNumber: "111",
		},
		Error: nil,
	}
	mockUpdateTicketDetail := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockRoomRepositoryQuery.On("FindOneQueueByUserId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockQueueByUser))
	suite.mockOrderRepositoryQuery.On("FindBankTicketByParam", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankTicketByParam))
	suite.mockTicketRepositoryQuery.On("FindTicketByEventId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockTicketByEvent))
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockUserById))
	suite.mockOrderRepositoryCommand.On("UpdateBankTicket", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateBankTicket))
	suite.mockTicketRepositoryCommand.On("UpdateOneTicketDetail", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateTicketDetail))

	_, err := suite.usecase.CreateOrderTicket(suite.ctx, payload)
	assert.NoError(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateOrderTicketErrEvent() {
	payload := request.OrderReq{
		UserId:     "id",
		TicketType: "type",
		EventId:    "id",
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag: "tag",
		},
		Error: errors.BadRequest("error"),
	}
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))

	_, err := suite.usecase.CreateOrderTicket(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateOrderTicketErrEventNil() {
	payload := request.OrderReq{
		UserId:     "id",
		TicketType: "type",
		EventId:    "id",
	}

	mockEventById := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))

	_, err := suite.usecase.CreateOrderTicket(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateOrderTicketErrEventParse() {
	payload := request.OrderReq{
		UserId:     "id",
		TicketType: "type",
		EventId:    "id",
	}

	mockEventById := helpers.Result{
		Data: &entity.Country{
			Name: "name",
		},
		Error: nil,
	}
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))

	_, err := suite.usecase.CreateOrderTicket(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateOrderTicketErrQueue() {
	payload := request.OrderReq{
		UserId:     "id",
		TicketType: "type",
		EventId:    "id",
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag: "tag",
		},
		Error: nil,
	}
	mockQueueByUser := helpers.Result{
		Data: &roomEntity.QueueRoom{
			QueueId: "id",
		},
		Error: errors.BadRequest("error"),
	}
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockRoomRepositoryQuery.On("FindOneQueueByUserId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockQueueByUser))

	_, err := suite.usecase.CreateOrderTicket(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateOrderTicketErrQueueNil() {
	payload := request.OrderReq{
		UserId:     "id",
		TicketType: "type",
		EventId:    "id",
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag: "tag",
		},
		Error: nil,
	}
	mockQueueByUser := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockRoomRepositoryQuery.On("FindOneQueueByUserId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockQueueByUser))

	_, err := suite.usecase.CreateOrderTicket(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateOrderTicketErrQueueParse() {
	payload := request.OrderReq{
		UserId:     "id",
		TicketType: "type",
		EventId:    "id",
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag: "tag",
		},
		Error: nil,
	}
	mockQueueByUser := helpers.Result{
		Data: &entity.Country{
			Name: "name",
		},
		Error: nil,
	}
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockRoomRepositoryQuery.On("FindOneQueueByUserId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockQueueByUser))

	_, err := suite.usecase.CreateOrderTicket(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateOrderTicketErrBank() {
	payload := request.OrderReq{
		UserId:     "id",
		TicketType: "type",
		EventId:    "id",
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag: "tag",
		},
		Error: nil,
	}
	mockQueueByUser := helpers.Result{
		Data: &roomEntity.QueueRoom{
			QueueId: "id",
		},
		Error: nil,
	}
	mockBankTicketByParam := helpers.Result{
		Data:  nil,
		Error: errors.InternalServerError("error"),
	}
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockRoomRepositoryQuery.On("FindOneQueueByUserId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockQueueByUser))
	suite.mockOrderRepositoryQuery.On("FindBankTicketByParam", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankTicketByParam))

	res, _ := suite.usecase.CreateOrderTicket(suite.ctx, payload)
	assert.Nil(suite.T(), res)
}

func (suite *CommandUsecaseTestSuite) TestCreateOrderTicketErrBankExist() {
	payload := request.OrderReq{
		UserId:     "id",
		TicketType: "type",
		EventId:    "id",
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag: "tag",
		},
		Error: nil,
	}
	mockQueueByUser := helpers.Result{
		Data: &roomEntity.QueueRoom{
			QueueId: "id",
		},
		Error: nil,
	}
	mockBankTicketByParam := helpers.Result{
		Data: &entity.BankTicket{
			TicketNumber: "111",
		},
		Error: nil,
	}
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockRoomRepositoryQuery.On("FindOneQueueByUserId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockQueueByUser))
	suite.mockOrderRepositoryQuery.On("FindBankTicketByParam", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankTicketByParam))

	_, err := suite.usecase.CreateOrderTicket(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateOrderTicketOnline() {
	payload := request.OrderReq{
		UserId:     "id",
		TicketType: constants.Online,
		EventId:    "id",
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag: "tag",
		},
		Error: nil,
	}
	mockQueueByUser := helpers.Result{
		Data: &roomEntity.QueueRoom{
			QueueId: "id",
		},
		Error: nil,
	}
	mockBankTicketByParam := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	mockTicketByEvent := helpers.Result{
		Data: &ticketEntity.Ticket{
			TicketId:       "id",
			TicketPrice:    50,
			TotalRemaining: 10,
		},
		Error: nil,
	}
	mockUserById := helpers.Result{
		Data: &userEntity.User{
			Country: userEntity.Country{
				Code: "ID",
			},
		},
		Error: nil,
	}
	mockUpdateBankTicket := helpers.Result{
		Data: &entity.BankTicket{
			TicketId:     "id",
			EventId:      "id",
			TicketNumber: "111",
		},
		Error: nil,
	}
	mockUpdateTicketDetail := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	mockTotalOfflineTicket := helpers.Result{
		Data: &[]ticketEntity.AggregateTotalTicket{
			{
				Id:                   "id",
				TotalAvailableTicket: 0,
			},
		},
	}
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockRoomRepositoryQuery.On("FindOneQueueByUserId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockQueueByUser))
	suite.mockOrderRepositoryQuery.On("FindBankTicketByParam", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankTicketByParam))
	suite.mockTicketRepositoryQuery.On("FindTicketByEventId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockTicketByEvent))
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockUserById))
	suite.mockOrderRepositoryCommand.On("UpdateBankTicket", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateBankTicket))
	suite.mockTicketRepositoryCommand.On("UpdateOneTicketDetail", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateTicketDetail))
	suite.mockTicketRepositoryQuery.On("FindTotalAvalailableTicketByCountry", mock.Anything, mock.Anything).Return(mockChannel(mockTotalOfflineTicket))

	_, err := suite.usecase.CreateOrderTicket(suite.ctx, payload)
	assert.NoError(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateOrderTicketOnlineErr() {
	payload := request.OrderReq{
		UserId:     "id",
		TicketType: constants.Online,
		EventId:    "id",
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag: "tag",
		},
		Error: nil,
	}
	mockQueueByUser := helpers.Result{
		Data: &roomEntity.QueueRoom{
			QueueId: "id",
		},
		Error: nil,
	}
	mockBankTicketByParam := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	mockTicketByEvent := helpers.Result{
		Data: &ticketEntity.Ticket{
			TicketId:       "id",
			TicketPrice:    50,
			TotalRemaining: 10,
		},
		Error: nil,
	}
	mockUserById := helpers.Result{
		Data: &userEntity.User{
			Country: userEntity.Country{
				Code: "ID",
			},
		},
		Error: nil,
	}
	mockUpdateBankTicket := helpers.Result{
		Data: &entity.BankTicket{
			TicketId:     "id",
			EventId:      "id",
			TicketNumber: "111",
		},
		Error: nil,
	}
	mockUpdateTicketDetail := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	mockTotalOfflineTicket := helpers.Result{
		Data:  nil,
		Error: errors.BadRequest("error"),
	}
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockRoomRepositoryQuery.On("FindOneQueueByUserId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockQueueByUser))
	suite.mockOrderRepositoryQuery.On("FindBankTicketByParam", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankTicketByParam))
	suite.mockTicketRepositoryQuery.On("FindTicketByEventId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockTicketByEvent))
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockUserById))
	suite.mockOrderRepositoryCommand.On("UpdateBankTicket", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateBankTicket))
	suite.mockTicketRepositoryCommand.On("UpdateOneTicketDetail", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateTicketDetail))
	suite.mockTicketRepositoryQuery.On("FindTotalAvalailableTicketByCountry", mock.Anything, mock.Anything).Return(mockChannel(mockTotalOfflineTicket))

	_, err := suite.usecase.CreateOrderTicket(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateOrderTicketOnlineErrNil() {
	payload := request.OrderReq{
		UserId:     "id",
		TicketType: constants.Online,
		EventId:    "id",
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag: "tag",
		},
		Error: nil,
	}
	mockQueueByUser := helpers.Result{
		Data: &roomEntity.QueueRoom{
			QueueId: "id",
		},
		Error: nil,
	}
	mockBankTicketByParam := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	mockTicketByEvent := helpers.Result{
		Data: &ticketEntity.Ticket{
			TicketId:       "id",
			TicketPrice:    50,
			TotalRemaining: 10,
		},
		Error: nil,
	}
	mockUserById := helpers.Result{
		Data: &userEntity.User{
			Country: userEntity.Country{
				Code: "ID",
			},
		},
		Error: nil,
	}
	mockUpdateBankTicket := helpers.Result{
		Data: &entity.BankTicket{
			TicketId:     "id",
			EventId:      "id",
			TicketNumber: "111",
		},
		Error: nil,
	}
	mockUpdateTicketDetail := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	mockTotalOfflineTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockRoomRepositoryQuery.On("FindOneQueueByUserId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockQueueByUser))
	suite.mockOrderRepositoryQuery.On("FindBankTicketByParam", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankTicketByParam))
	suite.mockTicketRepositoryQuery.On("FindTicketByEventId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockTicketByEvent))
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockUserById))
	suite.mockOrderRepositoryCommand.On("UpdateBankTicket", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateBankTicket))
	suite.mockTicketRepositoryCommand.On("UpdateOneTicketDetail", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateTicketDetail))
	suite.mockTicketRepositoryQuery.On("FindTotalAvalailableTicketByCountry", mock.Anything, mock.Anything).Return(mockChannel(mockTotalOfflineTicket))

	_, err := suite.usecase.CreateOrderTicket(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateOrderTicketOnlineErrParse() {
	payload := request.OrderReq{
		UserId:     "id",
		TicketType: constants.Online,
		EventId:    "id",
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag: "tag",
		},
		Error: nil,
	}
	mockQueueByUser := helpers.Result{
		Data: &roomEntity.QueueRoom{
			QueueId: "id",
		},
		Error: nil,
	}
	mockBankTicketByParam := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	mockTicketByEvent := helpers.Result{
		Data: &ticketEntity.Ticket{
			TicketId:       "id",
			TicketPrice:    50,
			TotalRemaining: 10,
		},
		Error: nil,
	}
	mockUserById := helpers.Result{
		Data: &userEntity.User{
			Country: userEntity.Country{
				Code: "ID",
			},
		},
		Error: nil,
	}
	mockUpdateBankTicket := helpers.Result{
		Data: &entity.BankTicket{
			TicketId:     "id",
			EventId:      "id",
			TicketNumber: "111",
		},
		Error: nil,
	}
	mockUpdateTicketDetail := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	mockTotalOfflineTicket := helpers.Result{
		Data: &entity.Country{
			Name: "name",
		},
		Error: nil,
	}
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockRoomRepositoryQuery.On("FindOneQueueByUserId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockQueueByUser))
	suite.mockOrderRepositoryQuery.On("FindBankTicketByParam", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankTicketByParam))
	suite.mockTicketRepositoryQuery.On("FindTicketByEventId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockTicketByEvent))
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockUserById))
	suite.mockOrderRepositoryCommand.On("UpdateBankTicket", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateBankTicket))
	suite.mockTicketRepositoryCommand.On("UpdateOneTicketDetail", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateTicketDetail))
	suite.mockTicketRepositoryQuery.On("FindTotalAvalailableTicketByCountry", mock.Anything, mock.Anything).Return(mockChannel(mockTotalOfflineTicket))

	_, err := suite.usecase.CreateOrderTicket(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateOrderTicketOnlineFalse() {
	payload := request.OrderReq{
		UserId:     "id",
		TicketType: constants.Online,
		EventId:    "id",
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag: "tag",
		},
		Error: nil,
	}
	mockQueueByUser := helpers.Result{
		Data: &roomEntity.QueueRoom{
			QueueId: "id",
		},
		Error: nil,
	}
	mockBankTicketByParam := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	mockTicketByEvent := helpers.Result{
		Data: &ticketEntity.Ticket{
			TicketId:       "id",
			TicketPrice:    50,
			TotalRemaining: 10,
		},
		Error: nil,
	}
	mockUserById := helpers.Result{
		Data: &userEntity.User{
			Country: userEntity.Country{
				Code: "ID",
			},
		},
		Error: nil,
	}
	mockUpdateBankTicket := helpers.Result{
		Data: &entity.BankTicket{
			TicketId:     "id",
			EventId:      "id",
			TicketNumber: "111",
		},
		Error: nil,
	}
	mockUpdateTicketDetail := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	mockTotalOfflineTicket := helpers.Result{
		Data: &[]ticketEntity.AggregateTotalTicket{
			{
				Id:                   "id",
				TotalAvailableTicket: 1,
			},
		},
	}
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockRoomRepositoryQuery.On("FindOneQueueByUserId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockQueueByUser))
	suite.mockOrderRepositoryQuery.On("FindBankTicketByParam", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankTicketByParam))
	suite.mockTicketRepositoryQuery.On("FindTicketByEventId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockTicketByEvent))
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockUserById))
	suite.mockOrderRepositoryCommand.On("UpdateBankTicket", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateBankTicket))
	suite.mockTicketRepositoryCommand.On("UpdateOneTicketDetail", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateTicketDetail))
	suite.mockTicketRepositoryQuery.On("FindTotalAvalailableTicketByCountry", mock.Anything, mock.Anything).Return(mockChannel(mockTotalOfflineTicket))

	_, err := suite.usecase.CreateOrderTicket(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateOrderTicketErrDetail() {
	payload := request.OrderReq{
		UserId:     "id",
		TicketType: "type",
		EventId:    "id",
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag: "tag",
		},
		Error: nil,
	}
	mockQueueByUser := helpers.Result{
		Data: &roomEntity.QueueRoom{
			QueueId: "id",
		},
		Error: nil,
	}
	mockBankTicketByParam := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	mockTicketByEvent := helpers.Result{
		Data:  nil,
		Error: errors.BadRequest("error"),
	}
	mockUserById := helpers.Result{
		Data: &userEntity.User{
			Country: userEntity.Country{
				Code: "ID",
			},
		},
		Error: nil,
	}
	mockUpdateBankTicket := helpers.Result{
		Data: &entity.BankTicket{
			TicketId:     "id",
			EventId:      "id",
			TicketNumber: "111",
		},
		Error: nil,
	}
	mockUpdateTicketDetail := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockRoomRepositoryQuery.On("FindOneQueueByUserId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockQueueByUser))
	suite.mockOrderRepositoryQuery.On("FindBankTicketByParam", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankTicketByParam))
	suite.mockTicketRepositoryQuery.On("FindTicketByEventId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockTicketByEvent))
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockUserById))
	suite.mockOrderRepositoryCommand.On("UpdateBankTicket", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateBankTicket))
	suite.mockTicketRepositoryCommand.On("UpdateOneTicketDetail", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateTicketDetail))

	_, err := suite.usecase.CreateOrderTicket(suite.ctx, payload)
	assert.NoError(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateOrderTicketErrDetailNil() {
	payload := request.OrderReq{
		UserId:     "id",
		TicketType: "type",
		EventId:    "id",
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag: "tag",
		},
		Error: nil,
	}
	mockQueueByUser := helpers.Result{
		Data: &roomEntity.QueueRoom{
			QueueId: "id",
		},
		Error: nil,
	}
	mockBankTicketByParam := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	mockTicketByEvent := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	mockUserById := helpers.Result{
		Data: &userEntity.User{
			Country: userEntity.Country{
				Code: "ID",
			},
		},
		Error: nil,
	}
	mockUpdateBankTicket := helpers.Result{
		Data: &entity.BankTicket{
			TicketId:     "id",
			EventId:      "id",
			TicketNumber: "111",
		},
		Error: nil,
	}
	mockUpdateTicketDetail := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockRoomRepositoryQuery.On("FindOneQueueByUserId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockQueueByUser))
	suite.mockOrderRepositoryQuery.On("FindBankTicketByParam", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankTicketByParam))
	suite.mockTicketRepositoryQuery.On("FindTicketByEventId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockTicketByEvent))
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockUserById))
	suite.mockOrderRepositoryCommand.On("UpdateBankTicket", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateBankTicket))
	suite.mockTicketRepositoryCommand.On("UpdateOneTicketDetail", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateTicketDetail))

	_, err := suite.usecase.CreateOrderTicket(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateOrderTicketErrDetailParse() {
	payload := request.OrderReq{
		UserId:     "id",
		TicketType: "type",
		EventId:    "id",
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag: "tag",
		},
		Error: nil,
	}
	mockQueueByUser := helpers.Result{
		Data: &roomEntity.QueueRoom{
			QueueId: "id",
		},
		Error: nil,
	}
	mockBankTicketByParam := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	mockTicketByEvent := helpers.Result{
		Data: &entity.Country{
			Name: "name",
		},
		Error: nil,
	}
	mockUserById := helpers.Result{
		Data: &userEntity.User{
			Country: userEntity.Country{
				Code: "ID",
			},
		},
		Error: nil,
	}
	mockUpdateBankTicket := helpers.Result{
		Data: &entity.BankTicket{
			TicketId:     "id",
			EventId:      "id",
			TicketNumber: "111",
		},
		Error: nil,
	}
	mockUpdateTicketDetail := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockRoomRepositoryQuery.On("FindOneQueueByUserId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockQueueByUser))
	suite.mockOrderRepositoryQuery.On("FindBankTicketByParam", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankTicketByParam))
	suite.mockTicketRepositoryQuery.On("FindTicketByEventId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockTicketByEvent))
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockUserById))
	suite.mockOrderRepositoryCommand.On("UpdateBankTicket", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateBankTicket))
	suite.mockTicketRepositoryCommand.On("UpdateOneTicketDetail", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateTicketDetail))

	_, err := suite.usecase.CreateOrderTicket(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateOrderTicketSold() {
	payload := request.OrderReq{
		UserId:     "id",
		TicketType: "type",
		EventId:    "id",
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag: "tag",
		},
		Error: nil,
	}
	mockQueueByUser := helpers.Result{
		Data: &roomEntity.QueueRoom{
			QueueId: "id",
		},
		Error: nil,
	}
	mockBankTicketByParam := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	mockTicketByEvent := helpers.Result{
		Data: &ticketEntity.Ticket{
			TicketId:       "id",
			TicketPrice:    50,
			TotalRemaining: 0,
		},
		Error: nil,
	}
	mockUserById := helpers.Result{
		Data: &userEntity.User{
			Country: userEntity.Country{
				Code: "ID",
			},
		},
		Error: nil,
	}
	mockUpdateBankTicket := helpers.Result{
		Data: &entity.BankTicket{
			TicketId:     "id",
			EventId:      "id",
			TicketNumber: "111",
		},
		Error: nil,
	}
	mockUpdateTicketDetail := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockRoomRepositoryQuery.On("FindOneQueueByUserId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockQueueByUser))
	suite.mockOrderRepositoryQuery.On("FindBankTicketByParam", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankTicketByParam))
	suite.mockTicketRepositoryQuery.On("FindTicketByEventId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockTicketByEvent))
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockUserById))
	suite.mockOrderRepositoryCommand.On("UpdateBankTicket", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateBankTicket))
	suite.mockTicketRepositoryCommand.On("UpdateOneTicketDetail", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateTicketDetail))

	_, err := suite.usecase.CreateOrderTicket(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateOrderTicketErrUser() {
	payload := request.OrderReq{
		UserId:     "id",
		TicketType: "type",
		EventId:    "id",
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag: "tag",
		},
		Error: nil,
	}
	mockQueueByUser := helpers.Result{
		Data: &roomEntity.QueueRoom{
			QueueId: "id",
		},
		Error: nil,
	}
	mockBankTicketByParam := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	mockTicketByEvent := helpers.Result{
		Data: &ticketEntity.Ticket{
			TicketId:       "id",
			TicketPrice:    50,
			TotalRemaining: 10,
		},
		Error: nil,
	}
	mockUserById := helpers.Result{
		Data: &userEntity.User{
			Country: userEntity.Country{
				Code: "ID",
			},
		},
		Error: errors.BadRequest("error"),
	}
	mockUpdateBankTicket := helpers.Result{
		Data: &entity.BankTicket{
			TicketId:     "id",
			EventId:      "id",
			TicketNumber: "111",
		},
		Error: nil,
	}
	mockUpdateTicketDetail := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockRoomRepositoryQuery.On("FindOneQueueByUserId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockQueueByUser))
	suite.mockOrderRepositoryQuery.On("FindBankTicketByParam", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankTicketByParam))
	suite.mockTicketRepositoryQuery.On("FindTicketByEventId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockTicketByEvent))
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockUserById))
	suite.mockOrderRepositoryCommand.On("UpdateBankTicket", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateBankTicket))
	suite.mockTicketRepositoryCommand.On("UpdateOneTicketDetail", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateTicketDetail))

	_, err := suite.usecase.CreateOrderTicket(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateOrderTicketErrUserNil() {
	payload := request.OrderReq{
		UserId:     "id",
		TicketType: "type",
		EventId:    "id",
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag: "tag",
		},
		Error: nil,
	}
	mockQueueByUser := helpers.Result{
		Data: &roomEntity.QueueRoom{
			QueueId: "id",
		},
		Error: nil,
	}
	mockBankTicketByParam := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	mockTicketByEvent := helpers.Result{
		Data: &ticketEntity.Ticket{
			TicketId:       "id",
			TicketPrice:    50,
			TotalRemaining: 10,
		},
		Error: nil,
	}
	mockUserById := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	mockUpdateBankTicket := helpers.Result{
		Data: &entity.BankTicket{
			TicketId:     "id",
			EventId:      "id",
			TicketNumber: "111",
		},
		Error: nil,
	}
	mockUpdateTicketDetail := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockRoomRepositoryQuery.On("FindOneQueueByUserId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockQueueByUser))
	suite.mockOrderRepositoryQuery.On("FindBankTicketByParam", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankTicketByParam))
	suite.mockTicketRepositoryQuery.On("FindTicketByEventId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockTicketByEvent))
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockUserById))
	suite.mockOrderRepositoryCommand.On("UpdateBankTicket", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateBankTicket))
	suite.mockTicketRepositoryCommand.On("UpdateOneTicketDetail", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateTicketDetail))

	_, err := suite.usecase.CreateOrderTicket(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateOrderTicketErrUserParse() {
	payload := request.OrderReq{
		UserId:     "id",
		TicketType: "type",
		EventId:    "id",
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag: "tag",
		},
		Error: nil,
	}
	mockQueueByUser := helpers.Result{
		Data: &roomEntity.QueueRoom{
			QueueId: "id",
		},
		Error: nil,
	}
	mockBankTicketByParam := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	mockTicketByEvent := helpers.Result{
		Data: &ticketEntity.Ticket{
			TicketId:       "id",
			TicketPrice:    50,
			TotalRemaining: 10,
		},
		Error: nil,
	}
	mockUserById := helpers.Result{
		Data: &entity.Country{
			Name: "name",
		},
		Error: nil,
	}
	mockUpdateBankTicket := helpers.Result{
		Data: &entity.BankTicket{
			TicketId:     "id",
			EventId:      "id",
			TicketNumber: "111",
		},
		Error: nil,
	}
	mockUpdateTicketDetail := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockRoomRepositoryQuery.On("FindOneQueueByUserId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockQueueByUser))
	suite.mockOrderRepositoryQuery.On("FindBankTicketByParam", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankTicketByParam))
	suite.mockTicketRepositoryQuery.On("FindTicketByEventId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockTicketByEvent))
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockUserById))
	suite.mockOrderRepositoryCommand.On("UpdateBankTicket", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateBankTicket))
	suite.mockTicketRepositoryCommand.On("UpdateOneTicketDetail", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateTicketDetail))

	_, err := suite.usecase.CreateOrderTicket(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateOrderTicketErrUpdateBank() {
	payload := request.OrderReq{
		UserId:     "id",
		TicketType: "type",
		EventId:    "id",
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag: "tag",
		},
		Error: nil,
	}
	mockQueueByUser := helpers.Result{
		Data: &roomEntity.QueueRoom{
			QueueId: "id",
		},
		Error: nil,
	}
	mockBankTicketByParam := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	mockTicketByEvent := helpers.Result{
		Data: &ticketEntity.Ticket{
			TicketId:       "id",
			TicketPrice:    50,
			TotalRemaining: 10,
		},
		Error: nil,
	}
	mockUserById := helpers.Result{
		Data: &userEntity.User{
			Country: userEntity.Country{
				Code: "ID",
			},
		},
		Error: nil,
	}
	mockUpdateBankTicket := helpers.Result{
		Data: &entity.BankTicket{
			TicketId:     "id",
			EventId:      "id",
			TicketNumber: "111",
		},
		Error: errors.BadRequest("error"),
	}
	mockUpdateTicketDetail := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockRoomRepositoryQuery.On("FindOneQueueByUserId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockQueueByUser))
	suite.mockOrderRepositoryQuery.On("FindBankTicketByParam", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankTicketByParam))
	suite.mockTicketRepositoryQuery.On("FindTicketByEventId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockTicketByEvent))
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockUserById))
	suite.mockOrderRepositoryCommand.On("UpdateBankTicket", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateBankTicket))
	suite.mockTicketRepositoryCommand.On("UpdateOneTicketDetail", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateTicketDetail))

	_, err := suite.usecase.CreateOrderTicket(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateOrderTicketErrUpdateBankNil() {
	payload := request.OrderReq{
		UserId:     "id",
		TicketType: "type",
		EventId:    "id",
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag: "tag",
		},
		Error: nil,
	}
	mockQueueByUser := helpers.Result{
		Data: &roomEntity.QueueRoom{
			QueueId: "id",
		},
		Error: nil,
	}
	mockBankTicketByParam := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	mockTicketByEvent := helpers.Result{
		Data: &ticketEntity.Ticket{
			TicketId:       "id",
			TicketPrice:    50,
			TotalRemaining: 10,
		},
		Error: nil,
	}
	mockUserById := helpers.Result{
		Data: &userEntity.User{
			Country: userEntity.Country{
				Code: "ID",
			},
		},
		Error: nil,
	}
	mockUpdateBankTicket := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	mockUpdateTicketDetail := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockRoomRepositoryQuery.On("FindOneQueueByUserId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockQueueByUser))
	suite.mockOrderRepositoryQuery.On("FindBankTicketByParam", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankTicketByParam))
	suite.mockTicketRepositoryQuery.On("FindTicketByEventId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockTicketByEvent))
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockUserById))
	suite.mockOrderRepositoryCommand.On("UpdateBankTicket", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateBankTicket))
	suite.mockTicketRepositoryCommand.On("UpdateOneTicketDetail", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateTicketDetail))

	_, err := suite.usecase.CreateOrderTicket(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateOrderTicketErrUpdateBankParse() {
	payload := request.OrderReq{
		UserId:     "id",
		TicketType: "type",
		EventId:    "id",
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag: "tag",
		},
		Error: nil,
	}
	mockQueueByUser := helpers.Result{
		Data: &roomEntity.QueueRoom{
			QueueId: "id",
		},
		Error: nil,
	}
	mockBankTicketByParam := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	mockTicketByEvent := helpers.Result{
		Data: &ticketEntity.Ticket{
			TicketId:       "id",
			TicketPrice:    50,
			TotalRemaining: 10,
		},
		Error: nil,
	}
	mockUserById := helpers.Result{
		Data: &userEntity.User{
			Country: userEntity.Country{
				Code: "ID",
			},
		},
		Error: nil,
	}
	mockUpdateBankTicket := helpers.Result{
		Data: &entity.Country{
			Name: "name",
		},
		Error: nil,
	}
	mockUpdateTicketDetail := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockRoomRepositoryQuery.On("FindOneQueueByUserId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockQueueByUser))
	suite.mockOrderRepositoryQuery.On("FindBankTicketByParam", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankTicketByParam))
	suite.mockTicketRepositoryQuery.On("FindTicketByEventId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockTicketByEvent))
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockUserById))
	suite.mockOrderRepositoryCommand.On("UpdateBankTicket", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateBankTicket))
	suite.mockTicketRepositoryCommand.On("UpdateOneTicketDetail", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateTicketDetail))

	_, err := suite.usecase.CreateOrderTicket(suite.ctx, payload)
	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateOrderTicketErrUpdateDetail() {
	payload := request.OrderReq{
		UserId:     "id",
		TicketType: "type",
		EventId:    "id",
	}

	mockEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Name:    "name",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag: "tag",
		},
		Error: nil,
	}
	mockQueueByUser := helpers.Result{
		Data: &roomEntity.QueueRoom{
			QueueId: "id",
		},
		Error: nil,
	}
	mockBankTicketByParam := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	mockTicketByEvent := helpers.Result{
		Data: &ticketEntity.Ticket{
			TicketId:       "id",
			TicketPrice:    50,
			TotalRemaining: 10,
		},
		Error: nil,
	}
	mockUserById := helpers.Result{
		Data: &userEntity.User{
			Country: userEntity.Country{
				Code: "ID",
			},
		},
		Error: nil,
	}
	mockUpdateBankTicket := helpers.Result{
		Data: &entity.BankTicket{
			TicketId:     "id",
			EventId:      "id",
			TicketNumber: "111",
		},
		Error: nil,
	}
	mockUpdateTicketDetail := helpers.Result{
		Data:  nil,
		Error: errors.BadRequest("error"),
	}
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockEventById))
	suite.mockRoomRepositoryQuery.On("FindOneQueueByUserId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockQueueByUser))
	suite.mockOrderRepositoryQuery.On("FindBankTicketByParam", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockBankTicketByParam))
	suite.mockTicketRepositoryQuery.On("FindTicketByEventId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockTicketByEvent))
	suite.mockUserRepositoryQuery.On("FindOneUserId", mock.Anything, mock.Anything).Return(mockChannel(mockUserById))
	suite.mockOrderRepositoryCommand.On("UpdateBankTicket", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateBankTicket))
	suite.mockTicketRepositoryCommand.On("UpdateOneTicketDetail", mock.Anything, mock.Anything).Return(mockChannel(mockUpdateTicketDetail))

	_, err := suite.usecase.CreateOrderTicket(suite.ctx, payload)
	assert.Error(suite.T(), err)
}
