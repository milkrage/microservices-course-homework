package order

import (
	"context"
	"slices"
	"sync"

	"github.com/samber/lo"

	"github.com/milkrage/microservices-course-homework/order/internal/model"
	"github.com/milkrage/microservices-course-homework/order/internal/repository/converter"
	repoModel "github.com/milkrage/microservices-course-homework/order/internal/repository/model"
)

type Repository struct {
	storage map[string]repoModel.Order
	mu      sync.RWMutex
}

func NewOrderRepository() *Repository {
	return &Repository{
		storage: make(map[string]repoModel.Order),
	}
}

func (r *Repository) Get(_ context.Context, orderID string) (model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.storage[orderID]
	if !ok {
		return model.Order{}, model.ErrOrderNotFound
	}

	var transactionUUID *string
	if order.TransactionID != nil {
		transactionUUID = lo.ToPtr(*order.TransactionID)
	}

	var paymentMethod *repoModel.OrderPaymentMethod
	if order.PaymentMethod != nil {
		paymentMethod = lo.ToPtr(*order.PaymentMethod)
	}

	result := repoModel.Order{
		OrderID:       order.OrderID,
		UserID:        order.UserID,
		PartIDs:       slices.Clone(order.PartIDs),
		TotalPrice:    order.TotalPrice,
		TransactionID: transactionUUID,
		PaymentMethod: paymentMethod,
		Status:        order.Status,
	}

	return converter.OrderToModel(result), nil
}

func (r *Repository) Upsert(_ context.Context, order model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.storage[order.OrderID] = converter.OrderToRepo(order)

	return nil
}
