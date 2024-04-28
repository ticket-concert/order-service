package usecases

import (
	"context"
	"fmt"
	"order-service/configs"
	"order-service/internal/modules/event"
	"order-service/internal/modules/room"
	"order-service/internal/modules/room/models/entity"
	"order-service/internal/modules/room/models/request"
	"order-service/internal/modules/room/models/response"
	"order-service/internal/modules/ticket"
	"order-service/internal/pkg/constants"
	"order-service/internal/pkg/errors"
	"order-service/internal/pkg/helpers"
	"order-service/internal/pkg/log"
	"order-service/internal/pkg/redis"
	"strconv"
	"time"

	eventEntity "order-service/internal/modules/event/models/entity"
	ticketEntity "order-service/internal/modules/ticket/models/entity"

	"go.elastic.co/apm"
)

type commandUsecase struct {
	roomRepositoryQuery   room.MongodbRepositoryQuery
	roomRepositoryCommand room.MongodbRepositoryCommand
	ticketRepositoryQuery ticket.MongodbRepositoryQuery
	eventRepositoryQuery  event.MongodbRepositoryQuery
	logger                log.Logger
	redis                 redis.Collections
}

func NewCommandUsecase(
	rmq room.MongodbRepositoryQuery, rmc room.MongodbRepositoryCommand,
	trq ticket.MongodbRepositoryQuery, emq event.MongodbRepositoryQuery, log log.Logger, rc redis.Collections) room.UsecaseCommand {
	return commandUsecase{
		roomRepositoryQuery:   rmq,
		roomRepositoryCommand: rmc,
		ticketRepositoryQuery: trq,
		eventRepositoryQuery:  emq,
		logger:                log,
		redis:                 rc,
	}
}

func (c commandUsecase) CreateQueueRoom(origCtx context.Context, payload request.QueueReq) (*response.QueueResp, error) {
	domain := "roomUsecase-CreateQueueRoom"
	span, ctx := apm.StartSpanOptions(origCtx, domain, "function", apm.SpanOptions{
		Start:  time.Now(),
		Parent: apm.TraceContext{},
	})
	defer span.End()

	if configs.GetConfig().DayFlag {
		day := time.Now().Weekday()
		if day != time.Saturday && day != time.Sunday {
			msg := "this day not Saturday or Sunday"
			c.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
			return nil, errors.BadRequest("this day not Saturday or Sunday")
		}
	}

	eventData := <-c.eventRepositoryQuery.FindEventById(ctx, payload.EventId)
	if eventData.Error != nil {
		msg := "Error DB connection FindEventById"
		c.logger.Error(ctx, msg, fmt.Sprintf("%+v", eventData.Error))
		return nil, eventData.Error
	}

	if eventData.Data == nil {
		msg := "event not found"
		c.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		return nil, errors.BadRequest("event not found")
	}

	event, ok := eventData.Data.(*eventEntity.Event)
	if !ok {
		msg := "cannot parsing data event"
		c.logger.Error(ctx, msg, fmt.Sprintf("%+v", eventData.Data))
		return nil, errors.InternalServerError("cannot parsing data event")
	}

	queueRoom := <-c.roomRepositoryQuery.FindOneQueueByUserId(ctx, payload.UserId, payload.EventId)
	if queueRoom.Error != nil {
		msg := "Error DB connection FindOneQueueByUserId"
		c.logger.Error(ctx, msg, fmt.Sprintf("%+v", queueRoom.Error))
		return nil, queueRoom.Error
	}

	if queueRoom.Data != nil {
		msg := "user already in the queue"
		c.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		return nil, errors.BadRequest("user already in the queue")
	}

	lastQueue := <-c.roomRepositoryQuery.FindOneLastQueue(ctx, payload.EventId)
	if lastQueue.Error != nil {
		msg := "Error DB connection FindOneLastQueue"
		c.logger.Error(ctx, msg, fmt.Sprintf("%+v", lastQueue.Error))
		return nil, lastQueue.Error
	}

	state := 1

	if lastQueue.Data != nil {
		queue, ok := lastQueue.Data.(*entity.QueueRoom)
		if !ok {
			msg := "cannot parsing data last queue"
			c.logger.Error(ctx, msg, fmt.Sprintf("%+v", lastQueue.Data))
			return nil, errors.InternalServerError("cannot parsing data")
		}
		state = queue.QueueNumber + 1
	}

	var queueLimit int

	checkedLimit, _ := c.redis.Get(ctx, fmt.Sprintf("%s:%s:%s:%s", constants.ORDER, constants.QueueLimit, event.EventId, event.Tag)).Result()
	if checkedLimit == "" {
		totalTicket := <-c.ticketRepositoryQuery.FindTotalAvalailableTicket(ctx, event.Country.Code, event.Tag)
		if totalTicket.Error != nil {
			msg := "Error DB connection FindTotalAvalailableTicket"
			c.logger.Error(ctx, msg, fmt.Sprintf("%+v", totalTicket.Error))
			return nil, totalTicket.Error
		}

		aggregateTicket, ok := totalTicket.Data.(*[]ticketEntity.AggregateTotalTicket)
		if !ok {
			msg := "cannot parsing data total ticket"
			c.logger.Error(ctx, msg, fmt.Sprintf("%+v", totalTicket.Data))
			return nil, errors.InternalServerError("cannot parsing data")
		}

		availableTicket := (*aggregateTicket)[0]

		quartal := helpers.GetCurrentQuartal()
		if quartal != helpers.Q4 {
			queueLimit = availableTicket.TotalAvailableTicket / 4
		} else {
			queueLimit = availableTicket.TotalAvailableTicket
		}
		c.redis.Set(ctx, fmt.Sprintf("%s:%s:%s:%s", constants.ORDER, constants.QueueLimit, event.EventId, event.Tag), queueLimit, 4*30*24*time.Hour)
	} else {
		limit, err := strconv.Atoi(checkedLimit)
		queueLimit = limit
		if err != nil {
			msg := "cannot parsing redis data"
			c.logger.Error(ctx, msg, fmt.Sprintf("%+v", checkedLimit))
			return nil, errors.InternalServerError("cannot parsing redis data")
		}
	}

	if state > queueLimit {
		msg := "queue is full"
		c.logger.Error(ctx, msg, fmt.Sprintf("%+v", payload))
		return nil, errors.BadRequest("queue is full")
	}

	data := entity.QueueRoom{
		UserId:      payload.UserId,
		EventId:     event.EventId,
		QueueNumber: state,
		CountryCode: event.Country.Code,
	}

	respQueue := <-c.roomRepositoryCommand.InsertOneRoom(ctx, data)
	if respQueue.Error != nil {
		msg := "Error DB connection InsertOneRoom"
		c.logger.Error(ctx, msg, fmt.Sprintf("%+v", respQueue.Error))
		return nil, respQueue.Error
	}

	return &response.QueueResp{
		UserId:      data.UserId,
		QueueNumber: data.QueueNumber,
		CountryCode: data.CountryCode,
	}, nil
}
