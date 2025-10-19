package order_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	grpcMocks "github.com/milkrage/microservices-course-homework/order/internal/client/grpc/mocks"
	repoMocks "github.com/milkrage/microservices-course-homework/order/internal/repository/mocks"
	orderService "github.com/milkrage/microservices-course-homework/order/internal/service/order"
)

type TestSuite struct {
	suite.Suite

	ctx context.Context // nolint: containedctx

	orderRepository *repoMocks.OrderRepository
	paymentClient   *grpcMocks.PaymentClient
	inventoryClient *grpcMocks.InventoryClient

	service *orderService.Service
}

func (ts *TestSuite) SetupTest() {
	ts.ctx = context.Background()
	ts.orderRepository = repoMocks.NewOrderRepository(ts.T())
	ts.paymentClient = grpcMocks.NewPaymentClient(ts.T())
	ts.inventoryClient = grpcMocks.NewInventoryClient(ts.T())
	ts.service = orderService.NewOrderService(ts.orderRepository, ts.paymentClient, ts.inventoryClient)
}

func (ts *TestSuite) TearDownTest() {
}

func TestOrderService(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
