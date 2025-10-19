package part

import (
	"context"

	"github.com/milkrage/microservices-course-homework/inventory/internal/model"
)

func (s *Service) Get(ctx context.Context, partID string) (model.Part, error) {
	return s.partRepository.Get(ctx, partID)
}
