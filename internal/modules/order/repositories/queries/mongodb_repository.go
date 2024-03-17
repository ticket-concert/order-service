package queries

import (
	"context"
	"order-service/internal/modules/order"
	"order-service/internal/modules/order/models/entity"
	"order-service/internal/modules/order/models/request"
	"order-service/internal/pkg/databases/mongodb"
	wrapper "order-service/internal/pkg/helpers"
	"order-service/internal/pkg/log"

	"go.mongodb.org/mongo-driver/bson"
)

type queryMongodbRepository struct {
	mongoDb mongodb.Collections
	logger  log.Logger
}

func NewQueryMongodbRepository(mongodb mongodb.Collections, log log.Logger) order.MongodbRepositoryQuery {
	return &queryMongodbRepository{
		mongoDb: mongodb,
		logger:  log,
	}
}

func (q queryMongodbRepository) FindBankTicketByParam(ctx context.Context, queueId string, userId string) <-chan wrapper.Result {
	var bankTicket entity.BankTicket
	output := make(chan wrapper.Result)

	go func() {
		resp := <-q.mongoDb.FindOne(mongodb.FindOne{
			Result:         &bankTicket,
			CollectionName: "bank-ticket",
			Filter: bson.M{
				"$and": []interface{}{
					bson.M{"queueId": queueId},
					bson.M{"userId": userId},
				},
			},
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}

func (q queryMongodbRepository) FindOrderByUser(ctx context.Context, payload request.OrderList) <-chan wrapper.Result {
	var orders []entity.Order
	var countData int64
	output := make(chan wrapper.Result)

	go func() {
		resp := <-q.mongoDb.FindAllData(mongodb.FindAllData{
			Result:         &orders,
			CountData:      &countData,
			CollectionName: "order",
			Filter:         bson.M{"userId": payload.UserId},
			Sort: &mongodb.Sort{
				FieldName: "createdAt",
				By:        mongodb.SortDescending,
			},
			Page: payload.Page,
			Size: payload.Size,
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}

func (q queryMongodbRepository) FindBankTicketByUser(ctx context.Context, payload request.PreOrderList) <-chan wrapper.Result {
	var bankTicket []entity.BankTicket
	var countData int64
	output := make(chan wrapper.Result)

	go func() {
		resp := <-q.mongoDb.FindAllData(mongodb.FindAllData{
			Result:         &bankTicket,
			CountData:      &countData,
			CollectionName: "bank-ticket",
			Filter:         bson.M{"userId": payload.UserId},
			Sort: &mongodb.Sort{
				FieldName: "updatedAt",
				By:        mongodb.SortDescending,
			},
			Page: payload.Page,
			Size: payload.Size,
		}, ctx)
		output <- resp
		close(output)
	}()

	return output
}
