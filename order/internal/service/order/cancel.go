package order

import (
	"context"

	"github.com/milkrage/microservices-course-homework/order/internal/model"
)

func (s *Service) Cancel(ctx context.Context, orderID string) error {
	order, err := s.orderRepository.Get(ctx, orderID)
	if err != nil {
		return err
	}

	if !order.IsCancelable() {
		return model.ErrOrderNotCancelable
	}

	order.Status = model.OrderStatusCanceled

	err = s.orderRepository.Upsert(ctx, order)
	if err != nil {
		return err
	}

	return nil
}
