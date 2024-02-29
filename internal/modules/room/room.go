package room

import (
	"context"
	"order-service/internal/modules/room/models/entity"
	"order-service/internal/modules/room/models/request"
	"order-service/internal/modules/room/models/response"
	wrapper "order-service/internal/pkg/helpers"
)

type UsecaseCommand interface {
	CreateQueueRoom(origCtx context.Context, payload request.QueueReq) (*response.QueueResp, error)
}

type UsecaseQuery interface {
}

type MongodbRepositoryQuery interface {
	FindOneLastQueue(ctx context.Context, eventId string) <-chan wrapper.Result
	FindOneQueueByUserId(ctx context.Context, userId string, eventId string) <-chan wrapper.Result
}

type MongodbRepositoryCommand interface {
	InsertOneRoom(ctx context.Context, room entity.QueueRoom) <-chan wrapper.Result
}
