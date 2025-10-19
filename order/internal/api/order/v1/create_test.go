package v1_test

import (
	"bytes"
	"encoding/json"
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

func TestCreateOrder(t *testing.T) {
	testCases := []struct {
		name           string
		serviceError   error
		respStatusCode int
	}{
		{name: "not found all parts", serviceError: model.ErrNotFoundAllParts, respStatusCode: http.StatusNotFound},
		{name: "other errors", serviceError: gofakeit.Error(), respStatusCode: http.StatusInternalServerError},
		{name: "success", serviceError: nil, respStatusCode: http.StatusOK},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			service := serviceMocks.NewOrderService(t)
			api := v1.NewOrderHandler(service)

			service.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(model.Order{}, test.serviceError)

			body, _ := json.Marshal(map[string]any{"user_uuid": "1", "part_uuids": []string{"1", "2"}})
			req := httptest.NewRequest(http.MethodPost, "/api/v1/orders", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			srv, _ := orderV1.NewServer(api)
			srv.ServeHTTP(rr, req)

			require.Equal(t, test.respStatusCode, rr.Code)
		})
	}
}
