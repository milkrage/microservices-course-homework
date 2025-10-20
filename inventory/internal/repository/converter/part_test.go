package converter_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/milkrage/microservices-course-homework/inventory/internal/model"
	"github.com/milkrage/microservices-course-homework/inventory/internal/repository/converter"
	repoModel "github.com/milkrage/microservices-course-homework/inventory/internal/repository/model"
)

func TestPartToModel(t *testing.T) {
	metadataV, err := repoModel.NewPartMetadataValue("dd")
	require.NoError(t, err)

	repoPart := repoModel.Part{
		ID:            "id",
		Name:          "name",
		Description:   "desc",
		Price:         1.234,
		StockQuantity: 567,
		Category:      repoModel.NewPartCategory(1, "cat"),
		Dimensions:    &repoModel.PartDimensions{Weight: 1, Width: 2, Height: 3, Length: 4},
		Manufacturer:  &repoModel.PartManufacturer{Name: "1", Country: "2", Website: "3"},
		Tags:          []string{"1", "2", "3"},
		Metadata:      map[string]repoModel.PartMetadataValue{"s": metadataV},
	}

	result, err := converter.PartToModel(repoPart)
	require.NoError(t, err)
	require.Equal(t, repoPart.ID, result.ID)
	require.Equal(t, repoPart.Name, result.Name)
	require.Equal(t, repoPart.Description, result.Description)
	require.Equal(t, repoPart.Price, result.Price)
	require.Equal(t, repoPart.StockQuantity, result.StockQuantity)
	require.Equal(t, repoPart.Category.Name(), result.Category.Name())
	require.Equal(t, repoPart.Category.Number(), result.Category.Number())
	require.Equal(t, repoPart.Dimensions.Weight, result.Dimensions.Weight)
	require.Equal(t, repoPart.Dimensions.Width, result.Dimensions.Width)
	require.Equal(t, repoPart.Dimensions.Height, result.Dimensions.Height)
	require.Equal(t, repoPart.Dimensions.Length, result.Dimensions.Length)
	require.Equal(t, repoPart.Manufacturer.Name, result.Manufacturer.Name)
	require.Equal(t, repoPart.Manufacturer.Country, result.Manufacturer.Country)
	require.Equal(t, repoPart.Manufacturer.Website, result.Manufacturer.Website)
	require.Equal(t, repoPart.Tags, result.Tags)
	require.Equal(t, repoPart.Metadata["s"].Get(), result.Metadata["s"].Get())
}

func TestPartFilterToRepoModel(t *testing.T) {
	in := model.PartFilter{
		IDs:                   []string{"1"},
		Names:                 []string{"2"},
		Categories:            []model.PartCategory{model.NewPartCategory(1, "2")},
		ManufacturerCountries: []string{"3"},
		Tags:                  []string{"4"},
	}

	out := converter.PartFilterToRepoModel(in)
	require.Equal(t, in.IDs, out.IDs)
	require.Equal(t, in.Names, out.Names)
	require.Equal(t, in.ManufacturerCountries, out.ManufacturerCountries)
	require.Equal(t, in.Tags, out.Tags)
	require.Equal(t, in.Categories[0].Name(), out.Categories[0].Name())
	require.Equal(t, in.Categories[0].Number(), out.Categories[0].Number())
}
