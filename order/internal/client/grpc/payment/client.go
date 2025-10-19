package payment

import (
	paymentV1 "github.com/milkrage/microservices-course-homework/shared/pkg/proto/payment/v1"
)

type Client struct {
	payment paymentV1.PaymentServiceClient
}

func NewPaymentClient(payment paymentV1.PaymentServiceClient) *Client {
	return &Client{
		payment: payment,
	}
}
