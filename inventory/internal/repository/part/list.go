package part

import (
	"context"
	"fmt"
	"maps"
	"slices"

	"github.com/milkrage/microservices-course-homework/inventory/internal/model"
	"github.com/milkrage/microservices-course-homework/inventory/internal/repository/converter"
	repo "github.com/milkrage/microservices-course-homework/inventory/internal/repository/model"
)

func (r *Repository) List(ctx context.Context, partFilter model.PartFilter) ([]model.Part, error) {
	filter := converter.PartFilterToRepoModel(partFilter)

	r.mu.RLock()
	defer r.mu.RUnlock()

	// Return all parts if all filters is empty.
	if r.isEmptyFilter(filter) {
		parts, err := converter.PartsToModel(slices.Collect(maps.Values(r.storage)))
		if err != nil {
			return nil, fmt.Errorf("failed to convert part to model: %w", err)
		}

		return parts, nil
	}

	// Find parts by filter.
	parts, err := converter.PartsToModel(r.findParts(filter))
	if err != nil {
		return nil, fmt.Errorf("failed to convert part to model: %w", err)
	}

	return parts, nil
}

func (r *Repository) isEmptyFilter(f repo.PartFilter) bool {
	return len(f.IDs)+len(f.Names)+len(f.Categories)+len(f.ManufacturerCountries)+len(f.Tags) == 0
}

// findParts - find parts by filter.
// WARNING: This method does not use a mutex. You must wrap calls to this method with a mutex.
func (r *Repository) findParts(filter repo.PartFilter) []repo.Part {
	result := make([]repo.Part, 0)

	for part := range maps.Values(r.storage) {
		if slices.Contains(filter.IDs, part.ID) ||
			slices.Contains(filter.Categories, part.Category) ||
			slices.Contains(filter.ManufacturerCountries, part.Manufacturer.Country) ||
			slices.Contains(filter.Names, part.Name) {
			result = append(result, part)
			continue
		}

		for _, tag := range filter.Tags {
			if slices.Contains(part.Tags, tag) {
				result = append(result, part)
				break
			}
		}
	}

	return result
}
