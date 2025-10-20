package payment

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/milkrage/microservices-course-homework/payment/internal/model"
)

func (s *Service) Pay(_ context.Context, _, _, _ string) (model.Payment, error) {
	payment := model.Payment{TransactionID: uuid.New().String()}
	log.Printf("payment was successful, transaction_uuid: %s", payment.TransactionID)

	return payment, nil
}
