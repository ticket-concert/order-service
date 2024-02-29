package queries

import (
	"context"
	"order-service/internal/modules/event"
	"order-service/internal/modules/event/models/entity"
	"order-service/internal/pkg/databases/mongodb"
	wrapper "order-service/internal/pkg/helpers"
	"order-service/internal/pkg/log"

	"go.mongodb.org/mongo-driver/bson"
)

type queryMongodbRepository struct {
	mongoDb mongodb.Collections
	logger  log.Logger
}

func NewQueryMongodbRepository(mongodb mongodb.Collections, log log.Logger) event.MongodbRepositoryQuery {
	return &queryMongodbRepository{
		mongoDb: mongodb,
		logger:  log,
	}
}

func (q queryMongodbRepository) FindEventById(ctx context.Context, eventId string) <-chan wrapper.Result {
	var event entity.Event
	output := make(chan wrapper.Result)

	go func() {
		resp := <-q.mongoDb.FindOne(mongodb.FindOne{
			Result:         &event,
			CollectionName: "event",
			Filter: bson.M{
				"eventId": eventId,
			},
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}
