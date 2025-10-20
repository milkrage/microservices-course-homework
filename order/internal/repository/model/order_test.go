package model_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/milkrage/microservices-course-homework/order/internal/repository/model"
)

func TestNewOrderPaymentMethod(t *testing.T) {
	card := model.NewOrderPaymentMethod("CARD")
	require.Equal(t, int32(1), card.Number())

	unknown := model.NewOrderPaymentMethod("!CARD!")
	require.Equal(t, int32(0), unknown.Number())
}
