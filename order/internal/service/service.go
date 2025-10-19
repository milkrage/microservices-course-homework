package service

import (
	"context"

	"github.com/milkrage/microservices-course-homework/order/internal/model"
)

type OrderService interface {
	Get(ctx context.Context, orderID string) (model.Order, error)
	Create(ctx context.Context, userID string, partIDs []string) (model.Order, error)
	Pay(ctx context.Context, orderID string, paymentMethod model.OrderPaymentMethod) (string, error)
	Cancel(ctx context.Context, orderID string) error
}
