package memory

import (
	"slices"
	"sync"
)

type Order struct {
	OrderUUID       string
	UserUUID        string
	PartUUIDs       []string
	TotalPrice      float64
	TransactionUUID *string
	PaymentMethod   *string
	Status          string
}

type OrderMemoryStorage struct {
	storage map[string]Order
	mu      sync.RWMutex
}

func NewOrderMemoryStorage() *OrderMemoryStorage {
	return &OrderMemoryStorage{
		storage: make(map[string]Order),
	}
}

func (o *OrderMemoryStorage) Get(uuid string) (Order, bool) {
	o.mu.RLock()
	order, ok := o.storage[uuid]
	o.mu.RUnlock()

	if !ok {
		return Order{}, false
	}

	var transactionUUID *string
	if order.TransactionUUID != nil {
		tmp := *order.TransactionUUID
		transactionUUID = &tmp
	}

	var paymentMethod *string
	if order.PaymentMethod != nil {
		tmp := *order.PaymentMethod
		paymentMethod = &tmp
	}

	result := Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUUIDs:       slices.Clone(order.PartUUIDs),
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          order.Status,
	}

	return result, true
}

func (o *OrderMemoryStorage) Upsert(order Order) {
	o.mu.Lock()
	o.storage[order.OrderUUID] = order
	o.mu.Unlock()
}
