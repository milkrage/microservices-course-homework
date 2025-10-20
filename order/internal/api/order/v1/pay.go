package v1

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/milkrage/microservices-course-homework/order/internal/model"
	orderV1 "github.com/milkrage/microservices-course-homework/shared/pkg/openapi/order/v1"
)

func (h *Handler) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	err := req.Validate()
	if err != nil {
		return &orderV1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}, nil
	}

	transactionID, err := h.service.Pay(ctx, params.OrderUUID, *model.NewOrderPaymentMethod(string(req.PaymentMethod)))
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: fmt.Sprintf("Order %s not found", params.OrderUUID),
			}, nil
		}

		log.Printf("failed to pay order: %v", err)

		return &orderV1.InternalServerError{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
		}, nil
	}

	return &orderV1.PayOrderResponse{TransactionUUID: orderV1.NewOptString(transactionID)}, nil
}
