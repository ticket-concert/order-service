package handlers

import (
	"fmt"
	"order-service/internal/modules/order"
	"order-service/internal/modules/order/models/request"
	"order-service/internal/pkg/errors"
	"order-service/internal/pkg/helpers"
	"order-service/internal/pkg/log"
	"order-service/internal/pkg/redis"

	middlewares "order-service/configs/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type OrderHttpHandler struct {
	OrderUsecaseCommand order.UsecaseCommand
	OrderUsecaseQuery   order.UsecaseQuery
	Logger              log.Logger
	Validator           *validator.Validate
}

func InitOrderHttpHandler(app *fiber.App, ouc order.UsecaseCommand, ouq order.UsecaseQuery, log log.Logger, redisClient redis.Collections) {
	handler := &OrderHttpHandler{
		OrderUsecaseCommand: ouc,
		OrderUsecaseQuery:   ouq,
		Logger:              log,
		Validator:           validator.New(),
	}
	middlewares := middlewares.NewMiddlewares(redisClient)
	route := app.Group("/api/order")

	route.Post("/v1/create-order", middlewares.VerifyBearer(), handler.CreateOrder)
	route.Get("/v1/list", middlewares.VerifyBearer(), handler.GetOrderList)
	route.Get("/v1/preorder-list", middlewares.VerifyBearer(), handler.GetPreOrderList)
}

func (t OrderHttpHandler) CreateOrder(c *fiber.Ctx) error {
	req := new(request.OrderReq)
	if err := c.BodyParser(req); err != nil {
		return helpers.RespError(c, t.Logger, errors.BadRequest("bad request"))
	}

	userId := c.Locals("userId").(string)
	req.UserId = userId

	if err := t.Validator.Struct(req); err != nil {
		fmt.Println(err)
		return helpers.RespError(c, t.Logger, errors.BadRequest(err.Error()))
	}
	resp, err := t.OrderUsecaseCommand.CreateOrderTicket(c.Context(), *req)
	if err != nil {
		return helpers.RespCustomError(c, t.Logger, err)
	}
	return helpers.RespSuccess(c, t.Logger, resp, "Create order success")
}

func (t OrderHttpHandler) GetOrderList(c *fiber.Ctx) error {
	req := new(request.OrderList)
	if err := c.QueryParser(req); err != nil {
		return helpers.RespError(c, t.Logger, errors.BadRequest("bad request"))
	}

	userId := c.Locals("userId").(string)
	req.UserId = userId
	if err := t.Validator.Struct(req); err != nil {
		return helpers.RespError(c, t.Logger, errors.BadRequest(err.Error()))
	}

	resp, err := t.OrderUsecaseQuery.FindOrderList(c.Context(), *req)
	if err != nil {
		return helpers.RespCustomError(c, t.Logger, err)
	}
	return helpers.RespPagination(c, t.Logger, resp.CollectionData, resp.MetaData, "Get order list success")
}

func (t OrderHttpHandler) GetPreOrderList(c *fiber.Ctx) error {
	req := new(request.PreOrderList)
	if err := c.QueryParser(req); err != nil {
		return helpers.RespError(c, t.Logger, errors.BadRequest("bad request"))
	}

	userId := c.Locals("userId").(string)
	req.UserId = userId
	if err := t.Validator.Struct(req); err != nil {
		return helpers.RespError(c, t.Logger, errors.BadRequest(err.Error()))
	}

	resp, err := t.OrderUsecaseQuery.FindPreOrderList(c.Context(), *req)
	if err != nil {
		return helpers.RespCustomError(c, t.Logger, err)
	}
	return helpers.RespPagination(c, t.Logger, resp.CollectionData, resp.MetaData, "Get preorder list success")
}
