package part

import (
	"context"

	"github.com/milkrage/microservices-course-homework/inventory/internal/model"
)

func (s *Service) List(ctx context.Context, filter model.PartFilter) ([]model.Part, error) {
	return s.partRepository.List(ctx, filter)
}
