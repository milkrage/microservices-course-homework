package repository

import (
	"context"

	"github.com/milkrage/microservices-course-homework/order/internal/model"
)

type OrderRepository interface {
	Get(ctx context.Context, orderID string) (model.Order, error)
	Upsert(ctx context.Context, order model.Order) error
}
