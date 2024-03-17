// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	context "context"
	entity "order-service/internal/modules/ticket/models/entity"
	helpers "order-service/internal/pkg/helpers"

	mock "github.com/stretchr/testify/mock"
)

// MongodbRepositoryCommand is an autogenerated mock type for the MongodbRepositoryCommand type
type MongodbRepositoryCommand struct {
	mock.Mock
}

// UpdateOneTicketDetail provides a mock function with given fields: ctx, payload
func (_m *MongodbRepositoryCommand) UpdateOneTicketDetail(ctx context.Context, payload entity.Ticket) <-chan helpers.Result {
	ret := _m.Called(ctx, payload)

	if len(ret) == 0 {
		panic("no return value specified for UpdateOneTicketDetail")
	}

	var r0 <-chan helpers.Result
	if rf, ok := ret.Get(0).(func(context.Context, entity.Ticket) <-chan helpers.Result); ok {
		r0 = rf(ctx, payload)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan helpers.Result)
		}
	}

	return r0
}

// NewMongodbRepositoryCommand creates a new instance of MongodbRepositoryCommand. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMongodbRepositoryCommand(t interface {
	mock.TestingT
	Cleanup(func())
}) *MongodbRepositoryCommand {
	mock := &MongodbRepositoryCommand{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}