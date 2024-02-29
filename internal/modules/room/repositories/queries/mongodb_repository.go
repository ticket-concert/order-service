package queries

import (
	"context"
	room "order-service/internal/modules/room"
	"order-service/internal/modules/room/models/entity"
	"order-service/internal/pkg/databases/mongodb"
	wrapper "order-service/internal/pkg/helpers"
	"order-service/internal/pkg/log"

	"go.mongodb.org/mongo-driver/bson"
)

type queryMongodbRepository struct {
	mongoDb mongodb.Collections
	logger  log.Logger
}

func NewQueryMongodbRepository(mongodb mongodb.Collections, log log.Logger) room.MongodbRepositoryQuery {
	return &queryMongodbRepository{
		mongoDb: mongodb,
		logger:  log,
	}
}

func (q queryMongodbRepository) FindOneLastQueue(ctx context.Context, eventId string) <-chan wrapper.Result {
	var room entity.QueueRoom
	output := make(chan wrapper.Result)

	go func() {
		resp := <-q.mongoDb.FindOne(mongodb.FindOne{
			Result:         &room,
			CollectionName: "queue-room",
			Filter: bson.M{
				"eventId": eventId,
			},
			Sort: &mongodb.Sort{
				FieldName: "seatNumber",
				By:        mongodb.SortDescending,
			},
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}

func (q queryMongodbRepository) FindOneQueueByUserId(ctx context.Context, userId string, eventId string) <-chan wrapper.Result {
	var room entity.QueueRoom
	output := make(chan wrapper.Result)

	go func() {
		resp := <-q.mongoDb.FindOne(mongodb.FindOne{
			Result:         &room,
			CollectionName: "queue-room",
			Filter: bson.M{
				"userId":  userId,
				"eventId": eventId,
			},
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}
