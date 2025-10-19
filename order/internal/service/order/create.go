package order

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/milkrage/microservices-course-homework/order/internal/model"
)

func (s *Service) Create(ctx context.Context, userID string, partIDs []string) (model.Order, error) {
	listParts, err := s.inventoryClient.ListParts(ctx, model.PartFilter{IDs: partIDs})
	if err != nil {
		return model.Order{}, err
	}

	if !s.allPartsExists(listParts.Parts, partIDs) {
		return model.Order{}, model.ErrNotFoundAllParts
	}

	order := model.Order{
		OrderID:       uuid.NewString(),
		UserID:        userID,
		PartIDs:       partIDs,
		TotalPrice:    s.calcTotalPrice(listParts.Parts),
		TransactionID: nil,
		PaymentMethod: nil,
		Status:        model.OrderStatusPending,
	}

	err = s.orderRepository.Upsert(ctx, order)
	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}

func (s *Service) allPartsExists(source []model.Part, subset []string) bool {
	collection := make([]string, 0, len(source))
	for _, part := range source {
		collection = append(collection, part.ID)
	}

	return lo.Every(collection, subset)
}

func (s *Service) calcTotalPrice(parts []model.Part) float64 {
	var totalPrice float64

	for _, part := range parts {
		totalPrice += part.Price
	}

	return totalPrice
}
