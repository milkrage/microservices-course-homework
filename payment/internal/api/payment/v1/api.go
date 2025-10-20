package v1

import (
	"github.com/milkrage/microservices-course-homework/payment/internal/service"
	paymentV1 "github.com/milkrage/microservices-course-homework/shared/pkg/proto/payment/v1"
)

type PaymentHandler struct {
	paymentService service.PaymentService
	paymentV1.UnimplementedPaymentServiceServer
}

func NewPaymentHandler(paymentService service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}
