package v1

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/milkrage/microservices-course-homework/order/internal/converter"
	"github.com/milkrage/microservices-course-homework/order/internal/model"
	orderV1 "github.com/milkrage/microservices-course-homework/shared/pkg/openapi/order/v1"
)

func (h *Handler) GetOrder(ctx context.Context, params orderV1.GetOrderParams) (orderV1.GetOrderRes, error) {
	order, err := h.service.Get(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: fmt.Sprintf("Order %s not found", params.OrderUUID),
			}, nil
		}

		log.Printf("failed to cancel order: %v", err)

		return &orderV1.InternalServerError{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
		}, nil
	}

	return converter.OrderToGetOrderResponse(order), nil
}
