package part

import (
	"context"
	"fmt"

	"github.com/milkrage/microservices-course-homework/inventory/internal/model"
	"github.com/milkrage/microservices-course-homework/inventory/internal/repository/converter"
)

func (r *Repository) Get(_ context.Context, partID string) (model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	part, ok := r.storage[partID]
	if !ok {
		return model.Part{}, model.ErrPartNotFound
	}

	out, err := converter.PartToModel(part)
	if err != nil {
		return model.Part{}, fmt.Errorf("failed to convert part to model: %w", err)
	}

	return out, nil
}
