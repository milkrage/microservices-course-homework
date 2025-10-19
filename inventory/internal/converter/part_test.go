package converter_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/milkrage/microservices-course-homework/inventory/internal/converter"
	"github.com/milkrage/microservices-course-homework/inventory/internal/model"
)

func TestPartToProto(t *testing.T) {
	metavalue, err := model.NewPartMetadataValue("32")
	require.NoError(t, err)

	in := model.Part{
		ID:            "id",
		Name:          "name",
		Description:   "desc",
		Price:         123.4,
		StockQuantity: 567,
		Category:      model.NewPartCategory(2, "FUEL"),
		Dimensions:    &model.PartDimensions{Width: 1, Height: 2, Length: 3, Weight: 4},
		Manufacturer:  &model.PartManufacturer{Name: "1", Country: "2", Website: "3"},
		Tags:          []string{"1"},
		Metadata:      map[string]model.PartMetadataValue{"1": metavalue},
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	out := converter.PartToProto(in)
	require.Equal(t, in.ID, out.Uuid)
	require.Equal(t, in.Name, out.Name)
	require.Equal(t, in.Description, out.Description)
	require.Equal(t, in.Price, out.Price)
	require.Equal(t, in.StockQuantity, out.StockQuantity)
	require.Equal(t, in.Category.Name(), out.Category.String())
	require.Equal(t, in.Category.Number(), int32(out.Category.Number()))
}
