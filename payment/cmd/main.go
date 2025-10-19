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

	paymentHandler "github.com/milkrage/microservices-course-homework/payment/internal/api/payment/v1"
	paymentService "github.com/milkrage/microservices-course-homework/payment/internal/service/payment"
	paymentV1 "github.com/milkrage/microservices-course-homework/shared/pkg/proto/payment/v1"
)

const grpcPort = 50052

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	service := paymentService.NewPaymentService()
	handler := paymentHandler.NewPaymentHandler(service)

	reflection.Register(s)
	paymentV1.RegisterPaymentServiceServer(s, handler)

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
