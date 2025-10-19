package converter

import (
	"github.com/milkrage/microservices-course-homework/order/internal/model"
	inventoryV1 "github.com/milkrage/microservices-course-homework/shared/pkg/proto/inventory/v1"
)

func PartFilterToListPartsRequest(in model.PartFilter) *inventoryV1.ListPartsRequest {
	return &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{
			Uuids: in.IDs,
		},
	}
}

func ListPartsResponseToModel(in *inventoryV1.ListPartsResponse) model.ListParts {
	out := model.ListParts{
		Parts: make([]model.Part, 0),
	}

	for _, part := range in.Parts {
		out.Parts = append(out.Parts, model.Part{ID: part.Uuid, Price: part.Price})
	}

	return out
}
