package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/milkrage/microservices-course-homework/inventory/internal/storage/memory"
	inventoryV1 "github.com/milkrage/microservices-course-homework/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50051

type inventoryService struct {
	inventoryV1.UnimplementedInventoryServiceServer
	storage *memory.InventoryStorage
}

func (i *inventoryService) GetPart(_ context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	part, ok := i.storage.GetPart(req.Uuid)
	if !ok {
		return nil, status.Errorf(codes.NotFound, "part with uuid %s not found", req.Uuid)
	}

	return &inventoryV1.GetPartResponse{Part: part}, nil
}

func (i *inventoryService) ListParts(_ context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	if req.Filter == nil {
		return nil, status.Errorf(codes.InvalidArgument, "filter is required key")
	}

	parts := i.storage.ListParts(req.Filter)

	return &inventoryV1.ListPartsResponse{Parts: parts}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	storage := memory.NewInventoryStorage()
	service := &inventoryService{storage: storage}

	reflection.Register(s)

	inventoryV1.RegisterInventoryServiceServer(s, service)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Fatalf("failed to serve: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
