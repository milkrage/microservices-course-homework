package order

import (
	"github.com/milkrage/microservices-course-homework/order/internal/client/grpc"
	"github.com/milkrage/microservices-course-homework/order/internal/repository"
)

type Service struct {
	orderRepository repository.OrderRepository
	paymentClient   grpc.PaymentClient
	inventoryClient grpc.InventoryClient
}

func NewOrderService(
	orderRepository repository.OrderRepository,
	paymentClient grpc.PaymentClient,
	inventoryClient grpc.InventoryClient,
) *Service {
	return &Service{
		orderRepository: orderRepository,
		paymentClient:   paymentClient,
		inventoryClient: inventoryClient,
	}
}
