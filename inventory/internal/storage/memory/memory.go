package memory

import (
	"maps"
	"slices"
	"sync"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	inventoryV1 "github.com/milkrage/microservices-course-homework/shared/pkg/proto/inventory/v1"
)

type InventoryStorage struct {
	storage map[string]*inventoryV1.Part
	mu      sync.RWMutex
}

func NewInventoryStorage() *InventoryStorage {
	inventoryStorage := &InventoryStorage{
		storage: generateData(10),
	}

	return inventoryStorage
}

func (i *InventoryStorage) GetPart(uuid string) (*inventoryV1.Part, bool) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	part, ok := i.storage[uuid]
	if !ok {
		return nil, false
	}

	return part, true
}

func (i *InventoryStorage) ListParts(filter *inventoryV1.PartsFilter) []*inventoryV1.Part {
	i.mu.RLock()
	defer i.mu.RUnlock()

	// Return all parts if all filters is empty.
	if i.isEmptyFilter(filter) {
		return slices.Collect(maps.Values(i.storage))
	}

	// This map is an implementation of a unique set (by keys).
	result := make(map[*inventoryV1.Part]any)

	for part := range maps.Values(i.storage) {
		if slices.Contains(filter.Uuids, part.Uuid) {
			result[part] = nil

			continue
		}

		if slices.Contains(filter.Categories, part.Category) {
			result[part] = nil

			continue
		}

		if slices.Contains(filter.ManufacturerCountries, part.Manufacturer.Country) {
			result[part] = nil

			continue
		}

		if slices.Contains(filter.Names, part.Name) {
			result[part] = nil

			continue
		}

		for _, tag := range filter.Tags {
			if slices.Contains(part.Tags, tag) {
				result[part] = nil

				// Exit from loop of tags.
				break
			}
		}
	}

	return slices.Collect(maps.Keys(result))
}

func (i *InventoryStorage) isEmptyFilter(filter *inventoryV1.PartsFilter) bool {
	filters := len(filter.Uuids) +
		len(filter.Names) +
		len(filter.Categories) +
		len(filter.ManufacturerCountries) +
		len(filter.Tags)

	if filters == 0 {
		return true
	}

	return false
}

func generateData(count int) map[string]*inventoryV1.Part {
	storage := make(map[string]*inventoryV1.Part, count)

	categories := map[int]inventoryV1.Category{
		1: inventoryV1.Category_ENGINE,
		2: inventoryV1.Category_FUEL,
		3: inventoryV1.Category_PORTHOLE,
		4: inventoryV1.Category_WING,
	}

	for range count {
		tags := make([]string, 0)

		for range gofakeit.Number(1, 10) {
			tags = append(tags, gofakeit.Bird())
		}

		part := &inventoryV1.Part{
			Uuid:         uuid.NewString(),
			Category:     categories[gofakeit.Number(1, 4)],
			Manufacturer: &inventoryV1.Manufacturer{Country: gofakeit.Country()},
			Name:         gofakeit.Name(),
			Tags:         tags,
			Price:        gofakeit.Price(1, 1000),
		}

		storage[part.Uuid] = part
	}

	return storage
}
