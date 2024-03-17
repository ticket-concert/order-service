package commands

import (
	"context"
	room "order-service/internal/modules/room"
	"order-service/internal/modules/room/models/entity"
	"order-service/internal/pkg/databases/mongodb"
	wrapper "order-service/internal/pkg/helpers"
	"order-service/internal/pkg/log"
	"time"

	"github.com/google/uuid"
)

type commandMongodbRepository struct {
	mongoDb mongodb.Collections
	logger  log.Logger
}

func NewCommandMongodbRepository(mongodb mongodb.Collections, log log.Logger) room.MongodbRepositoryCommand {
	return &commandMongodbRepository{
		mongoDb: mongodb,
		logger:  log,
	}
}

func (c commandMongodbRepository) InsertOneRoom(ctx context.Context, room entity.QueueRoom) <-chan wrapper.Result {
	output := make(chan wrapper.Result)
	room.QueueId = uuid.NewString()
	room.CreatedAt = time.Now()
	room.UpdatedAt = time.Now()

	go func() {
		resp := <-c.mongoDb.InsertOne(mongodb.InsertOne{
			CollectionName: "queue-room",
			Document:       room,
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}

// func (c commandMongodbRepository) DeleteOneQueue(ctx context.Context, queueId string) <-chan wrapper.Result {
// 	output := make(chan wrapper.Result)

// 	go func() {
// 		resp := <-c.mongoDb.DeleteOne(mongodb.DeleteOne{
// 			CollectionName: "queue-room",
// 			Filter: bson.M{
// 				"queueId": queueId,
// 			},
// 		}, ctx)
// 		output <- resp
// 		close(output)
// 	}()

// 	return output
// }
