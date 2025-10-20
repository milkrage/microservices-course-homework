package payment

import (
	"context"
	"fmt"
	"time"

	"github.com/milkrage/microservices-course-homework/order/internal/model"
	paymentV1 "github.com/milkrage/microservices-course-homework/shared/pkg/proto/payment/v1"
)

func (c *Client) Pay(ctx context.Context, userId, orderID string, paymentMethod model.OrderPaymentMethod) (string, error) {
	req := &paymentV1.PayOrderRequest{
		OrderUuid:     orderID,
		UserUuid:      userId,
		PaymentMethod: paymentV1.PaymentMethod(paymentMethod.Number()),
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	resp, err := c.payment.PayOrder(ctxWithTimeout, req)
	if err != nil {
		return "", fmt.Errorf("%w: %w", model.ErrInternal, err)
	}

	return resp.TransactionUuid, nil
}
