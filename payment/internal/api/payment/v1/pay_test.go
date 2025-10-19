package v1_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	v1 "github.com/milkrage/microservices-course-homework/payment/internal/api/payment/v1"
	"github.com/milkrage/microservices-course-homework/payment/internal/model"
	"github.com/milkrage/microservices-course-homework/payment/internal/service/mocks"
	paymentV1 "github.com/milkrage/microservices-course-homework/shared/pkg/proto/payment/v1"
)

type TestSuite struct {
	suite.Suite

	ctx     context.Context // nolint: containedctx
	service *mocks.PaymentService
	api     *v1.PaymentHandler
}

func TestPaymentHandler(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (ts *TestSuite) SetupTest() {
	ts.ctx = context.Background()
	ts.service = mocks.NewPaymentService(ts.T())
	ts.api = v1.NewPaymentHandler(ts.service)
}

func (ts *TestSuite) TearDownTest() {
}

func (ts *TestSuite) TestPayOrderSuccess() {
	transactionID := uuid.NewString()

	ts.service.
		On("Pay", context.Background(), mock.Anything, mock.Anything, mock.Anything).
		Return(model.Payment{TransactionID: transactionID}, nil)

	resp, err := ts.api.PayOrder(ts.ctx, &paymentV1.PayOrderRequest{
		OrderUuid:     gofakeit.UUID(),
		UserUuid:      gofakeit.UUID(),
		PaymentMethod: paymentV1.PaymentMethod_CARD,
	})

	ts.Require().NoError(err)
	ts.Require().Equal(transactionID, resp.TransactionUuid)
}

func (ts *TestSuite) TestPayOrderServiceError() {
	expected := status.Errorf(codes.Internal, "internal service error")

	ts.service.
		On("Pay", context.Background(), mock.Anything, mock.Anything, mock.Anything).
		Return(model.Payment{}, errors.New("random error"))

	resp, err := ts.api.PayOrder(ts.ctx, &paymentV1.PayOrderRequest{
		OrderUuid:     gofakeit.UUID(),
		UserUuid:      gofakeit.UUID(),
		PaymentMethod: paymentV1.PaymentMethod_CARD,
	})

	ts.Require().Nil(resp)
	ts.Require().Error(err)
	ts.Require().Equal(expected, err)
}
