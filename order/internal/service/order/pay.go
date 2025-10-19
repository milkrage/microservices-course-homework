package order

import (
	"context"

	"github.com/milkrage/microservices-course-homework/order/internal/model"
)

func (s *Service) Pay(ctx context.Context, orderID string, paymentMethod model.OrderPaymentMethod) (string, error) {
	order, err := s.orderRepository.Get(ctx, orderID)
	if err != nil {
		return "", err
	}

	transactionID, err := s.paymentClient.Pay(ctx, order.UserID, order.OrderID, paymentMethod)
	if err != nil {
		return "", err
	}

	order.TransactionID = &transactionID
	order.PaymentMethod = &paymentMethod
	order.Status = model.OrderStatusPaid

	err = s.orderRepository.Upsert(ctx, order)
	if err != nil {
		return "", err
	}

	return transactionID, nil
}
