package converter

import (
	"fmt"
	"slices"

	"github.com/milkrage/microservices-course-homework/inventory/internal/model"
	repoModel "github.com/milkrage/microservices-course-homework/inventory/internal/repository/model"
)

func PartToModel(in repoModel.Part) (model.Part, error) {
	var dimensions *model.PartDimensions
	if in.Dimensions != nil {
		dimensions = &model.PartDimensions{
			Height: in.Dimensions.Height,
			Length: in.Dimensions.Length,
			Weight: in.Dimensions.Weight,
			Width:  in.Dimensions.Width,
		}
	}

	var manufacturer *model.PartManufacturer
	if in.Manufacturer != nil {
		manufacturer = &model.PartManufacturer{
			Name:    in.Manufacturer.Name,
			Country: in.Manufacturer.Country,
			Website: in.Manufacturer.Website,
		}
	}

	var metadata map[string]model.PartMetadataValue
	if len(in.Metadata) > 0 {
		metadata = make(map[string]model.PartMetadataValue, len(in.Metadata))

		for k, v := range in.Metadata {
			metadataValue, err := model.NewPartMetadataValue(v.Get())
			if err != nil {
				return model.Part{}, fmt.Errorf("failed to convert metadata value: %w", err)
			}

			metadata[k] = metadataValue
		}
	}

	out := model.Part{
		ID:            in.ID,
		Name:          in.Name,
		Description:   in.Description,
		Price:         in.Price,
		StockQuantity: in.StockQuantity,
		Category:      model.NewPartCategory(in.Category.Number(), in.Category.Name()),
		Dimensions:    dimensions,
		Manufacturer:  manufacturer,
		Tags:          slices.Clone(in.Tags),
		Metadata:      metadata,
	}

	return out, nil
}

func PartsToModel(in []repoModel.Part) ([]model.Part, error) {
	parts := make([]model.Part, 0, len(in))

	for _, v := range in {
		part, err := PartToModel(v)
		if err != nil {
			return nil, err
		}

		parts = append(parts, part)
	}

	return parts, nil
}

func PartFilterToRepoModel(filter model.PartFilter) repoModel.PartFilter {
	categories := make([]repoModel.PartCategory, 0, len(filter.Categories))

	for _, v := range filter.Categories {
		categories = append(categories, repoModel.NewPartCategory(v.Number(), v.Name()))
	}

	return repoModel.PartFilter{
		IDs:                   slices.Clone(filter.IDs),
		Names:                 slices.Clone(filter.Names),
		Categories:            categories,
		ManufacturerCountries: slices.Clone(filter.ManufacturerCountries),
		Tags:                  slices.Clone(filter.Tags),
	}
}
