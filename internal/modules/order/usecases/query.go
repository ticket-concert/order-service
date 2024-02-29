package usecases

import (
	"context"
	"order-service/internal/modules/order"
	"order-service/internal/modules/order/models/entity"
	"order-service/internal/modules/order/models/request"
	"order-service/internal/modules/order/models/response"
	"order-service/internal/pkg/constants"
	"order-service/internal/pkg/errors"
	"order-service/internal/pkg/helpers"
	"order-service/internal/pkg/log"
	"time"

	"go.elastic.co/apm"
)

type queryUsecase struct {
	orderRepositoryQuery order.MongodbRepositoryQuery
	logger               log.Logger
}

func NewQueryUsecase(omq order.MongodbRepositoryQuery, log log.Logger) order.UsecaseQuery {
	return queryUsecase{
		orderRepositoryQuery: omq,
		logger:               log,
	}
}

func (q queryUsecase) FindOrderList(origCtx context.Context, payload request.OrderList) (*response.OrderListResp, error) {
	domain := "orderUsecase-FindOrderList"
	span, ctx := apm.StartSpanOptions(origCtx, domain, "function", apm.SpanOptions{
		Start:  time.Now(),
		Parent: apm.TraceContext{},
	})
	defer span.End()

	orderData := <-q.orderRepositoryQuery.FindOrderByUser(ctx, payload)
	if orderData.Error != nil {
		return nil, orderData.Error
	}

	if orderData.Data == nil {
		return nil, errors.BadRequest("order not found")
	}

	orders, ok := orderData.Data.(*[]entity.Order)
	if !ok {
		return nil, errors.InternalServerError("cannot parsing data order")
	}

	var collectionData = make([]response.OrderList, 0)
	for _, value := range *orders {
		collectionData = append(collectionData, response.OrderList{
			FullName:     value.FullName,
			TicketNumber: value.TicketNumber,
			TicketType:   value.TicketType,
			TicketPrice:  value.Amount,
			SeatNumber:   value.SeatNumber,
			EventName:    value.EventName,
			EventTime:    value.DateTime,
			EventPlace:   value.Country.Place,
			OrderTime:    value.OrderTime,
			EventId:      value.EventId,
			TicketId:     value.TicketId,
		})
	}

	return &response.OrderListResp{
		CollectionData: collectionData,
		MetaData:       helpers.GenerateMetaData(orderData.Count, int64(len(*orders)), payload.Page, payload.Size),
	}, nil
}

func (q queryUsecase) FindPreOrderList(origCtx context.Context, payload request.PreOrderList) (*response.PreOrderListResp, error) {
	domain := "orderUsecase-FindOrderList"
	span, ctx := apm.StartSpanOptions(origCtx, domain, "function", apm.SpanOptions{
		Start:  time.Now(),
		Parent: apm.TraceContext{},
	})
	defer span.End()

	bankTicketData := <-q.orderRepositoryQuery.FindBankTicketByUser(ctx, payload)
	if bankTicketData.Error != nil {
		return nil, bankTicketData.Error
	}

	if bankTicketData.Data == nil {
		return nil, errors.BadRequest("order not found")
	}

	bankTicket, ok := bankTicketData.Data.(*[]entity.BankTicket)
	if !ok {
		return nil, errors.InternalServerError("cannot parsing data order")
	}

	var collectionData = make([]response.PreOrderList, 0)
	for _, value := range *bankTicket {
		var maxWaitTime string
		if value.PaymentStatus == constants.Pending {
			count := 15
			then := value.UpdatedAt.Local().Add(time.Duration(+count) * time.Minute)
			maxWaitTime = then.Format("2006-01-02 15:04")
		}
		collectionData = append(collectionData, response.PreOrderList{
			TicketNumber: value.TicketNumber,
			TicketType:   value.TicketType,
			TicketPrice:  value.Price,
			OrderTime:    value.UpdatedAt.Local(),
			UserId:       value.UserId,
			EventId:      value.EventId,
			MaxWaitTime:  maxWaitTime,
		})
	}

	return &response.PreOrderListResp{
		CollectionData: collectionData,
		MetaData:       helpers.GenerateMetaData(bankTicketData.Count, int64(len(*bankTicket)), payload.Page, payload.Size),
	}, nil
}
