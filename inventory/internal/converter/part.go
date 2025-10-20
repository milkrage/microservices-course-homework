package converter

import (
	"slices"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/milkrage/microservices-course-homework/inventory/internal/model"
	inventoryV1 "github.com/milkrage/microservices-course-homework/shared/pkg/proto/inventory/v1"
)

func PartToProto(in model.Part) *inventoryV1.Part {
	var dimensions *inventoryV1.Dimensions
	if in.Dimensions != nil {
		dimensions = &inventoryV1.Dimensions{
			Length: in.Dimensions.Length,
			Width:  in.Dimensions.Width,
			Weight: in.Dimensions.Weight,
			Height: in.Dimensions.Height,
		}
	}

	var manufacturer *inventoryV1.Manufacturer
	if in.Manufacturer != nil {
		manufacturer = &inventoryV1.Manufacturer{
			Name:    in.Manufacturer.Name,
			Country: in.Manufacturer.Country,
			Website: in.Manufacturer.Website,
		}
	}

	var metadata map[string]*inventoryV1.Value
	if len(in.Metadata) > 0 {
		metadata = make(map[string]*inventoryV1.Value, len(in.Metadata))

		for k, v := range in.Metadata {
			metadata[k] = PartMetadataValueToProto(v)
		}
	}

	out := &inventoryV1.Part{
		Uuid:          in.ID,
		Name:          in.Name,
		Description:   in.Description,
		Price:         in.Price,
		StockQuantity: in.StockQuantity,
		Category:      inventoryV1.Category(in.Category.Number()),
		Dimensions:    dimensions,
		Manufacturer:  manufacturer,
		Tags:          slices.Clone(in.Tags),
		Metadata:      metadata,
		CreatedAt:     timestamppb.New(in.CreatedAt),
		UpdatedAt:     timestamppb.New(in.UpdatedAt),
	}

	return out
}

func PartMetadataValueToProto(in model.PartMetadataValue) *inventoryV1.Value {
	v := &inventoryV1.Value{}

	switch in.Get().(type) {
	case string:
		v.Value = &inventoryV1.Value_StringValue{StringValue: in.Get().(string)}
	case int64:
		v.Value = &inventoryV1.Value_Int64Value{Int64Value: in.Get().(int64)}
	case float64:
		v.Value = &inventoryV1.Value_DoubleValue{DoubleValue: in.Get().(float64)}
	case bool:
		v.Value = &inventoryV1.Value_BoolValue{BoolValue: in.Get().(bool)}
	}

	return v
}

func PartsToProto(in []model.Part) []*inventoryV1.Part {
	parts := make([]*inventoryV1.Part, 0, len(in))

	for _, v := range in {
		parts = append(parts, PartToProto(v))
	}

	return parts
}

func PartFilterToModel(in *inventoryV1.PartsFilter) model.PartFilter {
	if in == nil {
		return model.PartFilter{}
	}

	categories := make([]model.PartCategory, 0, len(in.Categories))

	for _, v := range in.Categories {
		categories = append(categories, model.NewPartCategory(int32(v.Number()), v.String()))
	}

	out := model.PartFilter{
		IDs:                   slices.Clone(in.Uuids),
		Names:                 slices.Clone(in.Names),
		Categories:            categories,
		ManufacturerCountries: slices.Clone(in.ManufacturerCountries),
		Tags:                  slices.Clone(in.Tags),
	}

	return out
}
