package ticket

import (
	"context"
	"order-service/internal/modules/ticket/models/entity"
	"order-service/internal/modules/ticket/models/request"
	wrapper "order-service/internal/pkg/helpers"
)

type MongodbRepositoryQuery interface {
	FindTotalAvalailableTicket(ctx context.Context, countryCode string, tag string) <-chan wrapper.Result
	FindTotalAvalailableTicketByCountry(ctx context.Context, payload request.TicketReq) <-chan wrapper.Result
	FindTicketByEventId(ctx context.Context, eventId string, ticketType string) <-chan wrapper.Result
}

type MongodbRepositoryCommand interface {
	UpdateOneTicketDetail(ctx context.Context, payload entity.Ticket) <-chan wrapper.Result
}
