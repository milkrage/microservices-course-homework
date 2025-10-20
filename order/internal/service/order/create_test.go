package order_test

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"

	"github.com/milkrage/microservices-course-homework/order/internal/model"
)

func (ts *TestSuite) TestCreateWithoutAllParts() {
	ts.inventoryClient.On("ListParts", ts.ctx, mock.Anything).Return(model.ListParts{Parts: []model.Part{{ID: "1"}}}, nil)

	order, err := ts.service.Create(ts.ctx, gofakeit.UUID(), []string{"1", "2"})

	ts.Require().Equal(model.Order{}, order)
	ts.Require().ErrorIs(err, model.ErrNotFoundAllParts)
}

func (ts *TestSuite) TestCreateSuccess() {
	partIDs := []string{"1"}
	userID := gofakeit.UUID()

	ts.inventoryClient.On("ListParts", ts.ctx, mock.Anything).Return(model.ListParts{Parts: []model.Part{{ID: "1"}}}, nil)
	ts.orderRepository.On("Upsert", ts.ctx, mock.Anything).Return(nil)

	order, err := ts.service.Create(ts.ctx, userID, partIDs)

	ts.Require().NoError(err)
	ts.Require().Equal(model.OrderStatusPending, order.Status)
	ts.Require().Equal(partIDs, order.PartIDs)
	ts.Require().Equal(userID, order.UserID)
}
