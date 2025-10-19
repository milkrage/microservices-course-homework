package service

import (
	"context"

	"github.com/milkrage/microservices-course-homework/inventory/internal/model"
)

type PartService interface {
	Get(ctx context.Context, partID string) (model.Part, error)
	List(ctx context.Context, filter model.PartFilter) ([]model.Part, error)
}
