package handlers

import (
	"fmt"
	"order-service/internal/modules/room"
	"order-service/internal/modules/room/models/request"
	"order-service/internal/pkg/errors"
	"order-service/internal/pkg/helpers"
	"order-service/internal/pkg/log"
	"order-service/internal/pkg/redis"

	middlewares "order-service/configs/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type RoomHttpHandler struct {
	RoomUsecaseCommand room.UsecaseCommand
	Logger             log.Logger
	Validator          *validator.Validate
}

func InitRoomHttpHandler(app *fiber.App, ruc room.UsecaseCommand, log log.Logger, redisClient redis.Collections) {
	handler := &RoomHttpHandler{
		RoomUsecaseCommand: ruc,
		Logger:             log,
		Validator:          validator.New(),
	}
	middlewares := middlewares.NewMiddlewares(redisClient)
	route := app.Group("/api/room")

	route.Post("/v1/create-queue", middlewares.VerifyBearer(), handler.CreateQueueRoom)
}

func (t RoomHttpHandler) CreateQueueRoom(c *fiber.Ctx) error {
	req := new(request.QueueReq)
	if err := c.BodyParser(req); err != nil {
		return helpers.RespError(c, t.Logger, errors.BadRequest("bad request"))
	}

	userId := c.Locals("userId").(string)
	req.UserId = userId

	if err := t.Validator.Struct(req); err != nil {
		fmt.Println(err)
		return helpers.RespError(c, t.Logger, errors.BadRequest(err.Error()))
	}
	resp, err := t.RoomUsecaseCommand.CreateQueueRoom(c.Context(), *req)
	if err != nil {
		return helpers.RespCustomError(c, t.Logger, err)
	}
	return helpers.RespSuccess(c, t.Logger, resp, "Create queue room success")
}
