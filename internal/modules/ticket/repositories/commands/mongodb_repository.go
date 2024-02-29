package commands

import (
	"context"
	"order-service/internal/modules/ticket"
	"order-service/internal/modules/ticket/models/entity"
	"order-service/internal/pkg/databases/mongodb"
	wrapper "order-service/internal/pkg/helpers"
	"order-service/internal/pkg/log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type commandMongodbRepository struct {
	mongoDb mongodb.Collections
	logger  log.Logger
}

func NewCommandMongodbRepository(mongodb mongodb.Collections, log log.Logger) ticket.MongodbRepositoryCommand {
	return &commandMongodbRepository{
		mongoDb: mongodb,
		logger:  log,
	}
}

func (c commandMongodbRepository) UpdateOneTicketDetail(ctx context.Context, payload entity.Ticket) <-chan wrapper.Result {
	output := make(chan wrapper.Result)

	go func() {
		resp := <-c.mongoDb.UpdateOne(mongodb.UpdateOne{
			CollectionName: "ticket-detail",
			Filter: bson.M{
				"ticketId": payload.TicketId,
				"eventId":  payload.EventId,
			},
			Document: bson.M{
				"totalRemaining": payload.TotalRemaining,
				"updatedAt":      time.Now(),
			},
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}
