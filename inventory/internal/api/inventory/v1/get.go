package v1

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/milkrage/microservices-course-homework/inventory/internal/converter"
	"github.com/milkrage/microservices-course-homework/inventory/internal/model"
	inventoryV1 "github.com/milkrage/microservices-course-homework/shared/pkg/proto/inventory/v1"
)

func (p *PartHandler) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	part, err := p.partService.Get(ctx, req.Uuid)
	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			return nil, status.Errorf(codes.NotFound, "part with uuid %s not found", req.Uuid)
		}

		log.Printf("failed to get part, id: %s, err: %v", req.Uuid, err)

		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &inventoryV1.GetPartResponse{Part: converter.PartToProto(part)}, nil
}
