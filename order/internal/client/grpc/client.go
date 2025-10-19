package grpc

import (
	"context"

	"github.com/milkrage/microservices-course-homework/order/internal/model"
)

type PaymentClient interface {
	Pay(ctx context.Context, userId, OrderID string, paymentMethod model.OrderPaymentMethod) (string, error)
}

type InventoryClient interface {
	ListParts(ctx context.Context, filter model.PartFilter) (model.ListParts, error)
}
