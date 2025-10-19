package v1

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/milkrage/microservices-course-homework/inventory/internal/converter"
	inventoryV1 "github.com/milkrage/microservices-course-homework/shared/pkg/proto/inventory/v1"
)

func (p *PartHandler) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	parts, err := p.partService.List(ctx, converter.PartFilterToModel(req.Filter))
	if err != nil {
		log.Printf("failed to list parts: %v", err)
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	return &inventoryV1.ListPartsResponse{Parts: converter.PartsToProto(parts)}, nil
}
