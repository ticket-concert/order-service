package usecases_test

import (
	"context"
	"order-service/internal/modules/room"
	"order-service/internal/pkg/errors"
	"order-service/internal/pkg/helpers"
	"testing"

	eventEntity "order-service/internal/modules/event/models/entity"
	roomEntity "order-service/internal/modules/room/models/entity"
	"order-service/internal/modules/room/models/request"
	uc "order-service/internal/modules/room/usecases"
	ticketEntity "order-service/internal/modules/ticket/models/entity"
	mockcertEvent "order-service/mocks/modules/event"
	mockcert "order-service/mocks/modules/room"
	mockcertTicket "order-service/mocks/modules/ticket"
	mocklog "order-service/mocks/pkg/log"
	mockredis "order-service/mocks/pkg/redis"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CommandUsecaseTestSuite struct {
	suite.Suite
	mockRoomRepositoryQuery   *mockcert.MongodbRepositoryQuery
	mockRoomRepositoryCommand *mockcert.MongodbRepositoryCommand
	mockTicketRepositoryQuery *mockcertTicket.MongodbRepositoryQuery
	mockEventRepositoryQuery  *mockcertEvent.MongodbRepositoryQuery
	mockLogger                *mocklog.Logger
	mockRedis                 *mockredis.Collections
	usecase                   room.UsecaseCommand
	ctx                       context.Context
}

func (suite *CommandUsecaseTestSuite) SetupTest() {
	suite.mockRoomRepositoryQuery = &mockcert.MongodbRepositoryQuery{}
	suite.mockRoomRepositoryCommand = &mockcert.MongodbRepositoryCommand{}
	suite.mockTicketRepositoryQuery = &mockcertTicket.MongodbRepositoryQuery{}
	suite.mockEventRepositoryQuery = &mockcertEvent.MongodbRepositoryQuery{}
	suite.mockLogger = &mocklog.Logger{}
	suite.mockRedis = &mockredis.Collections{}
	suite.ctx = context.Background()
	suite.usecase = uc.NewCommandUsecase(
		suite.mockRoomRepositoryQuery,
		suite.mockRoomRepositoryCommand,
		suite.mockTicketRepositoryQuery,
		suite.mockEventRepositoryQuery,
		suite.mockLogger,
		suite.mockRedis,
	)
}

func TestCommandUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(CommandUsecaseTestSuite))
}

func (suite *CommandUsecaseTestSuite) TestCreateQueueRoom() {
	payload := request.QueueReq{
		UserId:  "id",
		EventId: "id",
	}
	mockFindEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag: "tag",
		},
		Error: nil,
	}
	mockFindOneQueueByUserId := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockFindOneLastQueue := helpers.Result{
		Data: &roomEntity.QueueRoom{
			QueueId:     "id",
			QueueNumber: 1,
		},
		Error: nil,
	}

	mockInsertOneRoom := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	data := roomEntity.QueueRoom{
		UserId:      "id",
		EventId:     "id",
		QueueNumber: 2,
		CountryCode: "code",
	}

	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockFindEventById))
	suite.mockRoomRepositoryQuery.On("FindOneQueueByUserId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockFindOneQueueByUserId))
	suite.mockRoomRepositoryQuery.On("FindOneLastQueue", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockFindOneLastQueue))
	suite.mockRedis.On("Get", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(redis.NewStringResult("5", nil))
	suite.mockRoomRepositoryCommand.On("InsertOneRoom", suite.ctx, data).Return(mockChannel(mockInsertOneRoom))

	_, err := suite.usecase.CreateQueueRoom(suite.ctx, payload)

	assert.NoError(suite.T(), err)

	mockFindEventById3 := helpers.Result{
		Data:  nil,
		Error: nil,
	}
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockFindEventById3))
	suite.mockRoomRepositoryQuery.On("FindOneQueueByUserId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockFindOneQueueByUserId))
	suite.mockRoomRepositoryQuery.On("FindOneLastQueue", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockFindOneLastQueue))
	suite.mockRedis.On("Get", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(redis.NewStringResult("5", nil))
	suite.mockRoomRepositoryCommand.On("InsertOneRoom", suite.ctx, data).Return(mockChannel(mockInsertOneRoom))

	_, err3 := suite.usecase.CreateQueueRoom(suite.ctx, payload)

	assert.Error(suite.T(), err3)
}

