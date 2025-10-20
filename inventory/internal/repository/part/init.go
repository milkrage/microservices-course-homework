package part

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"

	"github.com/milkrage/microservices-course-homework/inventory/internal/repository/model"
)

func (r *Repository) init() {
	r.mu.Lock()
	defer r.mu.Unlock()

	categories := []model.PartCategory{
		model.NewPartCategory(1, "ENGINE"),
		model.NewPartCategory(2, "FUEL"),
		model.NewPartCategory(3, "PORTHOLE"),
		model.NewPartCategory(4, "WING"),
	}

	for range 100 {
		tags := make([]string, 0)

		for range gofakeit.Number(1, 10) {
			tags = append(tags, gofakeit.Bird())
		}

		part := model.Part{
			ID:           uuid.NewString(),
			Category:     categories[gofakeit.Number(0, len(categories)-1)],
			Manufacturer: &model.PartManufacturer{Country: gofakeit.Country()},
			Name:         gofakeit.Name(),
			Tags:         tags,
			Price:        gofakeit.Price(1, 1000),
		}

		r.storage[part.ID] = part
	}
}
