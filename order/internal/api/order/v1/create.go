package v1

import (
	"context"
	"log"
	"net/http"

	"github.com/go-faster/errors"

	"github.com/milkrage/microservices-course-homework/order/internal/converter"
	"github.com/milkrage/microservices-course-homework/order/internal/model"
	orderV1 "github.com/milkrage/microservices-course-homework/shared/pkg/openapi/order/v1"
)

func (h *Handler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	err := req.Validate()
	if err != nil {
		return &orderV1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}, nil
	}

	order, err := h.service.Create(ctx, req.UserUUID, req.PartUuids)
	if err != nil {
		if errors.Is(err, model.ErrNotFoundAllParts) {
			return &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: "Not all requested parts were found.",
			}, nil
		}

		log.Printf("failed to create order: %v", err)

		return &orderV1.InternalServerError{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
		}, nil
	}

	return converter.OrderToCreateOrderResponse(order), nil
}
