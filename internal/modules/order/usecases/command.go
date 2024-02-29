package usecases

import (
	"context"
	"order-service/configs"
	"order-service/internal/modules/event"
	eventEntity "order-service/internal/modules/event/models/entity"
	"order-service/internal/modules/order"
	"order-service/internal/modules/order/models/entity"
	"order-service/internal/modules/order/models/request"
	"order-service/internal/modules/order/models/response"
	"order-service/internal/modules/room"
	roomEntity "order-service/internal/modules/room/models/entity"
	"order-service/internal/modules/ticket"
	ticketEntity "order-service/internal/modules/ticket/models/entity"
	ticketRequest "order-service/internal/modules/ticket/models/request"
	"order-service/internal/modules/user"
	userEntity "order-service/internal/modules/user/models/entity"
	"order-service/internal/pkg/constants"
	"order-service/internal/pkg/errors"
	"order-service/internal/pkg/log"
	"order-service/internal/pkg/redis"
	"time"

	"go.elastic.co/apm"
)

type commandUsecase struct {
	orderRepositoryCommand  order.MongodbRepositoryCommand
	orderRepositoryQuery    order.MongodbRepositoryQuery
	roomRepositoryQuery     room.MongodbRepositoryQuery
	ticketRepositoryQuery   ticket.MongodbRepositoryQuery
	ticketRepositoryCommand ticket.MongodbRepositoryCommand
	eventRepositoryQuery    event.MongodbRepositoryQuery
	userRepositoryQuery     user.MongodbRepositoryQuery
	logger                  log.Logger
	redis                   redis.Collections
}

func NewCommandUsecase(
	omc order.MongodbRepositoryCommand, omq order.MongodbRepositoryQuery, rmq room.MongodbRepositoryQuery,
	trq ticket.MongodbRepositoryQuery, trc ticket.MongodbRepositoryCommand,
	emq event.MongodbRepositoryQuery, umq user.MongodbRepositoryQuery, log log.Logger, rc redis.Collections) order.UsecaseCommand {
	return commandUsecase{
		orderRepositoryCommand:  omc,
		orderRepositoryQuery:    omq,
		roomRepositoryQuery:     rmq,
		ticketRepositoryQuery:   trq,
		ticketRepositoryCommand: trc,
		eventRepositoryQuery:    emq,
		userRepositoryQuery:     umq,
		logger:                  log,
		redis:                   rc,
	}
}

