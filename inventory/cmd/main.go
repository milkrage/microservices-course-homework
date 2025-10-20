package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	partHandler "github.com/milkrage/microservices-course-homework/inventory/internal/api/inventory/v1"
	partRepository "github.com/milkrage/microservices-course-homework/inventory/internal/repository/part"
	partService "github.com/milkrage/microservices-course-homework/inventory/internal/service/part"
	inventoryV1 "github.com/milkrage/microservices-course-homework/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50051

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	repository := partRepository.NewPartRepository()
	service := partService.NewPartService(repository)
	handler := partHandler.NewPartHandler(service)

	reflection.Register(s)
	inventoryV1.RegisterInventoryServiceServer(s, handler)

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
