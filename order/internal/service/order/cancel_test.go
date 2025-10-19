package order_test

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"

	"github.com/milkrage/microservices-course-homework/order/internal/model"
)

func (ts *TestSuite) TestOrderNotCancelable() {
	ts.orderRepository.On("Get", ts.ctx, mock.Anything).Return(model.Order{Status: "-"}, nil)

	err := ts.service.Cancel(ts.ctx, gofakeit.UUID())

	ts.Require().Error(err)
	ts.Require().ErrorIs(err, model.ErrOrderNotCancelable)
}

func (ts *TestSuite) TestCancelOrderSuccess() {
	ts.orderRepository.On("Get", ts.ctx, mock.Anything).Return(model.Order{Status: model.OrderStatusPending}, nil)
	ts.orderRepository.On("Upsert", ts.ctx, model.Order{Status: model.OrderStatusCanceled}).Return(nil)

	err := ts.service.Cancel(ts.ctx, gofakeit.UUID())

	ts.Require().NoError(err)
}
