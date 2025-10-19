package v1

import (
	"context"
	"net/http"

	"github.com/milkrage/microservices-course-homework/order/internal/service"
	orderV1 "github.com/milkrage/microservices-course-homework/shared/pkg/openapi/order/v1"
)

type Handler struct {
	service service.OrderService
}

func NewOrderHandler(service service.OrderService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) NewError(_ context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    orderV1.NewOptInt(http.StatusInternalServerError),
			Message: orderV1.NewOptString(err.Error()),
		},
	}
}
