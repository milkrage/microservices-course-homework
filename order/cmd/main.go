package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderHandler "github.com/milkrage/microservices-course-homework/order/internal/api/order/v1"
	"github.com/milkrage/microservices-course-homework/order/internal/client/grpc/inventory"
	"github.com/milkrage/microservices-course-homework/order/internal/client/grpc/payment"
	orderRepo "github.com/milkrage/microservices-course-homework/order/internal/repository/order"
	orderService "github.com/milkrage/microservices-course-homework/order/internal/service/order"
	orderV1 "github.com/milkrage/microservices-course-homework/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/milkrage/microservices-course-homework/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/milkrage/microservices-course-homework/shared/pkg/proto/payment/v1"
)

const (
	httpPort             = "8080"
	inventoryServicePort = "50051"
	paymentServicePort   = "50052"

	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func main() {
	inventoryServiceClient, err := grpc.NewClient(
		net.JoinHostPort("127.0.0.1", inventoryServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect to inventory service client: %v\n", err)
		return
	}
	defer func() {
		if err := inventoryServiceClient.Close(); err != nil {
			log.Printf("failed to close connect of inventory service client: %v", err)
		}
	}()

	paymentServiceClient, err := grpc.NewClient(
		net.JoinHostPort("127.0.0.1", paymentServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect to payment service client: %v\n", err)
		return
	}
	defer func() {
		if err := paymentServiceClient.Close(); err != nil {
			log.Printf("failed to close connect of payment service client: %v", err)
		}
	}()

	repository := orderRepo.NewOrderRepository()
	paymentClient := payment.NewPaymentClient(paymentV1.NewPaymentServiceClient(paymentServiceClient))
	inventoryClient := inventory.NewInventoryClient(inventoryV1.NewInventoryServiceClient(inventoryServiceClient))
	service := orderService.NewOrderService(repository, paymentClient, inventoryClient)
	handler := orderHandler.NewOrderHandler(service)

	orderServer, err := orderV1.NewServer(handler)
	if err != nil {
		log.Printf("failed to create OpenAPI Server: %v", err)
		return
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("🚀 HTTP server listening on %s\n", httpPort)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ failed to serve: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down HTTP server...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Failed to shutdown HTTP server: %v\n", err)
		return
	}

	log.Println("✅ Server stopped")
}
