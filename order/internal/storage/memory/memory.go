package memory

import "sync"

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

	return order, ok
}

func (o *OrderMemoryStorage) Upsert(order Order) {
	o.mu.Lock()
	o.storage[order.OrderUUID] = order
	o.mu.Unlock()
}
