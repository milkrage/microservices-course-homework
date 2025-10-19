package service

import (
	"context"

	"github.com/milkrage/microservices-course-homework/payment/internal/model"
)

type PaymentService interface {
	Pay(ctx context.Context, orderID, userID, paymentMethod string) (model.Payment, error)
}
