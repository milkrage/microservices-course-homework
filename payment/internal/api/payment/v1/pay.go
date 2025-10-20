package v1

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	paymentV1 "github.com/milkrage/microservices-course-homework/shared/pkg/proto/payment/v1"
)

func (h *PaymentHandler) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	payment, err := h.paymentService.Pay(ctx, req.OrderUuid, req.UserUuid, req.PaymentMethod.String())
	if err != nil {
		log.Printf("failed to pay order: %v", err)

		return nil, status.Errorf(codes.Internal, "internal service error")
	}

	return &paymentV1.PayOrderResponse{TransactionUuid: payment.TransactionID}, nil
}
