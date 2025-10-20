package v1

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-faster/errors"

	"github.com/milkrage/microservices-course-homework/order/internal/model"
	orderV1 "github.com/milkrage/microservices-course-homework/shared/pkg/openapi/order/v1"
)

func (h *Handler) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	err := h.service.Cancel(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotCancelable) {
			return &orderV1.ConflictError{
				Code:    http.StatusConflict,
				Message: fmt.Sprintf("Order %s has already been paid and cannot be cancelled", params.OrderUUID),
			}, nil
		}

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

	return &orderV1.CancelOrderNoContent{}, nil
}
