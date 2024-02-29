package queries

import (
	"context"
	address "order-service/internal/modules/ticket"
	"order-service/internal/modules/ticket/models/entity"
	"order-service/internal/modules/ticket/models/request"
	"order-service/internal/pkg/databases/mongodb"
	wrapper "order-service/internal/pkg/helpers"
	"order-service/internal/pkg/log"

	"go.mongodb.org/mongo-driver/bson"
)

type queryMongodbRepository struct {
	mongoDb mongodb.Collections
	logger  log.Logger
}

func NewQueryMongodbRepository(mongodb mongodb.Collections, log log.Logger) address.MongodbRepositoryQuery {
	return &queryMongodbRepository{
		mongoDb: mongodb,
		logger:  log,
	}
}

func (q queryMongodbRepository) FindTotalAvalailableTicket(ctx context.Context, countryCode string, tag string) <-chan wrapper.Result {
	var ticket []entity.AggregateTotalTicket
	output := make(chan wrapper.Result)

	go func() {
		resp := <-q.mongoDb.Aggregate(mongodb.Aggregate{
			Result:         &ticket,
			CollectionName: "ticket-detail",
			Filter: []bson.M{
				{
					"$match": bson.M{"country.code": countryCode, "tag": tag},
				},
				{
					"$group": bson.M{"_id": "$country.code", "totalAvailableTicket": bson.M{"$sum": "$totalRemaining"}},
				},
			},
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}

func (q queryMongodbRepository) FindTicketByEventId(ctx context.Context, eventId string, ticketType string) <-chan wrapper.Result {
	var tickets entity.Ticket
	output := make(chan wrapper.Result)

	go func() {
		resp := <-q.mongoDb.FindOne(mongodb.FindOne{
			Result:         &tickets,
			CollectionName: "ticket-detail",
			Filter: bson.M{
				"eventId":    eventId,
				"ticketType": ticketType,
			},
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}

func (q queryMongodbRepository) FindTotalAvalailableTicketByCountry(ctx context.Context, payload request.TicketReq) <-chan wrapper.Result {
	var ticket []entity.AggregateTotalTicket
	output := make(chan wrapper.Result)

	go func() {
		resp := <-q.mongoDb.Aggregate(mongodb.Aggregate{
			Result:         &ticket,
			CollectionName: "ticket-detail",
			Filter: []bson.M{
				{
					"$match": bson.M{
						"tag":          payload.Tag,
						"country.code": payload.CountryCode,
						"ticketType":   bson.M{"$ne": "Online"},
					},
				},
				{
					"$group": bson.M{
						"_id":                  "$country.code",
						"totalAvailableTicket": bson.M{"$sum": "$totalRemaining"},
					},
				},
				{
					"$sort": bson.M{
						"totalAvailableTicket": 1,
					},
				},
			},
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}