func (c commandUsecase) CreateOrderTicket(origCtx context.Context, payload request.OrderReq) (*response.OrderResp, error) {
	domain := "orderUsecase-CreateOrderTicket"
	span, ctx := apm.StartSpanOptions(origCtx, domain, "function", apm.SpanOptions{
		Start:  time.Now(),
		Parent: apm.TraceContext{},
	})
	defer span.End()

	if configs.GetConfig().DayFlag {
		day := time.Now().Weekday()
		if day != time.Saturday && day != time.Sunday {
			return nil, errors.BadRequest("This day not Saturday or Sunday")
		}
	}

	eventData := <-c.eventRepositoryQuery.FindEventById(ctx, payload.EventId)
	if eventData.Error != nil {
		return nil, eventData.Error
	}

	if eventData.Data == nil {
		return nil, errors.BadRequest("event not found")
	}

	event, ok := eventData.Data.(*eventEntity.Event)
	if !ok {
		return nil, errors.InternalServerError("cannot parsing data event")
	}

	queueRoom := <-c.roomRepositoryQuery.FindOneQueueByUserId(ctx, payload.UserId, payload.EventId)
	if queueRoom.Error != nil {
		return nil, queueRoom.Error
	}

	if queueRoom.Data == nil {
		return nil, errors.BadRequest("user not in the queue")
	}

	queueData, ok := queueRoom.Data.(*roomEntity.QueueRoom)
	if !ok {
		return nil, errors.InternalServerError("cannot parsing data queue")
	}

	currentTicket := <-c.orderRepositoryQuery.FindBankTicketByParam(ctx, queueData.QueueId, payload.UserId)
	if currentTicket.Error != nil {
		return nil, queueRoom.Error
	}

	if currentTicket.Data != nil {
		return nil, errors.BadRequest("user already order ticket")
	}

	if payload.TicketType == constants.Online {
		ticketReq := ticketRequest.TicketReq{
			CountryCode: event.Country.Code,
			Tag:         event.Tag,
		}
		totalAvailableTicket := <-c.ticketRepositoryQuery.FindTotalAvalailableTicketByCountry(ctx, ticketReq)
		if totalAvailableTicket.Error != nil {
			return nil, totalAvailableTicket.Error
		}

		if totalAvailableTicket.Data == nil {
			return nil, errors.BadRequest("country not found")
		}

		totalAvailable, ok := totalAvailableTicket.Data.(*[]ticketEntity.AggregateTotalTicket)
		if !ok {
			return nil, errors.InternalServerError("cannot parsing data")
		}
		eligibleBuy := true
		for i, v := range *totalAvailable {
			if i == 0 && v.TotalAvailableTicket != 0 {
				eligibleBuy = false
			}
		}
		if !eligibleBuy {
			return nil, errors.BadRequest("offline ticket still ready")
		}
	}

	ticketDetailData := <-c.ticketRepositoryQuery.FindTicketByEventId(ctx, event.EventId, payload.TicketType)
	if ticketDetailData.Error != nil {
		return nil, queueRoom.Error
	}

	if ticketDetailData.Data == nil {
		return nil, errors.BadRequest("ticket detail not found")
	}

	ticketDetail, ok := ticketDetailData.Data.(*ticketEntity.Ticket)
	if !ok {
		return nil, errors.InternalServerError("cannot parsing data")
	}

	if ticketDetail.TotalRemaining == 0 {
		return nil, errors.BadRequest("ticket category sold out")
	}

	userData := <-c.userRepositoryQuery.FindOneUserId(ctx, payload.UserId)
	if userData.Error != nil {
		return nil, userData.Error
	}

	if userData.Data == nil {
		return nil, errors.BadRequest("user not found")
	}

	user, ok := userData.Data.(*userEntity.User)
	if !ok {
		return nil, errors.InternalServerError("cannot parsing data user")
	}

	price := ticketDetail.TicketPrice
	if user.Country.Code != event.Country.Code && payload.TicketType != constants.Online {
		// discount 20% if buy different country
		price = price * 80 / 100
	}

	bankTicketReq := request.UpdateBankTicketReq{
		CountryCode:   event.Country.Code,
		TicketType:    payload.TicketType,
		Price:         price,
		UserId:        payload.UserId,
		QueueId:       queueData.QueueId,
		TicketId:      ticketDetail.TicketId,
		EventId:       event.EventId,
		PaymentStatus: constants.Pending,
		UpdatedAt:     time.Now(),
	}

	bankTicket := <-c.orderRepositoryCommand.UpdateBankTicket(ctx, bankTicketReq)
	if bankTicket.Error != nil {
		return nil, bankTicket.Error
	}

	if bankTicket.Data == nil {
		return nil, errors.BadRequest("failed to process order")
	}

	ticket, ok := bankTicket.Data.(*entity.BankTicket)
	if !ok {
		return nil, errors.InternalServerError("cannot parsing data")
	}

	ticketPayload := ticketEntity.Ticket{
		TicketId:       ticket.TicketId,
		EventId:        ticket.EventId,
		TotalRemaining: ticketDetail.TotalRemaining - 1,
	}

	ticketResp := <-c.ticketRepositoryCommand.UpdateOneTicketDetail(ctx, ticketPayload)
	if ticketResp.Error != nil {
		return nil, ticketResp.Error
	}

	return &response.OrderResp{
		TicketNumber: ticket.TicketNumber,
		QueueId:      ticket.QueueId,
		UserId:       ticket.UserId,
		EventId:      ticket.EventId,
		TicketType:   ticket.TicketType,
		CountryCode:  ticket.CountryCode,
		Price:        ticket.Price,
		OrderTime:    ticket.UpdatedAt,
	}, nil
}
