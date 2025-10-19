package converter

import (
	"slices"

	"github.com/samber/lo"

	"github.com/milkrage/microservices-course-homework/order/internal/model"
	repoModel "github.com/milkrage/microservices-course-homework/order/internal/repository/model"
)

func OrderToModel(in repoModel.Order) model.Order {
	var transactionID *string
	if in.TransactionID != nil {
		transactionID = lo.ToPtr(*in.TransactionID)
	}

	var paymentMethod *model.OrderPaymentMethod
	if in.PaymentMethod != nil {
		paymentMethod = model.NewOrderPaymentMethod(in.PaymentMethod.Name())
	}

	return model.Order{
		OrderID:       in.OrderID,
		UserID:        in.UserID,
		PartIDs:       slices.Clone(in.PartIDs),
		TotalPrice:    in.TotalPrice,
		TransactionID: transactionID,
		PaymentMethod: paymentMethod,
		Status:        in.Status,
	}
}

func OrderToRepo(in model.Order) repoModel.Order {
	var transactionID *string
	if in.TransactionID != nil {
		transactionID = lo.ToPtr(*in.TransactionID)
	}

	var paymentMethod *repoModel.OrderPaymentMethod
	if in.PaymentMethod != nil {
		paymentMethod = repoModel.NewOrderPaymentMethod(in.PaymentMethod.Name())
	}

	return repoModel.Order{
		OrderID:       in.OrderID,
		UserID:        in.UserID,
		PartIDs:       slices.Clone(in.PartIDs),
		TotalPrice:    in.TotalPrice,
		TransactionID: transactionID,
		PaymentMethod: paymentMethod,
		Status:        in.Status,
	}
}
