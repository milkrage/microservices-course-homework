package v1

import (
	"github.com/milkrage/microservices-course-homework/inventory/internal/service"
	inventoryV1 "github.com/milkrage/microservices-course-homework/shared/pkg/proto/inventory/v1"
)

type PartHandler struct {
	partService service.PartService
	inventoryV1.UnimplementedInventoryServiceServer
}

func NewPartHandler(partService service.PartService) *PartHandler {
	return &PartHandler{
		partService: partService,
	}
}