func (suite *CommandUsecaseTestSuite) TestCreateQueueRoomErrDay() {
	payload := request.QueueReq{
		UserId:  "id",
		EventId: "id",
	}

	mockFindEventById := helpers.Result{
		Data:  nil,
		Error: errors.BadRequest("error"),
	}
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockFindEventById))

	_, err := suite.usecase.CreateQueueRoom(suite.ctx, payload)

	assert.Error(suite.T(), err)

	mockFindEventById2 := helpers.Result{
		Data: &roomEntity.QueueRoom{
			QueueId: "id",
		},
		Error: nil,
	}
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockFindEventById2))

	_, err2 := suite.usecase.CreateQueueRoom(suite.ctx, payload)

	assert.Error(suite.T(), err2)
}

func (suite *CommandUsecaseTestSuite) TestCreateQueueRoomErr() {
	payload := request.QueueReq{
		UserId:  "id",
		EventId: "id",
	}

	mockFindEventById := helpers.Result{
		Data:  nil,
		Error: errors.BadRequest("error"),
	}
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockFindEventById))

	_, err := suite.usecase.CreateQueueRoom(suite.ctx, payload)

	assert.Error(suite.T(), err)

	mockFindEventById2 := helpers.Result{
		Data: &roomEntity.QueueRoom{
			QueueId: "id",
		},
		Error: nil,
	}
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockFindEventById2))

	_, err2 := suite.usecase.CreateQueueRoom(suite.ctx, payload)

	assert.Error(suite.T(), err2)
}

