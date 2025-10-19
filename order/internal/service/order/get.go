package order

import (
	"context"

	"github.com/milkrage/microservices-course-homework/order/internal/model"
)

func (s *Service) Get(ctx context.Context, orderID string) (model.Order, error) {
	return s.orderRepository.Get(ctx, orderID)
}
