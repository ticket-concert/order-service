// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	context "context"
	helpers "order-service/internal/pkg/helpers"

	mock "github.com/stretchr/testify/mock"

	request "order-service/internal/modules/order/models/request"
)

// MongodbRepositoryQuery is an autogenerated mock type for the MongodbRepositoryQuery type
type MongodbRepositoryQuery struct {
	mock.Mock
}

// FindBankTicketByParam provides a mock function with given fields: ctx, queueId, userId
func (_m *MongodbRepositoryQuery) FindBankTicketByParam(ctx context.Context, queueId string, userId string) <-chan helpers.Result {
	ret := _m.Called(ctx, queueId, userId)

	if len(ret) == 0 {
		panic("no return value specified for FindBankTicketByParam")
	}

	var r0 <-chan helpers.Result
	if rf, ok := ret.Get(0).(func(context.Context, string, string) <-chan helpers.Result); ok {
		r0 = rf(ctx, queueId, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan helpers.Result)
		}
	}

	return r0
}

// FindBankTicketByUser provides a mock function with given fields: ctx, payload
func (_m *MongodbRepositoryQuery) FindBankTicketByUser(ctx context.Context, payload request.PreOrderList) <-chan helpers.Result {
	ret := _m.Called(ctx, payload)

	if len(ret) == 0 {
		panic("no return value specified for FindBankTicketByUser")
	}

	var r0 <-chan helpers.Result
	if rf, ok := ret.Get(0).(func(context.Context, request.PreOrderList) <-chan helpers.Result); ok {
		r0 = rf(ctx, payload)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan helpers.Result)
		}
	}

	return r0
}

// FindOrderByUser provides a mock function with given fields: ctx, payload
func (_m *MongodbRepositoryQuery) FindOrderByUser(ctx context.Context, payload request.OrderList) <-chan helpers.Result {
	ret := _m.Called(ctx, payload)

	if len(ret) == 0 {
		panic("no return value specified for FindOrderByUser")
	}

	var r0 <-chan helpers.Result
	if rf, ok := ret.Get(0).(func(context.Context, request.OrderList) <-chan helpers.Result); ok {
		r0 = rf(ctx, payload)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan helpers.Result)
		}
	}

	return r0
}

// NewMongodbRepositoryQuery creates a new instance of MongodbRepositoryQuery. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMongodbRepositoryQuery(t interface {
	mock.TestingT
	Cleanup(func())
}) *MongodbRepositoryQuery {
	mock := &MongodbRepositoryQuery{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
