package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/milkrage/microservices-course-homework/order/internal/storage/memory"
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

// TODO(milkrage): remove global variables in week2.
var paymentMethodsMap = map[string]int32{
	"UNKNOWN":        0,
	"CARD":           1,
	"SBP":            2,
	"CREDIT_CARD":    3,
	"INVESTOR_MONEY": 4,
}

type OderHandler struct {
	payment   paymentV1.PaymentServiceClient
	inventory inventoryV1.InventoryServiceClient
	storage   *memory.OrderMemoryStorage
}

func NewOrderHandler(
	payment paymentV1.PaymentServiceClient,
	inventory inventoryV1.InventoryServiceClient,
	storage *memory.OrderMemoryStorage,
) *OderHandler {
	return &OderHandler{
		payment:   payment,
		inventory: inventory,
		storage:   storage,
	}
}

func (o *OderHandler) GetOrder(_ context.Context, params orderV1.GetOrderParams) (orderV1.GetOrderRes, error) {
	order, ok := o.storage.Get(params.OrderUUID)
	if !ok {
		return &orderV1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf("Order %s not found", params.OrderUUID),
		}, nil
	}

	resp := &orderV1.GetOrderResponse{
		OrderUUID:  orderV1.NewOptString(order.OrderUUID),
		UserUUID:   orderV1.NewOptString(order.UserUUID),
		PartUuids:  order.PartUUIDs,
		TotalPrice: orderV1.NewOptFloat64(order.TotalPrice),
		Status:     orderV1.NewOptOrderStatus(orderV1.OrderStatus(order.Status)),
	}

	if order.TransactionUUID != nil {
		resp.TransactionUUID = orderV1.NewOptString(*order.TransactionUUID)
	}

	if order.PaymentMethod != nil {
		resp.PaymentMethod = orderV1.NewOptPaymentMethod(orderV1.PaymentMethod(*order.PaymentMethod))
	}

	return resp, nil
}

func (o *OderHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	err := req.Validate()
	if err != nil {
		return &orderV1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}, nil
	}

	parts, err := o.inventory.ListParts(ctx, &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{Uuids: req.PartUuids},
	})
	if err != nil {
		log.Printf("failed to list parts via grpc client: %v", err)

		return &orderV1.InternalServerError{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
		}, nil
	}

	inventoryUUIDs := make([]string, 0, len(parts.Parts))

	// Extract UUIDs of parts from response (inventory service)
	for _, part := range parts.Parts {
		inventoryUUIDs = append(inventoryUUIDs, part.Uuid)
	}

	notFoundUUIDs := make([]string, 0)

	// Verify that all UUIDs from the current request are in the response UUID from the inventory service.
	// We do not use a length check, as the request may contain duplicates.
	for _, partUUID := range req.PartUuids {
		if !slices.Contains(inventoryUUIDs, partUUID) {
			notFoundUUIDs = append(notFoundUUIDs, partUUID)
		}
	}

	if len(notFoundUUIDs) > 0 {
		return &orderV1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf("The following parts were not found: %v", notFoundUUIDs),
		}, nil
	}

	var totalPrice float64

	for _, part := range parts.Parts {
		totalPrice += part.Price
	}

	order := memory.Order{
		OrderUUID:       uuid.NewString(),
		UserUUID:        req.UserUUID,
		PartUUIDs:       req.PartUuids,
		TotalPrice:      totalPrice,
		TransactionUUID: nil,
		PaymentMethod:   nil,
		Status:          string(orderV1.OrderStatusPENDINGPAYMENT),
	}

	o.storage.Upsert(order)

	return &orderV1.CreateOrderResponse{OrderUUID: orderV1.NewOptString(order.OrderUUID), TotalPrice: orderV1.NewOptFloat64(totalPrice)}, nil
}

func (o *OderHandler) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	err := req.Validate()
	if err != nil {
		return &orderV1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}, nil
	}

	order, ok := o.storage.Get(params.OrderUUID)
	if !ok {
		return &orderV1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf("Order %s not found", params.OrderUUID),
		}, nil
	}

	paymentMethod := string(req.PaymentMethod)

	paymentMethodCode, ok := paymentMethodsMap[paymentMethod]
	if !ok {
		// Neglecting proper error handling.
		paymentMethodCode = int32(0)
	}

	paymentOrder, err := o.payment.PayOrder(ctx, &paymentV1.PayOrderRequest{
		UserUuid:      order.UserUUID,
		OrderUuid:     order.OrderUUID,
		PaymentMethod: paymentV1.PaymentMethod(paymentMethodCode),
	})
	if err != nil {
		log.Printf("failed to pay order via grpc client: %v", err)

		return &orderV1.InternalServerError{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
		}, nil
	}

	order.TransactionUUID = &paymentOrder.TransactionUuid
	order.PaymentMethod = &paymentMethod
	order.Status = string(orderV1.OrderStatusPAID)

	o.storage.Upsert(order)

	return &orderV1.PayOrderResponse{TransactionUUID: orderV1.NewOptString(*order.TransactionUUID)}, nil
}

func (o *OderHandler) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	order, ok := o.storage.Get(params.OrderUUID)
	if !ok {
		return &orderV1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf("Order %s not found", params.OrderUUID),
		}, nil
	}

	switch order.Status {
	case string(orderV1.OrderStatusPAID):
		return &orderV1.NotFoundError{
			Code:    http.StatusConflict,
			Message: fmt.Sprintf("Order %s has already been paid and cannot be cancelled", params.OrderUUID),
		}, nil
	case string(orderV1.OrderStatusPENDINGPAYMENT):
		order.Status = string(orderV1.OrderStatusCANCELLED)

		o.storage.Upsert(order)

		return &orderV1.CancelOrderNoContent{}, nil
	case string(orderV1.OrderStatusCANCELLED):
		return &orderV1.NotFoundError{
			Code:    http.StatusConflict,
			Message: fmt.Sprintf("Order %s has already been canceled", params.OrderUUID),
		}, nil
	default:
		log.Printf("ERROR (CancelOrder) unknown order status %s in order %s\n", order.Status, order.OrderUUID)

		return &orderV1.InternalServerError{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
		}, nil
	}
}

func (o *OderHandler) NewError(_ context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    orderV1.NewOptInt(http.StatusInternalServerError),
			Message: orderV1.NewOptString(err.Error()),
		},
	}
}

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

	storage := memory.NewOrderMemoryStorage()

	handlers := NewOrderHandler(
		paymentV1.NewPaymentServiceClient(paymentServiceClient),
		inventoryV1.NewInventoryServiceClient(inventoryServiceClient),
		storage,
	)

	orderServer, err := orderV1.NewServer(handlers)
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
