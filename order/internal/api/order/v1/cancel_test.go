package v1_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	v1 "github.com/milkrage/microservices-course-homework/order/internal/api/order/v1"
	"github.com/milkrage/microservices-course-homework/order/internal/model"
	serviceMocks "github.com/milkrage/microservices-course-homework/order/internal/service/mocks"
	orderV1 "github.com/milkrage/microservices-course-homework/shared/pkg/openapi/order/v1"
)

func TestCancelOrder(t *testing.T) {
	testCases := []struct {
		name           string
		serviceError   error
		respStatusCode int
	}{
		{name: "order not cancelable", serviceError: model.ErrOrderNotCancelable, respStatusCode: http.StatusConflict},
		{name: "order not found", serviceError: model.ErrOrderNotFound, respStatusCode: http.StatusNotFound},
		{name: "other errors", serviceError: gofakeit.Error(), respStatusCode: http.StatusInternalServerError},
		{name: "success", serviceError: nil, respStatusCode: http.StatusNoContent},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			service := serviceMocks.NewOrderService(t)
			api := v1.NewOrderHandler(service)
			orderID := gofakeit.UUID()

			service.On("Cancel", mock.Anything, orderID).Return(test.serviceError)

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/orders/%s/cancel", orderID), nil)
			rr := httptest.NewRecorder()

			srv, _ := orderV1.NewServer(api)
			srv.ServeHTTP(rr, req)

			require.Equal(t, test.respStatusCode, rr.Code)
		})
	}
}
