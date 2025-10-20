package order_test

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"

	"github.com/milkrage/microservices-course-homework/order/internal/model"
)

func (ts *TestSuite) TestPayOrderSuccess() {
	orderID := gofakeit.UUID()
	userID := gofakeit.UUID()
	paymentMethod := *model.NewOrderPaymentMethod("CARD")
	transactionID := gofakeit.UUID()

	checkStatusFn := func(arg model.Order) bool {
		return arg.Status == model.OrderStatusPaid
	}

	ts.orderRepository.On("Get", ts.ctx, orderID).Return(model.Order{OrderID: orderID, UserID: userID}, nil)
	ts.paymentClient.On("Pay", ts.ctx, userID, orderID, paymentMethod).Return(transactionID, nil)
	ts.orderRepository.On("Upsert", ts.ctx, mock.MatchedBy(checkStatusFn)).Return(nil)

	resp, err := ts.service.Pay(ts.ctx, orderID, paymentMethod)

	ts.Require().NoError(err)
	ts.Require().Equal(transactionID, resp)
}
