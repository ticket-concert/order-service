package order

import (
	"context"
	"order-service/internal/modules/order/models/request"
	"order-service/internal/modules/order/models/response"
	wrapper "order-service/internal/pkg/helpers"
)

type UsecaseCommand interface {
	CreateOrderTicket(origCtx context.Context, payload request.OrderReq) (*response.OrderResp, error)
}

type UsecaseQuery interface {
	FindOrderList(origCtx context.Context, payload request.OrderList) (*response.OrderListResp, error)
	FindPreOrderList(origCtx context.Context, payload request.PreOrderList) (*response.PreOrderListResp, error)
}

type MongodbRepositoryQuery interface {
	FindBankTicketByParam(ctx context.Context, queueId string, userId string) <-chan wrapper.Result
	FindOrderByUser(ctx context.Context, payload request.OrderList) <-chan wrapper.Result
	FindBankTicketByUser(ctx context.Context, payload request.PreOrderList) <-chan wrapper.Result
}

type MongodbRepositoryCommand interface {
	UpdateBankTicket(ctx context.Context, payload request.UpdateBankTicketReq) <-chan wrapper.Result
}
