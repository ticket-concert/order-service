package commands

import (
	"context"
	"order-service/internal/modules/order"
	"order-service/internal/modules/order/models/entity"
	"order-service/internal/modules/order/models/request"
	"order-service/internal/pkg/databases/mongodb"
	wrapper "order-service/internal/pkg/helpers"
	"order-service/internal/pkg/log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type commandMongodbRepository struct {
	mongoDb mongodb.Collections
	logger  log.Logger
}

func NewCommandMongodbRepository(mongodb mongodb.Collections, log log.Logger) order.MongodbRepositoryCommand {
	return &commandMongodbRepository{
		mongoDb: mongodb,
		logger:  log,
	}
}

func (c commandMongodbRepository) UpdateBankTicket(ctx context.Context, payload request.UpdateBankTicketReq) <-chan wrapper.Result {
	output := make(chan wrapper.Result)
	var bankTicket entity.BankTicket

	go func() {
		resp := <-c.mongoDb.FindOneAndUpdate(mongodb.FindOneAndUpdate{
			CollectionName: "bank-ticket",
			Result:         &bankTicket,
			Filter: bson.M{
				"isUsed":     false,
				"eventId":    payload.EventId,
				"ticketType": payload.TicketType,
			},
			Update: bson.M{
				"$set": bson.M{
					"isUsed":        true,
					"userId":        payload.UserId,
					"price":         payload.Price,
					"queueId":       payload.QueueId,
					"ticketId":      payload.TicketId,
					"eventId":       payload.EventId,
					"paymentStatus": payload.PaymentStatus,
					"updatedAt":     payload.UpdatedAt,
				},
			},
			Upsert: false,
		}, options.After, ctx)
		output <- resp
		close(output)
	}()

	return output
}