func (suite *CommandUsecaseTestSuite) TestCreateQueueRoomErrParse() {
	payload := request.QueueReq{
		UserId:  "id",
		EventId: "id",
	}

	mockFindEventById := helpers.Result{
		Data: &roomEntity.QueueRoom{
			QueueId: "id",
		},
		Error: nil,
	}
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockFindEventById))

	_, err := suite.usecase.CreateQueueRoom(suite.ctx, payload)

	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateQueueRoomErrQueue() {
	payload := request.QueueReq{
		UserId:  "id",
		EventId: "id",
	}

	mockFindEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag: "tag",
		},
		Error: nil,
	}
	mockFindOneQueueByUserId := helpers.Result{
		Data:  nil,
		Error: errors.BadRequest("error"),
	}
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockFindEventById))
	suite.mockRoomRepositoryQuery.On("FindOneQueueByUserId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockFindOneQueueByUserId))

	_, err := suite.usecase.CreateQueueRoom(suite.ctx, payload)

	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateQueueRoomErrLastQueue() {
	payload := request.QueueReq{
		UserId:  "id",
		EventId: "id",
	}
	mockFindEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag: "tag",
		},
		Error: nil,
	}
	mockFindOneQueueByUserId := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockFindOneLastQueue := helpers.Result{
		Data:  nil,
		Error: errors.BadRequest("error"),
	}

	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockFindEventById))
	suite.mockRoomRepositoryQuery.On("FindOneQueueByUserId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockFindOneQueueByUserId))
	suite.mockRoomRepositoryQuery.On("FindOneLastQueue", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockFindOneLastQueue))
	suite.mockRedis.On("Get", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(redis.NewStringResult("5", nil))

	_, err := suite.usecase.CreateQueueRoom(suite.ctx, payload)

	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateQueueRoomErrLastQueueParse() {
	payload := request.QueueReq{
		UserId:  "id",
		EventId: "id",
	}
	mockFindEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag: "tag",
		},
		Error: nil,
	}
	mockFindOneQueueByUserId := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockFindOneLastQueue := helpers.Result{
		Data: &eventEntity.Country{
			Name: "name",
		},
		Error: nil,
	}

	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockFindEventById))
	suite.mockRoomRepositoryQuery.On("FindOneQueueByUserId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockFindOneQueueByUserId))
	suite.mockRoomRepositoryQuery.On("FindOneLastQueue", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockFindOneLastQueue))
	suite.mockRedis.On("Get", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(redis.NewStringResult("5", nil))

	_, err := suite.usecase.CreateQueueRoom(suite.ctx, payload)

	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateQueueRoomErrNotNil() {
	payload := request.QueueReq{
		UserId:  "id",
		EventId: "id",
	}

	mockFindEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag: "tag",
		},
		Error: nil,
	}
	mockFindOneQueueByUserId := helpers.Result{
		Data: &roomEntity.QueueRoom{
			QueueId: "id",
		},
		Error: nil,
	}
	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockFindEventById))
	suite.mockRoomRepositoryQuery.On("FindOneQueueByUserId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockFindOneQueueByUserId))

	_, err := suite.usecase.CreateQueueRoom(suite.ctx, payload)

	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateQueueRoomEmptyRedis() {
	payload := request.QueueReq{
		UserId:  "id",
		EventId: "id",
	}
	mockFindEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag: "tag",
		},
		Error: nil,
	}
	mockFindOneQueueByUserId := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockFindOneLastQueue := helpers.Result{
		Data: &roomEntity.QueueRoom{
			QueueId:     "id",
			QueueNumber: 1,
		},
		Error: nil,
	}

	mockInsertOneRoom := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	data := roomEntity.QueueRoom{
		UserId:      "id",
		EventId:     "id",
		QueueNumber: 2,
		CountryCode: "code",
	}
	totalAvailableTicket := helpers.Result{
		Data: &[]ticketEntity.AggregateTotalTicket{
			{
				Id:                   "1",
				TotalAvailableTicket: 15,
			},
		},
		Error: nil,
	}

	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockFindEventById))
	suite.mockRoomRepositoryQuery.On("FindOneQueueByUserId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockFindOneQueueByUserId))
	suite.mockRoomRepositoryQuery.On("FindOneLastQueue", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockFindOneLastQueue))
	suite.mockRedis.On("Get", mock.Anything, mock.Anything).Return(redis.NewStringResult("", nil))
	suite.mockTicketRepositoryQuery.On("FindTotalAvalailableTicket", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(totalAvailableTicket))
	suite.mockRoomRepositoryCommand.On("InsertOneRoom", suite.ctx, data).Return(mockChannel(mockInsertOneRoom))
	suite.mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	_, err := suite.usecase.CreateQueueRoom(suite.ctx, payload)

	assert.NoError(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateQueueRoomEmptyRedisErr() {
	payload := request.QueueReq{
		UserId:  "id",
		EventId: "id",
	}
	mockFindEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag: "tag",
		},
		Error: nil,
	}
	mockFindOneQueueByUserId := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockFindOneLastQueue := helpers.Result{
		Data: &roomEntity.QueueRoom{
			QueueId:     "id",
			QueueNumber: 1,
		},
		Error: nil,
	}

	mockInsertOneRoom := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	data := roomEntity.QueueRoom{
		UserId:      "id",
		EventId:     "id",
		QueueNumber: 2,
		CountryCode: "code",
	}
	totalAvailableTicket := helpers.Result{
		Data:  nil,
		Error: errors.BadRequest("error"),
	}

	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockFindEventById))
	suite.mockRoomRepositoryQuery.On("FindOneQueueByUserId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockFindOneQueueByUserId))
	suite.mockRoomRepositoryQuery.On("FindOneLastQueue", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockFindOneLastQueue))
	suite.mockRedis.On("Get", mock.Anything, mock.Anything).Return(redis.NewStringResult("", nil))
	suite.mockTicketRepositoryQuery.On("FindTotalAvalailableTicket", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(totalAvailableTicket))
	suite.mockRoomRepositoryCommand.On("InsertOneRoom", suite.ctx, data).Return(mockChannel(mockInsertOneRoom))
	suite.mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	_, err := suite.usecase.CreateQueueRoom(suite.ctx, payload)

	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateQueueRoomEmptyRedisErrTotal() {
	payload := request.QueueReq{
		UserId:  "id",
		EventId: "id",
	}
	mockFindEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag: "tag",
		},
		Error: nil,
	}
	mockFindOneQueueByUserId := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockFindOneLastQueue := helpers.Result{
		Data: &roomEntity.QueueRoom{
			QueueId:     "id",
			QueueNumber: 1,
		},
		Error: nil,
	}

	mockInsertOneRoom := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	data := roomEntity.QueueRoom{
		UserId:      "id",
		EventId:     "id",
		QueueNumber: 2,
		CountryCode: "code",
	}
	totalAvailableTicket := helpers.Result{
		Data: &eventEntity.Country{
			Name: "name",
		},
		Error: nil,
	}

	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockFindEventById))
	suite.mockRoomRepositoryQuery.On("FindOneQueueByUserId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockFindOneQueueByUserId))
	suite.mockRoomRepositoryQuery.On("FindOneLastQueue", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockFindOneLastQueue))
	suite.mockRedis.On("Get", mock.Anything, mock.Anything).Return(redis.NewStringResult("", nil))
	suite.mockTicketRepositoryQuery.On("FindTotalAvalailableTicket", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(totalAvailableTicket))
	suite.mockRoomRepositoryCommand.On("InsertOneRoom", suite.ctx, data).Return(mockChannel(mockInsertOneRoom))
	suite.mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	_, err := suite.usecase.CreateQueueRoom(suite.ctx, payload)

	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateQueueRoomEmptyRedisErrParse() {
	payload := request.QueueReq{
		UserId:  "id",
		EventId: "id",
	}
	mockFindEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag: "tag",
		},
		Error: nil,
	}
	mockFindOneQueueByUserId := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockFindOneLastQueue := helpers.Result{
		Data: &roomEntity.QueueRoom{
			QueueId:     "id",
			QueueNumber: 1,
		},
		Error: nil,
	}

	mockInsertOneRoom := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	data := roomEntity.QueueRoom{
		UserId:      "id",
		EventId:     "id",
		QueueNumber: 2,
		CountryCode: "code",
	}

	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockFindEventById))
	suite.mockRoomRepositoryQuery.On("FindOneQueueByUserId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockFindOneQueueByUserId))
	suite.mockRoomRepositoryQuery.On("FindOneLastQueue", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockFindOneLastQueue))
	suite.mockRedis.On("Get", mock.Anything, mock.Anything).Return(redis.NewStringResult("tes", nil))
	suite.mockRoomRepositoryCommand.On("InsertOneRoom", suite.ctx, data).Return(mockChannel(mockInsertOneRoom))

	_, err := suite.usecase.CreateQueueRoom(suite.ctx, payload)

	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateQueueRoomErrLimit() {
	payload := request.QueueReq{
		UserId:  "id",
		EventId: "id",
	}
	mockFindEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag: "tag",
		},
		Error: nil,
	}
	mockFindOneQueueByUserId := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockFindOneLastQueue := helpers.Result{
		Data: &roomEntity.QueueRoom{
			QueueId:     "id",
			QueueNumber: 1,
		},
		Error: nil,
	}

	mockInsertOneRoom := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	data := roomEntity.QueueRoom{
		UserId:      "id",
		EventId:     "id",
		QueueNumber: 2,
		CountryCode: "code",
	}

	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockFindEventById))
	suite.mockRoomRepositoryQuery.On("FindOneQueueByUserId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockFindOneQueueByUserId))
	suite.mockRoomRepositoryQuery.On("FindOneLastQueue", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockFindOneLastQueue))
	suite.mockRedis.On("Get", mock.Anything, mock.Anything).Return(redis.NewStringResult("1", nil))
	suite.mockRoomRepositoryCommand.On("InsertOneRoom", suite.ctx, data).Return(mockChannel(mockInsertOneRoom))

	_, err := suite.usecase.CreateQueueRoom(suite.ctx, payload)

	assert.Error(suite.T(), err)
}

