package converter

import (
	"github.com/milkrage/microservices-course-homework/order/internal/model"
	orderV1 "github.com/milkrage/microservices-course-homework/shared/pkg/openapi/order/v1"
)

func OrderToCreateOrderResponse(in model.Order) orderV1.CreateOrderRes {
	return &orderV1.CreateOrderResponse{
		OrderUUID:  orderV1.NewOptString(in.OrderID),
		TotalPrice: orderV1.NewOptFloat64(in.TotalPrice),
	}
}

func OrderToGetOrderResponse(in model.Order) orderV1.GetOrderRes {
	out := &orderV1.GetOrderResponse{
		OrderUUID:  orderV1.NewOptString(in.OrderID),
		UserUUID:   orderV1.NewOptString(in.UserID),
		PartUuids:  in.PartIDs,
		TotalPrice: orderV1.NewOptFloat64(in.TotalPrice),
		Status:     orderV1.NewOptOrderStatus(orderV1.OrderStatus(in.Status)),
	}

	if in.TransactionID != nil {
		out.TransactionUUID = orderV1.NewOptString(*in.TransactionID)
	}

	if in.PaymentMethod != nil {
		out.PaymentMethod = orderV1.NewOptPaymentMethod(orderV1.PaymentMethod(in.PaymentMethod.Name()))
	}

	return out
}
