package event

import (
	"context"
	wrapper "order-service/internal/pkg/helpers"
)

type MongodbRepositoryQuery interface {
	FindEventById(ctx context.Context, eventId string) <-chan wrapper.Result
}