func (suite *CommandUsecaseTestSuite) TestCreateQueueRoomErrInsert() {
	payload := request.QueueReq{
		UserId:  "id",
		EventId: "id",
	}
	mockFindEventById := helpers.Result{
		Data: &eventEntity.Event{
			EventId: "id",
			Country: eventEntity.Country{
				Code: "code",
			},
			Tag: "tag",
		},
		Error: nil,
	}
	mockFindOneQueueByUserId := helpers.Result{
		Data:  nil,
		Error: nil,
	}

	mockFindOneLastQueue := helpers.Result{
		Data: &roomEntity.QueueRoom{
			QueueId:     "id",
			QueueNumber: 1,
		},
		Error: nil,
	}

	mockInsertOneRoom := helpers.Result{
		Data:  nil,
		Error: errors.BadRequest("error"),
	}

	data := roomEntity.QueueRoom{
		UserId:      "id",
		EventId:     "id",
		QueueNumber: 2,
		CountryCode: "code",
	}

	suite.mockEventRepositoryQuery.On("FindEventById", mock.Anything, mock.Anything).Return(mockChannel(mockFindEventById))
	suite.mockRoomRepositoryQuery.On("FindOneQueueByUserId", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockFindOneQueueByUserId))
	suite.mockRoomRepositoryQuery.On("FindOneLastQueue", mock.Anything, mock.Anything, mock.Anything).Return(mockChannel(mockFindOneLastQueue))
	suite.mockRedis.On("Get", mock.Anything, mock.Anything).Return(redis.NewStringResult("5", nil))
	suite.mockRoomRepositoryCommand.On("InsertOneRoom", suite.ctx, data).Return(mockChannel(mockInsertOneRoom))

	_, err := suite.usecase.CreateQueueRoom(suite.ctx, payload)

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
